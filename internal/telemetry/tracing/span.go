package tracing

import (
	"context"

	grpc_errs "github.com/caos/zitadel/internal/api/grpc/errors"
	label "go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
	s.span.SetAttributes(label.Uint32("grpc_code", uint32(code)), label.String("grpc_msg", msg), label.String("error_id", id))
}
