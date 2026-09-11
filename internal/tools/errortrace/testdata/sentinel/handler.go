// Package sentinel is a tracer fixture pinning the empty-ID skip: a
// zerrors.Throw* call with an empty ID literal is the errors.Is() sentinel
// idiom (constructed only for a type comparison, never actually returned —
// see internal/command/org_domain.go for the real-code pattern this
// mirrors), and must not be recorded as a real error site. A real,
// non-empty-ID throw right next to it must still be recorded, proving the
// skip doesn't swallow the rest of the function.
package sentinel

import (
	"context"
	"errors"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type Handler struct{}

func (h *Handler) Lookup(ctx context.Context) error {
	notFoundSentinel := zerrors.ThrowNotFound(nil, "", "sentinel, never returned")
	if errors.Is(fetch(ctx), notFoundSentinel) {
		return zerrors.ThrowNotFound(nil, "TEST-s3nt1n", "not found")
	}
	return nil
}

func fetch(ctx context.Context) error { return nil }
