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
	userAgg := UserAggregateFromWriteModel(&existingEmail.WriteModel)
	changedEvent, hasChanged := existingEmail.NewChangedEvent(ctx, userAgg, email.EmailAddress)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2b7fM", "Errors.User.Email.NotChanged")
	}

	events := []eventstore.EventPusher{changedEvent}

	if email.IsEmailVerified {
		events = append(events, user.NewHumanEmailVerifiedEvent(ctx, userAgg))
	} else {
		emailCode, err := domain.NewEmailCode(r.emailVerificationCode)
		if err != nil {
			return nil, err
		}
		events = append(events, user.NewHumanEmailCodeAddedEvent(ctx, userAgg, emailCode.Code, emailCode.Expiry))
	}

	pushedEvents, err := r.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingEmail, pushedEvents...)
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
		return caos_errs.ThrowNotFound(nil, "COMMAND-3n8ud", "Errors.User.Code.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingCode.WriteModel)
	err = crypto.VerifyCode(existingCode.CodeCreationDate, existingCode.CodeExpiry, existingCode.Code, code, r.emailVerificationCode)
	if err == nil {
		_, err = r.eventstore.PushEvents(ctx, user.NewHumanEmailVerifiedEvent(ctx, userAgg))
		return err
	}

	_, err = r.eventstore.PushEvents(ctx, user.NewHumanEmailVerificationFailedEvent(ctx, userAgg))
	logging.LogWithFields("COMMAND-Dg2z5", "userID", userAgg.ID).OnError(err).Error("NewHumanEmailVerificationFailedEvent push failed")
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
	_, err = r.eventstore.PushEvents(ctx, user.NewHumanEmailCodeAddedEvent(ctx, userAgg, emailCode.Code, emailCode.Expiry))
	return err
}

func (r *CommandSide) HumanEmailVerificationCodeSent(ctx context.Context, orgID, userID string) (err error) {
	existingEmail, err := r.emailWriteModel(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if existingEmail.UserState == domain.UserStateUnspecified || existingEmail.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-6n8uH", "Errors.User.Email.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingEmail.WriteModel)
	_, err = r.eventstore.PushEvents(ctx, user.NewHumanEmailCodeSentEvent(ctx, userAgg))
	return err
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
