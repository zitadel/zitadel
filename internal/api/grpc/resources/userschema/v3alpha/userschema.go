package userschema

import (
	"context"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/api/authz"
	resource_object "github.com/zitadel/zitadel/internal/api/grpc/resources/object/v3alpha"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	schema "github.com/zitadel/zitadel/pkg/grpc/resources/userschema/v3alpha"
)

func (s *Server) CreateUserSchema(ctx context.Context, req *schema.CreateUserSchemaRequest) (*schema.CreateUserSchemaResponse, error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	instanceID := authz.GetInstance(ctx).InstanceID()
	userSchema, err := createUserSchemaToCommand(req, instanceID)
	if err != nil {
		return nil, err
	}

	if err := s.command.CreateUserSchema(ctx, userSchema); err != nil {
		return nil, err
	}
	return &schema.CreateUserSchemaResponse{
		Details: resource_object.DomainToDetailsPb(userSchema.Details, object.OwnerType_OWNER_TYPE_INSTANCE, instanceID),
	}, nil
}

func (s *Server) PatchUserSchema(ctx context.Context, req *schema.PatchUserSchemaRequest) (*schema.PatchUserSchemaResponse, error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	instanceID := authz.GetInstance(ctx).InstanceID()
	userSchema, err := patchUserSchemaToCommand(req, instanceID)
	if err != nil {
		return nil, err
	}
	if err := s.command.ChangeUserSchema(ctx, userSchema); err != nil {
		return nil, err
	}
	return &schema.PatchUserSchemaResponse{
		Details: resource_object.DomainToDetailsPb(userSchema.Details, object.OwnerType_OWNER_TYPE_INSTANCE, instanceID),
	}, nil
}

func (s *Server) DeactivateUserSchema(ctx context.Context, req *schema.DeactivateUserSchemaRequest) (*schema.DeactivateUserSchemaResponse, error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	instanceID := authz.GetInstance(ctx).InstanceID()
	details, err := s.command.DeactivateUserSchema(ctx, req.GetId(), instanceID)
	if err != nil {
		return nil, err
	}
	return &schema.DeactivateUserSchemaResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_INSTANCE, instanceID),
	}, nil
}

func (s *Server) ReactivateUserSchema(ctx context.Context, req *schema.ReactivateUserSchemaRequest) (*schema.ReactivateUserSchemaResponse, error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	instanceID := authz.GetInstance(ctx).InstanceID()
	details, err := s.command.ReactivateUserSchema(ctx, req.GetId(), instanceID)
	if err != nil {
		return nil, err
	}
	return &schema.ReactivateUserSchemaResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_INSTANCE, instanceID),
	}, nil
}

func (s *Server) DeleteUserSchema(ctx context.Context, req *schema.DeleteUserSchemaRequest) (*schema.DeleteUserSchemaResponse, error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	instanceID := authz.GetInstance(ctx).InstanceID()
	details, err := s.command.DeleteUserSchema(ctx, req.GetId(), instanceID)
	if err != nil {
		return nil, err
	}
	return &schema.DeleteUserSchemaResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_INSTANCE, instanceID),
	}, nil
}

func createUserSchemaToCommand(req *schema.CreateUserSchemaRequest, resourceOwner string) (*command.CreateUserSchema, error) {
	schema, err := req.GetUserSchema().GetSchema().MarshalJSON()
	if err != nil {
		return nil, err
	}
	return &command.CreateUserSchema{
		ResourceOwner:          resourceOwner,
		Type:                   req.GetUserSchema().GetType(),
		Schema:                 schema,
		PossibleAuthenticators: authenticatorsToDomain(req.GetUserSchema().GetPossibleAuthenticators()),
	}, nil
}

func patchUserSchemaToCommand(req *schema.PatchUserSchemaRequest, resourceOwner string) (*command.ChangeUserSchema, error) {
	schema, err := req.GetUserSchema().GetSchema().MarshalJSON()
	if err != nil {
		return nil, err
	}

	var ty *string
	if req.GetUserSchema() != nil && req.GetUserSchema().GetType() != "" {
		ty = gu.Ptr(req.GetUserSchema().GetType())
	}
	return &command.ChangeUserSchema{
		ID:                     req.GetId(),
		ResourceOwner:          resourceOwner,
		Type:                   ty,
		Schema:                 schema,
		PossibleAuthenticators: authenticatorsToDomain(req.GetUserSchema().GetPossibleAuthenticators()),
	}, nil
}

func authenticatorsToDomain(authenticators []schema.AuthenticatorType) []domain.AuthenticatorType {
	if authenticators == nil {
		return nil
	}
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
