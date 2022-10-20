package logstorage

import (
	"github.com/zitadel/zitadel/internal/logstorage/config"
)

type Config struct {
	Access    *config.Config
	Execution interface{}
}
