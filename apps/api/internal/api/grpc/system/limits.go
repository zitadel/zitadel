package system

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/object"
	objectpb "github.com/zitadel/zitadel/pkg/grpc/object"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

func (s *Server) SetLimits(ctx context.Context, req *system.SetLimitsRequest) (*system.SetLimitsResponse, error) {
	details, err := s.command.SetLimits(ctx, setInstanceLimitsPbToCommand(req))
	if err != nil {
		return nil, err
	}
	return &system.SetLimitsResponse{
		Details: object.AddToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}

func (s *Server) BulkSetLimits(ctx context.Context, req *system.BulkSetLimitsRequest) (*system.BulkSetLimitsResponse, error) {
	details, targetDetails, err := s.command.SetInstanceLimitsBulk(ctx, bulkSetInstanceLimitsPbToCommand(req))
	if err != nil {
		return nil, err
	}
	resp := &system.BulkSetLimitsResponse{
		Details:       object.AddToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
		TargetDetails: make([]*objectpb.ObjectDetails, len(targetDetails)),
	}
	for i := range targetDetails {
		resp.TargetDetails[i] = object.AddToDetailsPb(targetDetails[i].Sequence, targetDetails[i].EventDate, targetDetails[i].ResourceOwner)
	}
	return resp, nil
}

func (s *Server) ResetLimits(ctx context.Context, _ *system.ResetLimitsRequest) (*system.ResetLimitsResponse, error) {
	details, err := s.command.ResetLimits(ctx)
	if err != nil {
		return nil, err
	}
	return &system.ResetLimitsResponse{
		Details: object.ChangeToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}
