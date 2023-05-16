package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

func (s *Server) RegisterPasskey(ctx context.Context, req *user.RegisterPasskeyRequest) (*user.RegisterPasskeyResponse, error) {
	generator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypePasswordlessInitCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	ctxData := authz.GetCtxData(ctx)
	platform := passkeyAuthenticatorToDomain(req.GetAuthenticator())
	webAuthNToken, err := s.command.HumanAddPasswordlessSetupInitCode(ctx, req.GetUserId(), ctxData.ResourceOwner, req.GetCode().GetId(), req.GetCode().GetCode(), platform, generator)
	if err != nil {
		return nil, err
	}

	return &user.RegisterPasskeyResponse{
		PublicKeyCredential: webAuthNToken.PublicKey,
	}, nil
}

func passkeyAuthenticatorToDomain(pa user.PasskeyAuthenticator) domain.AuthenticatorAttachment {
	switch pa {
	case user.PasskeyAuthenticator_PASSKEY_AUTHENTICATOR_UNSPECIFIED:
		return domain.AuthenticatorAttachmentUnspecified
	case user.PasskeyAuthenticator_PASSKEY_AUTHENTICATOR_PLATFORM:
		return domain.AuthenticatorAttachmentPlattform
	case user.PasskeyAuthenticator_PASSKEY_AUTHENTICATOR_CROSS_PLATFORM:
		return domain.AuthenticatorAttachmentCrossPlattform
	default:
		return domain.AuthenticatorAttachmentUnspecified
	}
}

func (s *Server) VerifyPasskeyRegistration(ctx context.Context, req *user.VerifyPasskeyRegistrationRequest) (*user.VerifyPasskeyRegistrationResponse, error) {
	resourceOwner := authz.GetCtxData(ctx).ResourceOwner
	objectDetails, err := s.command.HumanHumanPasswordlessSetup(ctx, req.GetUserId(), resourceOwner, req.GetPasskeyName(), "", req.GetPublicKeyCredential())
	if err != nil {
		return nil, err
	}
	return &user.VerifyPasskeyRegistrationResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) CreatePasskeyRegistrationLink(ctx context.Context, req *user.CreatePasskeyRegistrationLinkRequest) (resp *user.CreatePasskeyRegistrationLinkResponse, err error) {
	resourceOwner := authz.GetCtxData(ctx).ResourceOwner

	switch medium := req.Medium.(type) {
	case nil:
		return passkeyDetailsToPb(
			s.command.AddUserPasskeyCode(ctx, req.GetUserId(), resourceOwner, s.userCodeAlg),
		)
	case *user.CreatePasskeyRegistrationLinkRequest_SendLink:
		return passkeyDetailsToPb(
			s.command.AddUserPasskeyCodeURLTemplate(ctx, req.GetUserId(), resourceOwner, s.userCodeAlg, medium.SendLink.GetUrlTemplate()),
		)
	case *user.CreatePasskeyRegistrationLinkRequest_ReturnCode:
		return passkeyCodeDetailsToPb(
			s.command.AddUserPasskeyCodeReturn(ctx, req.GetUserId(), resourceOwner, s.userCodeAlg),
		)
	default:
		return nil, caos_errs.ThrowUnimplementedf(nil, "USERv2-gaD8y", "verification oneOf %T in method CreatePasskeyRegistrationLink not implemented", medium)
	}
}

func passkeyDetailsToPb(details *domain.ObjectDetails, err error) (*user.CreatePasskeyRegistrationLinkResponse, error) {
	if err != nil {
		return nil, err
	}
	return &user.CreatePasskeyRegistrationLinkResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func passkeyCodeDetailsToPb(details *domain.PasskeyCodeDetails, err error) (*user.CreatePasskeyRegistrationLinkResponse, error) {
	if err != nil {
		return nil, err
	}
	return &user.CreatePasskeyRegistrationLinkResponse{
		Details: object.DomainToDetailsPb(details.ObjectDetails),
		Code: &user.PasskeyRegistrationCode{
			Id:   details.CodeID,
			Code: details.Code,
		},
	}, nil
}
