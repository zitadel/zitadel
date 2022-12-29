package logstore_test

import (
	"time"

	"github.com/zitadel/zitadel/internal/logstore"
)

type option func(config *logstore.EmitterConfig)

func emitterConfig(options ...option) *logstore.EmitterConfig {
	cfg := &logstore.EmitterConfig{
		Enabled:         true,
		Keep:            time.Hour,
		CleanupInterval: time.Hour,
		Debounce: &logstore.DebouncerConfig{
			MinFrequency: 0,
			MaxBulkSize:  0,
		},
	}
	for _, opt := range options {
		opt(cfg)
	}
	return cfg
}

func withDebouncerConfig(config *logstore.DebouncerConfig) option {
	return func(c *logstore.EmitterConfig) {
		c.Debounce = config
	}
}

func withDisabled() option {
	return func(c *logstore.EmitterConfig) {
		c.Enabled = false
	}
}

func withCleanupping(keep, interval time.Duration) option {
	return func(c *logstore.EmitterConfig) {
		c.Keep = keep
		c.CleanupInterval = interval
	}
}

func repeat(value, times int) []int {
	ints := make([]int, times)
	for i := 0; i < times; i++ {
		ints[i] = value
	}
	return ints
}
