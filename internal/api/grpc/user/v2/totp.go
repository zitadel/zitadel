package user

import (
	"context"

	user "github.com/zitadel/zitadel/pkg/grpc/user/v2"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/domain"
)

func (s *Server) RegisterTOTP(ctx context.Context, req *user.RegisterTOTPRequest) (*user.RegisterTOTPResponse, error) {
	return totpDetailsToPb(
		s.command.AddUserTOTP(ctx, req.GetUserId(), ""),
	)
}

func totpDetailsToPb(totp *domain.TOTP, err error) (*user.RegisterTOTPResponse, error) {
	if err != nil {
		return nil, err
	}
	return &user.RegisterTOTPResponse{
		Details: object.DomainToDetailsPb(totp.ObjectDetails),
		Uri:     totp.URI,
		Secret:  totp.Secret,
	}, nil
}

func (s *Server) VerifyTOTPRegistration(ctx context.Context, req *user.VerifyTOTPRegistrationRequest) (*user.VerifyTOTPRegistrationResponse, error) {
	objectDetails, err := s.command.CheckUserTOTP(ctx, req.GetUserId(), req.GetCode(), "")
	if err != nil {
		return nil, err
	}
	return &user.VerifyTOTPRegistrationResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) RemoveTOTP(ctx context.Context, req *user.RemoveTOTPRequest) (*user.RemoveTOTPResponse, error) {
	objectDetails, err := s.command.HumanRemoveTOTP(ctx, req.GetUserId(), "")
	if err != nil {
		return nil, err
	}
	return &user.RemoveTOTPResponse{Details: object.DomainToDetailsPb(objectDetails)}, nil
}
