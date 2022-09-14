package actions

import (
	"context"
	"time"

	"github.com/dop251/goja_nodejs/require"
	"github.com/zitadel/logging"
)

const (
	maxPrepareTimeout = 5 * time.Second
)

type runConfig struct {
	allowedToFail bool
	timeout,
	prepareTimeout time.Duration
	modules map[string]require.ModuleLoader
	end     time.Time
}

func newRunConfig(ctx context.Context, opts ...Option) *runConfig {
	deadline, ok := ctx.Deadline()
	if !ok {
		logging.Warn("no timeout set on action run")
	}

	config := &runConfig{
		timeout:        time.Until(deadline),
		prepareTimeout: maxPrepareTimeout,
		modules:        map[string]require.ModuleLoader{},
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

type Option func(*runConfig)

func WithAllowedToFail() Option {
	return func(c *runConfig) {
		c.allowedToFail = true
	}
}
