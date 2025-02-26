package command

import (
	"context"
	"io"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// RequestPasswordReset generates a code
// and triggers a notification e-mail with the default confirmation URL format.
func (c *Commands) RequestPasswordReset(ctx context.Context, userID string) (*domain.ObjectDetails, *string, error) {
	return c.requestPasswordReset(ctx, userID, false, "", domain.NotificationTypeEmail)
}

// RequestPasswordResetURLTemplate generates a code
// and triggers a notification e-mail with the confirmation URL rendered from the passed urlTmpl.
// urlTmpl must be a valid [tmpl.Template].
func (c *Commands) RequestPasswordResetURLTemplate(ctx context.Context, userID, urlTmpl string, notificationType domain.NotificationType) (*domain.ObjectDetails, *string, error) {
	if err := domain.RenderConfirmURLTemplate(io.Discard, urlTmpl, userID, "code", "orgID"); err != nil {
		return nil, nil, err
	}
	return c.requestPasswordReset(ctx, userID, false, urlTmpl, notificationType)
}

// RequestPasswordResetReturnCode generates a code and does not send a notification email.
// The generated plain text code will be returned.
func (c *Commands) RequestPasswordResetReturnCode(ctx context.Context, userID string) (*domain.ObjectDetails, *string, error) {
	return c.requestPasswordReset(ctx, userID, true, "", 0)
}

// requestPasswordReset creates a code for a password change.
// returnCode controls if the plain text version of the code will be set in the return object.
// When the plain text code is returned, no notification e-mail will be sent to the user.
// urlTmpl allows changing the target URL that is used by the e-mail and should be a validated Go template, if used.
func (c *Commands) requestPasswordReset(ctx context.Context, userID string, returnCode bool, urlTmpl string, notificationType domain.NotificationType) (_ *domain.ObjectDetails, plainCode *string, err error) {
	if userID == "" {
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-SAFdda", "Errors.User.IDMissing")
	}
	model, err := c.getHumanWriteModelByID(ctx, userID, "")
	if err != nil {
		return nil, nil, err
	}
	if !model.UserState.Exists() {
		return nil, nil, zerrors.ThrowNotFound(nil, "COMMAND-SAF4f", "Errors.User.NotFound")
	}
	if model.UserState == domain.UserStateInitial {
		return nil, nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Sfe4g", "Errors.User.NotInitialised")
	}
	if err = c.checkPermissionUpdateUser(ctx, model.ResourceOwner, userID); err != nil {
		return nil, nil, err
	}
	var passwordCode *EncryptedCode
	var generatorID string
	if notificationType == domain.NotificationTypeSms {
		passwordCode, generatorID, err = c.newPhoneCode(ctx, c.eventstore.Filter, domain.SecretGeneratorTypePasswordResetCode, c.userEncryption, c.defaultSecretGenerators.PasswordVerificationCode) //nolint:staticcheck
	} else {
		passwordCode, err = c.newEncryptedCode(ctx, c.eventstore.Filter, domain.SecretGeneratorTypePasswordResetCode, c.userEncryption) //nolint:staticcheck
	}
	if err != nil {
		return nil, nil, err
	}
	cmd := user.NewHumanPasswordCodeAddedEventV2(ctx, UserAggregateFromWriteModelCtx(ctx, &model.WriteModel), passwordCode.CryptedCode(), passwordCode.CodeExpiry(), notificationType, urlTmpl, returnCode, generatorID)

	if returnCode {
		plainCode = &passwordCode.Plain
	}
	if err = c.pushAppendAndReduce(ctx, model, cmd); err != nil {
		return nil, nil, err
	}
	return writeModelToObjectDetails(&model.WriteModel), plainCode, nil
}
