package tracing

import (
	grpc_errs "github.com/caos/zitadel/internal/api/grpc/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Span struct {
	span trace.Span
	opts []trace.SpanEndOption
}

func CreateSpan(span trace.Span) *Span {
	return &Span{span: span, opts: []trace.SpanEndOption{}}
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
		// trace.WithErrorStatus(codes.Error)
		s.span.RecordError(err)
		s.span.SetAttributes(
			attribute.KeyValue{},
		)
	}

	code, msg, id, _ := grpc_errs.ExtractCaosError(err)
	s.span.SetAttributes(attribute.Int("grpc_code", int(code)), attribute.String("grpc_msg", msg), attribute.String("error_id", id))
}
