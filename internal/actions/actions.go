package actions

import (
	"context"
	"errors"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"

	z_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
)

var (
	ErrHalt = errors.New("interrupt")
)

type jsAction func(*Context, *API) error

func Run(ctx context.Context, runtimeCtx *Context, api *API, script, name string, opts ...Option) error {
	config := newRunConfig(ctx, opts...)
	if config.timeout == 0 {
		return z_errs.ThrowInternal(nil, "ACTIO-uCpCx", "Errrors.Internal")
	}

	vm, err := prepareRun(script, config)
	if err != nil {
		return err
	}

	var fn jsAction
	jsFn := vm.Get(name)
	if jsFn == nil {
		return errors.New("function not found")
	}
	err = vm.ExportTo(jsFn, &fn)
	if err != nil {
		return err
	}

	t := setInterrupt(vm, config.timeout)
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

		err = fn(runtimeCtx, api)
		if err != nil && !config.allowedToFail {
			errCh <- err
			return
		}
		errCh <- nil
	}()
	return <-errCh
}

func newRuntime(config *runConfig) *goja.Runtime {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())

	registry := new(require.Registry)
	registry.Enable(vm)

	for name, loader := range config.modules {
		registry.RegisterNativeModule(name, loader)
	}

	return vm
}

func prepareRun(script string, config *runConfig) (*goja.Runtime, error) {
	vm := newRuntime(config)
	t := setInterrupt(vm, config.prepareTimeout)
	defer func() {
		t.Stop()
	}()
	errCh := make(chan error)
	go func() {
		defer func() {
			r := recover()
			if r != nil {
				errCh <- r.(error)
				return
			}
		}()
		_, err := vm.RunString(script)
		if err != nil {
			errCh <- err
			return
		}
		errCh <- nil
	}()
	return vm, <-errCh
}

func setInterrupt(vm *goja.Runtime, timeout time.Duration) *time.Timer {
	vm.ClearInterrupt()
	return time.AfterFunc(timeout, func() {
		vm.Interrupt(ErrHalt)
	})
}

func ActionToOptions(a *query.Action) []Option {
	opts := make([]Option, 0, 1)
	if a.AllowedToFail {
		opts = append(opts, WithAllowedToFail())
	}
	return opts
}
