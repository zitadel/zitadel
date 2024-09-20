package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user/schemauser"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type ChangeSchemaUserEmail struct {
	ResourceOwner string
	ID            string

	Email      *Email
	ReturnCode string
}

func (s *ChangeSchemaUserEmail) Valid() (err error) {
	if s.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-0oj2PquNGA", "Errors.IDMissing")
	}
	if s.Email != nil && s.Email.Address != "" {
		if err := s.Email.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func existingSchemaUserEmail(ctx context.Context, c *Commands, resourceOwner, userID string) (*UserV3WriteModel, error) {
	writeModel, err := c.getSchemaUserEmailWriteModelByID(ctx, resourceOwner, userID)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func (c *Commands) ChangeSchemaUserEmail(ctx context.Context, user *ChangeSchemaUserEmail, alg crypto.EncryptionAlgorithm) (_ *domain.ObjectDetails, err error) {
	if err := user.Valid(); err != nil {
		return nil, err
	}

	writeModel, err := existingSchemaUserEmail(ctx, c, user.ResourceOwner, user.ID)
	if err != nil {
		return nil, err
	}
	events, plainCode, err := c.updateSchemaUserEmail(ctx, writeModel, nil, user.Email, alg)
	if err != nil {
		return nil, err
	}
	if plainCode != "" {
		user.ReturnCode = plainCode
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

func (c *Commands) VerifySchemaUserEmail(ctx context.Context, resourceOwner, id, code string, alg crypto.EncryptionAlgorithm) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-y3n4Sdu8j5", "Errors.IDMissing")
	}
	writeModel, err := existingSchemaUserEmail(ctx, c, resourceOwner, id)
	if err != nil {
		return nil, err
	}

	events, err := writeModel.NewEmailVerify(ctx,
		func(creationDate time.Time, expiry time.Duration, cryptoCode *crypto.CryptoValue) error {
			return crypto.VerifyCode(creationDate, expiry, cryptoCode, code, alg)
		},
	)
	if err != nil {
		return nil, err
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

type ResendSchemaUserEmailCode struct {
	ResourceOwner string
	ID            string

	URLTemplate string
	ReturnCode  bool
	PlainCode   string
}

func (c *Commands) ResendSchemaUserEmailCode(ctx context.Context, user *ResendSchemaUserEmailCode, alg crypto.EncryptionAlgorithm) (*domain.ObjectDetails, error) {
	if user.ID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-KvPc5o9GeJ", "Errors.IDMissing")
	}
	writeModel, err := existingSchemaUserEmail(ctx, c, user.ResourceOwner, user.ID)
	if err != nil {
		return nil, err
	}

	events, plainCode, err := writeModel.NewResendEmailCode(ctx,
		func(ctx context.Context) (*EncryptedCode, error) {
			return c.newEmailCode(ctx, c.eventstore.Filter, alg) //nolint:staticcheck
		},
		user.URLTemplate,
		user.ReturnCode,
	)
	if err != nil {
		return nil, err
	}
	if plainCode != "" {
		user.PlainCode = plainCode
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

type EmailUpdate interface {
	EmailCodeGenerate
	GetAddress() domain.EmailAddress
	IsVerified() bool
}

func (c *Commands) updateSchemaUserEmail(ctx context.Context, writeModel *UserV3WriteModel, events []eventstore.Command, email *Email, alg crypto.EncryptionAlgorithm) (_ []eventstore.Command, plainCode string, err error) {
	if events == nil {
		events = []eventstore.Command{}
	}

	updateEvents, plainCode, err := writeModel.NewEmailUpdatedEvents(ctx,
		email,
		func(ctx context.Context) (*EncryptedCode, error) {
			return c.newEmailCode(ctx, c.eventstore.Filter, alg) //nolint:staticcheck
		},
	)
	if err != nil {
		return nil, "", err
	}
	return append(events, updateEvents...), plainCode, nil
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
