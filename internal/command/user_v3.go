package command

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	domain_schema "github.com/zitadel/zitadel/internal/domain/schema"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user/schemauser"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type CreateSchemaUser struct {
	Details       *domain.ObjectDetails
	ResourceOwner string

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

func (s *CreateSchemaUser) Valid(ctx context.Context, c *Commands) (err error) {
	if s.ResourceOwner == "" {
		return zerrors.ThrowInvalidArgument(nil, "TODO", "Errors.ResourceOwnerMissing")
	}
	if s.SchemaID == "" {
		return zerrors.ThrowInvalidArgument(nil, "TODO", "Errors.UserSchema.User.Type.Missing")
	}

	schemaWriteModel, err := c.getSchemaWriteModelByID(ctx, "", s.SchemaID)
	if err != nil {
		return err
	}
	if !schemaWriteModel.Exists() {
		return zerrors.ThrowPreconditionFailed(nil, "TODO", "TODO")
	}
	s.schemaRevision = schemaWriteModel.Revision

	if s.ID == "" {
		s.ID, err = c.idGenerator.Next()
		if err != nil {
			return err
		}
	}

	role := domain_schema.RoleUnspecified
	// if register then the role is directly self, and check if it is permitted to register is checked before
	if s.Register {
		role = domain_schema.RoleSelf
	} else {
		// get role for permission check in schema through extension
		role, err = c.getSchemaRoleForWrite(ctx, s.ResourceOwner, s.ID)
		if err != nil {
			return err
		}
	}

	schema, err := domain_schema.NewSchema(role, bytes.NewReader(schemaWriteModel.Schema))
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

func (c *Commands) getSchemaRoleForWrite(ctx context.Context, resourceOwner, userID string) (domain_schema.Role, error) {
	if userID == authz.GetCtxData(ctx).UserID {
		return domain_schema.RoleSelf, nil
	}
	if err := c.checkPermission(ctx, domain.PermissionUserWrite, resourceOwner, userID); err != nil {
		return domain_schema.RoleUnspecified, err
	}
	return domain_schema.RoleOwner, nil
}

func (c *Commands) CreateSchemaUser(ctx context.Context, user *CreateSchemaUser, emailCodeGenerator, phoneCodeGenerator crypto.Generator) (err error) {
	if err := user.Valid(ctx, c); err != nil {
		return err
	}

	writeModel, err := c.getSchemaUserWriteModelByID(ctx, user.ResourceOwner, user.ID)
	if err != nil {
		return err
	}
	if writeModel.Exists() {
		return zerrors.ThrowPreconditionFailed(nil, "TODO", "TODO")
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
