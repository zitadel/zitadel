package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) SetProjectMetadata(ctx context.Context, projectID, resourceOwner string, metadata *domain.Metadata) (_ *domain.Metadata, err error) {
	_, err = c.checkProjectExists(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}

	setMetadata := NewProjectMetadataWriteModel(projectID, resourceOwner, metadata.Key)
	projectAgg := ProjectAggregateFromWriteModel(&setMetadata.WriteModel)

	event, err := c.setProjectMetadata(ctx, projectAgg, metadata)
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

	return writeModelToProjectMetadata(setMetadata), nil
}

func (c *Commands) BulkSetProjectMetadata(ctx context.Context, projectID, resourceOwner string, metadatas ...*domain.Metadata) (_ *domain.ObjectDetails, err error) {
	if len(metadatas) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "META-zUTude", "Errors.Metadata.NoData")
	}

	_, err = c.checkProjectExists(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}

	events := make([]eventstore.Command, len(metadatas))
	setMetadata := NewProjectMetadataListWriteModel(projectID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&setMetadata.WriteModel)
	for i, data := range metadatas {
		event, err := c.setProjectMetadata(ctx, projectAgg, data)
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

func (c *Commands) setProjectMetadata(ctx context.Context, projectAgg *eventstore.Aggregate, metadata *domain.Metadata) (command eventstore.Command, err error) {
	if !metadata.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "META-cls90f", "Errors.Metadata.Invalid")
	}

	return project.NewMetadataSetEvent(
		ctx,
		projectAgg,
		metadata.Key,
		metadata.Value,
	), nil
}

func (c *Commands) RemoveProjectMetadata(ctx context.Context, projectID, resourceOwner, metadataKey string) (_ *domain.ObjectDetails, err error) {
	if metadataKey == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "META-y7tSu8", "Errors.Metadata.Invalid")
	}

	_, err = c.checkProjectExists(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}

	removeMetadata, err := c.getProjectMetadataModelByID(ctx, projectID, resourceOwner, metadataKey)
	if err != nil {
		return nil, err
	}

	if !removeMetadata.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "META-ux4yJT", "Errors.Metadata.NotFound")
	}

	projectAgg := ProjectAggregateFromWriteModel(&removeMetadata.WriteModel)
	event, err := c.removeProjectMetadata(ctx, projectAgg, metadataKey)
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

func (c *Commands) BulkRemoveProjectMetadata(ctx context.Context, projectID, resourceOwner string, metadataKeys ...string) (_ *domain.ObjectDetails, err error) {
	if len(metadataKeys) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "META-5j9jYR", "Errors.Metadata.NoData")
	}

	_, err = c.checkProjectExists(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}

	events := make([]eventstore.Command, len(metadataKeys))
	removeMetadata, err := c.getProjectMetadataListModelByID(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}

	projectAgg := ProjectAggregateFromWriteModel(&removeMetadata.WriteModel)
	for i, key := range metadataKeys {
		if key == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-idzMSu", "Errors.Metadata.Invalid")
		}
		if _, found := removeMetadata.metadataList[key]; !found {
			return nil, zerrors.ThrowNotFound(nil, "META-aoxqwl", "Errors.Metadata.KeyNotExisting")
		}
		event, err := c.removeProjectMetadata(ctx, projectAgg, key)
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

func (c *Commands) removeProjectMetadata(ctx context.Context, projectAgg *eventstore.Aggregate, metadataKey string) (command eventstore.Command, err error) {
	command = project.NewMetadataRemovedEvent(
		ctx,
		projectAgg,
		metadataKey,
	)

	return command, nil
}

func (c *Commands) getProjectMetadataModelByID(ctx context.Context, projectID, resourceOwner, key string) (*ProjectMetadataWriteModel, error) {
	projectMetadataWriteModel := NewProjectMetadataWriteModel(projectID, resourceOwner, key)
	err := c.eventstore.FilterToQueryReducer(ctx, projectMetadataWriteModel)
	if err != nil {
		return nil, err
	}

	return projectMetadataWriteModel, nil
}

func (c *Commands) getProjectMetadataListModelByID(ctx context.Context, projectID, resourceOwner string) (*ProjectMetadataListWriteModel, error) {
	projectMetadataWriteModel := NewProjectMetadataListWriteModel(projectID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, projectMetadataWriteModel)
	if err != nil {
		return nil, err
	}

	return projectMetadataWriteModel, nil
}

func writeModelToProjectMetadata(wm *ProjectMetadataWriteModel) *domain.Metadata {
	return &domain.Metadata{
		ObjectRoot: writeModelToObjectRoot(wm.WriteModel),
		Key:        wm.Key,
		Value:      wm.Value,
		State:      wm.State,
	}
}
