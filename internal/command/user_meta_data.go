package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/user"
)

func (c *Commands) SetUserMetaData(ctx context.Context, metaData *domain.MetaData, userID, resourceOwner string) (_ *domain.MetaData, err error) {
	err = c.checkUserExists(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	setMetaData := NewUserMetaDataWriteModel(userID, resourceOwner, metaData.Key)
	userAgg := UserAggregateFromWriteModel(&setMetaData.WriteModel)
	event, err := c.setUserMetaData(ctx, userAgg, metaData)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, event)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(setMetaData, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToUserMetaData(setMetaData), nil
}

func (c *Commands) BulkSetUserMetaData(ctx context.Context, userID, resourceOwner string, metaDatas ...*domain.MetaData) (_ *domain.ObjectDetails, err error) {
	if len(metaDatas) == 0 {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "META-9mm2d", "Errors.MetaData.NoData")
	}
	err = c.checkUserExists(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}

	events := make([]eventstore.EventPusher, len(metaDatas))

	setMetaData := NewUserMetaDataListWriteModel(userID, resourceOwner)
	userAgg := UserAggregateFromWriteModel(&setMetaData.WriteModel)
	for i, data := range metaDatas {
		event, err := c.setUserMetaData(ctx, userAgg, data)
		if err != nil {
			return nil, err
		}
		events[i] = event
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(setMetaData, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&setMetaData.WriteModel), nil
}

func (c *Commands) setUserMetaData(ctx context.Context, userAgg *eventstore.Aggregate, metaData *domain.MetaData) (pusher eventstore.EventPusher, err error) {
	if !metaData.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-2m00f", "Errors.MetaData.Invalid")
	}
	pusher = user.NewMetaDataSetEvent(
		ctx,
		userAgg,
		metaData.Key,
		metaData.Value,
	)
	return pusher, nil
}

func (c *Commands) RemoveUserMetaData(ctx context.Context, metaDataKey, userID, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	err = c.checkUserExists(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}

	removeMetaData, err := c.getUserMetaDataModelByID(ctx, userID, resourceOwner, metaDataKey)
	if err != nil {
		return nil, err
	}
	if !removeMetaData.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "META-ncnw3", "Errors.MetaData.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&removeMetaData.WriteModel)
	event, err := c.removeUserMetaData(ctx, userAgg, metaDataKey)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, event)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(removeMetaData, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&removeMetaData.WriteModel), nil
}

func (c *Commands) BulkRemoveUserMetaData(ctx context.Context, userID, resourceOwner string, metaDataKeys ...string) (_ *domain.ObjectDetails, err error) {
	if len(metaDataKeys) == 0 {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "META-9mm2d", "Errors.MetaData.NoData")
	}
	err = c.checkUserExists(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}

	events := make([]eventstore.EventPusher, len(metaDataKeys))
	setMetaData, err := c.getUserMetaDataListModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	userAgg := UserAggregateFromWriteModel(&setMetaData.WriteModel)
	for i, data := range metaDataKeys {
		event, err := c.removeUserMetaData(ctx, userAgg, data)
		if err != nil {
			return nil, err
		}
		events[i] = event
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(setMetaData, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&setMetaData.WriteModel), nil
}

func (c *Commands) removeUserMetaData(ctx context.Context, userAgg *eventstore.Aggregate, metaDataKey string) (pusher eventstore.EventPusher, err error) {
	if metaDataKey == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-2m00f", "Errors.MetaData.Invalid")
	}
	pusher = user.NewMetaDataRemovedEvent(
		ctx,
		userAgg,
		metaDataKey,
	)
	return pusher, nil
}

func (c *Commands) getUserMetaDataModelByID(ctx context.Context, userID, resourceOwner, key string) (*UserMetaDataWriteModel, error) {
	userMetaDataWriteModel := NewUserMetaDataWriteModel(userID, resourceOwner, key)
	err := c.eventstore.FilterToQueryReducer(ctx, userMetaDataWriteModel)
	if err != nil {
		return nil, err
	}
	return userMetaDataWriteModel, nil
}

func (c *Commands) getUserMetaDataListModelByID(ctx context.Context, userID, resourceOwner string) (*UserMetaDataListWriteModel, error) {
	userMetaDataWriteModel := NewUserMetaDataListWriteModel(userID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, userMetaDataWriteModel)
	if err != nil {
		return nil, err
	}
	return userMetaDataWriteModel, nil
}
