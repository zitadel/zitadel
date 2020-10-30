package tracing

import (
	"context"

	grpc_errs "github.com/caos/zitadel/internal/api/grpc/errors"
	apitrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
)

type Span struct {
	span apitrace.Span
	opts []apitrace.SpanOption
}

func CreateSpan(span apitrace.Span) *Span {
	return &Span{span: span, opts: []apitrace.SpanOption{}}
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
		s.span.RecordError(context.TODO(), err, apitrace.WithErrorStatus(codes.Error))
	}

	code, msg, _, _ := grpc_errs.ExtractCaosError(err)
	s.span.SetAttributes(label.Uint32("grpc_code", uint32(code)), label.String("grpc_msg", msg))
}
