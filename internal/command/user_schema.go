package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/user/schema"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type CreateUserSchema struct {
	ResourceOwner          string
	Type                   string
	Schema                 map[string]any
	PossibleAuthenticators []domain.AuthenticatorType
}

func (s *CreateUserSchema) Valid() error {
	if s.Type == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMA-DGFj3", "Errors.UserSchema.Type.Missing") // TODO: i18n
	}
	if err := validateUserSchema(s.Schema); err != nil {
		return err
	}
	for _, authenticator := range s.PossibleAuthenticators {
		if authenticator == domain.AuthenticatorTypeUnspecified {
			return zerrors.ThrowInvalidArgument(nil, "COMMA-DGFj3", "Errors.UserSchema.Authenticator.Invalid") // TODO: i18n
		}
	}
	return nil
}

func (c *Commands) CreateUserSchema(ctx context.Context, userSchema *CreateUserSchema) (string, *domain.ObjectDetails, error) {
	if err := userSchema.Valid(); err != nil {
		return "", nil, err
	}
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewUserSchemaWriteModel(id, userSchema.ResourceOwner)
	err = c.pushAppendAndReduce(ctx, writeModel,
		schema.NewCreatedEvent(ctx,
			&schema.NewAggregate(id, userSchema.ResourceOwner).Aggregate,
			userSchema.Type, userSchema.Schema, userSchema.PossibleAuthenticators,
		),
	)
	if err != nil {
		return "", nil, err
	}
	return id, writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func validateUserSchema(schema map[string]any) error {
	return nil // TODO: impl
}
