package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) SetOneTimePassword(ctx context.Context, orgID, userID, passwordString string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	existingPassword, err := r.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return err
	}
	password := &domain.Password{
		SecretString:   passwordString,
		ChangeRequired: true,
	}
	userAgg := UserAggregateFromWriteModel(&existingPassword.WriteModel)
	passwordEvent, err := r.changePassword(ctx, "", password, userAgg, existingPassword)
	if err != nil {
		return err
	}
	_, err = r.eventstore.PushEvents(ctx, passwordEvent)
	return err
}

func (r *CommandSide) SetPassword(ctx context.Context, orgID, userID, code, passwordString, userAgentID string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	existingCode, err := r.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return err
	}

	if existingCode.Code == nil || existingCode.UserState == domain.UserStateUnspecified || existingCode.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-2M9fs", "Errors.User.Code.NotFound")
	}

	err = crypto.VerifyCode(existingCode.CodeCreationDate, existingCode.CodeExpiry, existingCode.Code, code, r.emailVerificationCode)
	if err != nil {
		return err
	}

	password := &domain.Password{
		SecretString:   passwordString,
		ChangeRequired: false,
	}
	userAgg := UserAggregateFromWriteModel(&existingCode.WriteModel)
	passwordEvent, err := r.changePassword(ctx, userAgentID, password, userAgg, existingCode)
	if err != nil {
		return err
	}
	_, err = r.eventstore.PushEvents(ctx, passwordEvent)
	return err
}

func (r *CommandSide) ChangePassword(ctx context.Context, orgID, userID, oldPassword, newPassword, userAgentID string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	existingPassword, err := r.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if existingPassword.Secret == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Fds3s", "Errors.User.Password.Empty")
	}
	ctx, spanPasswordComparison := tracing.NewNamedSpan(ctx, "crypto.CompareHash")
	err = crypto.CompareHash(existingPassword.Secret, []byte(oldPassword), r.userPasswordAlg)
	spanPasswordComparison.EndWithError(err)

	if err != nil {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-3M0fs", "Errors.User.Password.Invalid")
	}
	password := &domain.Password{
		SecretString:   newPassword,
		ChangeRequired: false,
	}

	userAgg := UserAggregateFromWriteModel(&existingPassword.WriteModel)
	eventPusher, err := r.changePassword(ctx, userAgentID, password, userAgg, existingPassword)
	if err != nil {
		return err
	}
	_, err = r.eventstore.PushEvents(ctx, eventPusher)
	return err
}

func (r *CommandSide) changePassword(ctx context.Context, userAgentID string, password *domain.Password, userAgg *eventstore.Aggregate, existingPassword *HumanPasswordWriteModel) (event eventstore.EventPusher, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if existingPassword.UserState == domain.UserStateUnspecified || existingPassword.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-G8dh3", "Errors.User.Email.NotFound")
	}
	if existingPassword.UserState == domain.UserStateInitial {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-M9dse", "Errors.User.NotInitialised")
	}
	pwPolicy, err := r.getOrgPasswordComplexityPolicy(ctx, userAgg.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if err := password.HashPasswordIfExisting(pwPolicy, r.userPasswordAlg); err != nil {
		return nil, err
	}
	return user.NewHumanPasswordChangedEvent(ctx, userAgg, password.SecretCrypto, password.ChangeRequired, userAgentID), nil
}

func (r *CommandSide) RequestSetPassword(ctx context.Context, userID, resourceOwner string, notifyType domain.NotificationType) (err error) {
	existingHuman, err := r.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if existingHuman.UserState == domain.UserStateUnspecified || existingHuman.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-Hj9ds", "Errors.User.NotFound")
	}
	if existingHuman.UserState == domain.UserStateInitial {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M9sd", "Errors.User.NotInitialised")
	}
	userAgg := UserAggregateFromWriteModel(&existingHuman.WriteModel)
	passwordCode, err := domain.NewPasswordCode(r.passwordVerificationCode)
	if err != nil {
		return err
	}
	_, err = r.eventstore.PushEvents(ctx, user.NewHumanPasswordCodeAddedEvent(ctx, userAgg, passwordCode.Code, passwordCode.Expiry, notifyType))
	return err
}

func (r *CommandSide) PasswordCodeSent(ctx context.Context, orgID, userID string) (err error) {
	existingPassword, err := r.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if existingPassword.UserState == domain.UserStateUnspecified || existingPassword.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-3n77z", "Errors.User.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingPassword.WriteModel)
	_, err = r.eventstore.PushEvents(ctx, user.NewHumanPasswordCodeSentEvent(ctx, userAgg))
	return err
}

func (r *CommandSide) HumanCheckPassword(ctx context.Context, orgID, userID, password string, authRequest *domain.AuthRequest) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if password == "" {
		return caos_errs.ThrowNotFound(nil, "COMMAND-3n8fs", "Errors.User.Password.Empty")
	}

	existingPassword, err := r.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if existingPassword.UserState == domain.UserStateUnspecified || existingPassword.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-3n77z", "Errors.User.NotFound")
	}

	if existingPassword.Secret == nil {
		return caos_errs.ThrowNotFound(nil, "COMMAND-3n77z", "Errors.User.Password.NotSet")
	}

	userAgg := UserAggregateFromWriteModel(&existingPassword.WriteModel)
	ctx, spanPasswordComparison := tracing.NewNamedSpan(ctx, "crypto.CompareHash")
	err = crypto.CompareHash(existingPassword.Secret, []byte(password), r.userPasswordAlg)
	spanPasswordComparison.EndWithError(err)
	if err == nil {
		_, err = r.eventstore.PushEvents(ctx, user.NewHumanPasswordCheckSucceededEvent(ctx, userAgg, authRequestDomainToAuthRequestInfo(authRequest)))
		return err
	}
	_, err = r.eventstore.PushEvents(ctx, user.NewHumanPasswordCheckFailedEvent(ctx, userAgg, authRequestDomainToAuthRequestInfo(authRequest)))
	logging.Log("COMMAND-9fj7s").OnError(err).Error("error create password check failed event")
	return caos_errs.ThrowInvalidArgument(nil, "COMMAND-452ad", "Errors.User.Password.Invalid")
}

func (r *CommandSide) passwordWriteModel(ctx context.Context, userID, resourceOwner string) (writeModel *HumanPasswordWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanPasswordWriteModel(userID, resourceOwner)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
