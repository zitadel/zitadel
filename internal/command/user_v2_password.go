package command

import (
	"context"
	"io"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/repository/user"
)

// RequestPasswordReset generates a code
// and triggers a notification e-mail with the default confirmation URL format.
func (c *Commands) RequestPasswordReset(ctx context.Context, userID, resourceOwner string) (*domain.ObjectDetails, *string, error) {
	return c.requestPasswordReset(ctx, userID, resourceOwner, false, "")
}

// RequestPasswordResetURLTemplate generates a code
// and triggers a notification e-mail with the confirmation URL rendered from the passed urlTmpl.
// urlTmpl must be a valid [tmpl.Template].
func (c *Commands) RequestPasswordResetURLTemplate(ctx context.Context, userID, resourceOwner string, urlTmpl string) (*domain.ObjectDetails, *string, error) {
	if err := domain.RenderConfirmURLTemplate(io.Discard, urlTmpl, userID, "code", "orgID"); err != nil {
		return nil, nil, err
	}
	return c.requestPasswordReset(ctx, userID, resourceOwner, false, urlTmpl)
}

// RequestPasswordResetReturnCode generates a code and does not send a notification email.
// The generated plain text code will be returned.
func (c *Commands) RequestPasswordResetReturnCode(ctx context.Context, userID, resourceOwner string) (*domain.ObjectDetails, *string, error) {
	return c.requestPasswordReset(ctx, userID, resourceOwner, true, "")
}

//
//// requestPasswordReset creates a code for a password change.
//// returnCode controls if the plain text version of the code will be set in the return object.
//// When the plain text code is returned, no notification e-mail will be sent to the user.
//// urlTmpl allows changing the target URL that is used by the e-mail and should be a validated Go template, if used.
//func (c *Commands) requestPasswordReset(ctx context.Context, userID, resourceOwner string, alg crypto.EncryptionAlgorithm, returnCode bool, urlTmpl string) (*domain.ObjectDetails, *string, error) {
//	config, err := secretGeneratorConfig(ctx, c.eventstore.Filter, domain.SecretGeneratorTypePasswordResetCode)
//	if err != nil {
//		return nil, nil, err
//	}
//	gen := crypto.NewEncryptionGenerator(*config, alg)
//	cmd, err := c.NewUserPasswordResetEvent(ctx, userID, resourceOwner, alg, returnCode, urlTmpl)
//	if err != nil {
//		return nil, nil, err
//	}
//	//if err = cmd.AddGeneratedCode(ctx, gen, urlTmpl, returnCode); err != nil {
//	//	return nil, nil, err
//	//}
//	if err = c.pushAppendAndReduce(ctx, cmd.model, cmd.event); err != nil {
//		return nil, nil, err
//	}
//	return writeModelToObjectDetails(&cmd.model.WriteModel), cmd.plainCode, nil
//}

// requestPasswordReset creates a code for a password change.
// returnCode controls if the plain text version of the code will be set in the return object.
// When the plain text code is returned, no notification e-mail will be sent to the user.
// urlTmpl allows changing the target URL that is used by the e-mail and should be a validated Go template, if used.
func (c *Commands) requestPasswordReset(ctx context.Context, userID, resourceOwner string, returnCode bool, urlTmpl string) (_ *domain.ObjectDetails, plainCode *string, err error) {
	if userID == "" {
		return nil, nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-SAFdda", "Errors.User.IDMissing")
	}
	model, err := c.getHumanWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, nil, err
	}
	if !model.UserState.Exists() {
		return nil, nil, caos_errs.ThrowNotFound(nil, "COMMAND-SAF4f", "Errors.User.NotFound")
	}
	if model.UserState == domain.UserStateInitial {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Sfe4g", "Errors.User.NotInitialised")
	}
	if authz.GetCtxData(ctx).UserID != userID {
		if err = c.checkPermission(ctx, domain.PermissionUserWrite, model.ResourceOwner, userID); err != nil {
			return nil, nil, err
		}
	}
	code, err := c.newCode(ctx, c.eventstore.Filter, domain.SecretGeneratorTypePasswordResetCode, c.userEncryption)
	if err != nil {
		return nil, nil, err
	}
	cmd := user.NewHumanPasswordCodeAddedEventV2(ctx, UserAggregateFromWriteModel(&model.WriteModel), code.Crypted, code.Expiry, domain.NotificationTypeEmail, urlTmpl, returnCode) //TODO: notification type

	if returnCode {
		plainCode = &code.Plain
	}
	if err = c.pushAppendAndReduce(ctx, model, cmd); err != nil {
		return nil, nil, err
	}
	return writeModelToObjectDetails(&model.WriteModel), plainCode, nil
}

//
//// UserPasswordResetEvent allows step-by-step additions of events,
//// operating on the Human Password Model.
//type UserPasswordResetEvent struct {
//	aggregate *eventstore.Aggregate
//	event     *user.HumanPasswordCodeAddedEvent
//	model     *HumanWriteModel
//
//	plainCode *string
//}
//
//// NewUserPasswordResetEvent constructs a UserPasswordResetEvents with a HumanWriteModel,
//// filtered by userID and resourceOwner.
//// If a model cannot be found, or it's state is invalid and error is returned.
//func (c *Commands) NewUserPasswordResetEvent(ctx context.Context, userID, resourceOwner string, alg crypto.EncryptionAlgorithm, returnCode bool, urlTmpl string) (*UserPasswordResetEvent, error) {
//	if userID == "" {
//		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-SAFdda", "Errors.User.IDMissing")
//	}
//
//	model, err := c.getHumanWriteModelByID(ctx, userID, resourceOwner)
//	if err != nil {
//		return nil, err
//	}
//	if model.UserState == domain.UserStateUnspecified || model.UserState == domain.UserStateDeleted {
//		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-SAF4f", "Errors.User.NotFound")
//	}
//	if model.UserState == domain.UserStateInitial {
//		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Sfe4g", "Errors.User.NotInitialised")
//	}
//	if authz.GetCtxData(ctx).UserID != userID {
//		if err = c.checkPermission(ctx, domain.PermissionUserWrite, model.ResourceOwner, userID); err != nil {
//			return nil, err
//		}
//	}
//	gen, _, err := secretGenerator(ctx, c.eventstore.Filter, domain.SecretGeneratorTypePasswordResetCode, alg)
//	if err != nil {
//		return nil, err
//	}
//	value, plain, err := crypto.NewCode(gen)
//	if err != nil {
//		return nil, err
//	}
//
//	e := &UserPasswordResetEvent{
//		model: model,
//	}
//	e.event = user.NewHumanPasswordCodeAddedEventV2(ctx, UserAggregateFromWriteModel(&model.WriteModel), value, gen.Expiry(), domain.NotificationTypeEmail, urlTmpl, returnCode) //TODO: notification type
//	if returnCode {
//		e.plainCode = &plain
//	}
//	return e, nil
//}
//
//// AddGeneratedCode generates a new encrypted code and sets it to the email address.
//// When returnCode a plain text of the code will be returned from Push.
//func (c *UserPasswordResetEvents) AddGeneratedCode(ctx context.Context, gen crypto.Generator, urlTmpl string, returnCode bool) error {
//
//	return nil
//}
