package actions

import (
	"errors"
	"time"

	"github.com/caos/logging"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
)

var ErrHalt = errors.New("interrupt")

type jsAction func(*Context, *API) error

func Run(ctx *Context, api *API, script, name string, timeout time.Duration, allowedToFail bool) error {
	if timeout <= 0 || timeout > 20 {
		timeout = 20 * time.Second
	}
	prepareTimeout := timeout
	if prepareTimeout > 5 {
		prepareTimeout = 5 * time.Second
	}
	vm, err := prepareRun(script, prepareTimeout)
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
	t := setInterrupt(vm, timeout)
	defer func() {
		t.Stop()
	}()
	errCh := make(chan error)
	go func() {
		defer func() {
			r := recover()
			if r != nil && !allowedToFail {
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
		if err != nil && !allowedToFail {
			errCh <- err
			return
		}
		errCh <- nil
	}()
	return <-errCh
}

func newRuntime() *goja.Runtime {
	vm := goja.New()

	printer := console.PrinterFunc(func(s string) {
		logging.Log("ACTIONS-dfgg2").Debug(s)
	})
	registry := new(require.Registry)
	registry.Enable(vm)
	registry.RegisterNativeModule("console", console.RequireWithPrinter(printer))
	console.Enable(vm)

	return vm
}

func prepareRun(script string, timeout time.Duration) (*goja.Runtime, error) {
	vm := newRuntime()
	t := setInterrupt(vm, timeout)
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
