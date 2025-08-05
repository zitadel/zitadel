package record

import (
	"time"

	"github.com/sirupsen/logrus"
)

type ExecutionLog struct {
	LogDate    time.Time              `json:"logDate"`
	Took       time.Duration          `json:"took"`
	Message    string                 `json:"message"`
	LogLevel   logrus.Level           `json:"logLevel"`
	InstanceID string                 `json:"instanceId"`
	ActionID   string                 `json:"actionId,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

func (e ExecutionLog) Normalize() *ExecutionLog {
	e.Message = cutString(e.Message, 2000)
	return &e
}
