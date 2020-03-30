package tracing

import (
	"fmt"
	"strconv"

	"go.opencensus.io/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/caos/zitadel/internal/errors"
)

type Span struct {
	span       *trace.Span
	attributes []trace.Attribute
}

func CreateSpan(span *trace.Span) *Span {
	return &Span{span: span, attributes: []trace.Attribute{}}
}

func (s *Span) End() {
	if s.span == nil {
		return
	}
	s.span.AddAttributes(s.attributes...)
	s.span.End()
}

func (s *Span) EndWithError(err error) {
	s.SetStatusByError(err)
	s.End()
}

func (s *Span) SetStatusByError(err error) {
	if s.span == nil {
		return
	}
	s.span.SetStatus(statusFromError(err))
}

func statusFromError(err error) trace.Status {
	if statusErr, ok := status.FromError(err); ok {
		return trace.Status{Code: int32(statusErr.Code()), Message: statusErr.Message()}
	}
	return trace.Status{Code: int32(codes.Unknown), Message: "Unknown"}
}

// AddAnnotation creates an annotation. The annotation will not be added to the tracing use Annotate(msg) afterwards
func (s *Span) AddAnnotation(key string, value interface{}) *Span {
	attribute, err := toTraceAttribute(key, value)
	if err != nil {
		return s
	}
	s.attributes = append(s.attributes, attribute)
	return s
}

// Annotate creates an annotation in tracing. Before added annotations will be set
func (s *Span) Annotate(message string) *Span {
	if s.span == nil {
		return s
	}
	s.span.Annotate(s.attributes, message)
	s.attributes = []trace.Attribute{}
	return s
}

func (s *Span) Annotatef(format string, addiations ...interface{}) *Span {
	s.Annotate(fmt.Sprintf(format, addiations...))
	return s
}

func toTraceAttribute(key string, value interface{}) (attr trace.Attribute, err error) {
	switch value := value.(type) {
	case bool:
		return trace.BoolAttribute(key, value), nil
	case string:
		return trace.StringAttribute(key, value), nil
	}
	if valueInt, err := convertToInt64(value); err == nil {
		return trace.Int64Attribute(key, valueInt), nil
	}
	return attr, errors.ThrowInternal(nil, "TRACE-jlq3s", "Attribute is not of type bool, string or int64")
}

func convertToInt64(value interface{}) (int64, error) {
	valueString := fmt.Sprintf("%v", value)
	return strconv.ParseInt(valueString, 10, 64)
}
