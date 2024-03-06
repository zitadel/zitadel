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
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	userSchema := createUserSchemaToCommand(req)
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
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	userSchema := updateUserSchemaToCommand(req)
	details, err := s.command.UpdateUserSchema(ctx, userSchema)
	if err != nil {
		return nil, err
	}
	return &schema.UpdateUserSchemaResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}
func (s *Server) DeactivateUserSchema(ctx context.Context, req *schema.DeactivateUserSchemaRequest) (*schema.DeactivateUserSchemaResponse, error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.DeactivateUserSchema(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &schema.DeactivateUserSchemaResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}
func (s *Server) ReactivateUserSchema(ctx context.Context, req *schema.ReactivateUserSchemaRequest) (*schema.ReactivateUserSchemaResponse, error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.ReactivateUserSchema(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &schema.ReactivateUserSchemaResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}
func (s *Server) DeleteUserSchema(ctx context.Context, req *schema.DeleteUserSchemaRequest) (*schema.DeleteUserSchemaResponse, error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.DeleteUserSchema(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &schema.DeleteUserSchemaResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func checkUserSchemaEnabled(ctx context.Context) error {
	if authz.GetInstance(ctx).Features().UserSchema {
		return nil
	}
	return zerrors.ThrowPreconditionFailed(nil, "SCHEMA-SFjk3", "Errors.UserSchema.NotEnabled")
}

func createUserSchemaToCommand(req *schema.CreateUserSchemaRequest) *command.CreateUserSchema {
	return &command.CreateUserSchema{
		Type:                   req.GetType(),
		Schema:                 req.GetSchema().AsMap(),
		PossibleAuthenticators: authenticatorsToDomain(req.GetPossibleAuthenticators()),
	}
}

func updateUserSchemaToCommand(req *schema.UpdateUserSchemaRequest) *command.UpdateUserSchema {
	return &command.UpdateUserSchema{
		ID:                     req.GetId(),
		Type:                   req.Type,
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
