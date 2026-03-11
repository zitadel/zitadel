package detection

import (
	"context"

	"github.com/zitadel/zitadel/internal/signals"
)

// SignalRecorder persists a signal together with its associated findings.
// The implementation may enrich or transform the data before writing.
type SignalRecorder interface {
	Record(ctx context.Context, signal signals.Signal, findings []Finding) error
}
