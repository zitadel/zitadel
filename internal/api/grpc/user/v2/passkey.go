package user

import (
	"context"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

func (s *Server) RegisterPasskey(ctx context.Context, req *user.RegisterPasskeyRequest) (resp *user.RegisterPasskeyResponse, err error) {
	var (
		resourceOwner = authz.GetCtxData(ctx).ResourceOwner
		authenticator = passkeyAuthenticatorToDomain(req.GetAuthenticator())
	)
	if code := req.GetCode(); code != nil {
		return passkeyRegistrationDetailsToPb(
			s.command.RegisterUserPasskeyWithCode(ctx, req.GetUserId(), resourceOwner, authenticator, code.Id, code.Code, s.userCodeAlg),
		)
	}
	return passkeyRegistrationDetailsToPb(
		s.command.RegisterUserPasskey(ctx, req.GetUserId(), resourceOwner, authenticator),
	)
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

func passkeyRegistrationDetailsToPb(details *domain.PasskeyRegistrationDetails, err error) (*user.RegisterPasskeyResponse, error) {
	if err != nil {
		return nil, err
	}
	options := new(structpb.Struct)
	if err := protojson.Unmarshal(details.PublicKeyCredentialCreationOptions, options); err != nil {
		return nil, caos_errs.ThrowInternal(err, "USERv2-Dohr6", "Errors.Internal")
	}
	return &user.RegisterPasskeyResponse{
		Details:                            object.DomainToDetailsPb(details.ObjectDetails),
		PasskeyId:                          details.PasskeyID,
		PublicKeyCredentialCreationOptions: options,
	}, nil
}

func (s *Server) VerifyPasskeyRegistration(ctx context.Context, req *user.VerifyPasskeyRegistrationRequest) (*user.VerifyPasskeyRegistrationResponse, error) {
	resourceOwner := authz.GetCtxData(ctx).ResourceOwner
	pkc, err := protojson.Marshal(req.GetPublicKeyCredential())
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "USERv2-Pha2o", "Errors.Internal")
	}
	objectDetails, err := s.command.HumanHumanPasswordlessSetup(ctx, req.GetUserId(), resourceOwner, req.GetPasskeyName(), "", pkc)
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
