package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/domain"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

func (s *Server) RegisterPasskey(ctx context.Context, req *user.RegisterPasskeyRequest) (*user.RegisterPasskeyResponse, error) {
	generator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypePasswordlessInitCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	ctxData := authz.GetCtxData(ctx)
	platform := passkeyAuthenticatorToDomain(req.GetAuthenticator())
	webAuthNToken, err := s.command.HumanAddPasswordlessSetupInitCode(ctx, req.GetUserId(), ctxData.ResourceOwner, req.GetCodeId(), req.GetCode(), platform, generator)
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

func (s *Server) CreatePasskeyRegistrationLink(ctx context.Context, req *user.CreatePasskeyRegistrationLinkRequest) (*user.CreatePasskeyRegistrationLinkResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	passwordlessInitCode, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypePasswordlessInitCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	initCode, err := s.command.HumanSendPasswordlessInitCode(ctx, req.UserId, ctxData.OrgID, passwordlessInitCode)
	if err != nil {
		return nil, err
	}
	return &user.CreatePasskeyRegistrationLinkResponse{
		// TODO: Details: ...,
		CodeId: &initCode.CodeID,
		Code:   &initCode.Code,
	}, nil
}
