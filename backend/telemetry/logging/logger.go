package logging

import "log/slog"

type Logger struct {
	*slog.Logger
}

func New(l *slog.Logger) *Logger {
	return &Logger{Logger: l}
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{l.Logger.With(args...)}
}
