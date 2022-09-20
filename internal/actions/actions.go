package actions

import (
	"context"
	"errors"

	"github.com/dop251/goja_nodejs/require"
	z_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
)

type Config struct {
	HTTP HTTPConfig
}

var (
	ErrHalt = errors.New("interrupt")
)

type jsAction func(parameter, parameter) error

func Run(ctx context.Context, script, name string, opts ...Option) error {
	config, err := prepareRun(ctx, script, opts)
	if err != nil {
		return err
	}

	var fn jsAction
	jsFn := config.vm.Get(name)
	if jsFn == nil {
		return errors.New("function not found")
	}
	err = config.vm.ExportTo(jsFn, &fn)
	if err != nil {
		return err
	}

	t := config.Start()
	defer func() {
		t.Stop()
	}()
	errCh := make(chan error)
	defer close(errCh)

	go func() {
		defer func() {
			r := recover()
			if r != nil && !config.allowedToFail {
				err, ok := r.(error)
				if !ok {
					e, ok := r.(string)
					if ok {
						err = errors.New(e)
					}
				}
				errCh <- err
				return
			}
		}()

		err = fn(config.ctxParam.parameter, config.apiParam.parameter)
		if err != nil && !config.allowedToFail {
			errCh <- err
			return
		}
		errCh <- nil
	}()
	return <-errCh
}

func prepareRun(ctx context.Context, script string, opts []Option) (*runConfig, error) {
	config := newRunConfig(ctx, opts...)
	if config.timeout == 0 {
		return nil, z_errs.ThrowInternal(nil, "ACTIO-uCpCx", "Errrors.Internal")
	}
	t := config.Prepare()
	defer func() {
		t.Stop()
	}()

	registry := new(require.Registry)
	registry.Enable(config.vm)

	for name, loader := range config.modules {
		registry.RegisterNativeModule(name, loader)
	}

	errCh := make(chan error)
	//load function in seperate go routine to recover panics
	go func() {
		defer func() {
			r := recover()
			if r != nil {
				errCh <- r.(error)
				return
			}
		}()
		_, err := config.vm.RunString(script)
		if err != nil {
			errCh <- err
			return
		}
		errCh <- nil
	}()
	return config, <-errCh
}

func ActionToOptions(a *query.Action) []Option {
	opts := make([]Option, 0, 1)
	if a.AllowedToFail {
		opts = append(opts, WithAllowedToFail())
	}
	return opts
}
