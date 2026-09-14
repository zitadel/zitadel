// Package sentinel is a tracer fixture pinning the errors.Is() sentinel
// idiom: a zerrors.Throw* call passed directly as an argument to
// errors.Is(...) is only ever constructed for a type comparison, never
// actually returned (see internal/command/org_domain.go and
// internal/api/oidc/token_refresh.go for the two real occurrences of this
// pattern), so it must not be recorded as a real error site. A real throw
// right next to it must still be recorded, proving the skip doesn't
// swallow the rest of the function.
//
// The signal is being a direct argument to errors.Is(...), not having an
// empty ID: internal/api/grpc/user/v2/user.go's CreateUser and UpdateUser
// both directly return a real, empty-ID error from a default switch case,
// and internal/api/oidc/token_refresh.go's real sentinel has a non-empty
// ID. Either one would be handled wrong by a rule keyed on ID emptiness
// alone.
package sentinel

import (
	"context"
	"errors"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type Handler struct{}

func (h *Handler) Lookup(ctx context.Context) error {
	err := fetch(ctx)
	if errors.Is(err, zerrors.ThrowNotFound(nil, "", "sentinel, never returned")) {
		return zerrors.ThrowNotFound(nil, "TEST-s3nt1n", "not found")
	}
	return nil
}

// DirectEmptyID mirrors CreateUser/UpdateUser's real default-case shape: a
// zerrors.Throw* call with an empty ID that's genuinely returned, not
// compared. Recording it (even with an empty id) is the correct outcome,
// not skipping it as if it were a sentinel.
func (h *Handler) DirectEmptyID(ctx context.Context, kind int) error {
	switch kind {
	case 1:
		return nil
	default:
		return zerrors.ThrowInternal(nil, "", "kind is not implemented")
	}
}

func fetch(ctx context.Context) error { return nil }
