package meta

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/stretchr/testify/require"
)

func Test_ModulesPackageName(t *testing.T) {
	r := require.New(t)
	tmp := os.TempDir()
	envy.Set(envy.GO111MODULE, "on")

	r.NoError(os.Chdir(tmp))

	tcases := []struct {
		Content     string
		PackageName string
	}{
		{"module github.com/wawandco/zekito", "github.com/wawandco/zekito"},
		{"module zekito", "zekito"},
		{"module gopkg.in/some/friday.v2", "gopkg.in/some/friday.v2"},
		{"", "zekito"},
	}

	for _, tcase := range tcases {
		envy.Set("GOPATH", tmp)

		t.Run(tcase.Content, func(st *testing.T) {
			r := require.New(st)

			r.NoError(ioutil.WriteFile("go.mod", []byte(tcase.Content), 0644))

			a := New(filepath.Join(tmp, "zekito"))
			r.Equal(tcase.PackageName, a.PackagePkg)
		})
	}
}

func Test_App_Encoding(t *testing.T) {
	r := require.New(t)

	a := New(".")
	bb := &bytes.Buffer{}
	r.NoError(a.Encode(bb))

	b := App{}

	r.NoError((&b).Decode(bb))

	r.Equal(a.String(), b.String())
}

func Test_App_IsZero(t *testing.T) {
	r := require.New(t)

	app := App{}
	r.True(app.IsZero())
	app = New(".")
	r.False(app.IsZero())
}

func Test_App_PackageRoot(t *testing.T) {
	r := require.New(t)

	app := App{}

	app.PackageRoot("foo.com/bar")
	r.Equal("foo.com/bar/actions", app.ActionsPkg)
	r.Equal("foo.com/bar/models", app.ModelsPkg)
	r.Equal("foo.com/bar/grifts", app.GriftsPkg)
}

func Test_App_HasNodeJsScript(t *testing.T) {
	r := require.New(t)

	const pJSON = `
{
    "name": "buffalo",
    "version": "1.0.0",
    "main": "index.js",
    "license": "MIT",
    "scripts": {
        "dev": "webpack --watch",
		"build": "webpack -p --progress"
    },
    "dependencies": {
    "bootstrap-sass": "~3.3.7",
    "font-awesome": "~4.7.0",
    "highlightjs": "^9.12.0",
    "jquery": "~3.2.1",
    "jquery-ujs": "~1.2.2"
  },
  "devDependencies": {
    "@babel/cli": "^7.0.0",
    "@babel/core": "^7.0.0",
    "@babel/preset-env": "^7.0.0",
    "babel-loader": "^8.0.0-beta.6",
    "copy-webpack-plugin": "~4.5.2",
    "css-loader": "~1.0.0",
    "expose-loader": "~0.7.5",
    "file-loader": "~2.0.0",
    "gopherjs-loader": "^0.0.1",
    "image-webpack-loader": "^4.5.0",
    "imagemin": "^6.0.0",
    "mini-css-extract-plugin": "^0.4.2",
    "node-sass": "~4.9.0",
    "npm-install-webpack-plugin": "4.0.5",
    "sass-loader": "~7.1.0",
    "style-loader": "~0.23.0",
    "uglifyjs-webpack-plugin": "~1.3.0",
    "url-loader": "~1.1.1",
    "webpack": "~4.5.0",
    "webpack-clean-obsolete-chunks": "^0.4.0",
    "webpack-cli": "3.1.0",
    "webpack-livereload-plugin": "2.1.1",
    "webpack-manifest-plugin": "~2.0.0"
  }
}
	`

	tmp, err := ioutil.TempDir("", "")
	r.NoError(err)
	defer os.RemoveAll(tmp)
	cd, err := os.Getwd()
	r.NoError(err)
	defer os.Chdir(cd)
	r.NoError(os.Chdir(tmp))
	r.NoError(ioutil.WriteFile("package.json", []byte(pJSON), 0644))
	a := New(".")

	s, err := a.NodeScript("dev")
	r.NoError(err)
	r.Equal("webpack --watch", s)
	s, err = a.NodeScript("test")
	r.EqualError(err, "node script test not found")

	r.NoError(os.Remove("package.json"))
	a = New(".")
	s, err = a.NodeScript("dev")
	r.EqualError(err, "package.json not found")
}
