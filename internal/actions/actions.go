package actions

import (
	"context"
	"errors"
	"fmt"

	z_errs "github.com/zitadel/zitadel/internal/errors"

	"github.com/sirupsen/logrus"

	"github.com/zitadel/logging"

	"github.com/dop251/goja_nodejs/require"
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

func Run(ctx context.Context, ctxParam contextFields, apiParam apiFields, script, name string, opts ...Option) error {

	config := newRunConfig(ctx, opts...)
	if config.functionTimeout == 0 {
		return z_errs.ThrowInternal(nil, "ACTIO-uCpCx", "Errrors.Internal")
	}

	doLimit, remaining, err := logstoreService.Limit(ctx, config.instanceID)
	if err != nil {
		logging.Warnf("failed to check whether action executions should be limited: %s", err.Error())
		err = nil
	}

	config.cutTimeouts(remaining)

	if doLimit {
		err = errors.New("action execution seconds exhausted")
		if config.allowedToFail {
			config.logger.log(actionFailedMessage(err), logrus.ErrorLevel, true)
			return nil
		}
		return err
	}

	config.logger.Log(actionStartedMessage)

	if err = executeScript(config, ctxParam, apiParam, script); err != nil {
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

	t := config.StartFunction()
	defer func() {
		t.Stop()
	}()

	return executeFn(config, fn)
}

// TODO: Why does this return non nil errors even though config.allowedToFail is true?
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
			config.logger.log(actionFailedMessage(err), logrus.ErrorLevel, true)
			if config.allowedToFail {
				err = nil
			}
			return
		}

		e, ok := r.(string)
		if ok {
			err = errors.New(e)
			config.logger.log(actionFailedMessage(err), logrus.ErrorLevel, true)
			if config.allowedToFail {
				err = nil
			}
			return
		}
		err = fmt.Errorf("unknown error occured: %v", r)
		config.logger.log(actionFailedMessage(err), logrus.ErrorLevel, true)
		if config.allowedToFail {
			err = nil
		}
	}()

	err = fn(config.ctxParam.fields, config.apiParam.fields)
	if err != nil {
		config.logger.log(actionFailedMessage(err), logrus.ErrorLevel, true)
		if config.allowedToFail {
			return nil
		}
		return err
	}
	config.logger.log(actionSucceededMessage, logrus.InfoLevel, true)
	return nil
}

func ActionToOptions(a *query.Action) []Option {
	opts := make([]Option, 0, 1)
	if a.AllowedToFail {
		opts = append(opts, WithAllowedToFail())
	}
	return opts
}
