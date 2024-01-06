package system

import (
	"context"
	"errors"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	objectpb "github.com/zitadel/zitadel/pkg/grpc/object"
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

func (s *Server) BulkSetLimits(ctx context.Context, req *system.BulkSetLimitsRequest) (*system.BulkSetLimitsResponse, error) {
	details := make([]*objectpb.ObjectDetails, 0, len(req.Limits))
	var errs error
	for _, limit := range req.Limits {
		detail, err := s.command.SetLimits(
			authz.WithInstanceID(ctx, limit.GetInstanceId()),
			limit.GetInstanceId(),
			instanceLimitsPbToCommand(limit),
		)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		details = append(details, object.AddToDetailsPb(detail.Sequence, detail.EventDate, detail.ResourceOwner))
	}
	return &system.BulkSetLimitsResponse{
		Details: details,
	}, errs
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
