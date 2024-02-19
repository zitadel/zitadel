package schema

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
	schema "github.com/zitadel/zitadel/pkg/grpc/user/schema/v3alpha"
)

func (s *Server) CreateUserSchema(ctx context.Context, req *schema.CreateUserSchemaRequest) (*schema.CreateUserSchemaResponse, error) {
	userSchema := createUserSchemaToCommand(ctx, req)
	id, details, err := s.command.CreateUserSchema(ctx, userSchema)
	if err != nil {
		return nil, err
	}
	return &schema.CreateUserSchemaResponse{
		Id:      id,
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) UpdateUserSchema(ctx context.Context, req *schema.UpdateUserSchemaRequest) (*schema.UpdateUserSchemaResponse, error) {
	return nil, zerrors.ThrowUnimplemented(nil, "", "")
}
func (s *Server) DeactivateUserSchema(ctx context.Context, req *schema.DeactivateUserSchemaRequest) (*schema.DeactivateUserSchemaResponse, error) {
	return nil, zerrors.ThrowUnimplemented(nil, "", "")
}
func (s *Server) ReactivateUserSchema(ctx context.Context, req *schema.ReactivateUserSchemaRequest) (*schema.ReactivateUserSchemaResponse, error) {
	return nil, zerrors.ThrowUnimplemented(nil, "", "")
}
func (s *Server) DeleteUserSchema(ctx context.Context, req *schema.DeleteUserSchemaRequest) (*schema.DeleteUserSchemaResponse, error) {
	return nil, zerrors.ThrowUnimplemented(nil, "", "")
}

func createUserSchemaToCommand(ctx context.Context, req *schema.CreateUserSchemaRequest) *command.CreateUserSchema {
	return &command.CreateUserSchema{
		ResourceOwner:          authz.GetInstance(ctx).InstanceID(),
		Type:                   req.GetType(),
		Schema:                 req.GetSchema().AsMap(),
		PossibleAuthenticators: authenticatorsToDomain(req.GetPossibleAuthenticators()),
	}
}

func authenticatorsToDomain(authenticators []schema.AuthenticatorType) []domain.AuthenticatorType {
	types := make([]domain.AuthenticatorType, len(authenticators))
	for i, authenticator := range authenticators {
		types[i] = authenticatorTypeToDomain(authenticator)
	}
	return types
}

func authenticatorTypeToDomain(authenticator schema.AuthenticatorType) domain.AuthenticatorType {
	switch authenticator {
	case schema.AuthenticatorType_AUTHENTICATOR_TYPE_UNSPECIFIED:
		return domain.AuthenticatorTypeUnspecified
	case schema.AuthenticatorType_AUTHENTICATOR_TYPE_USERNAME:
		return domain.AuthenticatorTypeUsername
	case schema.AuthenticatorType_AUTHENTICATOR_TYPE_PASSWORD:
		return domain.AuthenticatorTypePassword
	case schema.AuthenticatorType_AUTHENTICATOR_TYPE_WEBAUTHN:
		return domain.AuthenticatorTypeWebAuthN
	case schema.AuthenticatorType_AUTHENTICATOR_TYPE_TOTP:
		return domain.AuthenticatorTypeTOTP
	case schema.AuthenticatorType_AUTHENTICATOR_TYPE_OTP_EMAIL:
		return domain.AuthenticatorTypeOTPEmail
	case schema.AuthenticatorType_AUTHENTICATOR_TYPE_OTP_SMS:
		return domain.AuthenticatorTypeOTPSMS
	case schema.AuthenticatorType_AUTHENTICATOR_TYPE_AUTHENTICATION_KEY:
		return domain.AuthenticatorTypeAuthenticationKey
	case schema.AuthenticatorType_AUTHENTICATOR_TYPE_IDENTITY_PROVIDER:
		return domain.AuthenticatorTypeIdentityProvider
	default:
		return domain.AuthenticatorTypeUnspecified
	}
}
