package command

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/domain"
	domain_schema "github.com/zitadel/zitadel/internal/domain/schema"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/user/schema"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type CreateUserSchema struct {
	ResourceOwner          string
	Type                   string
	Schema                 json.RawMessage
	PossibleAuthenticators []domain.AuthenticatorType
}

func (s *CreateUserSchema) Valid() error {
	if s.Type == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMA-DGFj3", "Errors.UserSchema.Type.Missing")
	}
	if err := validateUserSchema(s.Schema); err != nil {
		return err
	}
	for _, authenticator := range s.PossibleAuthenticators {
		if authenticator == domain.AuthenticatorTypeUnspecified {
			return zerrors.ThrowInvalidArgument(nil, "COMMA-Gh652", "Errors.UserSchema.Authenticator.Invalid")
		}
	}
	return nil
}

type UpdateUserSchema struct {
	ID                     string
	ResourceOwner          string
	Type                   *string
	Schema                 json.RawMessage
	PossibleAuthenticators []domain.AuthenticatorType
}

func (s *UpdateUserSchema) Valid() error {
	if s.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMA-H5421", "Errors.IDMissing")
	}
	if s.Type != nil && *s.Type == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMA-G43gn", "Errors.UserSchema.Type.Missing")
	}
	if err := validateUserSchema(s.Schema); err != nil {
		return err
	}
	for _, authenticator := range s.PossibleAuthenticators {
		if authenticator == domain.AuthenticatorTypeUnspecified {
			return zerrors.ThrowInvalidArgument(nil, "COMMA-WF4hg", "Errors.UserSchema.Authenticator.Invalid")
		}
	}
	return nil
}

func (c *Commands) CreateUserSchema(ctx context.Context, userSchema *CreateUserSchema) (string, *domain.ObjectDetails, error) {
	if err := userSchema.Valid(); err != nil {
		return "", nil, err
	}
	if userSchema.ResourceOwner == "" {
		return "", nil, zerrors.ThrowInvalidArgument(nil, "COMMA-J3hhj", "Errors.ResourceOwnerMissing")
	}
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewUserSchemaWriteModel(id, userSchema.ResourceOwner)
	err = c.pushAppendAndReduce(ctx, writeModel,
		schema.NewCreatedEvent(ctx,
			UserSchemaAggregateFromWriteModel(&writeModel.WriteModel),
			userSchema.Type, userSchema.Schema, userSchema.PossibleAuthenticators,
		),
	)
	if err != nil {
		return "", nil, err
	}
	return id, writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) UpdateUserSchema(ctx context.Context, userSchema *UpdateUserSchema) (*domain.ObjectDetails, error) {
	if err := userSchema.Valid(); err != nil {
		return nil, err
	}
	writeModel := NewUserSchemaWriteModel(userSchema.ID, userSchema.ResourceOwner)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	if writeModel.State != domain.UserSchemaStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMA-HB3e1", "Errors.UserSchema.NotActive")
	}
	updatedEvent := writeModel.NewUpdatedEvent(
		ctx,
		UserSchemaAggregateFromWriteModel(&writeModel.WriteModel),
		userSchema.Type,
		userSchema.Schema,
		userSchema.PossibleAuthenticators,
	)
	if updatedEvent == nil {
		return writeModelToObjectDetails(&writeModel.WriteModel), nil
	}
	if err := c.pushAppendAndReduce(ctx, writeModel, updatedEvent); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) DeactivateUserSchema(ctx context.Context, id, resourceOwner string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMA-Vvf3w", "Errors.IDMissing")
	}
	writeModel := NewUserSchemaWriteModel(id, resourceOwner)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	if writeModel.State != domain.UserSchemaStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMA-E4t4z", "Errors.UserSchema.NotActive")
	}
	err := c.pushAppendAndReduce(ctx, writeModel,
		schema.NewDeactivatedEvent(ctx, UserSchemaAggregateFromWriteModel(&writeModel.WriteModel)),
	)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) ReactivateUserSchema(ctx context.Context, id, resourceOwner string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMA-wq3Gw", "Errors.IDMissing")
	}
	writeModel := NewUserSchemaWriteModel(id, resourceOwner)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	if writeModel.State != domain.UserSchemaStateInactive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMA-DGzh5", "Errors.UserSchema.NotInactive")
	}
	err := c.pushAppendAndReduce(ctx, writeModel,
		schema.NewReactivatedEvent(ctx, UserSchemaAggregateFromWriteModel(&writeModel.WriteModel)),
	)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) DeleteUserSchema(ctx context.Context, id, resourceOwner string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMA-E22gg", "Errors.IDMissing")
	}
	writeModel := NewUserSchemaWriteModel(id, resourceOwner)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	if !writeModel.Exists() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMA-Grg41", "Errors.UserSchema.NotExists")
	}
	// TODO: check for users based on that schema; this is only possible with / after https://github.com/zitadel/zitadel/issues/7308
	err := c.pushAppendAndReduce(ctx, writeModel,
		schema.NewDeletedEvent(ctx, UserSchemaAggregateFromWriteModel(&writeModel.WriteModel), writeModel.SchemaType),
	)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func validateUserSchema(userSchema json.RawMessage) error {
	_, err := domain_schema.NewSchema(0, bytes.NewReader(userSchema))
	if err != nil {
		return zerrors.ThrowInvalidArgument(err, "COMMA-W21tg", "Errors.UserSchema.Schema.Invalid")
	}
	return nil
}
