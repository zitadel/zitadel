package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) GenerateRecoveryCodes(ctx context.Context, req *user.GenerateRecoveryCodesRequest) (*user.GenerateRecoveryCodesResponse, error) {
	details, err := s.command.GenerateRecoveryCodes(ctx, req.GetUserId(), int(req.GetCount()), "", nil)
	if err != nil {
		return nil, err
	}
	return &user.GenerateRecoveryCodesResponse{
		Details:       object.DomainToDetailsPb(&details.ObjectDetails),
		RecoveryCodes: details.RawCodes,
	}, nil
}

func (s *Server) RemoveRecoveryCodes(ctx context.Context, req *user.RemoveRecoveryCodesRequest) (*user.RemoveRecoveryCodesResponse, error) {
	objectDetails, err := s.command.RemoveRecoveryCodes(ctx, req.GetUserId(), "", nil)
	if err != nil {
		return nil, err
	}
	return &user.RemoveRecoveryCodesResponse{Details: object.DomainToDetailsPb(objectDetails)}, nil
}
