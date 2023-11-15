package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/repository/limits"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) SetInstanceLimits(ctx context.Context, req *admin.SetInstanceLimitsRequest) (*admin.SetInstanceLimitsResponse, error) {
	details, err := s.command.SetLimits(
		ctx,
		authz.GetInstance(ctx).InstanceID(),
		instanceLimitsPbToCommand(req),
	)
	if err != nil {
		return nil, err
	}
	return &admin.SetInstanceLimitsResponse{
		Details: object.AddToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}

func (s *Server) ResetInstanceLimits(ctx context.Context, _ *admin.ResetInstanceLimitsRequest) (*admin.ResetInstanceLimitsResponse, error) {
	details, err := s.command.ResetLimits(ctx, authz.GetInstance(ctx).InstanceID(), limits.ResetAllowPublicOrgRegistration)
	if err != nil {
		return nil, err
	}
	return &admin.ResetInstanceLimitsResponse{
		Details: object.ChangeToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}
