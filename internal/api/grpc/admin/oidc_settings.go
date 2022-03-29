package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/object"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetOIDCSettings(ctx context.Context, _ *admin_pb.GetOIDCSettingsRequest) (*admin_pb.GetOIDCSettingsResponse, error) {
	result, err := s.query.OIDCSettingsByAggID(ctx, authz.GetInstance(ctx).ID)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetOIDCSettingsResponse{
		Settings: OIDCSettingsToPb(result),
	}, nil
}

func (s *Server) UpdateOIDCSettings(ctx context.Context, req *admin_pb.UpdateOIDCSettingsRequest) (*admin_pb.UpdateOIDCSettingsResponse, error) {
	result, err := s.command.ChangeOIDCSettings(ctx, authz.GetInstance(ctx).ID, UpdateOIDCConfigToConfig(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateOIDCSettingsResponse{
		Details: object.DomainToChangeDetailsPb(result),
	}, nil
}
