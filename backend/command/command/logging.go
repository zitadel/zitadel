package command

import (
	"context"
	"log/slog"
	"time"

	"github.com/zitadel/zitadel/backend/telemetry/logging"
)

type logger struct {
	level slog.Level
	*logging.Logger
	cmd Command
}

// Activity decorates the commands execute method with logging.
// It logs the command name, duration, and success or failure of the command.
func Activity(l *logging.Logger, command Command) Command {
	return &logger{
		Logger: l.With(slog.String("type", "activity")),
		level:  slog.LevelInfo,
		cmd:    command,
	}
}

// Name implements [Command].
func (l *logger) Name() string {
	return l.cmd.Name()
}

func (l *logger) Execute(ctx context.Context) error {
	start := time.Now()
	log := l.Logger.With(slog.String("command", l.cmd.Name()))
	log.InfoContext(ctx, "execute")
	err := l.cmd.Execute(ctx)
	log = log.With(slog.Duration("took", time.Since(start)))
	if err != nil {
		log.Log(ctx, l.level, "failed", slog.Any("cause", err))
		return err
	}
	log.Log(ctx, l.level, "successful")
	return nil
}
