package system

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

func (s *Server) AddQuota(ctx context.Context, req *system.AddQuotaRequest) (*system.AddQuotaResponse, error) {
	details, err := s.command.AddQuota(
		ctx,
		instanceQuotaPbToCommand(req),
	)
	if err != nil {
		return nil, err
	}
	return &system.AddQuotaResponse{
		Details: object.AddToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}

func (s *Server) SetQuota(ctx context.Context, req *system.SetQuotaRequest) (*system.SetQuotaResponse, error) {
	details, err := s.command.SetQuota(
		ctx,
		instanceQuotaPbToCommand(req),
	)
	if err != nil {
		return nil, err
	}
	return &system.SetQuotaResponse{
		Details: object.AddToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}

func (s *Server) RemoveQuota(ctx context.Context, req *system.RemoveQuotaRequest) (*system.RemoveQuotaResponse, error) {
	details, err := s.command.RemoveQuota(ctx, instanceQuotaUnitPbToCommand(req.Unit))
	if err != nil {
		return nil, err
	}
	return &system.RemoveQuotaResponse{
		Details: object.ChangeToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}
