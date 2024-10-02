package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SetSchemaUserPassword struct {
	ResourceOwner string
	UserID        string

	Verification *SchemaUserPasswordVerification
	Password     *SchemaUserPassword
}

type SchemaUserPasswordVerification struct {
	CurrentPassword string
	Code            string
}

type SchemaUserPassword struct {
	Password            string
	EncodedPasswordHash string
	ChangeRequired      bool
}

func (p *SchemaUserPassword) Validate(hasher *crypto.Hasher) (err error) {
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

func (c *Commands) SetSchemaUserPassword(ctx context.Context, set *SetSchemaUserPassword) (*domain.ObjectDetails, error) {
	if set.UserID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-aS3Vz5t6BS", "Errors.IDMissing")
	}
	if err := set.Password.Validate(c.userPasswordHasher); err != nil {
		return nil, err
	}
	schemauser, err := existingSchemaUser(ctx, c, set.ResourceOwner, set.UserID)
	if err != nil {
		return nil, err
	}
	set.ResourceOwner = schemauser.ResourceOwner

	_, err = existingSchema(ctx, c, "", schemauser.SchemaID)
	if err != nil {
		return nil, err
	}
	// TODO check for possible authenticators

	writeModel, events, err := c.setSchemaUserPassword(ctx, set.ResourceOwner, set.UserID, set.Verification, set.Password)
	if err != nil {
		return nil, err
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

func (c *Commands) setSchemaUserPassword(ctx context.Context, resourceOwner, userID string, verification *SchemaUserPasswordVerification, set *SchemaUserPassword) (*PasswordV3WriteModel, []eventstore.Command, error) {
	if set == nil {
		return nil, nil, nil
	}
	schemaUser := &schemaUserPassword{
		Set:           true,
		ResourceOwner: resourceOwner,
		UserID:        userID,
		Verification:  verification,
		NewPassword:   set,
	}
	writeModel, err := c.getSchemaUserPasswordWithVerification(ctx, schemaUser)
	if err != nil {
		return nil, nil, err
	}

	// If password is provided, let's check if is compliant with the policy.
	// If only a encodedPassword is passed, we can skip this.
	if set.Password != "" {
		if err = c.checkPasswordComplexity(ctx, set.Password, writeModel.ResourceOwner); err != nil {
			return nil, nil, err
		}
	}
	encodedPassword := schemaUser.NewPassword.EncodedPasswordHash
	if encodedPassword == "" && set.Password != "" {
		encodedPassword, err = c.userPasswordHasher.Hash(set.Password)
		if err = convertPasswapErr(err); err != nil {
			return nil, nil, err
		}
	}
	events, err := writeModel.NewCreate(ctx,
		encodedPassword,
		set.ChangeRequired,
	)
	if err != nil {
		return nil, nil, err
	}
	return writeModel, events, nil
}

type RequestSchemaUserPasswordReset struct {
	ResourceOwner string
	UserID        string

	URLTemplate      string
	NotificationType domain.NotificationType
	PlainCode        *string
	ReturnCode       bool
}

func (c *Commands) RequestSchemaUserPasswordReset(ctx context.Context, user *RequestSchemaUserPasswordReset) (_ *domain.ObjectDetails, err error) {
	writeModel, err := existsSchemaUserPasswordWithPermission(ctx, c, user.ResourceOwner, user.UserID)
	if err != nil {
		return nil, err
	}

	events, plainCode, err := writeModel.NewAddCode(ctx,
		user.NotificationType,
		user.URLTemplate,
		user.ReturnCode,
		func(ctx context.Context, notifyType domain.NotificationType) (*EncryptedCode, string, error) {
			var passwordCode *EncryptedCode
			var generatorID string
			if notifyType == domain.NotificationTypeSms {
				passwordCode, generatorID, err = c.newPhoneCode(ctx, c.eventstore.Filter, domain.SecretGeneratorTypePasswordResetCode, c.userEncryption, c.defaultSecretGenerators.PasswordVerificationCode) //nolint:staticcheck
			} else {
				passwordCode, err = c.newEncryptedCode(ctx, c.eventstore.Filter, domain.SecretGeneratorTypePasswordResetCode, c.userEncryption) //nolint:staticcheck
			}
			return passwordCode, generatorID, err
		},
	)
	if err != nil {
		return nil, err
	}
	if plainCode != "" {
		user.PlainCode = &plainCode
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

func (c *Commands) DeleteSchemaUserPassword(ctx context.Context, resourceOwner, id string) (_ *domain.ObjectDetails, err error) {
	writeModel, err := existsSchemaUserPasswordWithPermission(ctx, c, resourceOwner, id)
	if err != nil {
		return nil, err
	}

	events, err := writeModel.NewDelete(ctx)
	if err != nil {
		return nil, err
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

type schemaUserPassword struct {
	Set           bool
	ResourceOwner string
	UserID        string
	Verification  *SchemaUserPasswordVerification
	NewPassword   *SchemaUserPassword
}

func (c *Commands) getSchemaUserPasswordWM(ctx context.Context, resourceOwner, id string) (*PasswordV3WriteModel, error) {
	writeModel := NewPasswordV3WriteModel(resourceOwner, id, c.checkPermission)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	return writeModel, nil
}

func existsSchemaUserPasswordWithPermission(ctx context.Context, c *Commands, resourceOwner, id string) (*PasswordV3WriteModel, error) {
	writeModel, err := c.getSchemaUserPasswordWithVerification(ctx, &schemaUserPassword{ResourceOwner: resourceOwner, UserID: id})
	if err != nil {
		return nil, err
	}
	return writeModel, writeModel.Exists()
}

func (c *Commands) getSchemaUserPasswordWithVerification(ctx context.Context, user *schemaUserPassword) (*PasswordV3WriteModel, error) {
	if user.UserID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-PoSU5BOZCi", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserPasswordWM(ctx, user.ResourceOwner, user.UserID)
	if err != nil {
		return nil, err
	}
	if err := writeModel.Exists(); !user.Set && err != nil {
		return nil, err
	}

	// if no verification is set, the user must have the permission to change the password
	verification := c.setSchemaUserPasswordWithPermission(writeModel.UserID, writeModel.ResourceOwner)
	if user.Verification != nil {
		// otherwise check the password code...
		if user.Verification.Code != "" {
			verification = c.setSchemaUserPasswordWithVerifyCode(writeModel.CodeCreationDate, writeModel.CodeExpiry, writeModel.Code, writeModel.GeneratorID, writeModel.VerificationID, user.Verification.Code)
		}
		// ...or old password
		if user.Verification.CurrentPassword != "" {
			verification = c.checkSchemaUserCurrentPassword(user.NewPassword.Password, user.NewPassword.EncodedPasswordHash, user.Verification.CurrentPassword, writeModel.EncodedHash)
		}
	}

	if verification != nil {
		newEncodedPassword, err := verification(ctx)
		if err != nil {
			return nil, err
		}
		// use the new hash from the verification in case there is one (e.g. existing pw check)
		if newEncodedPassword != "" {
			user.NewPassword.EncodedPasswordHash = newEncodedPassword
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
	codeCreationDate time.Time,
	codeExpiry time.Duration,
	encryptedCode *crypto.CryptoValue,
	codeProviderID string,
	codeVerificationID string,
	code string,
) setPasswordVerification {
	return func(ctx context.Context) (_ string, err error) {
		return "", schemaUserVerifyCode(ctx, codeCreationDate, codeExpiry, encryptedCode, codeProviderID, codeVerificationID, code, c.userEncryption, c.phoneCodeVerifier)
	}
}

// checkSchemaUserCurrentPassword returns a password check as [setPasswordVerification] implementation
func (c *Commands) checkSchemaUserCurrentPassword(
	newPassword, newEncodedPassword, currentPassword, currentEncodePassword string,
) setPasswordVerification {
	// in case the new password is already encoded, we only need to verify the current
	if newEncodedPassword != "" {
		return func(ctx context.Context) (string, error) {
			_, spanPasswap := tracing.NewNamedSpan(ctx, "passwap.Verify")
			_, err := c.userPasswordHasher.Verify(currentEncodePassword, currentPassword)
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
