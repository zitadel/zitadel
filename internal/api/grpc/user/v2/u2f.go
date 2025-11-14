package user

import (
	"context"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) RegisterU2F(ctx context.Context, req *connect.Request[user.RegisterU2FRequest]) (*connect.Response[user.RegisterU2FResponse], error) {
	return u2fRegistrationDetailsToPb(
		s.command.RegisterUserU2F(ctx, req.Msg.GetUserId(), "", req.Msg.GetDomain()),
	)
}

func u2fRegistrationDetailsToPb(details *domain.WebAuthNRegistrationDetails, err error) (*connect.Response[user.RegisterU2FResponse], error) {
	objectDetails, options, err := webAuthNRegistrationDetailsToPb(details, err)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.RegisterU2FResponse{
		Details:                            objectDetails,
		U2FId:                              details.ID,
		PublicKeyCredentialCreationOptions: options,
	}), nil
}

func (s *Server) VerifyU2FRegistration(ctx context.Context, req *connect.Request[user.VerifyU2FRegistrationRequest]) (*connect.Response[user.VerifyU2FRegistrationResponse], error) {
	pkc, err := req.Msg.GetPublicKeyCredential().MarshalJSON()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USERv2-IeTh4", "Errors.Internal")
	}
	objectDetails, err := s.command.HumanVerifyU2FSetup(ctx, req.Msg.GetUserId(), "", req.Msg.GetTokenName(), "", pkc)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.VerifyU2FRegistrationResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}), nil
}

func (s *Server) RemoveU2F(ctx context.Context, req *connect.Request[user.RemoveU2FRequest]) (*connect.Response[user.RemoveU2FResponse], error) {
	objectDetails, err := s.command.HumanRemoveU2F(ctx, req.Msg.GetUserId(), req.Msg.GetU2FId(), "")
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.RemoveU2FResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}), nil
}
