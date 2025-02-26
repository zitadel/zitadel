package command

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	domain_schema "github.com/zitadel/zitadel/internal/domain/schema"
	"github.com/zitadel/zitadel/internal/repository/user/schemauser"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type CreateSchemaUser struct {
	SchemaID       string
	schemaRevision uint64

	ResourceOwner string
	ID            string
	Data          json.RawMessage

	Email           *Email
	ReturnCodeEmail *string
	Phone           *Phone
	ReturnCodePhone *string
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

func (c *Commands) CreateSchemaUser(ctx context.Context, user *CreateSchemaUser) (*domain.ObjectDetails, error) {
	if err := user.Valid(ctx, c); err != nil {
		return nil, err
	}

	writeModel, err := c.getSchemaUserExists(ctx, user.ResourceOwner, user.ID)
	if err != nil {
		return nil, err
	}

	events, codeEmail, codePhone, err := writeModel.NewCreated(ctx,
		user.SchemaID,
		user.schemaRevision,
		user.Data,
		user.Email,
		user.Phone,
		func(ctx context.Context) (*EncryptedCode, error) {
			return c.newEmailCode(ctx, c.eventstore.Filter, c.userEncryption) //nolint:staticcheck
		},
		func(ctx context.Context) (*EncryptedCode, string, error) {
			return c.newPhoneCode(ctx, c.eventstore.Filter, domain.SecretGeneratorTypeVerifyPhoneCode, c.userEncryption, c.defaultSecretGenerators.PhoneVerificationCode) //nolint:staticcheck
		},
	)
	if err != nil {
		return nil, err
	}
	if codeEmail != "" {
		user.ReturnCodeEmail = &codeEmail
	}
	if codePhone != "" {
		user.ReturnCodePhone = &codePhone
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

func (c *Commands) DeleteSchemaUser(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Vs4wJCME7T", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserExists(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}

	events, err := writeModel.NewDelete(ctx)
	if err != nil {
		return nil, err
	}

	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

type ChangeSchemaUser struct {
	schemaWriteModel *UserSchemaWriteModel

	ResourceOwner string
	ID            string

	SchemaUser *SchemaUser

	Email           *Email
	ReturnCodeEmail *string
	Phone           *Phone
	ReturnCodePhone *string
}

type SchemaUser struct {
	SchemaID string
	Data     json.RawMessage
}

func (s *ChangeSchemaUser) Valid() (err error) {
	if s.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-gEJR1QOGHb", "Errors.IDMissing")
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

func (c *Commands) ChangeSchemaUser(ctx context.Context, user *ChangeSchemaUser) (*domain.ObjectDetails, error) {
	if err := user.Valid(); err != nil {
		return nil, err
	}

	writeModel, err := c.getSchemaUserWriteModelByID(ctx, user.ResourceOwner, user.ID)
	if err != nil {
		return nil, err
	}
	if !writeModel.Exists() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Nn8CRVlkeZ", "Errors.User.NotFound")
	}

	schemaID := writeModel.SchemaID
	if user.SchemaUser != nil && user.SchemaUser.SchemaID != "" {
		schemaID = user.SchemaUser.SchemaID
	}

	var schemaWM *UserSchemaWriteModel
	if user.SchemaUser != nil {
		schemaWriteModel, err := c.getSchemaWriteModelByID(ctx, "", schemaID)
		if err != nil {
			return nil, err
		}
		if !schemaWriteModel.Exists() {
			return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-VLDTtxT3If", "Errors.UserSchema.NotExists")
		}
		schemaWM = schemaWriteModel
	}

	events, codeEmail, codePhone, err := writeModel.NewUpdate(ctx,
		schemaWM,
		user.SchemaUser,
		user.Email,
		user.Phone,
		func(ctx context.Context) (*EncryptedCode, error) {
			return c.newEmailCode(ctx, c.eventstore.Filter, c.userEncryption) //nolint:staticcheck
		},
		func(ctx context.Context) (*EncryptedCode, string, error) {
			return c.newPhoneCode(ctx, c.eventstore.Filter, domain.SecretGeneratorTypeVerifyPhoneCode, c.userEncryption, c.defaultSecretGenerators.PhoneVerificationCode) //nolint:staticcheck
		},
	)
	if err != nil {
		return nil, err
	}

	if codeEmail != "" {
		user.ReturnCodeEmail = &codeEmail
	}
	if codePhone != "" {
		user.ReturnCodePhone = &codePhone
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
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
	writeModel, err := c.getSchemaUserExists(ctx, resourceOwner, id)
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

func (c *Commands) getSchemaUserExists(ctx context.Context, resourceOwner, id string) (*UserV3WriteModel, error) {
	writeModel := NewExistsUserV3WriteModel(resourceOwner, id, c.checkPermission)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	return writeModel, nil
}

func (c *Commands) getSchemaUserWriteModelByID(ctx context.Context, resourceOwner, id string) (*UserV3WriteModel, error) {
	writeModel := NewUserV3WriteModel(resourceOwner, id, c.checkPermission)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	return writeModel, nil
}

func (c *Commands) getSchemaUserEmailWriteModelByID(ctx context.Context, resourceOwner, id string) (*UserV3WriteModel, error) {
	writeModel := NewUserV3EmailWriteModel(resourceOwner, id, c.checkPermission)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	return writeModel, nil
}

func (c *Commands) getSchemaUserPhoneWriteModelByID(ctx context.Context, resourceOwner, id string) (*UserV3WriteModel, error) {
	writeModel := NewUserV3PhoneWriteModel(resourceOwner, id, c.checkPermission)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	return writeModel, nil
}
