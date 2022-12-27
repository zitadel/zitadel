package execution

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/zitadel/zitadel/internal/logstore"
)

var _ logstore.LogRecord = (*Record)(nil)

type Record struct {
	Timestamp  time.Time
	TookMS     int64
	Message    string
	LogLevel   logrus.Level
	InstanceID string
	ProjectID  string
	ActionID   string
	Metadata   map[string]interface{}
}

func (e *Record) RedactSecrets() logstore.LogRecord {
	// TODO implement?
	return e
}
