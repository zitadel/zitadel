package actions

import (
	"time"

	"github.com/dop251/goja_nodejs/require"
)

const (
	maxTimeout        = 20 * time.Second
	maxPrepareTimeout = 5 * time.Second
)

type runConfig struct {
	allowedToFail bool
	timeout,
	prepareTimeout time.Duration
	modules map[string]require.ModuleLoader
	end     time.Time
}

func newRunConfig(opts ...Option) *runConfig {
	config := &runConfig{
		timeout:        maxTimeout,
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

// WithTimeout sets the passed timeout for the execution
// timeout has to be between 0 and 20 seconds
// values out of range are ignored
func WithTimeout(timeout time.Duration) Option {
	return func(c *runConfig) {
		if timeout <= 0 || timeout > maxTimeout {
			return
		}
		c.timeout = timeout
		if timeout > maxPrepareTimeout {
			return
		}
		c.prepareTimeout = timeout
	}
}

func WithAllowedToFail() Option {
	return func(c *runConfig) {
		c.allowedToFail = true
	}
}
