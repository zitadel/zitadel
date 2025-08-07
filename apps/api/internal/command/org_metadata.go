package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) SetOrgMetadata(ctx context.Context, orgID string, metadata *domain.Metadata) (_ *domain.Metadata, err error) {
	err = c.checkOrgExists(ctx, orgID)
	if err != nil {
		return nil, err
	}
	setMetadata := NewOrgMetadataWriteModel(orgID, metadata.Key)
	orgAgg := OrgAggregateFromWriteModel(&setMetadata.WriteModel)
	event, err := c.setOrgMetadata(ctx, orgAgg, metadata)
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
	return writeModelToOrgMetadata(setMetadata), nil
}

func (c *Commands) BulkSetOrgMetadata(ctx context.Context, orgID string, metadatas ...*domain.Metadata) (_ *domain.ObjectDetails, err error) {
	if len(metadatas) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "META-9mm2d", "Errors.Metadata.NoData")
	}
	err = c.checkOrgExists(ctx, orgID)
	if err != nil {
		return nil, err
	}

	events := make([]eventstore.Command, len(metadatas))
	setMetadata := NewOrgMetadataListWriteModel(orgID)
	orgAgg := OrgAggregateFromWriteModel(&setMetadata.WriteModel)
	for i, data := range metadatas {
		event, err := c.setOrgMetadata(ctx, orgAgg, data)
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

func (c *Commands) setOrgMetadata(ctx context.Context, orgAgg *eventstore.Aggregate, metadata *domain.Metadata) (command eventstore.Command, err error) {
	if !metadata.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "META-2ml0f", "Errors.Metadata.Invalid")
	}
	return org.NewMetadataSetEvent(
		ctx,
		orgAgg,
		metadata.Key,
		metadata.Value,
	), nil
}

func (c *Commands) RemoveOrgMetadata(ctx context.Context, orgID, metadataKey string) (_ *domain.ObjectDetails, err error) {
	if metadataKey == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "META-2n0f1", "Errors.Metadata.Invalid")
	}
	err = c.checkOrgExists(ctx, orgID)
	if err != nil {
		return nil, err
	}
	removeMetadata, err := c.getOrgMetadataModelByID(ctx, orgID, metadataKey)
	if err != nil {
		return nil, err
	}
	if !removeMetadata.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "META-mcnw3", "Errors.Metadata.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&removeMetadata.WriteModel)
	event, err := c.removeOrgMetadata(ctx, orgAgg, metadataKey)
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

func (c *Commands) BulkRemoveOrgMetadata(ctx context.Context, orgID string, metadataKeys ...string) (_ *domain.ObjectDetails, err error) {
	if len(metadataKeys) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "META-9mw2d", "Errors.Metadata.NoData")
	}
	err = c.checkOrgExists(ctx, orgID)
	if err != nil {
		return nil, err
	}

	events := make([]eventstore.Command, len(metadataKeys))
	removeMetadata, err := c.getOrgMetadataListModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	orgAgg := OrgAggregateFromWriteModel(&removeMetadata.WriteModel)
	for i, key := range metadataKeys {
		if key == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-m19ds", "Errors.Metadata.Invalid")
		}
		if _, found := removeMetadata.metadataList[key]; !found {
			return nil, zerrors.ThrowNotFound(nil, "META-2npds", "Errors.Metadata.KeyNotExisting")
		}
		event, err := c.removeOrgMetadata(ctx, orgAgg, key)
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

func (c *Commands) removeOrgMetadata(ctx context.Context, orgAgg *eventstore.Aggregate, metadataKey string) (command eventstore.Command, err error) {
	command = org.NewMetadataRemovedEvent(
		ctx,
		orgAgg,
		metadataKey,
	)
	return command, nil
}

func (c *Commands) getOrgMetadataModelByID(ctx context.Context, orgID, key string) (*OrgMetadataWriteModel, error) {
	orgMetadataWriteModel := NewOrgMetadataWriteModel(orgID, key)
	err := c.eventstore.FilterToQueryReducer(ctx, orgMetadataWriteModel)
	if err != nil {
		return nil, err
	}
	return orgMetadataWriteModel, nil
}

func (c *Commands) getOrgMetadataListModelByID(ctx context.Context, orgID string) (*OrgMetadataListWriteModel, error) {
	orgMetadataWriteModel := NewOrgMetadataListWriteModel(orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, orgMetadataWriteModel)
	if err != nil {
		return nil, err
	}
	return orgMetadataWriteModel, nil
}

func writeModelToOrgMetadata(wm *OrgMetadataWriteModel) *domain.Metadata {
	return &domain.Metadata{
		ObjectRoot: writeModelToObjectRoot(wm.WriteModel),
		Key:        wm.Key,
		Value:      wm.Value,
		State:      wm.State,
	}
}
