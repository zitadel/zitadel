// Logs the operation of an instance
package log

import (
	"log/slog"
)

type Logger struct {
	*slog.Logger
}

func NewLogger(logger *slog.Logger) *Logger {
	return &Logger{logger}
}
