package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) ChangeHumanEmail(ctx context.Context, email *domain.Email) (*domain.Email, error) {
	if !email.IsValid() || email.AggregateID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4M9sf", "Errors.Email.Invalid")
	}

	existingEmail, err := r.emailWriteModel(ctx, email.AggregateID, email.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if existingEmail.UserState == domain.UserStateUnspecified || existingEmail.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-0Pe4r", "Errors.User.Email.NotFound")
	}
	changedEvent, hasChanged := existingEmail.NewChangedEvent(ctx, email.EmailAddress)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M9fs", "Errors.User.Email.NotChanged")
	}
	userAgg := UserAggregateFromWriteModel(&existingEmail.WriteModel)
	userAgg.PushEvents(changedEvent)

	if email.IsEmailVerified {
		userAgg.PushEvents(user.NewHumanEmailVerifiedEvent(ctx))
	} else {
		emailCode, err := domain.NewEmailCode(r.emailVerificationCode)
		if err != nil {
			return nil, err
		}
		userAgg.PushEvents(user.NewHumanEmailCodeAddedEvent(ctx, emailCode.Code, emailCode.Expiry))
	}

	err = r.eventstore.PushAggregate(ctx, existingEmail, userAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToEmail(existingEmail), nil
}

func (r *CommandSide) VerifyHumanEmail(ctx context.Context, userID, code, resourceowner string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4M0ds", "Errors.User.UserIDMissing")
	}
	if code == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-çm0ds", "Errors.User.Code.Empty")
	}

	existingCode, err := r.emailWriteModel(ctx, userID, resourceowner)
	if err != nil {
		return err
	}
	if existingCode.Code == nil || existingCode.UserState == domain.UserStateUnspecified || existingCode.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-2M9fs", "Errors.User.Code.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingCode.WriteModel)
	err = crypto.VerifyCode(existingCode.CodeCreationDate, existingCode.CodeExpiry, existingCode.Code, code, r.emailVerificationCode)
	if err == nil {
		userAgg.PushEvents(user.NewHumanEmailVerifiedEvent(ctx))
		return r.eventstore.PushAggregate(ctx, existingCode, userAgg)
	}
	userAgg.PushEvents(user.NewHumanEmailVerificationFailedEvent(ctx))
	err = r.eventstore.PushAggregate(ctx, existingCode, userAgg)
	logging.LogWithFields("COMMAND-Dg2z5", "userID", userAgg.ID()).OnError(err).Error("NewHumanEmailVerificationFailedEvent push failed")
	return caos_errs.ThrowInvalidArgument(err, "COMMAND-Gdsgs", "Errors.User.Code.Invalid")
}

func (r *CommandSide) CreateHumanEmailVerificationCode(ctx context.Context, userID, resourceOwner string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4M0ds", "Errors.User.UserIDMissing")
	}

	existingEmail, err := r.emailWriteModel(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if existingEmail.UserState == domain.UserStateUnspecified || existingEmail.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-0Pe4r", "Errors.User.Email.NotFound")
	}
	if existingEmail.UserState == domain.UserStateInitial {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-E3fbw", "Errors.User.NotInitialised")
	}
	if existingEmail.IsEmailVerified {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3M9ds", "Errors.User.Email.AlreadyVerified")
	}
	userAgg := UserAggregateFromWriteModel(&existingEmail.WriteModel)
	emailCode, err := domain.NewEmailCode(r.emailVerificationCode)
	if err != nil {
		return err
	}
	userAgg.PushEvents(user.NewHumanEmailCodeAddedEvent(ctx, emailCode.Code, emailCode.Expiry))

	return r.eventstore.PushAggregate(ctx, existingEmail, userAgg)
}

func (r *CommandSide) emailWriteModel(ctx context.Context, userID, resourceOwner string) (writeModel *HumanEmailWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanEmailWriteModel(userID, resourceOwner)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
