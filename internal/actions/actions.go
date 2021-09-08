package actions

import (
	"errors"
	"fmt"
	"time"

	"github.com/caos/oidc/pkg/oidc"
	"github.com/caos/zitadel/internal/domain"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
)

var ErrHalt = errors.New("interrupt")

type Context struct {
	User         *domain.Human
	ExternalUser *domain.ExternalUser
	Tokens       *oidc.Tokens
}

type API map[string]interface{}

func (a *API) set(name string, value interface{}) {
	map[string]interface{}(*a)[name] = value
}

type jsAction func(*Context, *API) error

type runOpt func(*API, *goja.Runtime) error

func Appender(list *[]UserGrant) runOpt {
	return func(a *API, _ *goja.Runtime) error {
		a.set("append", appendUserGrant(list))
		return nil
	}
}

func SetExternalUser(user *domain.ExternalUser) runOpt {
	return func(a *API, _ *goja.Runtime) error {
		a.set("User", user)
		return nil
	}
}

type User interface {
	SetFirstName(string)
	SetLastName(string)
	//SetEmail(string)
	//SetEmailVerified(bool
	//SetPhone(string)
	//SetPhoneVerified(bool)
}

type user struct {
	Firstname string
	lastname  string
}

func (u *user) SetFirstName(firstname string) {
	u.Firstname = firstname
}

func (u *user) SetLastName(lastname string) {
	u.lastname = lastname
}

//
//func (u *user) SetEmail(email string) {
//	panic("implement me")
//}
//
//func (u *user) SetEmailVerified(verified bool) {
//	panic("implement me")
//}
//
//func (u *user) SetPhone(phone string) {
//	panic("implement me")
//}
//
//func (u *user) SetPhoneVerified(verified bool) {
//	panic("implement me")
//}

func NewUser(
	firstname,
	lastname string,
	//email string,
	//emailVerified bool,
	//phone string,
	//phoneVerified bool,
) User {
	return &user{
		Firstname: firstname,
		lastname:  lastname,
	}
}

func SetUser(firsname func(string)) runOpt {
	return func(a *API, vm *goja.Runtime) error {
		a.set("setFirstname", firsname)
		return nil
	}
}

func SetMetadata(metadata []*domain.Metadata) runOpt {
	return func(a *API, _ *goja.Runtime) error {
		a.set("User.Metadatas", metadata)
		return nil
	}
}

func Run(ctx *Context, script, name string, timeout time.Duration, allowedToFail bool, opts ...runOpt) error {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	vm, err := prepareRun(script, timeout)
	if err != nil {
		return err
	}
	var fn jsAction
	err = vm.ExportTo(vm.Get(name), &fn)
	if err != nil {
		return err
	}
	api := &API{}
	for _, opt := range opts {
		if err := opt(api, vm); err != nil {
			return err
		}
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
		fmt.Println(s)
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
