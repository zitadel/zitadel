package streams

type LogFieldKey string

const (
	LogFieldKeyVersion    LogFieldKey = "version"
	LogFieldKeyStream     LogFieldKey = "stream"
	LogFieldKeyInstanceID LogFieldKey = "instance"
)

type Stream string

const (
	LogFieldValueStreamActivity Stream = "activity"
	LogFieldValueStreamVersion  Stream = "v1"
)
