package execution

import (
	"time"

	"github.com/zitadel/zitadel/internal/logstore"
)

var _ logstore.LogRecord = (*ExecutionLogRecord)(nil)

type ExecutionLogRecord struct {
	Timestamp      time.Time
	InstanceID     string
	OrganizationID string
	ActionID       string
	RunID          string
	Message        string
	FileDescriptor FileDescriptor
}

type FileDescriptor string

const (
	StdOut FileDescriptor = "stdout"
	StdErr FileDescriptor = "stderr"
)

func (e *ExecutionLogRecord) RedactSecrets() logstore.LogRecord {
	// TODO implement
	return e
}
