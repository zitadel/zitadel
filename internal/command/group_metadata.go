package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) SetGroupMetadata(ctx context.Context, metadata *domain.Metadata, groupID, resourceOwner string) (_ *domain.Metadata, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	err = c.checkGroupExists(ctx, groupID, resourceOwner)
	if err != nil {
		return nil, err
	}
	setMetadata := NewGroupMetadataWriteModel(groupID, resourceOwner, metadata.Key)
	groupAgg := GroupAggregateFromWriteModel(&setMetadata.WriteModel)
	event, err := c.setGroupMetadata(ctx, groupAgg, metadata)
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
	return writeModelToGroupMetadata(setMetadata), nil
}

func (c *Commands) BulkSetGroupMetadata(ctx context.Context, groupID, resourceOwner string, metadatas ...*domain.Metadata) (_ *domain.ObjectDetails, err error) {
	if len(metadatas) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "META-9mm2d", "Errors.Metadata.NoData")
	}
	err = c.checkGroupExists(ctx, groupID, resourceOwner)
	if err != nil {
		return nil, err
	}

	events := make([]eventstore.Command, len(metadatas))
	setMetadata := NewGroupMetadataListWriteModel(groupID, resourceOwner)
	groupAgg := GroupAggregateFromWriteModel(&setMetadata.WriteModel)
	for i, data := range metadatas {
		event, err := c.setGroupMetadata(ctx, groupAgg, data)
		if err != nil {
			return nil, err
		}
		events[i] = event
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

func (c *Commands) setGroupMetadata(ctx context.Context, groupAgg *eventstore.Aggregate, metadata *domain.Metadata) (command eventstore.Command, err error) {
	if !metadata.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "META-2m00f", "Errors.Metadata.Invalid")
	}
	return group.NewMetadataSetEvent(
		ctx,
		groupAgg,
		metadata.Key,
		metadata.Value,
	), nil
}

func (c *Commands) RemoveGroupMetadata(ctx context.Context, metadataKey, groupID, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	if metadataKey == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "META-2n0fs", "Errors.Metadata.Invalid")
	}
	err = c.checkGroupExists(ctx, groupID, resourceOwner)
	if err != nil {
		return nil, err
	}
	removeMetadata, err := c.getGroupMetadataModelByID(ctx, groupID, resourceOwner, metadataKey)
	if err != nil {
		return nil, err
	}
	if !removeMetadata.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "META-ncnw3", "Errors.Metadata.NotFound")
	}
	groupAgg := GroupAggregateFromWriteModel(&removeMetadata.WriteModel)
	event, err := c.removeGroupMetadata(ctx, groupAgg, metadataKey)
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

func (c *Commands) BulkRemoveGroupMetadata(ctx context.Context, groupID, resourceOwner string, metadataKeys ...string) (_ *domain.ObjectDetails, err error) {
	if len(metadataKeys) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "META-9mm2d", "Errors.Metadata.NoData")
	}
	err = c.checkGroupExists(ctx, groupID, resourceOwner)
	if err != nil {
		return nil, err
	}

	events := make([]eventstore.Command, len(metadataKeys))
	removeMetadata, err := c.getGroupMetadataListModelByID(ctx, groupID, resourceOwner)
	if err != nil {
		return nil, err
	}
	groupAgg := GroupAggregateFromWriteModel(&removeMetadata.WriteModel)
	for i, key := range metadataKeys {
		if key == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-m29ds", "Errors.Metadata.Invalid")
		}
		if _, found := removeMetadata.metadataList[key]; !found {
			return nil, zerrors.ThrowNotFound(nil, "META-2nnds", "Errors.Metadata.KeyNotExisting")
		}
		event, err := c.removeGroupMetadata(ctx, groupAgg, key)
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

// func (c *Commands) removeGroupMetadataFromOrg(ctx context.Context, resourceOwner string) ([]eventstore.Command, error) {
// 	existingGroupMetadata, err := c.getGroupMetadataByOrgListModelByID(ctx, resourceOwner)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(existingGroupMetadata.GroupMetadata) == 0 {
// 		return nil, nil
// 	}
// 	events := make([]eventstore.Command, 0)
// 	for key, value := range existingGroupMetadata.GroupMetadata {
// 		if len(value) == 0 {
// 			continue
// 		}
// 		events = append(events, group.NewMetadataRemovedAllEvent(ctx, &group.NewAggregate(key, resourceOwner).Aggregate))
// 	}
// 	return events, nil
// }

func (c *Commands) removeGroupMetadata(ctx context.Context, groupAgg *eventstore.Aggregate, metadataKey string) (command eventstore.Command, err error) {
	command = group.NewMetadataRemovedEvent(
		ctx,
		groupAgg,
		metadataKey,
	)
	return command, nil
}

func (c *Commands) getGroupMetadataModelByID(ctx context.Context, groupID, resourceOwner, key string) (*GroupMetadataWriteModel, error) {
	groupMetadataWriteModel := NewGroupMetadataWriteModel(groupID, resourceOwner, key)
	err := c.eventstore.FilterToQueryReducer(ctx, groupMetadataWriteModel)
	if err != nil {
		return nil, err
	}
	return groupMetadataWriteModel, nil
}

func (c *Commands) getGroupMetadataListModelByID(ctx context.Context, groupID, resourceOwner string) (*GroupMetadataListWriteModel, error) {
	groupMetadataWriteModel := NewGroupMetadataListWriteModel(groupID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, groupMetadataWriteModel)
	if err != nil {
		return nil, err
	}
	return groupMetadataWriteModel, nil
}

func (c *Commands) getGroupMetadataByOrgListModelByID(ctx context.Context, resourceOwner string) (*GroupMetadataByOrgListWriteModel, error) {
	groupMetadataWriteModel := NewGroupMetadataByOrgListWriteModel(resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, groupMetadataWriteModel)
	if err != nil {
		return nil, err
	}
	return groupMetadataWriteModel, nil
}
