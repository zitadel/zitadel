package system

import (
	"context"

	"github.com/zitadel/zitadel/internal/command"

	"github.com/zitadel/zitadel/internal/api/grpc/object"

	"github.com/zitadel/zitadel/internal/api/authz"

	"github.com/zitadel/zitadel/pkg/grpc/system"
	system_pb "github.com/zitadel/zitadel/pkg/grpc/system"
)

func (s *Server) AddQuota(ctx context.Context, req *system.AddQuotaRequest) (*system.AddQuotaResponse, error) {
	ctx = authz.WithInstanceID(ctx, req.InstanceId)

	details, err := s.command.AddInstanceQuota(
		ctx,
		instanceQuotaPbToQuota(req),
	)
	if err != nil {
		return nil, err
	}
	return &system_pb.AddQuotaResponse{
		Details: object.AddToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}

func (s *Server) RemoveQuota(ctx context.Context, req *system.RemoveQuotaRequest) (*system.RemoveQuotaResponse, error) {
	ctx = authz.WithInstanceID(ctx, req.InstanceId)
	details, err := s.command.RemoveInstanceQuota(ctx, command.QuotaUnit(req.Unit))
	if err != nil {
		return nil, err
	}
	return &system_pb.RemoveQuotaResponse{
		Details: object.ChangeToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}
