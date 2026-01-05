package logging

import (
	"io"
	"log/slog"
)

type commandErrorWriter struct {
	command string
}

func (w *commandErrorWriter) Write(errTxt []byte) (int, error) {
	logger := slog.Default()
	if w.command != "" {
		logger = logger.With("command", w.command)
	}
	logger.Error("failed to run Zitadel", "err", string(errTxt))
	return len(errTxt), nil
}

// CommandErrorWriter creates a new [io.Writer] that can be used in cobra commands
// to log errors using the global slog logger.
func CommandErrorWriter(command string) io.Writer {
	return &commandErrorWriter{command: command}
}
