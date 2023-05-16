package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/command"
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
	ctxData := authz.GetCtxData(ctx)
	objectDetails, err := s.command.HumanHumanPasswordlessSetup(ctx, req.GetUserId(), ctxData.ResourceOwner, req.GetPasskeyName(), "", req.GetPublicKeyCredential())
	// ------------------------------------------------------------------------------------------------------- ^ tokenName? ------ ^ UA? -----
	if err != nil {
		return nil, err
	}
	return &user.VerifyPasskeyRegistrationResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) CreatePasskeyRegistrationLink(ctx context.Context, req *user.CreatePasskeyRegistrationLinkRequest) (resp *user.CreatePasskeyRegistrationLinkResponse, err error) {
	var (
		userID        = req.GetUserId()
		resourceOwner = authz.GetCtxData(ctx).ResourceOwner
		details       *command.PasskeyCodeDetails
	)

	switch medium := req.Medium.(type) {
	case nil:
		details, err = s.command.AddUserPasskeyCode(ctx, userID, resourceOwner, s.userCodeAlg)
	case *user.CreatePasskeyRegistrationLinkRequest_SendLink:
		details, err = s.command.AddUserPasskeyCodeURLTemplate(ctx, userID, resourceOwner, s.userCodeAlg, medium.SendLink.GetUrlTemplate())
	case *user.CreatePasskeyRegistrationLinkRequest_ReturnCode:
		details, err = s.command.AddUserPasskeyCodeReturn(ctx, userID, resourceOwner, s.userCodeAlg)
	default:
		err = caos_errs.ThrowUnimplementedf(nil, "USERv2-gaD8y", "verification oneOf %T in method CreatePasskeyRegistrationLink not implemented", medium)
	}
	if err != nil {
		return nil, err
	}
	return passkeyCodeDetailsToPB(details), nil
}

func passkeyCodeDetailsToPB(details *command.PasskeyCodeDetails) *user.CreatePasskeyRegistrationLinkResponse {
	return &user.CreatePasskeyRegistrationLinkResponse{
		Details: object.DomainToDetailsPb(details.ObjectDetails),
		Code: &user.PasskeyRegistrationCode{
			Id:   *details.CodeID,
			Code: *details.Code,
		},
	}
}
