package actions

import (
	"context"
	"errors"
	"fmt"

	"github.com/dop251/goja_nodejs/require"
	"github.com/sirupsen/logrus"

	z_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
)

type Config struct {
	HTTP HTTPConfig
}

var ErrHalt = errors.New("interrupt")

type jsAction func(fields, fields) error

const (
	actionStartedMessage   = "action run started"
	actionSucceededMessage = "action run succeeded"
)

func actionFailedMessage(err error) string {
	return fmt.Sprintf("action run failed: %s", err.Error())
}

func Run(ctx context.Context, ctxParam contextFields, apiParam apiFields, script, name string, opts ...Option) (err error) {
	config := newRunConfig(ctx, append(opts, withLogger(ctx))...)
	if config.functionTimeout == 0 {
		return z_errs.ThrowInternal(nil, "ACTIO-uCpCx", "Errrors.Internal")
	}

	remaining := logstoreService.Limit(ctx, config.instanceID)
	config.cutTimeouts(remaining)

	config.logger.Log(actionStartedMessage)
	if remaining != nil && *remaining == 0 {
		return z_errs.ThrowResourceExhausted(nil, "ACTIO-f19Ii", "Errors.Quota.Execution.Exhausted")
	}

	defer func() {
		if err != nil {
			config.logger.log(actionFailedMessage(err), logrus.ErrorLevel, true)
		} else {
			config.logger.log(actionSucceededMessage, logrus.InfoLevel, true)
		}
		if config.allowedToFail {
			err = nil
		}
	}()

	if err := executeScript(config, ctxParam, apiParam, script); err != nil {
		return err
	}

	var fn jsAction
	jsFn := config.vm.Get(name)
	if jsFn == nil {
		return errors.New("function not found")
	}
	if err := config.vm.ExportTo(jsFn, &fn); err != nil {
		return err
	}

	t := config.StartFunction()
	defer func() {
		t.Stop()
	}()

	return executeFn(config, fn)
}

func executeScript(config *runConfig, ctxParam contextFields, apiParam apiFields, script string) (err error) {
	t := config.StartScript()
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
	return err
}

func executeFn(config *runConfig, fn jsAction) (err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		var ok bool
		if err, ok = r.(error); ok {
			return
		}

		e, ok := r.(string)
		if ok {
			err = errors.New(e)
			return
		}
		err = fmt.Errorf("unknown error occurred: %v", r)
	}()

	if err = fn(config.ctxParam.fields, config.apiParam.fields); err != nil {
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
