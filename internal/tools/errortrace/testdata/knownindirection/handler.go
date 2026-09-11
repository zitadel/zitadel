// Package knownindirection is a tracer fixture: a handler that calls a
// function-typed struct field, the same shape as command.PermissionCheck /
// domain.PermissionCheck (a field, not a method — so the walker's normal
// method-call path can't resolve it, and it has to go through the
// knownIndirections table instead). The test points a fixture
// knownIndirections entry at IndirectTarget below instead of the real
// authz.CheckPermission, so this pins the *mechanism* without coupling the
// test to that function's real, independently-evolving internals.
package knownindirection

import (
	"context"

	"github.com/zitadel/zitadel/internal/zerrors"
)

// PermCheck mirrors the shape of domain.PermissionCheck: a named function
// type assigned to a struct field, called dynamically.
type PermCheck func(ctx context.Context) error

type Handler struct {
	checkPermission PermCheck
}

func (h *Handler) ViaField(ctx context.Context) error {
	return h.checkPermission(ctx)
}

// IndirectTarget stands in for authz.CheckPermission.
func IndirectTarget(ctx context.Context) error {
	return zerrors.ThrowPermissionDenied(nil, "TEST-p3rm1t", "permission denied")
}
