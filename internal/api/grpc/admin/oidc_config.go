package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetOIDCConfig(ctx context.Context, _ *admin_pb.GetOIDCConfigRequest) (*admin_pb.GetOIDCConfigResponse, error) {
	result, err := s.query.OIDCConfigByAggID(ctx, domain.IAMID)
	if err != nil {
		return nil, err

	}
	return &admin_pb.GetOIDCConfigResponse{
		Config: OIDCConfigToPb(result),
	}, nil
}

func (s *Server) UpdateOIDCConfig(ctx context.Context, req *admin_pb.UpdateOIDCConfigRequest) (*admin_pb.UpdateOIDCConfigResponse, error) {
	result, err := s.command.ChangeOIDCConfig(ctx, UpdateOIDCConfigToConfig(req))
	if err != nil {
		return nil, err

	}
	return &admin_pb.UpdateOIDCConfigResponse{
		Details: object.DomainToChangeDetailsPb(result),
	}, nil
}
