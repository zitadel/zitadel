package logstore

import "time"

type FileDescriptor string

const (
	StdOut FileDescriptor = "stdout"
	StdErr FileDescriptor = "stderr"
)

type ExecutionLogRecord struct {
	Timestamp      time.Time
	InstanceID     string
	OrganizationID string
	ActionID       string
	RunID          string
	Message        string
	FileDescriptor FileDescriptor
}
