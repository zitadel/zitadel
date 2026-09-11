// Package direct is a tracer fixture: a handler method that throws directly,
// with no indirection at all. Pins the simplest possible case — the walker
// must find this without following anything.
package direct

import (
	"context"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type Handler struct{}

func (h *Handler) Direct(ctx context.Context) error {
	return zerrors.ThrowNotFound(nil, "TEST-d1r3ct", "not found")
}
