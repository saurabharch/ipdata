package cmd

import (
	"context"

	"github.com/gobuffalo/buffalo-plugins/genny/install"
	"github.com/gobuffalo/buffalo-plugins/plugins"
	"github.com/gobuffalo/events"
	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

func Listen(e events.Event) error {
	if e.Kind != "buffalo:setup:started" {
		return nil
	}

	run := genny.WetRunner(context.Background())

	opts := &install.Options{}
	gg, err := install.New(opts)
	if err != nil {
		return errors.WithStack(err)
	}
	run.WithGroup(gg)
	payload := e.Payload
	payload["plugins"] = opts.Plugins
	events.EmitPayload(plugins.EvtSetupStarted, payload)
	if err := run.Run(); err != nil {
		events.EmitError(plugins.EvtSetupErr, err, payload)
		return errors.WithStack(err)
	}
	events.EmitPayload(plugins.EvtSetupFinished, payload)
	return nil
}
