// Package unresolved is a tracer fixture: a handler calling a function-typed
// struct field with no knownIndirections entry at all. The walker must
// report this on w.unresolved (surfaced, never silently dropped), not
// pretend the call simply doesn't exist.
package unresolved

import "context"

// Op is an in-module named function type nothing maps in knownIndirections.
type Op func(ctx context.Context) error

type Handler struct {
	op Op
}

func (h *Handler) Dynamic(ctx context.Context) error {
	return h.op(ctx)
}
