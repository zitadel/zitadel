package command

import (
	"bytes"
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) SetUserMetadata(ctx context.Context, metadata *domain.Metadata, userID, resourceOwner string) (_ *domain.Metadata, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	userResourceOwner, err := c.checkUserExists(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}

	if err := c.checkPermissionUpdateUser(ctx, userResourceOwner, userID); err != nil {
		return nil, err
	}

	setMetadata, err := c.getUserMetadataModelByID(ctx, userID, userResourceOwner, metadata.Key)
	if err != nil {
		return nil, err
	}
	userAgg := UserAggregateFromWriteModel(&setMetadata.WriteModel)
	// return if no change in the metadata
	if bytes.Equal(setMetadata.Value, metadata.Value) {
		return writeModelToUserMetadata(setMetadata), nil
	}

	event, err := c.setUserMetadata(ctx, userAgg, metadata)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(setMetadata, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToUserMetadata(setMetadata), nil
}

func (c *Commands) BulkSetUserMetadata(ctx context.Context, userID, resourceOwner string, metadatas ...*domain.Metadata) (_ *domain.ObjectDetails, err error) {
	if len(metadatas) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "META-9mm2d", "Errors.Metadata.NoData")
	}
	userResourceOwner, err := c.checkUserExists(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}

	if err := c.checkPermissionUpdateUser(ctx, userResourceOwner, userID); err != nil {
		return nil, err
	}

	events := make([]eventstore.Command, 0)
	setMetadata, err := c.getUserMetadataListModelByID(ctx, userID, userResourceOwner)
	if err != nil {
		return nil, err
	}
	userAgg := UserAggregateFromWriteModel(&setMetadata.WriteModel)
	for _, data := range metadatas {
		// if no change to metadata no event has to be pushed
		if existingValue, ok := setMetadata.metadataList[data.Key]; ok && bytes.Equal(existingValue, data.Value) {
			continue
		}
		event, err := c.setUserMetadata(ctx, userAgg, data)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	// no changes for the metadata
	if len(events) == 0 {
		return writeModelToObjectDetails(&setMetadata.WriteModel), nil
	}

	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(setMetadata, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&setMetadata.WriteModel), nil
}

func (c *Commands) setUserMetadata(ctx context.Context, userAgg *eventstore.Aggregate, metadata *domain.Metadata) (command eventstore.Command, err error) {
	if !metadata.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "META-2m00f", "Errors.Metadata.Invalid")
	}
	return user.NewMetadataSetEvent(
		ctx,
		userAgg,
		metadata.Key,
		metadata.Value,
	), nil
}

func (c *Commands) RemoveUserMetadata(ctx context.Context, metadataKey, userID, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	if metadataKey == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "META-2n0fs", "Errors.Metadata.Invalid")
	}
	userResourceOwner, err := c.checkUserExists(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}

	if err := c.checkPermissionUpdateUser(ctx, userResourceOwner, userID); err != nil {
		return nil, err
	}

	removeMetadata, err := c.getUserMetadataModelByID(ctx, userID, userResourceOwner, metadataKey)
	if err != nil {
		return nil, err
	}
	if !removeMetadata.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "META-ncnw3", "Errors.Metadata.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&removeMetadata.WriteModel)
	event, err := c.removeUserMetadata(ctx, userAgg, metadataKey)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(removeMetadata, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&removeMetadata.WriteModel), nil
}

func (c *Commands) BulkRemoveUserMetadata(ctx context.Context, userID, resourceOwner string, metadataKeys ...string) (_ *domain.ObjectDetails, err error) {
	if len(metadataKeys) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "META-9mm2d", "Errors.Metadata.NoData")
	}
	userResourceOwner, err := c.checkUserExists(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}

	if err := c.checkPermissionUpdateUser(ctx, userResourceOwner, userID); err != nil {
		return nil, err
	}

	events := make([]eventstore.Command, len(metadataKeys))
	removeMetadata, err := c.getUserMetadataListModelByID(ctx, userID, userResourceOwner)
	if err != nil {
		return nil, err
	}
	userAgg := UserAggregateFromWriteModel(&removeMetadata.WriteModel)
	for i, key := range metadataKeys {
		if key == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-m29ds", "Errors.Metadata.Invalid")
		}
		if _, found := removeMetadata.metadataList[key]; !found {
			return nil, zerrors.ThrowNotFound(nil, "META-2nnds", "Errors.Metadata.KeyNotExisting")
		}
		event, err := c.removeUserMetadata(ctx, userAgg, key)
		if err != nil {
			return nil, err
		}
		events[i] = event
	}

	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(removeMetadata, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&removeMetadata.WriteModel), nil
}

func (c *Commands) removeUserMetadata(ctx context.Context, userAgg *eventstore.Aggregate, metadataKey string) (command eventstore.Command, err error) {
	command = user.NewMetadataRemovedEvent(
		ctx,
		userAgg,
		metadataKey,
	)
	return command, nil
}

func (c *Commands) getUserMetadataModelByID(ctx context.Context, userID, resourceOwner, key string) (*UserMetadataWriteModel, error) {
	userMetadataWriteModel := NewUserMetadataWriteModel(userID, resourceOwner, key)
	err := c.eventstore.FilterToQueryReducer(ctx, userMetadataWriteModel)
	if err != nil {
		return nil, err
	}
	return userMetadataWriteModel, nil
}

func (c *Commands) getUserMetadataListModelByID(ctx context.Context, userID, resourceOwner string) (*UserMetadataListWriteModel, error) {
	userMetadataWriteModel := NewUserMetadataListWriteModel(userID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, userMetadataWriteModel)
	if err != nil {
		return nil, err
	}
	return userMetadataWriteModel, nil
}
