package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) GetOIDCSettings(ctx context.Context, _ *admin_pb.GetOIDCSettingsRequest) (*admin_pb.GetOIDCSettingsResponse, error) {
	result, err := s.query.OIDCSettingsByAggID(ctx, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetOIDCSettingsResponse{
		Settings: OIDCSettingsToPb(result),
	}, nil
}

func (s *Server) AddOIDCSettings(ctx context.Context, req *admin_pb.AddOIDCSettingsRequest) (*admin_pb.AddOIDCSettingsResponse, error) {
	result, err := s.command.AddOIDCSettings(ctx, AddOIDCConfigToConfig(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddOIDCSettingsResponse{
		Details: object.DomainToChangeDetailsPb(result),
	}, nil
}

func (s *Server) UpdateOIDCSettings(ctx context.Context, req *admin_pb.UpdateOIDCSettingsRequest) (*admin_pb.UpdateOIDCSettingsResponse, error) {
	result, err := s.command.ChangeOIDCSettings(ctx, UpdateOIDCConfigToConfig(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateOIDCSettingsResponse{
		Details: object.DomainToChangeDetailsPb(result),
	}, nil
}
