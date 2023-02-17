package actions

import (
	"context"
	"errors"
	"fmt"

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

type jsAction func(fields, fields) error

func Run(ctx context.Context, ctxParam contextFields, apiParam apiFields, script, name string, opts ...Option) (err error) {
	var config *runConfig
	// only return the error if the function is not allowed to fail
	// in each case
	defer func() {
		if err != nil && config.allowedToFail {
			err = nil
		}
	}()
	config, err = prepareRun(ctx, ctxParam, apiParam, script, opts)

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

	return executeFn(config, fn)
}

func prepareRun(ctx context.Context, ctxParam contextFields, apiParam apiFields, script string, opts []Option) (config *runConfig, err error) {
	config = newRunConfig(ctx, opts...)
	if config.timeout == 0 {
		return nil, z_errs.ThrowInternal(nil, "ACTIO-uCpCx", "Errrors.Internal")
	}
	t := config.Prepare()
	defer func() {
		t.Stop()
	}()

	if ctxParam != nil {
		ctxParam(config.ctxParam)
	}
	if apiParam != nil {
		apiParam(config.apiParam)
	}

	registry := new(require.Registry)
	registry.Enable(config.vm)

	for name, loader := range config.modules {
		registry.RegisterNativeModule(name, loader)
	}

	// overload error if function panics
	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
			return
		}
	}()
	_, err = config.vm.RunString(script)
	return config, err
}

func executeFn(config *runConfig, fn jsAction) (err error) {
	defer func() {
		r := recover()
		if r != nil && !config.allowedToFail {
			var ok bool
			if err, ok = r.(error); ok {
				return
			}

			e, ok := r.(string)
			if ok {
				err = errors.New(e)
				return
			}
			err = fmt.Errorf("unknown error occured: %v", r)
		}
	}()
	err = fn(config.ctxParam.fields, config.apiParam.fields)
	if err != nil && !config.allowedToFail {
		return err
	}
	return nil
}

func ActionToOptions(a *query.Action) []Option {
	opts := make([]Option, 0, 1)
	if a.AllowedToFail {
		opts = append(opts, WithAllowedToFail())
	}
	return opts
}
