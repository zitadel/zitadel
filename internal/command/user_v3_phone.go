package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type ChangeSchemaUserPhone struct {
	ResourceOwner string
	ID            string

	Phone      *Phone
	ReturnCode *string
}

func (s *ChangeSchemaUserPhone) Valid() (err error) {
	if s.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-DkQ9aurv5u", "Errors.IDMissing")
	}
	if s.Phone != nil && s.Phone.Number != "" {
		if s.Phone.Number, err = s.Phone.Number.Normalize(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Commands) ChangeSchemaUserPhone(ctx context.Context, user *ChangeSchemaUserPhone) (_ *domain.ObjectDetails, err error) {
	if err := user.Valid(); err != nil {
		return nil, err
	}

	writeModel, err := c.getSchemaUserWMForPhone(ctx, user.ResourceOwner, user.ID)
	if err != nil {
		return nil, err
	}

	events, plainCode, err := writeModel.NewPhoneUpdate(ctx,
		user.Phone,
		func(ctx context.Context) (*EncryptedCode, string, error) {
			return c.newPhoneCode(ctx, c.eventstore.Filter, domain.SecretGeneratorTypeVerifyPhoneCode, c.userEncryption, c.defaultSecretGenerators.PhoneVerificationCode) //nolint:staticcheck
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

func (c *Commands) VerifySchemaUserPhone(ctx context.Context, resourceOwner, id, code string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-R4LKY44Ke3", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserWMForPhone(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}

	events, err := writeModel.NewPhoneVerify(ctx,
		func(ctx context.Context, creationDate time.Time, expiry time.Duration, cryptoCode *crypto.CryptoValue, generatorID string, verificationID string) error {
			return schemaUserVerifyCode(ctx, creationDate, expiry, cryptoCode, generatorID, verificationID, code, c.userEncryption, c.phoneCodeVerifier)
		},
	)
	if err != nil {
		return nil, err
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

type ResendSchemaUserPhoneCode struct {
	ResourceOwner string
	ID            string

	ReturnCode bool
	PlainCode  *string
}

func (c *Commands) ResendSchemaUserPhoneCode(ctx context.Context, user *ResendSchemaUserPhoneCode) (*domain.ObjectDetails, error) {
	if user.ID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-zmxIFR2nMo", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserWMForPhone(ctx, user.ResourceOwner, user.ID)
	if err != nil {
		return nil, err
	}

	events, plainCode, err := writeModel.NewResendPhoneCode(ctx,
		func(ctx context.Context) (*EncryptedCode, string, error) {
			return c.newPhoneCode(ctx, c.eventstore.Filter, domain.SecretGeneratorTypeVerifyPhoneCode, c.userEncryption, c.defaultSecretGenerators.PhoneVerificationCode) //nolint:staticcheck
		},
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

func (c *Commands) getSchemaUserWMForPhone(ctx context.Context, resourceOwner, id string) (*UserV3WriteModel, error) {
	writeModel := NewUserV3PhoneWriteModel(resourceOwner, id, c.checkPermission)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	return writeModel, nil
}
