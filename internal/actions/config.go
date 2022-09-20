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

func WithContextOptions(opts ...ContextOption) Option {
	return func(c *runConfig) {
		for _, opt := range opts {
			opt(c.ctxParam)
		}
	}
}

func WithAPIOptions(opts ...APIOption) Option {
	return func(c *runConfig) {
		for _, opt := range opts {
			opt(c.apiParam)
		}
	}
}

type runConfig struct {
	allowedToFail bool
	timeout,
	prepareTimeout time.Duration
	modules map[string]require.ModuleLoader
	end     time.Time

	vm       *goja.Runtime
	ctxParam *contextParam
	apiParam *apiParam
}

func newRunConfig(ctx context.Context, opts ...Option) *runConfig {
	deadline, ok := ctx.Deadline()
	if !ok {
		logging.Warn("no timeout set on action run")
	}

	vm := goja.New()
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())

	config := &runConfig{
		timeout:        time.Until(deadline),
		prepareTimeout: maxPrepareTimeout,
		modules:        map[string]require.ModuleLoader{},
		vm:             vm,
		ctxParam: &contextParam{
			runtime:   vm,
			parameter: parameter{},
		},
		apiParam: &apiParam{
			runtime:   vm,
			parameter: parameter{},
		},
	}

	for _, opt := range opts {
		opt(config)
	}

	if config.prepareTimeout > config.timeout {
		config.prepareTimeout = config.timeout
	}

	config.end = time.Now().Add(config.timeout)

	return config
}

func (c *runConfig) Start() *time.Timer {
	c.vm.ClearInterrupt()
	return time.AfterFunc(c.timeout, func() {
		c.vm.Interrupt(ErrHalt)
	})
}

func (c *runConfig) Prepare() *time.Timer {
	c.vm.ClearInterrupt()
	return time.AfterFunc(c.prepareTimeout, func() {
		c.vm.Interrupt(ErrHalt)
	})
}
