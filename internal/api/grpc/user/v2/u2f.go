package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

func (s *Server) RegisterU2F(ctx context.Context, req *user.RegisterU2FRequest) (*user.RegisterU2FResponse, error) {
	return u2fRegistrationDetailsToPb(
		s.command.RegisterUserU2F(ctx, req.GetUserId(), authz.GetCtxData(ctx).ResourceOwner, req.GetDomain()),
	)
}

func u2fRegistrationDetailsToPb(details *domain.WebAuthNRegistrationDetails, err error) (*user.RegisterU2FResponse, error) {
	objectDetails, options, err := webAuthNRegistrationDetailsToPb(details, err)
	if err != nil {
		return nil, err
	}
	return &user.RegisterU2FResponse{
		Details:                            objectDetails,
		U2FId:                              details.ID,
		PublicKeyCredentialCreationOptions: options,
	}, nil
}

func (s *Server) VerifyU2FRegistration(ctx context.Context, req *user.VerifyU2FRegistrationRequest) (*user.VerifyU2FRegistrationResponse, error) {
	resourceOwner := authz.GetCtxData(ctx).ResourceOwner
	pkc, err := req.GetPublicKeyCredential().MarshalJSON()
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "USERv2-IeTh4", "Errors.Internal")
	}
	objectDetails, err := s.command.HumanVerifyU2FSetup(ctx, req.GetUserId(), resourceOwner, req.GetTokenName(), "", pkc)
	if err != nil {
		return nil, err
	}
	return &user.VerifyU2FRegistrationResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}
