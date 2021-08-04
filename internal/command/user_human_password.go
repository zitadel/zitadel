package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/user"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) SetPassword(ctx context.Context, orgID, userID, passwordString string, oneTime bool) (objectDetails *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-3M0fs", "Errors.IDMissing")
	}
	existingPassword, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}
	if !existingPassword.UserState.Exists() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3M0fs", "Errors.User.NotFound")
	}
	password := &domain.Password{
		SecretString:   passwordString,
		ChangeRequired: oneTime,
	}
	userAgg := UserAggregateFromWriteModel(&existingPassword.WriteModel)
	passwordEvent, err := c.changePassword(ctx, "", password, userAgg, existingPassword)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, passwordEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPassword, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPassword.WriteModel), nil
}

func (c *Commands) SetPasswordWithVerifyCode(ctx context.Context, orgID, userID, code, passwordString, userAgentID string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if userID == "" {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-3M9fs", "Errors.IDMissing")
	}
	if passwordString == "" {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-Mf0sd", "Errors.User.Password.Empty")
	}
	existingCode, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return err
	}

	if existingCode.Code == nil || existingCode.UserState == domain.UserStateUnspecified || existingCode.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M9fs", "Errors.User.Code.NotFound")
	}

	err = crypto.VerifyCode(existingCode.CodeCreationDate, existingCode.CodeExpiry, existingCode.Code, code, c.passwordVerificationCode)
	if err != nil {
		return err
	}

	password := &domain.Password{
		SecretString:   passwordString,
		ChangeRequired: false,
	}
	userAgg := UserAggregateFromWriteModel(&existingCode.WriteModel)
	passwordEvent, err := c.changePassword(ctx, userAgentID, password, userAgg, existingCode)
	if err != nil {
		return err
	}
	_, err = c.eventstore.PushEvents(ctx, passwordEvent)
	return err
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
	existingPassword, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}
	if existingPassword.Secret == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Fds3s", "Errors.User.Password.Empty")
	}
	ctx, spanPasswordComparison := tracing.NewNamedSpan(ctx, "crypto.CompareHash")
	err = crypto.CompareHash(existingPassword.Secret, []byte(oldPassword), c.userPasswordAlg)
	spanPasswordComparison.EndWithError(err)

	if err != nil {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-3M0fs", "Errors.User.Password.Invalid")
	}
	password := &domain.Password{
		SecretString:   newPassword,
		ChangeRequired: false,
	}

	userAgg := UserAggregateFromWriteModel(&existingPassword.WriteModel)
	eventPusher, err := c.changePassword(ctx, userAgentID, password, userAgg, existingPassword)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, eventPusher)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPassword, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPassword.WriteModel), nil
}

func (c *Commands) changePassword(ctx context.Context, userAgentID string, password *domain.Password, userAgg *eventstore.Aggregate, existingPassword *HumanPasswordWriteModel) (event eventstore.EventPusher, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if existingPassword.UserState == domain.UserStateUnspecified || existingPassword.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-G8dh3", "Errors.User.Password.NotFound")
	}
	if existingPassword.UserState == domain.UserStateInitial {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-M9dse", "Errors.User.NotInitialised")
	}
	pwPolicy, err := c.getOrgPasswordComplexityPolicy(ctx, userAgg.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if err := password.HashPasswordIfExisting(pwPolicy, c.userPasswordAlg); err != nil {
		return nil, err
	}
	return user.NewHumanPasswordChangedEvent(ctx, userAgg, password.SecretCrypto, password.ChangeRequired, userAgentID), nil
}

func (c *Commands) RequestSetPassword(ctx context.Context, userID, resourceOwner string, notifyType domain.NotificationType) (objectDetails *domain.ObjectDetails, err error) {
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
	passwordCode, err := domain.NewPasswordCode(c.passwordVerificationCode)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, user.NewHumanPasswordCodeAddedEvent(ctx, userAgg, passwordCode.Code, passwordCode.Expiry, notifyType))
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
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-MM9fs", "Errors.User.UserIDMissing")
	}

	existingPassword, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if existingPassword.UserState == domain.UserStateUnspecified || existingPassword.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3n77z", "Errors.User.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingPassword.WriteModel)
	_, err = c.eventstore.PushEvents(ctx, user.NewHumanPasswordCodeSentEvent(ctx, userAgg))
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

	existingPassword, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if existingPassword.UserState == domain.UserStateUnspecified || existingPassword.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3n77z", "Errors.User.NotFound")
	}

	if existingPassword.Secret == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3n77z", "Errors.User.Password.NotSet")
	}

	userAgg := UserAggregateFromWriteModel(&existingPassword.WriteModel)
	ctx, spanPasswordComparison := tracing.NewNamedSpan(ctx, "crypto.CompareHash")
	err = crypto.CompareHash(existingPassword.Secret, []byte(password), c.userPasswordAlg)
	spanPasswordComparison.EndWithError(err)
	if err == nil {
		_, err = c.eventstore.PushEvents(ctx, user.NewHumanPasswordCheckSucceededEvent(ctx, userAgg, authRequestDomainToAuthRequestInfo(authRequest)))
		return err
	}
	events := make([]eventstore.EventPusher, 0)
	errMsg := "Errors.User.Password.Invalid"
	events = append(events, user.NewHumanPasswordCheckFailedEvent(ctx, userAgg, authRequestDomainToAuthRequestInfo(authRequest)))
	if lockoutPolicy != nil && lockoutPolicy.MaxPasswordAttempts > 0 {
		if existingPassword.PasswordCheckFailedCount+1 >= lockoutPolicy.MaxPasswordAttempts {
			events = append(events, user.NewUserLockedEvent(ctx, userAgg))
		}

	}
	_, err = c.eventstore.PushEvents(ctx, events...)
	logging.Log("COMMAND-9fj7s").OnError(err).Error("error create password check failed event")
	return caos_errs.ThrowInvalidArgument(nil, "COMMAND-452ad", errMsg)
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
