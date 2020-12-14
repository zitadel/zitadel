package tracing

import (
	grpc_errs "github.com/caos/zitadel/internal/api/grpc/errors"
	"go.opentelemetry.io/otel/label"
	api_trace "go.opentelemetry.io/otel/trace"
)

type Span struct {
	span api_trace.Span
	opts []api_trace.SpanOption
}

func CreateSpan(span api_trace.Span) *Span {
	return &Span{span: span, opts: []api_trace.SpanOption{}}
}

func (s *Span) End() {
	if s.span == nil {
		return
	}
	s.span.End(s.opts...)
}

func (s *Span) EndWithError(err error) {
	s.SetStatusByError(err)
	s.End()
}

func (s *Span) SetStatusByError(err error) {
	if s.span == nil {
		return
	}
	if err != nil {
		s.span.RecordError(err)
	}

	code, msg, id, _ := grpc_errs.ExtractCaosError(err)
	s.span.SetAttributes(
		label.Uint32("grpc_code", uint32(code)),
		label.String("grpc_msg", msg),
		label.String("error_id", id),
	)
}
