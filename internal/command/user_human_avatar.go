package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/user"
)

func (c *Commands) ChangeAvatar(ctx context.Context, orgID, userID string, avatar []byte) (objectDetails *domain.ObjectDetails, err error) {
	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-2M9fs", "Errors.IDMissing")
	}
	existingAvatar, err := c.avatarWriteModel(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}
	if !existingAvatar.UserState.Exists() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-11Fp0", "Errors.User.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingAvatar.WriteModel)
	assetID, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	changeAvatar := user.NewHumanAvatarChangedEvent(ctx, userAgg, assetID, avatar)
	pushedEvents, err := c.eventstore.PushEvents(ctx, changeAvatar)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingAvatar, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingAvatar.WriteModel), nil
}

func (c *Commands) avatarWriteModel(ctx context.Context, userID, resourceOwner string) (writeModel *HumanAvatarWriteModel, err error) {
	writeModel = NewHumanAvatarWriteModel(userID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
