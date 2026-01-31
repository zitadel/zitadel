package queue

import (
	"context"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
)

type logMiddleware struct {
	river.MiddlewareDefaults
	logger *slog.Logger
}

func newLogMiddleware() rivertype.Middleware {
	return &logMiddleware{
		logger: logging.New(logging.StreamQueue),
	}
}

func (m *logMiddleware) InsertMany(
	ctx context.Context,
	manyParams []*rivertype.JobInsertParams,
	doInner func(context.Context) ([]*rivertype.JobInsertResult, error),
) ([]*rivertype.JobInsertResult, error) {
	start := time.Now()
	ctx = logging.ToCtx(ctx, m.logger)
	results, err := doInner(ctx)
	if err != nil {
		logging.WithError(ctx, err).Error("insert many error")
		return results, err
	}
	logging.Info(ctx, "jobs inserted",
		slog.Int("count", len(results)),
		slog.Duration("duration", time.Since(start)),
	)

	// Only do expensive operations if debug is enabled
	if m.logger.Enabled(ctx, slog.LevelDebug) {
		for _, result := range results {
			logging.Debug(ctx, "inserted job details", attributesFromJobInsertResult(result)...)
		}
	}
	return results, err
}

func (m *logMiddleware) Work(
	ctx context.Context,
	job *rivertype.JobRow,
	doInner func(context.Context) error,
) error {
	start := time.Now()
	ctx = logging.ToCtx(ctx, m.logger)
	ctx = logging.With(ctx, attributesFromJobRow(job)...)
	if err := doInner(ctx); err != nil {
		logging.WithError(ctx, err).Warn("job processing error")
		return err
	}
	logging.Info(ctx, "job processed successfully",
		slog.Duration("duration", time.Since(start)),
	)
	return nil
}

func attributesFromJobRow(j *rivertype.JobRow) []any {
	attributes := make([]any, 0, 14)
	attributes = append(attributes,
		slog.String("queue", j.Queue),
		slog.Int64("job_id", j.ID),
		slog.String("kind", j.Kind),
		slog.Int("priority", j.Priority),
		slog.Int("max_attempts", j.MaxAttempts),
		slog.String("state", string(j.State)),
		slog.String("unique_key", hex.EncodeToString(j.UniqueKey)),
	)
	if j.AttemptedAt != nil {
		attributes = append(attributes,
			slog.Time("created_at", j.CreatedAt),
			slog.Int("attempt", j.Attempt),
			slog.Time("attempted_at", *j.AttemptedAt),
			slog.Any("attempted_by", j.AttemptedBy),
		)
	}
	if j.FinalizedAt != nil {
		attributes = append(attributes, slog.Time("finalized_at", *j.FinalizedAt))
	}
	if !j.ScheduledAt.IsZero() {
		attributes = append(attributes, slog.Time("scheduled_at", j.ScheduledAt))
	}
	if len(j.Tags) > 0 {
		attributes = append(attributes, slog.Any("tags", j.Tags))
	}
	return attributes
}

func attributesFromJobInsertResult(j *rivertype.JobInsertResult) []any {
	attributes := make([]any, 0, 15)
	attributes = append(attributes, attributesFromJobRow(j.Job)...)
	attributes = append(attributes, slog.Bool("unique_skipped_as_duplicate", j.UniqueSkippedAsDuplicate))
	return attributes
}
