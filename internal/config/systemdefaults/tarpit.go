package systemdefaults

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type TarpitConfig struct {
	// After how many failed attempts, the tarpit should start.
	MinFailedAttempts uint64
	// The seconds that will be added per step.
	StepDuration time.Duration
	// The failed attempts that are needed to increase the tarpit by one step.
	StepSize uint64
	// The maximum duration the tarpit can reach.
	MaxDuration time.Duration
}

// Tarpit returns a function which delays the response after a failed check,
// depending on the number of failed attempts.
// A delay is recorded as span and logged, so it's visible in traces and logs.
func (t *TarpitConfig) Tarpit() func(ctx context.Context, failedCount uint64) {
	return func(ctx context.Context, failedCount uint64) {
		duration := t.duration(failedCount)
		if duration == 0 {
			return
		}
		ctx, span := tracing.NewNamedSpan(ctx, "tarpit")
		defer span.End()
		trace.SpanFromContext(ctx).SetAttributes(
			attribute.Int64("tarpit.failed_attempts", int64(failedCount)),
			attribute.String("tarpit.duration", duration.String()),
		)
		logging.Info(ctx, "tarpit", "failed_attempts", failedCount, "duration", duration)
		time.Sleep(duration)
	}
}

func (t *TarpitConfig) duration(failedCount uint64) time.Duration {
	if failedCount < t.MinFailedAttempts {
		return 0
	}
	// calculate the step we are at
	// every StepSize failed attempts increase the step by one
	step := (failedCount - t.MinFailedAttempts) / t.StepSize
	duration := time.Duration(step) * t.StepDuration
	if duration < t.MaxDuration {
		return duration
	}
	return t.MaxDuration
}
