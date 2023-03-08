package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (c *Commands) ChangeHumanProfile(ctx context.Context, profile *domain.Profile) (*domain.Profile, error) {
	if profile.AggregateID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-AwbEB", "Errors.User.Profile.IDMissing")
	}
	if err := profile.Validate(); err != nil {
		return nil, err
	}
	existingProfile, err := c.profileWriteModelByID(ctx, profile.AggregateID, profile.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if existingProfile.UserState == domain.UserStateUnspecified || existingProfile.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3M9sd", "Errors.User.Profile.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingProfile.WriteModel)
	changedEvent, hasChanged, err := existingProfile.NewChangedEvent(ctx, userAgg, profile.FirstName, profile.LastName, profile.NickName, profile.DisplayName, profile.PreferredLanguage, profile.Gender)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M0fs", "Errors.User.Profile.NotChanged")
	}

	events, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingProfile, events...)
	if err != nil {
		return nil, err
	}

	return writeModelToProfile(existingProfile), nil
}

func (c *Commands) profileWriteModelByID(ctx context.Context, userID, resourceOwner string) (writeModel *HumanProfileWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanProfileWriteModel(userID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
