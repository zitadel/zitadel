package execution

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/zitadel/zitadel/internal/logstore"
)

var _ logstore.LogRecord = (*Record)(nil)

type Record struct {
	LogDate    time.Time              `json:"logDate"`
	Took       time.Duration          `json:"took"`
	Message    string                 `json:"message"`
	LogLevel   logrus.Level           `json:"logLevel"`
	InstanceID string                 `json:"instanceId"`
	ProjectID  string                 `json:"projectId"`
	ActionID   string                 `json:"actionId,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

func (e Record) Normalize() logstore.LogRecord {
	return &e
}
