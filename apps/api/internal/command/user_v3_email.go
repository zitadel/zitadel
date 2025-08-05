package command

import (
	"context"
	"io"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type ChangeSchemaUserEmail struct {
	ResourceOwner string
	ID            string

	Email      *Email
	ReturnCode *string
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
	if s.Email != nil && s.Email.URLTemplate != "" {
		if err := domain.RenderConfirmURLTemplate(io.Discard, s.Email.URLTemplate, s.ID, "code", "orgID"); err != nil {
			return err
		}
	}
	return nil
}

func (c *Commands) ChangeSchemaUserEmail(ctx context.Context, user *ChangeSchemaUserEmail) (_ *domain.ObjectDetails, err error) {
	if err := user.Valid(); err != nil {
		return nil, err
	}

	writeModel, err := c.getSchemaUserEmailWriteModelByID(ctx, user.ResourceOwner, user.ID)
	if err != nil {
		return nil, err
	}

	events, plainCode, err := writeModel.NewEmailUpdate(ctx,
		user.Email,
		func(ctx context.Context) (*EncryptedCode, error) {
			return c.newEmailCode(ctx, c.eventstore.Filter, c.userEncryption) //nolint:staticcheck
		},
	)
	if err != nil {
		return nil, err
	}
	if plainCode != "" {
		user.ReturnCode = &plainCode
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

func (c *Commands) VerifySchemaUserEmail(ctx context.Context, resourceOwner, id, code string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-y3n4Sdu8j5", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserEmailWriteModelByID(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}

	events, err := writeModel.NewEmailVerify(ctx,
		func(creationDate time.Time, expiry time.Duration, cryptoCode *crypto.CryptoValue) error {
			return crypto.VerifyCode(creationDate, expiry, cryptoCode, code, c.userEncryption)
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
	PlainCode   *string
}

func (c *Commands) ResendSchemaUserEmailCode(ctx context.Context, user *ResendSchemaUserEmailCode) (*domain.ObjectDetails, error) {
	if user.ID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-KvPc5o9GeJ", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserEmailWriteModelByID(ctx, user.ResourceOwner, user.ID)
	if err != nil {
		return nil, err
	}

	events, plainCode, err := writeModel.NewResendEmailCode(ctx,
		func(ctx context.Context) (*EncryptedCode, error) {
			return c.newEmailCode(ctx, c.eventstore.Filter, c.userEncryption) //nolint:staticcheck
		},
		user.URLTemplate,
		user.ReturnCode,
	)
	if err != nil {
		return nil, err
	}
	if plainCode != "" {
		user.PlainCode = &plainCode
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}
