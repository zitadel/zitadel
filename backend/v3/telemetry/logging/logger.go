package logging

import "log/slog"

// Logger abstracts [slog.Logger] not sure if thats needed
type Logger struct {
	*slog.Logger
}

func NewLogger(logger *slog.Logger) *Logger {
	return &Logger{Logger: logger}
}
