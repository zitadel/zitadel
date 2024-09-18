package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user/schemauser"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type ChangeSchemaUserEmail struct {
	Details *domain.ObjectDetails

	ResourceOwner string
	ID            string

	Email      *Email
	ReturnCode string
}

func (s *ChangeSchemaUserEmail) Valid() (err error) {
	if s.Email != nil && s.Email.Address != "" {
		if err := s.Email.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func existingSchemaUserEmailWithPermission(ctx context.Context, c *Commands, resourceOwner, userID string) (*UserV3WriteModel, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-0oj2PquNGA", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserEmailWriteModelByID(ctx, resourceOwner, userID)
	if err != nil {
		return nil, err
	}
	if !writeModel.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-nJ0TQFuRmP", "Errors.User.NotFound")
	}
	if err := c.checkPermissionUpdateUser(ctx, writeModel.ResourceOwner, writeModel.AggregateID); err != nil {
		return nil, err
	}
	return writeModel, nil
}

func (c *Commands) ChangeSchemaUserEmail(ctx context.Context, user *ChangeSchemaUserEmail, alg crypto.EncryptionAlgorithm) (err error) {
	if err := user.Valid(); err != nil {
		return err
	}

	writeModel, err := existingSchemaUserEmailWithPermission(ctx, c, user.ResourceOwner, user.ID)
	if err != nil {
		return err
	}
	// when there is no change in the address, we don't want to change anything on the verification state
	if user.Email == nil || string(user.Email.Address) == writeModel.Email {
		user.Details = writeModelToObjectDetails(&writeModel.WriteModel)
		return nil
	}

	events := make([]eventstore.Command, 0)
	events, user.ReturnCode, err = c.updateSchemaUserEmail(ctx, writeModel, events, user.Email, alg)
	if err != nil {
		return err
	}
	if err := c.pushAppendAndReduce(ctx, writeModel, events...); err != nil {
		return err
	}
	user.Details = writeModelToObjectDetails(&writeModel.WriteModel)
	return nil
}

func (c *Commands) VerifySchemaUserEmail(ctx context.Context, resourceOwner, id, code string, alg crypto.EncryptionAlgorithm) (*domain.ObjectDetails, error) {
	writeModel, err := existingSchemaUserEmailWithPermission(ctx, c, resourceOwner, id)
	if err != nil {
		return nil, err
	}

	if writeModel.EmailCode == nil {
		return writeModelToObjectDetails(&writeModel.WriteModel), nil
	}
	if err := crypto.VerifyCode(writeModel.EmailCode.CreationDate, writeModel.EmailCode.Expiry, writeModel.EmailCode.Code, code, alg); err != nil {
		return nil, err
	}
	if err := c.pushAppendAndReduce(ctx, writeModel,
		schemauser.NewEmailVerifiedEvent(ctx, UserV3AggregateFromWriteModel(&writeModel.WriteModel)),
	); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

type ResendSchemaUserEmailCode struct {
	Details *domain.ObjectDetails

	ResourceOwner string
	ID            string

	URLTemplate string
	ReturnCode  bool
	PlainCode   string
}

func (r *ResendSchemaUserEmailCode) IsReturnCode() bool {
	return r.ReturnCode
}
func (r *ResendSchemaUserEmailCode) GetURLTemplate() string {
	return r.URLTemplate
}

func (c *Commands) ResendSchemaUserEmailCode(ctx context.Context, user *ResendSchemaUserEmailCode, alg crypto.EncryptionAlgorithm) error {
	writeModel, err := existingSchemaUserEmailWithPermission(ctx, c, user.ResourceOwner, user.ID)
	if err != nil {
		return err
	}
	if writeModel.EmailCode == nil {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-QRkNTBwF8q", "Errors.User.Code.Empty")
	}
	events := make([]eventstore.Command, 0)
	events, user.PlainCode, err = c.generateSchemaUserEmailCode(ctx, writeModel, events, user, alg)
	if err != nil {
		return err
	}
	if err := c.pushAppendAndReduce(ctx, writeModel, events...); err != nil {
		return err
	}
	user.Details = writeModelToObjectDetails(&writeModel.WriteModel)
	return nil
}

type EmailUpdate interface {
	EmailCodeGenerate
	GetAddress() domain.EmailAddress
	IsVerified() bool
}

func (c *Commands) updateSchemaUserEmail(ctx context.Context, existing *UserV3WriteModel, events []eventstore.Command, email EmailUpdate, alg crypto.EncryptionAlgorithm) (_ []eventstore.Command, plainCode string, err error) {
	if existing.Email == string(email.GetAddress()) {
		return events, plainCode, nil
	}
	events = append(events, schemauser.NewEmailUpdatedEvent(ctx,
		UserV3AggregateFromWriteModel(&existing.WriteModel),
		email.GetAddress(),
	))
	if email.IsVerified() {
		return append(events, schemauser.NewEmailVerifiedEvent(ctx, UserV3AggregateFromWriteModel(&existing.WriteModel))), "", nil
	}
	return c.generateSchemaUserEmailCode(ctx, existing, events, email, alg)
}

type EmailVerify interface {
	GetCode() string
}

func (c *Commands) verifySchemaUserEmail(ctx context.Context, existing *UserV3WriteModel, events []eventstore.Command, email EmailVerify, alg crypto.EncryptionAlgorithm) (_ []eventstore.Command, plainCode string, err error) {
	if err := crypto.VerifyCode(existing.EmailCode.CreationDate, existing.EmailCode.Expiry, existing.EmailCode.Code, email.GetCode(), alg); err != nil {
		return events, plainCode, err
	}
	return append(events, schemauser.NewEmailVerifiedEvent(ctx, UserV3AggregateFromWriteModel(&existing.WriteModel))), "", nil
}

type EmailCodeGenerate interface {
	IsReturnCode() bool
	GetURLTemplate() string
}

func (c *Commands) generateSchemaUserEmailCode(ctx context.Context, existing *UserV3WriteModel, events []eventstore.Command, email EmailCodeGenerate, alg crypto.EncryptionAlgorithm) (_ []eventstore.Command, plainCode string, err error) {
	cryptoCode, err := c.newEmailCode(ctx, c.eventstore.Filter, alg) //nolint:staticcheck
	if err != nil {
		return nil, "", err
	}
	if email.IsReturnCode() {
		plainCode = cryptoCode.Plain
	}
	events = append(events, schemauser.NewEmailCodeAddedEvent(ctx,
		UserV3AggregateFromWriteModel(&existing.WriteModel),
		cryptoCode.Crypted,
		cryptoCode.Expiry,
		email.GetURLTemplate(),
		email.IsReturnCode(),
	))
	return events, plainCode, nil
}
