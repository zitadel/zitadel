package actions

import (
	"context"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/zitadel/logging"
)

const (
	maxPrepareTimeout = 5 * time.Second
)

type Option func(*runConfig)

func WithAllowedToFail() Option {
	return func(c *runConfig) {
		c.allowedToFail = true
	}
}

type runConfig struct {
	allowedToFail bool
	functionTimeout,
	scriptTimeout time.Duration
	modules    map[string]require.ModuleLoader
	logger     *logger
	instanceID string
	vm         *goja.Runtime
	ctxParam   *ctxConfig
	apiParam   *apiConfig
}

func newRunConfig(ctx context.Context, opts ...Option) *runConfig {
	deadline, ok := ctx.Deadline()
	if !ok {
		logging.Warn("no timeout set on action run")
	}

	vm := goja.New()
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())

	config := &runConfig{
		functionTimeout: time.Until(deadline),
		scriptTimeout:   maxPrepareTimeout,
		modules:         map[string]require.ModuleLoader{},
		vm:              vm,
		ctxParam: &ctxConfig{
			FieldConfig: FieldConfig{
				Runtime: vm,
				fields:  fields{},
			},
		},
		apiParam: &apiConfig{
			FieldConfig: FieldConfig{
				Runtime: vm,
				fields:  fields{},
			},
		},
	}

	for _, opt := range opts {
		opt(config)
	}

	if config.scriptTimeout > config.functionTimeout {
		config.scriptTimeout = config.functionTimeout
	}

	return config
}

func (c *runConfig) StartFunction() *time.Timer {
	c.vm.ClearInterrupt()
	return time.AfterFunc(c.functionTimeout, func() {
		c.vm.Interrupt(ErrHalt)
	})
}

func (c *runConfig) StartScript() *time.Timer {
	c.vm.ClearInterrupt()
	return time.AfterFunc(c.scriptTimeout, func() {
		c.vm.Interrupt(ErrHalt)
	})
}

func (c *runConfig) cutTimeouts(remainingSeconds *uint64) {
	if remainingSeconds == nil {
		return
	}

	remainingDur := time.Duration(*remainingSeconds) * time.Second
	if c.functionTimeout > remainingDur {
		c.functionTimeout = remainingDur
	}
	if c.scriptTimeout > remainingDur {
		c.scriptTimeout = remainingDur
	}
}
