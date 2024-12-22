// Logs the operation of an instance
package log

import (
	"context"
	"log/slog"

	"github.com/zitadel/zitadel/internal/v3/instance"
	"github.com/zitadel/zitadel/internal/v3/storage"
)

var _ instance.InstanceStorage = (*InstanceLogger)(nil)

type InstanceLogger struct {
	*Logger
}

func NewInstanceLogger(logger *slog.Logger) *InstanceLogger {
	return &InstanceLogger{Logger: NewLogger(logger)}
}

// WriteInstanceAdded implements instance.InstanceStorage.
func (l *InstanceLogger) WriteInstanceAdded(ctx context.Context, tx storage.Transaction, instance *instance.AddInstanceRequest) error {
	tx.OnCommit(func(ctx context.Context) error {
		l.InfoContext(ctx, "Instance added", slog.Any("instance", instance))
		return nil
	})
	return nil
}
