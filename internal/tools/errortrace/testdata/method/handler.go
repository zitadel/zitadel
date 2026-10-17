// Package method is a tracer fixture: a handler method whose throw is one
// concrete method call away. Pins that the walker actually recurses into a
// same-package helper rather than only looking at the handler's own body.
package method

import (
	"context"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type Handler struct{}

func (h *Handler) ViaMethod(ctx context.Context) error {
	return h.helper(ctx)
}

func (h *Handler) helper(ctx context.Context) error {
	return zerrors.ThrowInvalidArgument(nil, "TEST-m3th0d", "bad input")
}
