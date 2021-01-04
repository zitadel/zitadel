package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/v2/domain"
)

func (r *CommandSide) ChangeHumanEmail(ctx context.Context, email *usr_model.Email) (*usr_model.Email, error) {
	if !email.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4M9sf", "Errors.Email.Invalid")
	}

	existingEmail, err := r.emailWriteModel(ctx, email.AggregateID)
	if err != nil {
		return nil, err
	}
	if existingEmail.UserState == domain.UserStateUnspecified || existingEmail.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowAlreadyExists(nil, "COMMAND-0Pe4r", "Errors.User.Email.NotFound")
	}
	changedEvent, hasChanged := existingEmail.NewChangedEvent(ctx, email.EmailAddress)
	if !hasChanged {
		return nil, caos_errs.ThrowAlreadyExists(nil, "COMMAND-2M9fs", "Errors.User.Email.NotChanged")
	}
	userAgg := UserAggregateFromWriteModel(&existingEmail.WriteModel)
	userAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingEmail, userAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToEmail(existingEmail), nil
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
