package command

import (
	"context"
	"errors"

	"github.com/zitadel/logging"
	"github.com/zitadel/passwap"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (c *Commands) SetPassword(ctx context.Context, orgID, userID, password string, oneTime bool) (objectDetails *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-3M0fs", "Errors.IDMissing")
	}
	wm, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}
	if !wm.UserState.Exists() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3M0fs", "Errors.User.NotFound")
	}
	if err = c.checkPermission(ctx, domain.PermissionUserWrite, wm.ResourceOwner, userID); err != nil {
		return nil, err
	}
	return c.setPassword(ctx, wm, password, oneTime)
}

func (c *Commands) SetPasswordWithVerifyCode(ctx context.Context, orgID, userID, code, password, userAgentID string) (objectDetails *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-3M9fs", "Errors.IDMissing")
	}
	if password == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-Mf0sd", "Errors.User.Password.Empty")
	}
	wm, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}

	if wm.Code == nil || wm.UserState == domain.UserStateUnspecified || wm.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M9fs", "Errors.User.Code.NotFound")
	}

	err = crypto.VerifyCodeWithAlgorithm(wm.CodeCreationDate, wm.CodeExpiry, wm.Code, code, c.userEncryption)
	if err != nil {
		return nil, err
	}

	return c.setPassword(ctx, wm, password, false)
}

func (c *Commands) setPassword(ctx context.Context, wm *HumanPasswordWriteModel, password string, changeRequired bool) (objectDetails *domain.ObjectDetails, err error) {
	command, err := c.setPasswordCommand(ctx, wm, password, changeRequired)
	if err != nil {
		return nil, err
	}
	err = c.pushAppendAndReduce(ctx, wm, command)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func (c *Commands) setPasswordCommand(ctx context.Context, wm *HumanPasswordWriteModel, password string, changeRequired bool) (_ eventstore.Command, err error) {
	if err = c.canUpdatePassword(ctx, password, wm); err != nil {
		return nil, err
	}
	ctx, span := tracing.NewNamedSpan(ctx, "passwap.Hash")
	encoded, err := c.userPasswordHasher.Hash(password)
	span.EndWithError(err)
	if err = convertPasswapErr(err); err != nil {
		return nil, err
	}
	return user.NewHumanPasswordChangedEvent(ctx, UserAggregateFromWriteModel(&wm.WriteModel), encoded, changeRequired, ""), nil
}

func (c *Commands) ChangePassword(ctx context.Context, orgID, userID, oldPassword, newPassword, userAgentID string) (objectDetails *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-3M0fs", "Errors.IDMissing")
	}
	if oldPassword == "" || newPassword == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-3M0fs", "Errors.User.Password.Empty")
	}
	wm, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}
	if wm.EncodedHash == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Fds3s", "Errors.User.Password.Empty")
	}
	if err = c.canUpdatePassword(ctx, newPassword, wm); err != nil {
		return nil, err
	}

	ctx, spanPasswap := tracing.NewNamedSpan(ctx, "passwap.VerifyAndUpdate")
	updated, err := c.userPasswordHasher.VerifyAndUpdate(wm.EncodedHash, oldPassword, newPassword)
	spanPasswap.EndWithError(err)
	if err = convertPasswapErr(err); err != nil {
		return nil, err
	}
	err = c.pushAppendAndReduce(ctx, wm,
		user.NewHumanPasswordChangedEvent(ctx, UserAggregateFromWriteModel(&wm.WriteModel), updated, false, userAgentID))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func (c *Commands) canUpdatePassword(ctx context.Context, newPassword string, wm *HumanPasswordWriteModel) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if wm.UserState == domain.UserStateUnspecified || wm.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-G8dh3", "Errors.User.Password.NotFound")
	}
	if wm.UserState == domain.UserStateInitial {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-M9dse", "Errors.User.NotInitialised")
	}
	policy, err := c.getOrgPasswordComplexityPolicy(ctx, wm.ResourceOwner)
	if err != nil {
		return err
	}

	if err := policy.Check(newPassword); err != nil {
		return err
	}
	return nil
}

func (c *Commands) RequestSetPassword(ctx context.Context, userID, resourceOwner string, notifyType domain.NotificationType, passwordVerificationCode crypto.Generator) (objectDetails *domain.ObjectDetails, err error) {
	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-M00oL", "Errors.User.UserIDMissing")
	}

	existingHuman, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingHuman.UserState == domain.UserStateUnspecified || existingHuman.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Hj9ds", "Errors.User.NotFound")
	}
	if existingHuman.UserState == domain.UserStateInitial {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M9sd", "Errors.User.NotInitialised")
	}
	userAgg := UserAggregateFromWriteModel(&existingHuman.WriteModel)
	passwordCode, err := domain.NewPasswordCode(passwordVerificationCode)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, user.NewHumanPasswordCodeAddedEvent(ctx, userAgg, passwordCode.Code, passwordCode.Expiry, notifyType))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingHuman, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingHuman.WriteModel), nil
}

func (c *Commands) PasswordCodeSent(ctx context.Context, orgID, userID string) (err error) {
	if userID == "" {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-meEfe", "Errors.User.UserIDMissing")
	}

	existingPassword, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if existingPassword.UserState == domain.UserStateUnspecified || existingPassword.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3n77z", "Errors.User.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingPassword.WriteModel)
	_, err = c.eventstore.Push(ctx, user.NewHumanPasswordCodeSentEvent(ctx, userAgg))
	return err
}

func (c *Commands) PasswordChangeSent(ctx context.Context, orgID, userID string) (err error) {
	if userID == "" {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-pqlm2n", "Errors.User.UserIDMissing")
	}

	existingPassword, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if existingPassword.UserState == domain.UserStateUnspecified || existingPassword.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-x902b2v", "Errors.User.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingPassword.WriteModel)
	_, err = c.eventstore.Push(ctx, user.NewHumanPasswordChangeSentEvent(ctx, userAgg))
	return err
}

func (c *Commands) HumanCheckPassword(ctx context.Context, orgID, userID, password string, authRequest *domain.AuthRequest, lockoutPolicy *domain.LockoutPolicy) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if userID == "" {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-4Mfsf", "Errors.User.UserIDMissing")
	}
	if password == "" {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-3n8fs", "Errors.User.Password.Empty")
	}

	loginPolicy, err := c.getOrgLoginPolicy(ctx, orgID)
	if err != nil {
		return caos_errs.ThrowPreconditionFailed(err, "COMMAND-Edf3g", "Errors.Org.LoginPolicy.NotFound")
	}
	if !loginPolicy.AllowUsernamePassword {
		return caos_errs.ThrowPreconditionFailed(err, "COMMAND-Dft32", "Errors.Org.LoginPolicy.UsernamePasswordNotAllowed")
	}

	wm, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if wm.UserState == domain.UserStateUnspecified || wm.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3n77z", "Errors.User.NotFound")
	}

	if wm.EncodedHash == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3n77z", "Errors.User.Password.NotSet")
	}

	userAgg := UserAggregateFromWriteModel(&wm.WriteModel)
	ctx, spanPasswordComparison := tracing.NewNamedSpan(ctx, "passwap.Verify")
	updated, err := c.userPasswordHasher.Verify(wm.EncodedHash, password)
	spanPasswordComparison.EndWithError(err)
	err = convertPasswapErr(err)

	commands := make([]eventstore.Command, 0, 2)
	if err == nil {
		commands = append(commands, user.NewHumanPasswordCheckSucceededEvent(ctx, userAgg, authRequestDomainToAuthRequestInfo(authRequest)))
		if updated != "" {
			commands = append(commands, user.NewHumanPasswordHashUpdatedEvent(ctx, userAgg, updated))
		}
		_, err = c.eventstore.Push(ctx, commands...)
		return err
	}

	commands = append(commands, user.NewHumanPasswordCheckFailedEvent(ctx, userAgg, authRequestDomainToAuthRequestInfo(authRequest)))
	if lockoutPolicy != nil && lockoutPolicy.MaxPasswordAttempts > 0 {
		if wm.PasswordCheckFailedCount+1 >= lockoutPolicy.MaxPasswordAttempts {
			commands = append(commands, user.NewUserLockedEvent(ctx, userAgg))
		}

	}
	_, pushErr := c.eventstore.Push(ctx, commands...)
	logging.OnError(pushErr).Error("error create password check failed event")
	return err
}

func (c *Commands) passwordWriteModel(ctx context.Context, userID, resourceOwner string) (writeModel *HumanPasswordWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanPasswordWriteModel(userID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func convertPasswapErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, passwap.ErrPasswordMismatch) {
		return caos_errs.ThrowInvalidArgument(err, "COMMAND-3M0fs", "Errors.User.Password.Invalid")
	}
	if errors.Is(err, passwap.ErrPasswordNoChange) {
		return caos_errs.ThrowPreconditionFailed(err, "COMMAND-Aesh5", "Errors.User.Password.NotChanged")
	}
	return caos_errs.ThrowInternal(err, "COMMAND-CahN2", "Errors.Internal")
}
