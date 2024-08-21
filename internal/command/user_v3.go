package command

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/domain"
	domain_schema "github.com/zitadel/zitadel/internal/domain/schema"
	"github.com/zitadel/zitadel/internal/repository/user/schemauser"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type CreateSchemaUser struct {
	SchemaType     string
	SchemaRevision uint64

	Data json.RawMessage
}

func (s *CreateSchemaUser) Valid(ctx context.Context, c *Commands, resourceOwner string) error {
	if s.SchemaType == "" {
		return zerrors.ThrowInvalidArgument(nil, "TODO", "Errors.UserSchema.User.Type.Missing")
	}
	if s.SchemaRevision == 0 {
		return zerrors.ThrowInvalidArgument(nil, "TODO", "Errors.UserSchema.User.Revision.Missing")
	}

	schemaWriteModel, err := c.getSchemaWriteModelByType(ctx, resourceOwner, s.SchemaType)
	if err != nil {
		return err
	}
	if !schemaWriteModel.Exists() {
		return zerrors.ThrowPreconditionFailed(nil, "TODO", "TODO")
	}

	schema, err := domain_schema.NewSchema(0, bytes.NewReader(schemaWriteModel.Schema))
	if err != nil {
		return err
	}

	var v interface{}
	if err := json.Unmarshal(s.Data, &v); err != nil {
		return zerrors.ThrowInvalidArgument(nil, "TODO", "TODO")
	}

	if err := schema.Validate(v); err != nil {
		return zerrors.ThrowPreconditionFailed(nil, "TODO", "TODO")
	}
	return nil
}

func (c *Commands) CreateSchemaUser(ctx context.Context, resourceOwner string, user *CreateSchemaUser) (string, *domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return "", nil, zerrors.ThrowInvalidArgument(nil, "TODO", "Errors.ResourceOwnerMissing")
	}
	if err := user.Valid(ctx, c, resourceOwner); err != nil {
		return "", nil, err
	}

	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel, err := c.getSchemaUserWriteModelByID(ctx, resourceOwner, id)
	if err != nil {
		return "", nil, err
	}

	if err := c.pushAppendAndReduce(ctx, writeModel,
		schemauser.NewCreatedEvent(ctx,
			UserV3AggregateFromWriteModel(&writeModel.WriteModel),
			user.SchemaType, user.SchemaRevision, "", "", user.Data,
		),
	); err != nil {
		return "", nil, err
	}
	return id, writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) DeleteSchemaUser(ctx context.Context, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "TODO", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserWriteModelByID(ctx, "", id)
	if err != nil {
		return nil, err
	}
	if !writeModel.Exists() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "TODO", "Errors.UserSchema.User.NotExists")
	}
	if err := c.pushAppendAndReduce(ctx, writeModel,
		schemauser.NewDeletedEvent(ctx, UserV3AggregateFromWriteModel(&writeModel.WriteModel)),
	); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) getSchemaUserWriteModelByID(ctx context.Context, resourceOwner, id string) (*UserV3WriteModel, error) {
	writeModel := NewUserV3WriteModel(resourceOwner, id)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	return writeModel, nil
}
