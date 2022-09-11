package actions

import (
	"errors"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

var (
	ErrHalt = errors.New("interrupt")
)

type jsAction func(*Context, *API) error

func Run(ctx *Context, api *API, script, name string, opts ...Option) error {
	config := newRunConfig(opts...)

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

		err = fn(ctx, api)
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
