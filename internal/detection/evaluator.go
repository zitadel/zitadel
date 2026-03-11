package detection

import (
	"context"

	"github.com/zitadel/zitadel/internal/signals"
)

// Evaluator performs synchronous detection evaluation for a given signal.
// Implementations decide whether the request should be allowed, blocked,
// or challenged based on rules, heuristics, and external classifiers.
type Evaluator interface {
	Evaluate(ctx context.Context, signal signals.Signal) (Decision, error)
}
