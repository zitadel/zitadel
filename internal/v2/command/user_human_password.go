package command

import (
	"context"
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
	return r.changePassword(ctx, orgID, userID, "", password, userAgg, existingPassword)
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
	return r.changePassword(ctx, orgID, userID, userAgentID, password, userAgg, existingCode)
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
	return r.changePassword(ctx, orgID, userID, userAgentID, password, userAgg, existingPassword)
}

func (r *CommandSide) changePassword(ctx context.Context, orgID, userID, userAgentID string, password *domain.Password, userAgg *user.Aggregate, existingPassword *HumanPasswordWriteModel) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-B8hY3", "Errors.User.UserIDMissing")
	}
	if existingPassword.UserState == domain.UserStateUnspecified || existingPassword.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-G8dh3", "Errors.User.Email.NotFound")
	}
	if existingPassword.UserState == domain.UserStateInitial {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-M9dse", "Errors.User.NotInitialised")
	}
	pwPolicy, err := r.getOrgPasswordComplexityPolicy(ctx, orgID)
	if err != nil {
		return err
	}
	if err := password.HashPasswordIfExisting(pwPolicy, r.userPasswordAlg); err != nil {
		return err
	}
	userAgg.PushEvents(user.NewHumanPasswordChangedEvent(ctx, password.SecretCrypto, password.ChangeRequired, userAgentID))
	return r.eventstore.PushAggregate(ctx, existingPassword, userAgg)
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
	userAgg.PushEvents(user.NewHumanPasswordCodeAddedEvent(ctx, passwordCode.Code, passwordCode.Expiry, notifyType))
	return r.eventstore.PushAggregate(ctx, existingHuman, userAgg)
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
