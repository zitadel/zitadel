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

	SchemaID       string
	schemaRevision uint64

	ID   string
	Data json.RawMessage

	Email           *Email
	ReturnCodeEmail *domain.EmailCode
	Phone           *Phone
	ReturnCodePhone *domain.PhoneCode
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
	if s.schemaRevision == 0 {
		return zerrors.ThrowInvalidArgument(nil, "TODO", "Errors.UserSchema.User.Revision.Missing")
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

	userAgg := UserV3AggregateFromWriteModel(&writeModel.WriteModel)
	events := []eventstore.Command{
		schemauser.NewCreatedEvent(ctx,
			userAgg,
			user.SchemaID, user.schemaRevision, user.Data,
		),
	}
	if user.Email != nil {
		events = append(events, schemauser.NewEmailChangedEvent(ctx,
			userAgg,
			user.Email.Address,
		))
		if user.Email.Verified {
			events = append(events, schemauser.NewEmailVerifiedEvent(ctx, userAgg))
		} else {
			user.ReturnCodeEmail, _, err = domain.NewEmailCode(emailCodeGenerator)
			if err != nil {
				return err
			}
			events = append(events, schemauser.NewEmailCodeAddedEvent(ctx, userAgg,
				user.ReturnCodeEmail.Code,
				user.ReturnCodeEmail.Expiry,
				user.Email.URLTemplate,
				user.Email.ReturnCode,
			))
		}
	}
	if user.Phone != nil {
		events = append(events, schemauser.NewPhoneChangedEvent(ctx,
			userAgg,
			user.Phone.Number,
		))
		if user.Phone.Verified {
			events = append(events, schemauser.NewPhoneVerifiedEvent(ctx, userAgg))
		} else {
			user.ReturnCodePhone, err = domain.NewPhoneCode(phoneCodeGenerator)
			if err != nil {
				return err
			}
			events = append(events, schemauser.NewPhoneCodeAddedEvent(ctx, userAgg,
				user.ReturnCodePhone.Code,
				user.ReturnCodePhone.Expiry,
				user.Phone.ReturnCode,
			))
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
