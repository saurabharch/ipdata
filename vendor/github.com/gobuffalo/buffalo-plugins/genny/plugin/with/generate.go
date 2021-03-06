package with

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/buffalo-plugins/genny/plugin"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/genny/new"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/pkg/errors"
)

func GenerateCmd(opts *plugin.Options) (*genny.Group, error) {
	gg := &genny.Group{}
	if err := opts.Validate(); err != nil {
		return gg, errors.WithStack(err)
	}

	g := genny.New()
	box := packr.New("./generate/templates", "./generate/templates")
	if err := g.Box(box); err != nil {
		return gg, errors.WithStack(err)
	}
	ctx := plush.NewContext()
	ctx.Set("opts", opts)
	g.Transformer(plushgen.Transformer(ctx))

	g.Transformer(genny.Replace("-shortName-", opts.ShortName))
	g.Transformer(genny.Dot())

	g.RunFn(func(r *genny.Runner) error {
		f, err := r.FindFile("cmd/available.go")
		if err != nil {
			return errors.WithStack(err)
		}
		const g = `Available.Add("generate", generateCmd)`
		const m = `Available.Mount(rootCmd)`
		body := strings.Replace(f.String(), m, fmt.Sprintf("\t%s\n%s", g, m), 1)
		return r.File(genny.NewFile(f.Name(), strings.NewReader(body)))
	})

	gg.Add(g)

	g, err := new.New(&new.Options{
		Name:   opts.ShortName,
		Prefix: "genny",
	})
	if err != nil {
		return gg, errors.WithStack(err)
	}
	gg.Add(g)

	return gg, nil
}
