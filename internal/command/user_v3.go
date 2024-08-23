package command

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	domain_schema "github.com/zitadel/zitadel/internal/domain/schema"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user/schemauser"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type CreateSchemaUser struct {
	Details *domain.ObjectDetails

	Register bool

	SchemaID       string
	schemaRevision uint64

	ID   string
	Data json.RawMessage

	Email           *Email
	ReturnCodeEmail string
	Phone           *Phone
	ReturnCodePhone string
}

func (s *CreateSchemaUser) Valid(ctx context.Context, c *Commands, resourceOwner string) error {
	if s.SchemaID == "" {
		return zerrors.ThrowInvalidArgument(nil, "TODO", "Errors.UserSchema.User.Type.Missing")
	}

	schemaWriteModel, err := c.getSchemaWriteModelByType(ctx, resourceOwner, s.SchemaID)
	if err != nil {
		return err
	}
	if !schemaWriteModel.Exists() {
		return zerrors.ThrowPreconditionFailed(nil, "TODO", "TODO")
	}
	s.schemaRevision = schemaWriteModel.Revision

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

	if s.Email != nil && s.Email.Address != "" {
		if err := s.Email.Validate(); err != nil {
			return err
		}
	}

	if s.Phone != nil && s.Phone.Number != "" {
		if s.Phone.Number, err = s.Phone.Number.Normalize(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Commands) CreateSchemaUser(ctx context.Context, resourceOwner string, user *CreateSchemaUser, emailCodeGenerator, phoneCodeGenerator crypto.Generator) (err error) {
	if resourceOwner == "" {
		return zerrors.ThrowInvalidArgument(nil, "TODO", "Errors.ResourceOwnerMissing")
	}
	if err := user.Valid(ctx, c, resourceOwner); err != nil {
		return err
	}

	if user.ID == "" {
		user.ID, err = c.idGenerator.Next()
		if err != nil {
			return err
		}
	}

	writeModel, err := c.getSchemaUserWriteModelByID(ctx, resourceOwner, user.ID)
	if err != nil {
		return err
	}
	if writeModel.Exists() {
		return zerrors.ThrowPreconditionFailed(nil, "TODO", "TODO")
	}
	if !user.Register {
		if err := c.checkPermission(ctx, domain.PermissionUserWrite, writeModel.ResourceOwner, writeModel.AggregateID); err != nil {
			return err
		}
	}

	userAgg := UserV3AggregateFromWriteModel(&writeModel.WriteModel)
	events := []eventstore.Command{
		schemauser.NewCreatedEvent(ctx,
			userAgg,
			user.SchemaID, user.schemaRevision, user.Data,
		),
	}
	if user.Email != nil {
		events, user.ReturnCodeEmail, err = c.updateSchemaUserEmail(ctx, events, userAgg, user.Email, emailCodeGenerator)
		if err != nil {
			return err
		}
	}
	if user.Phone != nil {
		events, user.ReturnCodePhone, err = c.updateSchemaUserPhone(ctx, events, userAgg, user.Phone, phoneCodeGenerator)
		if err != nil {
			return err
		}
	}

	if err := c.pushAppendAndReduce(ctx, writeModel, events...); err != nil {
		return err
	}
	user.Details = writeModelToObjectDetails(&writeModel.WriteModel)
	return nil
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
		return nil, zerrors.ThrowNotFound(nil, "TODO", "TODO")
	}
	if err := c.checkPermissionDeleteUser(ctx, writeModel.ResourceOwner, writeModel.AggregateID); err != nil {
		return nil, err
	}
	if err := c.pushAppendAndReduce(ctx, writeModel,
		schemauser.NewDeletedEvent(ctx, UserV3AggregateFromWriteModel(&writeModel.WriteModel)),
	); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) updateSchemaUserEmail(ctx context.Context, events []eventstore.Command, agg *eventstore.Aggregate, email *Email, codeGenerator crypto.Generator) (_ []eventstore.Command, plainCode string, err error) {
	events = append(events, schemauser.NewEmailChangedEvent(ctx,
		agg,
		email.Address,
	))
	if email.Verified {
		events = append(events, schemauser.NewEmailVerifiedEvent(ctx, agg))
	} else {
		returnCode, code, err := generateCode(codeGenerator, email.ReturnCode)
		if err != nil {
			return nil, "", err
		}
		plainCode = code
		events = append(events, schemauser.NewEmailCodeAddedEvent(ctx, agg,
			returnCode,
			codeGenerator.Expiry(),
			email.URLTemplate,
			email.ReturnCode,
		))
	}
	return events, plainCode, nil
}

func (c *Commands) updateSchemaUserPhone(ctx context.Context, events []eventstore.Command, agg *eventstore.Aggregate, phone *Phone, codeGenerator crypto.Generator) (_ []eventstore.Command, plainCode string, err error) {
	events = append(events, schemauser.NewPhoneChangedEvent(ctx,
		agg,
		phone.Number,
	))
	if phone.Verified {
		events = append(events, schemauser.NewPhoneVerifiedEvent(ctx, agg))
	} else {
		returnCode, code, err := generateCode(codeGenerator, phone.ReturnCode)
		if err != nil {
			return nil, "", err
		}
		plainCode = code
		events = append(events, schemauser.NewPhoneCodeAddedEvent(ctx, agg,
			returnCode,
			codeGenerator.Expiry(),
			phone.ReturnCode,
		))
	}
	return events, plainCode, nil
}

func generateCode(gen crypto.Generator, returnCode bool) (*crypto.CryptoValue, string, error) {
	value, plain, err := crypto.NewCode(gen)
	if err != nil {
		return nil, "", err
	}

	if returnCode {
		return value, plain, nil
	}
	return value, "", nil

}

func (c *Commands) getSchemaUserExists(ctx context.Context, resourceOwner, id string) (*UserV3WriteModel, error) {
	writeModel := NewExistsUserV3WriteModel(resourceOwner, id)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	return writeModel, nil
}

func (c *Commands) getSchemaUserWriteModelByID(ctx context.Context, resourceOwner, id string) (*UserV3WriteModel, error) {
	writeModel := NewUserV3WriteModel(resourceOwner, id)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	return writeModel, nil
}
