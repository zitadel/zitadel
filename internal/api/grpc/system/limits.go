package system

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

func (s *Server) SetLimits(ctx context.Context, req *system.SetLimitsRequest) (*system.SetLimitsResponse, error) {
	details, err := s.command.SetLimits(
		ctx,
		req.GetInstanceId(),
		instanceLimitsPbToCommand(req),
	)
	if err != nil {
		return nil, err
	}
	return &system.SetLimitsResponse{
		Details: object.AddToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}

func (s *Server) ResetLimits(ctx context.Context, req *system.ResetLimitsRequest) (*system.ResetLimitsResponse, error) {
	details, err := s.command.ResetLimits(ctx, req.GetInstanceId())
	if err != nil {
		return nil, err
	}
	return &system.ResetLimitsResponse{
		Details: object.ChangeToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}
