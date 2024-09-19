package command

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	"github.com/zitadel/zitadel/v2/internal/crypto"
	"github.com/zitadel/zitadel/v2/internal/domain"
	domain_schema "github.com/zitadel/zitadel/v2/internal/domain/schema"
	"github.com/zitadel/zitadel/v2/internal/eventstore"
	"github.com/zitadel/zitadel/v2/internal/repository/user/schemauser"
	"github.com/zitadel/zitadel/v2/internal/zerrors"
)

type CreateSchemaUser struct {
	Details *domain.ObjectDetails

	SchemaID       string
	schemaRevision uint64

	ResourceOwner string
	ID            string
	Data          json.RawMessage

	Email           *Email
	ReturnCodeEmail string
	Phone           *Phone
	ReturnCodePhone string
}

func (s *CreateSchemaUser) Valid(ctx context.Context, c *Commands) (err error) {
	if s.ResourceOwner == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-urEJKa1tJM", "Errors.ResourceOwnerMissing")
	}
	if s.SchemaID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-TFo06JgnF2", "Errors.UserSchema.ID.Missing")
	}

	schemaWriteModel, err := c.getSchemaWriteModelByID(ctx, "", s.SchemaID)
	if err != nil {
		return err
	}
	if !schemaWriteModel.Exists() {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-N9QOuN4F7o", "Errors.UserSchema.NotExists")
	}
	s.schemaRevision = schemaWriteModel.SchemaRevision

	if s.ID == "" {
		s.ID, err = c.idGenerator.Next()
		if err != nil {
			return err
		}
	}

	// get role for permission check in schema through extension
	role, err := c.getSchemaRoleForWrite(ctx, s.ResourceOwner, s.ID)
	if err != nil {
		return err
	}

	schema, err := domain_schema.NewSchema(role, bytes.NewReader(schemaWriteModel.Schema))
	if err != nil {
		return err
	}

	var v interface{}
	if err := json.Unmarshal(s.Data, &v); err != nil {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-7o3ZGxtXUz", "Errors.User.Invalid")
	}

	if err := schema.Validate(v); err != nil {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-SlKXqLSeL6", "Errors.UserSchema.Data.Invalid")
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

func (c *Commands) CreateSchemaUser(ctx context.Context, user *CreateSchemaUser, alg crypto.EncryptionAlgorithm) (err error) {
	if err := user.Valid(ctx, c); err != nil {
		return err
	}

	writeModel, err := c.getSchemaUserExists(ctx, user.ResourceOwner, user.ID)
	if err != nil {
		return err
	}
	if writeModel.Exists() {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-Nn8CRVlkeZ", "Errors.User.AlreadyExists")
	}

	userAgg := UserV3AggregateFromWriteModel(&writeModel.WriteModel)
	events := []eventstore.Command{
		schemauser.NewCreatedEvent(ctx,
			userAgg,
			user.SchemaID, user.schemaRevision, user.Data,
		),
	}
	if user.Email != nil {
		events, user.ReturnCodeEmail, err = c.updateSchemaUserEmail(ctx, writeModel, events, userAgg, user.Email, alg)
		if err != nil {
			return err
		}
	}
	if user.Phone != nil {
		events, user.ReturnCodePhone, err = c.updateSchemaUserPhone(ctx, writeModel, events, userAgg, user.Phone, alg)
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

func (c *Commands) DeleteSchemaUser(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Vs4wJCME7T", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserExists(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}
	if !writeModel.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-syHyCsGmvM", "Errors.User.NotFound")
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

type ChangeSchemaUser struct {
	Details *domain.ObjectDetails

	SchemaID         *string
	schemaWriteModel *UserSchemaWriteModel

	ResourceOwner string
	ID            string
	Data          json.RawMessage

	Email           *Email
	ReturnCodeEmail string
	Phone           *Phone
	ReturnCodePhone string
}

func (s *ChangeSchemaUser) Valid(ctx context.Context, c *Commands) (err error) {
	if s.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-gEJR1QOGHb", "Errors.IDMissing")
	}
	if s.SchemaID != nil {
		s.schemaWriteModel, err = c.getSchemaWriteModelByID(ctx, "", *s.SchemaID)
		if err != nil {
			return err
		}
		if !s.schemaWriteModel.Exists() {
			return zerrors.ThrowPreconditionFailed(nil, "COMMAND-VLDTtxT3If", "Errors.UserSchema.NotExists")
		}
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

func (s *ChangeSchemaUser) ValidData(ctx context.Context, c *Commands, existingUser *UserV3WriteModel) (err error) {
	// get role for permission check in schema through extension
	role, err := c.getSchemaRoleForWrite(ctx, existingUser.ResourceOwner, existingUser.AggregateID)
	if err != nil {
		return err
	}

	if s.schemaWriteModel == nil {
		s.schemaWriteModel, err = c.getSchemaWriteModelByID(ctx, "", existingUser.SchemaID)
		if err != nil {
			return err
		}
	}

	schema, err := domain_schema.NewSchema(role, bytes.NewReader(s.schemaWriteModel.Schema))
	if err != nil {
		return err
	}

	// if data not changed but a new schema or revision should be used
	data := s.Data
	if s.Data == nil {
		data = existingUser.Data
	}

	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-7o3ZGxtXUz", "Errors.User.Invalid")
	}

	if err := schema.Validate(v); err != nil {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-SlKXqLSeL6", "Errors.UserSchema.Data.Invalid")
	}
	return nil
}

func (c *Commands) ChangeSchemaUser(ctx context.Context, user *ChangeSchemaUser, alg crypto.EncryptionAlgorithm) (err error) {
	if err := user.Valid(ctx, c); err != nil {
		return err
	}

	writeModel, err := c.getSchemaUserWriteModelByID(ctx, user.ResourceOwner, user.ID)
	if err != nil {
		return err
	}
	if !writeModel.Exists() {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-Nn8CRVlkeZ", "Errors.User.NotFound")
	}

	userAgg := UserV3AggregateFromWriteModel(&writeModel.WriteModel)
	events := make([]eventstore.Command, 0)
	if user.Data != nil || user.SchemaID != nil {
		if err := user.ValidData(ctx, c, writeModel); err != nil {
			return err
		}
		updateEvent := writeModel.NewUpdatedEvent(ctx,
			userAgg,
			user.schemaWriteModel.AggregateID,
			user.schemaWriteModel.SchemaRevision,
			user.Data,
		)
		if updateEvent != nil {
			events = append(events, updateEvent)
		}
	}
	if user.Email != nil {
		events, user.ReturnCodeEmail, err = c.updateSchemaUserEmail(ctx, writeModel, events, userAgg, user.Email, alg)
		if err != nil {
			return err
		}
	}
	if user.Phone != nil {
		events, user.ReturnCodePhone, err = c.updateSchemaUserPhone(ctx, writeModel, events, userAgg, user.Phone, alg)
		if err != nil {
			return err
		}
	}
	if len(events) == 0 {
		user.Details = writeModelToObjectDetails(&writeModel.WriteModel)
		return nil
	}
	if err := c.pushAppendAndReduce(ctx, writeModel, events...); err != nil {
		return err
	}
	user.Details = writeModelToObjectDetails(&writeModel.WriteModel)
	return nil
}

func (c *Commands) checkPermissionUpdateUserState(ctx context.Context, resourceOwner, userID string) error {
	return c.checkPermission(ctx, domain.PermissionUserWrite, resourceOwner, userID)
}

func (c *Commands) LockSchemaUser(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Eu8I2VAfjF", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserExists(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}
	if !writeModel.Exists() || writeModel.Locked {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-G4LOrnjY7q", "Errors.User.NotFound")
	}
	if err := c.checkPermissionUpdateUserState(ctx, writeModel.ResourceOwner, writeModel.AggregateID); err != nil {
		return nil, err
	}
	if err := c.pushAppendAndReduce(ctx, writeModel,
		schemauser.NewLockedEvent(ctx, UserV3AggregateFromWriteModel(&writeModel.WriteModel)),
	); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) UnlockSchemaUser(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-krXtYscQZh", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserExists(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}
	if !writeModel.Exists() || !writeModel.Locked {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-gpBv46Lh9m", "Errors.User.NotFound")
	}
	if err := c.checkPermissionUpdateUserState(ctx, writeModel.ResourceOwner, writeModel.AggregateID); err != nil {
		return nil, err
	}
	if err := c.pushAppendAndReduce(ctx, writeModel,
		schemauser.NewUnlockedEvent(ctx, UserV3AggregateFromWriteModel(&writeModel.WriteModel)),
	); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) DeactivateSchemaUser(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-pjJhge86ZV", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserExists(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}
	if writeModel.State != domain.UserStateActive {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-Ob6lR5iFTe", "Errors.User.NotFound")
	}
	if err := c.checkPermissionUpdateUserState(ctx, writeModel.ResourceOwner, writeModel.AggregateID); err != nil {
		return nil, err
	}
	if err := c.pushAppendAndReduce(ctx, writeModel,
		schemauser.NewDeactivatedEvent(ctx, UserV3AggregateFromWriteModel(&writeModel.WriteModel)),
	); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) ActivateSchemaUser(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-17XupGvxBJ", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserExists(ctx, "", id)
	if err != nil {
		return nil, err
	}
	if writeModel.State != domain.UserStateInactive {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-rQjbBr4J3j", "Errors.User.NotFound")
	}
	if err := c.checkPermissionUpdateUserState(ctx, writeModel.ResourceOwner, writeModel.AggregateID); err != nil {
		return nil, err
	}
	if err := c.pushAppendAndReduce(ctx, writeModel,
		schemauser.NewActivatedEvent(ctx, UserV3AggregateFromWriteModel(&writeModel.WriteModel)),
	); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) updateSchemaUserEmail(ctx context.Context, existing *UserV3WriteModel, events []eventstore.Command, agg *eventstore.Aggregate, email *Email, alg crypto.EncryptionAlgorithm) (_ []eventstore.Command, plainCode string, err error) {
	if existing.Email == string(email.Address) {
		return events, plainCode, nil
	}

	events = append(events, schemauser.NewEmailUpdatedEvent(ctx,
		agg,
		email.Address,
	))
	if email.Verified {
		events = append(events, schemauser.NewEmailVerifiedEvent(ctx, agg))
	} else {
		cryptoCode, err := c.newEmailCode(ctx, c.eventstore.Filter, alg) //nolint:staticcheck
		if err != nil {
			return nil, "", err
		}
		if email.ReturnCode {
			plainCode = cryptoCode.Plain
		}
		events = append(events, schemauser.NewEmailCodeAddedEvent(ctx, agg,
			cryptoCode.Crypted,
			cryptoCode.Expiry,
			email.URLTemplate,
			email.ReturnCode,
		))
	}
	return events, plainCode, nil
}

func (c *Commands) updateSchemaUserPhone(ctx context.Context, existing *UserV3WriteModel, events []eventstore.Command, agg *eventstore.Aggregate, phone *Phone, alg crypto.EncryptionAlgorithm) (_ []eventstore.Command, plainCode string, err error) {
	if existing.Phone == string(phone.Number) {
		return events, plainCode, nil
	}

	events = append(events, schemauser.NewPhoneUpdatedEvent(ctx,
		agg,
		phone.Number,
	))
	if phone.Verified {
		events = append(events, schemauser.NewPhoneVerifiedEvent(ctx, agg))
	} else {
		cryptoCode, err := c.newPhoneCode(ctx, c.eventstore.Filter, alg) //nolint:staticcheck
		if err != nil {
			return nil, "", err
		}
		if phone.ReturnCode {
			plainCode = cryptoCode.Plain
		}
		events = append(events, schemauser.NewPhoneCodeAddedEvent(ctx, agg,
			cryptoCode.Crypted,
			cryptoCode.Expiry,
			phone.ReturnCode,
		))
	}
	return events, plainCode, nil
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
