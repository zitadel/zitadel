package otel

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"

	"github.com/zitadel/zitadel/cmd/build"
)

func ResourceWithService() (*resource.Resource, error) {
	attributes := []attribute.KeyValue{
		semconv.ServiceNameKey.String("ZITADEL"),
	}
	if build.Version() != "" {
		attributes = append(attributes, semconv.ServiceVersionKey.String(build.Version()))
	}
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			attributes...,
		),
	)
}
