package userv2

import (
	"go.opentelemetry.io/otel/trace"

	"github.com/zitadel/zitadel/backend/command/v2/domain"
)

type Server struct {
	tracer trace.Tracer
	domain *domain.Domain
}
