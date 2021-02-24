package admin

import (
	"context"

	"github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) Healthz(context.Context, *admin.HealthzRequest) (*admin.HealthzResponse, error) {
	return &admin.HealthzResponse{}, nil
}
