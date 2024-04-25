package command

import (
	"context"
	"errors"

	"github.com/zitadel/logging"
	"github.com/zitadel/passwap"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) SetPassword(ctx context.Context, orgID, userID, password string, oneTime bool) (objectDetails *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-3M0fs", "Errors.IDMissing")
	}
	wm, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}
	if !wm.UserState.Exists() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-3M0fs", "Errors.User.NotFound")
	}
	if err = c.checkPermission(ctx, domain.PermissionUserWrite, wm.ResourceOwner, userID); err != nil {
		return nil, err
	}
	return c.setPassword(ctx, wm, password, "", oneTime)
}

func (c *Commands) SetPasswordWithVerifyCode(ctx context.Context, orgID, userID, code, password, userAgentID string) (objectDetails *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-3M9fs", "Errors.IDMissing")
	}
	if password == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Mf0sd", "Errors.User.Password.Empty")
	}
	wm, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}

	if wm.Code == nil || wm.UserState == domain.UserStateUnspecified || wm.UserState == domain.UserStateDeleted {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-2M9fs", "Errors.User.Code.NotFound")
	}

	err = crypto.VerifyCode(wm.CodeCreationDate, wm.CodeExpiry, wm.Code, code, c.userEncryption)
	if err != nil {
		return nil, err
	}

	return c.setPassword(ctx, wm, password, userAgentID, false)
}

// setEncodedPassword add change event from already encoded password to HumanPasswordWriteModel and return the necessary object details for response
func (c *Commands) setEncodedPassword(ctx context.Context, wm *HumanPasswordWriteModel, password, userAgentID string, changeRequired bool) (objectDetails *domain.ObjectDetails, err error) {
	agg := user.NewAggregate(wm.AggregateID, wm.ResourceOwner)
	command, err := c.setPasswordCommand(ctx, &agg.Aggregate, wm.UserState, password, userAgentID, changeRequired, true)
	if err != nil {
		return nil, err
	}
	err = c.pushAppendAndReduce(ctx, wm, command)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

// setPassword add change event to HumanPasswordWriteModel and return the necessary object details for response
func (c *Commands) setPassword(ctx context.Context, wm *HumanPasswordWriteModel, password, userAgentID string, changeRequired bool) (objectDetails *domain.ObjectDetails, err error) {
	agg := user.NewAggregate(wm.AggregateID, wm.ResourceOwner)
	command, err := c.setPasswordCommand(ctx, &agg.Aggregate, wm.UserState, password, userAgentID, changeRequired, false)
	if err != nil {
		return nil, err
	}
	err = c.pushAppendAndReduce(ctx, wm, command)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func (c *Commands) setPasswordCommand(ctx context.Context, agg *eventstore.Aggregate, userState domain.UserState, password, userAgentID string, changeRequired, encoded bool) (_ eventstore.Command, err error) {
	if err = c.canUpdatePassword(ctx, password, agg.ResourceOwner, userState); err != nil {
		return nil, err
	}

	if !encoded {
		ctx, span := tracing.NewNamedSpan(ctx, "passwap.Hash")
		encodedPassword, err := c.userPasswordHasher.Hash(password)
		span.EndWithError(err)
		if err = convertPasswapErr(err); err != nil {
			return nil, err
		}
		return user.NewHumanPasswordChangedEvent(ctx, agg, encodedPassword, changeRequired, userAgentID), nil
	}
	return user.NewHumanPasswordChangedEvent(ctx, agg, password, changeRequired, userAgentID), nil
}

// ChangePassword change password of existing user
func (c *Commands) ChangePassword(ctx context.Context, orgID, userID, oldPassword, newPassword, userAgentID string) (objectDetails *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-3M0fs", "Errors.IDMissing")
	}
	if oldPassword == "" || newPassword == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-3M0fs", "Errors.User.Password.Empty")
	}
	wm, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}

	newPasswordHash, err := c.verifyAndUpdatePassword(ctx, wm.EncodedHash, oldPassword, newPassword)
	if err != nil {
		return nil, err
	}
	return c.setEncodedPassword(ctx, wm, newPasswordHash, userAgentID, false)
}

// verifyAndUpdatePassword verify if the old password is correct with the encoded hash and
// returns the hash of the new password if so
func (c *Commands) verifyAndUpdatePassword(ctx context.Context, encodedHash, oldPassword, newPassword string) (string, error) {
	if encodedHash == "" {
		return "", zerrors.ThrowPreconditionFailed(nil, "COMMAND-Fds3s", "Errors.User.Password.NotSet")
	}

	_, spanPasswap := tracing.NewNamedSpan(ctx, "passwap.Verify")
	updated, err := c.userPasswordHasher.VerifyAndUpdate(encodedHash, oldPassword, newPassword)
	spanPasswap.EndWithError(err)
	return updated, convertPasswapErr(err)
}

// canUpdatePassword checks uf the given password can be used to be the password of a user
func (c *Commands) canUpdatePassword(ctx context.Context, newPassword string, resourceOwner string, state domain.UserState) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if !isUserStateExists(state) {
		return zerrors.ThrowNotFound(nil, "COMMAND-G8dh3", "Errors.User.Password.NotFound")
	}
	if state == domain.UserStateInitial {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-M9dse", "Errors.User.NotInitialised")
	}
	policy, err := c.getOrgPasswordComplexityPolicy(ctx, resourceOwner)
	if err != nil {
		return err
	}

	if err := policy.Check(newPassword); err != nil {
		return err
	}
	return nil
}

// RequestSetPassword generate and send out new code to change password for a specific user
func (c *Commands) RequestSetPassword(ctx context.Context, userID, resourceOwner string, notifyType domain.NotificationType, passwordVerificationCode crypto.Generator, authRequestID string) (objectDetails *domain.ObjectDetails, err error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-M00oL", "Errors.User.UserIDMissing")
	}

	existingHuman, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingHuman.UserState) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Hj9ds", "Errors.User.NotFound")
	}
	if existingHuman.UserState == domain.UserStateInitial {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-2M9sd", "Errors.User.NotInitialised")
	}
	userAgg := UserAggregateFromWriteModel(&existingHuman.WriteModel)
	passwordCode, err := domain.NewPasswordCode(passwordVerificationCode)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, user.NewHumanPasswordCodeAddedEvent(ctx, userAgg, passwordCode.Code, passwordCode.Expiry, notifyType, authRequestID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingHuman, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingHuman.WriteModel), nil
}

// PasswordCodeSent notification send with code to change password
func (c *Commands) PasswordCodeSent(ctx context.Context, orgID, userID string) (err error) {
	if userID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-meEfe", "Errors.User.UserIDMissing")
	}

	existingPassword, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if existingPassword.UserState == domain.UserStateUnspecified || existingPassword.UserState == domain.UserStateDeleted {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-3n77z", "Errors.User.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingPassword.WriteModel)
	_, err = c.eventstore.Push(ctx, user.NewHumanPasswordCodeSentEvent(ctx, userAgg))
	return err
}

// PasswordChangeSent notification sent that user changed password
func (c *Commands) PasswordChangeSent(ctx context.Context, orgID, userID string) (err error) {
	if userID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-pqlm2n", "Errors.User.UserIDMissing")
	}

	existingPassword, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if existingPassword.UserState == domain.UserStateUnspecified || existingPassword.UserState == domain.UserStateDeleted {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-x902b2v", "Errors.User.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingPassword.WriteModel)
	_, err = c.eventstore.Push(ctx, user.NewHumanPasswordChangeSentEvent(ctx, userAgg))
	return err
}

// HumanCheckPassword check password for user with additional informations from authRequest
func (c *Commands) HumanCheckPassword(ctx context.Context, orgID, userID, password string, authRequest *domain.AuthRequest, lockoutPolicy *domain.LockoutPolicy) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if userID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-4Mfsf", "Errors.User.UserIDMissing")
	}
	if password == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-3n8fs", "Errors.User.Password.Empty")
	}

	loginPolicy, err := c.getOrgLoginPolicy(ctx, orgID)
	if err != nil {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-Edf3g", "Errors.Org.LoginPolicy.NotFound")
	}
	if !loginPolicy.AllowUsernamePassword {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-Dft32", "Errors.Org.LoginPolicy.UsernamePasswordNotAllowed")
	}

	wm, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return err
	}

	if !isUserStateExists(wm.UserState) {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-3n77z", "Errors.User.NotFound")
	}
	if wm.UserState == domain.UserStateLocked {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-JLK35", "Errors.User.Locked")
	}
	if wm.EncodedHash == "" {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-3nJ4t", "Errors.User.Password.NotSet")
	}

	userAgg := UserAggregateFromWriteModel(&wm.WriteModel)
	ctx, spanPasswordComparison := tracing.NewNamedSpan(ctx, "passwap.Verify")
	updated, err := c.userPasswordHasher.Verify(wm.EncodedHash, password)
	spanPasswordComparison.EndWithError(err)
	err = convertPasswapErr(err)
	commands := make([]eventstore.Command, 0, 2)

	// recheck for additional events (failed password checks or locks)
	recheckErr := c.eventstore.FilterToQueryReducer(ctx, wm)
	if recheckErr != nil {
		return recheckErr
	}
	if wm.UserState == domain.UserStateLocked {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-SFA3t", "Errors.User.Locked")
	}

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
		return zerrors.ThrowInvalidArgument(err, "COMMAND-3M0fs", "Errors.User.Password.Invalid")
	}
	if errors.Is(err, passwap.ErrPasswordNoChange) {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-Aesh5", "Errors.User.Password.NotChanged")
	}
	return zerrors.ThrowInternal(err, "COMMAND-CahN2", "Errors.Internal")
}
