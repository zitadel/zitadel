package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/v2/domain"
)

func (r *CommandSide) ChangeHumanProfile(ctx context.Context, profile *usr_model.Profile) (*usr_model.Profile, error) {
	if !profile.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-8io0d", "Errors.User.Profile.Invalid")
	}

	existingProfile, err := r.profileWriteModelByID(ctx, profile.AggregateID)
	if err != nil {
		return nil, err
	}
	if existingProfile.UserState == domain.UserStateUnspecified || existingProfile.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowAlreadyExists(nil, "COMMAND-3M9sd", "Errors.User.Profile.NotFound")
	}
	changedEvent, hasChanged := existingProfile.NewChangedEvent(ctx, profile.FirstName, profile.LastName, profile.NickName, profile.DisplayName, profile.PreferredLanguage, domain.Gender(profile.Gender))
	if !hasChanged {
		return nil, caos_errs.ThrowAlreadyExists(nil, "COMMAND-2M0fs", "Errors.User.Profile.NotChanged")
	}
	userAgg := UserAggregateFromWriteModel(&existingProfile.WriteModel)
	userAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingProfile, userAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToProfile(existingProfile), nil
}

func (r *CommandSide) profileWriteModelByID(ctx context.Context, userID string) (writeModel *HumanProfileWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanProfileWriteModel(userID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
