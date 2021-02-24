package management

import (
	"context"

	"github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) Healthz(context.Context, *management.HealthzRequest) (*management.HealthzResponse, error) {
	return &management.HealthzResponse{}, nil
}
