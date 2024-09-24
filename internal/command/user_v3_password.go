package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/user/authenticator"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SetSchemaUserPassword struct {
	ResourceOwner string
	UserID        string

	Password            string
	EncodedPasswordHash string
	ChangeRequired      bool

	CurrentPassword  string
	VerificationCode string
}

func (p *SetSchemaUserPassword) Validate(hasher *crypto.Hasher) (err error) {
	if p.UserID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-aS3Vz5t6BS", "Errors.IDMissing")
	}

	if p.EncodedPasswordHash != "" {
		if !hasher.EncodingSupported(p.EncodedPasswordHash) {
			return zerrors.ThrowInvalidArgument(nil, "COMMAND-oz74onzvqr", "Errors.User.Password.NotSupported")
		}
	}
	if p.Password == "" && p.EncodedPasswordHash == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-3klek4sbns", "Errors.User.Password.Empty")
	}

	return nil
}

func (c *Commands) SetSchemaUserPassword(ctx context.Context, user *SetSchemaUserPassword) (*domain.ObjectDetails, error) {
	if err := user.Validate(c.userPasswordHasher); err != nil {
		return nil, err
	}

	schemaUser := &schemaUserPassword{
		ResourceOwner:       user.ResourceOwner,
		UserID:              user.UserID,
		VerificationCode:    user.VerificationCode,
		CurrentPassword:     user.CurrentPassword,
		Password:            user.Password,
		EncodedPasswordHash: user.EncodedPasswordHash,
	}

	existing, err := c.getSchemaUserPasswordWithVerification(ctx, schemaUser)
	if err != nil {
		return nil, err
	}
	resourceOwner := existing.ResourceOwner
	// when no password was set yet
	if existing.EncodedHash == "" {
		existingUser, err := c.getSchemaUserExists(ctx, user.ResourceOwner, user.UserID)
		if err != nil {
			return nil, err
		}
		if !existingUser.Exists() {
			return nil, zerrors.ThrowNotFound(nil, "COMMAND-TODO", "Errors.User.Password.NotFound")
		}
		resourceOwner = existingUser.ResourceOwner
	}

	// If password is provided, let's check if is compliant with the policy.
	// If only a encodedPassword is passed, we can skip this.
	if user.Password != "" {
		if err = c.checkPasswordComplexity(ctx, user.Password, resourceOwner); err != nil {
			return nil, err
		}
	}

	encodedPassword := schemaUser.EncodedPasswordHash
	if encodedPassword == "" && user.Password != "" {
		encodedPassword, err = c.userPasswordHasher.Hash(user.Password)
		if err = convertPasswapErr(err); err != nil {
			return nil, err
		}
	}

	events, err := c.eventstore.Push(ctx,
		authenticator.NewPasswordCreatedEvent(ctx,
			&authenticator.NewAggregate(user.UserID, resourceOwner).Aggregate,
			existing.UserID,
			encodedPassword,
			user.ChangeRequired,
		),
	)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}

type RequestSchemaUserPasswordReset struct {
	ResourceOwner string
	UserID        string

	URLTemplate      string
	NotificationType domain.NotificationType
	PlainCode        string
	ReturnCode       bool
}

func (c *Commands) RequestSchemaUserPasswordReset(ctx context.Context, user *RequestSchemaUserPasswordReset) (_ *domain.ObjectDetails, err error) {
	existing, err := c.getSchemaUserPasswordExists(ctx, user.ResourceOwner, user.UserID)
	if err != nil {
		return nil, err
	}
	if existing.EncodedHash == "" {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-TODO", "Errors.User.Password.NotFound")
	}

	code, err := c.newEncryptedCode(ctx, c.eventstore.Filter, domain.SecretGeneratorTypePasswordResetCode, c.userEncryption) //nolint:staticcheck
	if err != nil {
		return nil, err
	}

	events, err := c.eventstore.Push(ctx,
		authenticator.NewPasswordCodeAddedEvent(ctx,
			&authenticator.NewAggregate(existing.UserID, existing.ResourceOwner).Aggregate,
			code.Crypted,
			code.Expiry,
			user.NotificationType,
			user.URLTemplate,
			user.ReturnCode,
		),
	)
	if err != nil {
		return nil, err
	}
	if user.ReturnCode {
		user.PlainCode = code.Plain
	}
	return pushedEventsToObjectDetails(events), nil
}

func (c *Commands) DeleteSchemaUserPassword(ctx context.Context, resourceOwner, id string) (_ *domain.ObjectDetails, err error) {
	existing, err := c.getSchemaUserPasswordExists(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}
	if existing.EncodedHash == "" {
		return nil, zerrors.ThrowNotFound(nil, "TODO", "TODO")
	}

	events, err := c.eventstore.Push(ctx,
		authenticator.NewPasswordDeletedEvent(ctx,
			&authenticator.NewAggregate(id, existing.ResourceOwner).Aggregate,
		),
	)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}

type schemaUserPassword struct {
	ResourceOwner       string
	UserID              string
	VerificationCode    string
	CurrentPassword     string
	Password            string
	EncodedPasswordHash string
}

func (c *Commands) getSchemaUserPasswordExists(ctx context.Context, resourceOwner, id string) (*PasswordV3WriteModel, error) {
	return c.getSchemaUserPasswordWithVerification(ctx, &schemaUserPassword{ResourceOwner: resourceOwner, UserID: id})
}

func (c *Commands) getSchemaUserPasswordWithVerification(ctx context.Context, user *schemaUserPassword) (*PasswordV3WriteModel, error) {
	if user.UserID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-PoSU5BOZCi", "Errors.IDMissing")
	}
	writeModel := NewPasswordV3WriteModel(user.ResourceOwner, user.UserID)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}

	// if no verification is set, the user must have the permission to change the password
	verification := c.setSchemaUserPasswordWithPermission(writeModel.UserID, writeModel.ResourceOwner)
	// otherwise check the password code...
	if user.VerificationCode != "" {
		verification = c.setSchemaUserPasswordWithVerifyCode(writeModel.CodeCreationDate, writeModel.CodeExpiry, writeModel.Code, user.VerificationCode)
	}
	// ...or old password
	if user.CurrentPassword != "" {
		verification = c.checkCurrentPassword(user.Password, user.EncodedPasswordHash, user.CurrentPassword, writeModel.EncodedHash)
	}

	if verification != nil {
		newEncodedPassword, err := verification(ctx)
		if err != nil {
			return nil, err
		}
		// use the new hash from the verification in case there is one (e.g. existing pw check)
		if newEncodedPassword != "" {
			user.EncodedPasswordHash = newEncodedPassword
		}
	}
	return writeModel, nil
}

// setSchemaUserPasswordWithPermission returns a permission check as [setPasswordVerification] implementation
func (c *Commands) setSchemaUserPasswordWithPermission(orgID, userID string) setPasswordVerification {
	return func(ctx context.Context) (_ string, err error) {
		return "", c.checkPermissionUpdateUser(ctx, orgID, userID)
	}
}

// setSchemaUserPasswordWithVerifyCode returns a password code check as [setPasswordVerification] implementation
func (c *Commands) setSchemaUserPasswordWithVerifyCode(
	passwordCodeCreationDate time.Time,
	passwordCodeExpiry time.Duration,
	passwordCode *crypto.CryptoValue,
	code string,
) setPasswordVerification {
	return func(ctx context.Context) (_ string, err error) {
		if passwordCode == nil {
			return "", zerrors.ThrowPreconditionFailed(nil, "COMMAND-TODO", "Errors.User.Code.NotFound")
		}
		_, spanCrypto := tracing.NewNamedSpan(ctx, "crypto.VerifyCode")
		defer func() {
			spanCrypto.EndWithError(err)
		}()
		return "", crypto.VerifyCode(passwordCodeCreationDate, passwordCodeExpiry, passwordCode, code, c.userEncryption)
	}
}

// checkSchemaUserCurrentPassword returns a password check as [setPasswordVerification] implementation
func (c *Commands) checkSchemaUserCurrentPassword(
	newPassword, newEncodedPassword, currentPassword, currentEncodePassword string,
) setPasswordVerification {
	// in case the new password is already encoded, we only need to verify the current
	if newEncodedPassword != "" {
		return func(ctx context.Context) (_ string, err error) {
			_, spanPasswap := tracing.NewNamedSpan(ctx, "passwap.Verify")
			_, err = c.userPasswordHasher.Verify(currentEncodePassword, currentPassword)
			spanPasswap.EndWithError(err)
			return "", convertPasswapErr(err)
		}
	}

	// otherwise let's directly verify and return the new generate hash, so we can reuse it in the event
	return func(ctx context.Context) (string, error) {
		return c.verifyAndUpdateSchemaUserPassword(ctx, currentEncodePassword, currentPassword, newPassword)
	}
}

// verifyAndUpdateSchemaUserPassword verify if the old password is correct with the encoded hash and
// returns the hash of the new password if so
func (c *Commands) verifyAndUpdateSchemaUserPassword(ctx context.Context, encodedHash, oldPassword, newPassword string) (string, error) {
	if encodedHash == "" {
		return "", zerrors.ThrowPreconditionFailed(nil, "COMMAND-TODO", "Errors.User.Password.NotSet")
	}

	_, spanPasswap := tracing.NewNamedSpan(ctx, "passwap.Verify")
	updated, err := c.userPasswordHasher.VerifyAndUpdate(encodedHash, oldPassword, newPassword)
	spanPasswap.EndWithError(err)
	return updated, convertPasswapErr(err)
}
