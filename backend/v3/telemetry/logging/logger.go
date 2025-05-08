package logging

import "log/slog"

type Logger struct {
	*slog.Logger
}

func NewLogger(logger *slog.Logger) *Logger {
	return &Logger{Logger: logger}
}
