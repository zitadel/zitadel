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
	s.schemaRevision = schemaWriteModel.Revision

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
		events, user.ReturnCodeEmail, err = c.updateSchemaUserEmail(ctx, events, userAgg, user.Email, alg)
		if err != nil {
			return err
		}
	}
	if user.Phone != nil {
		events, user.ReturnCodePhone, err = c.updateSchemaUserPhone(ctx, events, userAgg, user.Phone, alg)
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
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Vs4wJCME7T", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserExists(ctx, "", id)
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

type UpdateSchemaUser struct {
	Details *domain.ObjectDetails

	SchemaID         *string
	schemaWriteModel *UserSchemaWriteModel

	ID   string
	Data json.RawMessage

	Email           *Email
	ReturnCodeEmail string
	Phone           *Phone
	ReturnCodePhone string
}

func (s *UpdateSchemaUser) Valid(ctx context.Context, c *Commands) (err error) {
	if s.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-gEJR1QOGHb", "TODO")
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

func (s *UpdateSchemaUser) ValidData(ctx context.Context, c *Commands, existingUser *UserV3WriteModel) (err error) {
	if s.schemaWriteModel == nil {
		s.schemaWriteModel, err = c.getSchemaWriteModelByIDAndRevision(ctx, "", existingUser.SchemaID, existingUser.SchemaRevision)
		if err != nil {
			return err
		}
	}

	// get role for permission check in schema through extension
	role, err := c.getSchemaRoleForWrite(ctx, existingUser.ResourceOwner, existingUser.AggregateID)
	if err != nil {
		return err
	}

	schema, err := domain_schema.NewSchema(role, bytes.NewReader(s.schemaWriteModel.Schema))
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
}

func (c *Commands) UpdateSchemaUser(ctx context.Context, user *UpdateSchemaUser, alg crypto.EncryptionAlgorithm) (err error) {
	if err := user.Valid(ctx, c); err != nil {
		return err
	}

	writeModel, err := c.getSchemaUserWriteModelByID(ctx, "", user.ID)
	if err != nil {
		return err
	}
	if writeModel.Exists() {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-Nn8CRVlkeZ", "Errors.User.AlreadyExists")
	}

	userAgg := UserV3AggregateFromWriteModel(&writeModel.WriteModel)
	writeModel.NewUpdatedEvent(ctx,
		userAgg,
		user.schemaWriteModel.AggregateID,
		user.schemaWriteModel.Revision,
		user.Data,
	)
	events := []eventstore.Command{
		schemauser.NewCreatedEvent(ctx,
			userAgg,
		),
	}
	if user.Email != nil {
		events, user.ReturnCodeEmail, err = c.updateSchemaUserEmail(ctx, events, userAgg, user.Email, alg)
		if err != nil {
			return err
		}
	}
	if user.Phone != nil {
		events, user.ReturnCodePhone, err = c.updateSchemaUserPhone(ctx, events, userAgg, user.Phone, alg)
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

func (c *Commands) LockedSchemaUser(ctx context.Context, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Eu8I2VAfjF", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserExists(ctx, "", id)
	if err != nil {
		return nil, err
	}
	if !writeModel.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-G4LOrnjY7q", "Errors.User.NotFound")
	}
	if writeModel.State != domain.UserStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "TODO", "TODO")
	}
	if err := c.checkPermissionUpdateUser(ctx, writeModel.ResourceOwner, writeModel.AggregateID); err != nil {
		return nil, err
	}
	if err := c.pushAppendAndReduce(ctx, writeModel,
		schemauser.NewLockedEvent(ctx, UserV3AggregateFromWriteModel(&writeModel.WriteModel)),
	); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) UnlockSchemaUser(ctx context.Context, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-krXtYscQZh", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserExists(ctx, "", id)
	if err != nil {
		return nil, err
	}
	if !writeModel.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-gpBv46Lh9m", "Errors.User.NotFound")
	}
	if writeModel.State != domain.UserStateLocked {
		return nil, zerrors.ThrowPreconditionFailed(nil, "TODO", "TODO")
	}
	if err := c.checkPermissionUpdateUser(ctx, writeModel.ResourceOwner, writeModel.AggregateID); err != nil {
		return nil, err
	}
	if err := c.pushAppendAndReduce(ctx, writeModel,
		schemauser.NewUnlockedEvent(ctx, UserV3AggregateFromWriteModel(&writeModel.WriteModel)),
	); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) DeactiveSchemaUser(ctx context.Context, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-pjJhge86ZV", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserExists(ctx, "", id)
	if err != nil {
		return nil, err
	}
	if !writeModel.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-Ob6lR5iFTe", "Errors.User.NotFound")
	}
	if writeModel.State != domain.UserStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "TODO", "TODO")
	}
	if err := c.checkPermissionUpdateUser(ctx, writeModel.ResourceOwner, writeModel.AggregateID); err != nil {
		return nil, err
	}
	if err := c.pushAppendAndReduce(ctx, writeModel,
		schemauser.NewDeactivatedEvent(ctx, UserV3AggregateFromWriteModel(&writeModel.WriteModel)),
	); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) ReactiveSchemaUser(ctx context.Context, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-17XupGvxBJ", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserExists(ctx, "", id)
	if err != nil {
		return nil, err
	}
	if !writeModel.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-rQjbBr4J3j", "Errors.User.NotFound")
	}
	if writeModel.State != domain.UserStateInactive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "TODO", "TODO")
	}
	if err := c.checkPermissionUpdateUser(ctx, writeModel.ResourceOwner, writeModel.AggregateID); err != nil {
		return nil, err
	}
	if err := c.pushAppendAndReduce(ctx, writeModel,
		schemauser.NewReactivatedEvent(ctx, UserV3AggregateFromWriteModel(&writeModel.WriteModel)),
	); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) updateSchemaUserEmail(ctx context.Context, events []eventstore.Command, agg *eventstore.Aggregate, email *Email, alg crypto.EncryptionAlgorithm) (_ []eventstore.Command, plainCode string, err error) {

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

func (c *Commands) updateSchemaUserPhone(ctx context.Context, events []eventstore.Command, agg *eventstore.Aggregate, phone *Phone, alg crypto.EncryptionAlgorithm) (_ []eventstore.Command, plainCode string, err error) {
	events = append(events, schemauser.NewPhoneChangedEvent(ctx,
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
