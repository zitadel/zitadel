package config

import (
	"github.com/zitadel/zitadel/internal/logstorage/debouncer"
	"time"
)

type Config struct {
	Enabled   bool
	Retention time.Duration
	Debouncer *debouncer.Config
}
