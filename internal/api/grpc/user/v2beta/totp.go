package user

import (
	"context"

	"connectrpc.com/connect"

	object "github.com/zitadel/zitadel/internal/api/grpc/object/v2beta"
	"github.com/zitadel/zitadel/internal/domain"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

func (s *Server) RegisterTOTP(ctx context.Context, req *connect.Request[user.RegisterTOTPRequest]) (*connect.Response[user.RegisterTOTPResponse], error) {
	return totpDetailsToPb(
		s.command.AddUserTOTP(ctx, req.Msg.GetUserId(), ""),
	)
}

func totpDetailsToPb(totp *domain.TOTP, err error) (*connect.Response[user.RegisterTOTPResponse], error) {
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.RegisterTOTPResponse{
		Details: object.DomainToDetailsPb(totp.ObjectDetails),
		Uri:     totp.URI,
		Secret:  totp.Secret,
	}), nil
}

func (s *Server) VerifyTOTPRegistration(ctx context.Context, req *connect.Request[user.VerifyTOTPRegistrationRequest]) (*connect.Response[user.VerifyTOTPRegistrationResponse], error) {
	objectDetails, err := s.command.CheckUserTOTP(ctx, req.Msg.GetUserId(), req.Msg.GetCode(), "")
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.VerifyTOTPRegistrationResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}), nil
}

func (s *Server) RemoveTOTP(ctx context.Context, req *connect.Request[user.RemoveTOTPRequest]) (*connect.Response[user.RemoveTOTPResponse], error) {
	objectDetails, err := s.command.HumanRemoveTOTP(ctx, req.Msg.GetUserId(), "")
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.RemoveTOTPResponse{Details: object.DomainToDetailsPb(objectDetails)}), nil
}
