package auth

import (
	"context"

	"github.com/zitadel/zitadel/v2/pkg/grpc/auth"
)

func (s *Server) Healthz(context.Context, *auth.HealthzRequest) (*auth.HealthzResponse, error) {
	return &auth.HealthzResponse{}, nil
}
