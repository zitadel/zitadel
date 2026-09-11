// Package grpcstatus is a tracer fixture: a stub handler returning a raw
// gRPC status directly, never touching zerrors at all. Pins the
// status.Errorf(codes.X, ...) -> synthetic "GRPC-<CODE>" site path.
package grpcstatus

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct{}

func (h *Handler) Unimplemented(ctx context.Context) error {
	return status.Errorf(codes.Unimplemented, "method Unimplemented not implemented")
}

// FormattedMessage pins that a literal format string containing a real
// Printf verb is dropped rather than shown as-is — the verb's value only
// exists at runtime, so the raw literal ("instance %s not found") would be
// a false, half-templated message on a public docs page.
func (h *Handler) FormattedMessage(ctx context.Context, id string) error {
	return status.Errorf(codes.NotFound, "instance %s not found", id)
}
