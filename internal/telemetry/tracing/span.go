package tracing

import (
	"context"

	grpc_errs "github.com/caos/zitadel/internal/api/grpc/errors"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type Span struct {
	span trace.Span
	opts []trace.SpanOption
}

func CreateSpan(span trace.Span) *Span {
	return &Span{span: span, opts: []trace.SpanOption{}}
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
		s.span.RecordError(context.TODO(), err, trace.WithErrorStatus(codes.Error))
	}

	code, msg, id, _ := grpc_errs.ExtractCaosError(err)
	s.span.SetAttributes(
		attribute.Int64(("grpc_code", code),
		attribute.String("grpc_msg", msg),
		attribute.String("error_id", id))
}
