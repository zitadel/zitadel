package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/user"
)

func (c *Commands) AddHumanAvatar(ctx context.Context, orgID, userID, storageKey string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "USER-Ba5Ds", "Errors.IDMissing")
	}
	if storageKey == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "USER-1Xyud", "Errors.Assets.EmptyKey")
	}
	existingUser, err := c.userWriteModelByID(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}

	if existingUser.UserState == domain.UserStateUnspecified || existingUser.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "USER-vJ3fS", "Errors.Users.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingUser.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, user.NewHumanAvatarAddedEvent(ctx, userAgg, storageKey))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingUser, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingUser.WriteModel), nil
}

func (c *Commands) RemoveHumanAvatar(ctx context.Context, orgID, userID string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "USER-1B8sd", "Errors.IDMissing")
	}
	existingUser, err := c.getHumanWriteModelByID(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}
	if existingUser.UserState == domain.UserStateUnspecified || existingUser.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "USER-35N8f", "Errors.Users.NotFound")
	}
	err = c.RemoveAsset(ctx, orgID, existingUser.Avatar)
	if err != nil {
		return nil, err
	}
	userAgg := UserAggregateFromWriteModel(&existingUser.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, user.NewHumanAvatarRemovedEvent(ctx, userAgg, existingUser.Avatar))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingUser, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingUser.WriteModel), nil
}
