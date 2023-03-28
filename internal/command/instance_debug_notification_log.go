package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddDebugNotificationProviderLog(ctx context.Context, fileSystemProvider *fs.Config) (*domain.ObjectDetails, error) {
	writeModel := NewInstanceDebugNotificationLogWriteModel(ctx)
	instanceAgg := InstanceAggregateFromWriteModel(&writeModel.WriteModel)
	events, err := c.addDefaultDebugNotificationLog(ctx, instanceAgg, writeModel, fileSystemProvider)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(writeModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.DebugNotificationWriteModel.WriteModel), nil
}

func (c *Commands) addDefaultDebugNotificationLog(ctx context.Context, instanceAgg *eventstore.Aggregate, addedWriteModel *InstanceDebugNotificationLogWriteModel, fileSystemProvider *fs.Config) ([]eventstore.Command, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, addedWriteModel)
	if err != nil {
		return nil, err
	}
	if addedWriteModel.State.Exists() {
		return nil, caos_errs.ThrowAlreadyExists(nil, "INSTANCE-3h0fs", "Errors.IAM.DebugNotificationProvider.AlreadyExists")
	}

	events := []eventstore.Command{
		instance.NewDebugNotificationProviderLogAddedEvent(ctx,
			instanceAgg,
			fileSystemProvider.Compact),
	}
	return events, nil
}

func (c *Commands) ChangeDefaultNotificationLog(ctx context.Context, fileSystemProvider *fs.Config) (*domain.ObjectDetails, error) {
	writeModel := NewInstanceDebugNotificationLogWriteModel(ctx)
	instanceAgg := InstanceAggregateFromWriteModel(&writeModel.WriteModel)
	event, err := c.changeDefaultDebugNotificationProviderLog(ctx, instanceAgg, writeModel, fileSystemProvider)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(writeModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.DebugNotificationWriteModel.WriteModel), nil
}

func (c *Commands) changeDefaultDebugNotificationProviderLog(ctx context.Context, instanceAgg *eventstore.Aggregate, existingProvider *InstanceDebugNotificationLogWriteModel, fileSystemProvider *fs.Config) (eventstore.Command, error) {
	err := c.defaultDebugNotificationProviderLogWriteModelByID(ctx, existingProvider)
	if err != nil {
		return nil, err
	}
	if !existingProvider.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-2h0s3", "Errors.IAM.DebugNotificationProvider.NotFound")
	}
	changedEvent, hasChanged := existingProvider.NewChangedEvent(ctx,
		instanceAgg,
		fileSystemProvider.Compact)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "INSTANCE-fn9p3", "Errors.IAM.LoginPolicy.NotChanged")
	}
	return changedEvent, nil
}

func (c *Commands) RemoveDefaultNotificationLog(ctx context.Context) (*domain.ObjectDetails, error) {
	existingProvider := NewInstanceDebugNotificationLogWriteModel(ctx)
	instanceAgg := InstanceAggregateFromWriteModel(&existingProvider.WriteModel)
	err := c.defaultDebugNotificationProviderLogWriteModelByID(ctx, existingProvider)
	if err != nil {
		return nil, err
	}
	if !existingProvider.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-39lse", "Errors.IAM.DebugNotificationProvider.NotFound")
	}

	events, err := c.eventstore.Push(ctx, instance.NewDebugNotificationProviderLogRemovedEvent(ctx, instanceAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingProvider, events...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingProvider.DebugNotificationWriteModel.WriteModel), nil
}

func (c *Commands) defaultDebugNotificationProviderLogWriteModelByID(ctx context.Context, writeModel *InstanceDebugNotificationLogWriteModel) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return err
	}
	return nil
}
