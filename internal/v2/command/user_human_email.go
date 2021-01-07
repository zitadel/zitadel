package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) ChangeHumanEmail(ctx context.Context, email *domain.Email) (*domain.Email, error) {
	if !email.IsValid() || email.AggregateID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4M9sf", "Errors.Email.Invalid")
	}

	existingEmail, err := r.emailWriteModel(ctx, email.AggregateID)
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

func (r *CommandSide) CreateHumanEmailVerificationCode(ctx context.Context, userID string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4M0ds", "Errors.User.UserIDMissing")
	}

	existingEmail, err := r.emailWriteModel(ctx, userID)
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

func (r *CommandSide) emailWriteModel(ctx context.Context, userID string) (writeModel *HumanEmailWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanEmailWriteModel(userID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
