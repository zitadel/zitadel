package detection

import (
	"context"

	"github.com/zitadel/zitadel/internal/signals"
)

// FindingSink is the interface for forwarding findings to external systems.
// Implementations might push to a webhook, SIEM, OpenTelemetry collector,
// or any other downstream consumer.
//
// Sinks are invoked asynchronously after a detection evaluation completes.
// They receive the originating signal for correlation context alongside the
// findings the evaluation produced.
//
// Sinks MUST NOT block the detection evaluation path. Implementations
// should perform I/O in a bounded goroutine pool or buffer writes.
//
// This interface is defined for future extensibility — no concrete external
// sinks ship in the POC. Internal forwarding (emit to the signal stream,
// persist via FindingRecorder) remains the default.
type FindingSink interface {
	// Forward sends findings to the external system. The signal is provided
	// for correlation (instance, user, session, trace context). Errors are
	// logged but do not affect the detection decision.
	Forward(ctx context.Context, signal signals.Signal, findings []Finding) error
}

// MultiSink fans out findings to multiple sinks. If any sink returns an
// error, processing continues and all errors are collected.
type MultiSink struct {
	sinks []FindingSink
}

// NewMultiSink creates a sink that forwards to all provided sinks.
func NewMultiSink(sinks ...FindingSink) *MultiSink {
	return &MultiSink{sinks: sinks}
}

// Forward sends findings to all configured sinks, collecting errors.
func (m *MultiSink) Forward(ctx context.Context, signal signals.Signal, findings []Finding) error {
	if len(m.sinks) == 0 || len(findings) == 0 {
		return nil
	}
	var firstErr error
	for _, s := range m.sinks {
		if err := s.Forward(ctx, signal, findings); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
