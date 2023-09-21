package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	iam_repo "github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddDebugNotificationProviderFile(ctx context.Context, fileSystemProvider *fs.Config) (*domain.ObjectDetails, error) {
	writeModel := NewInstanceDebugNotificationFileWriteModel(ctx)
	instanceAgg := InstanceAggregateFromWriteModel(&writeModel.WriteModel)
	events, err := c.addDefaultDebugNotificationFile(ctx, instanceAgg, writeModel, fileSystemProvider)
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

func (c *Commands) addDefaultDebugNotificationFile(ctx context.Context, instanceAgg *eventstore.Aggregate, addedWriteModel *InstanceDebugNotificationFileWriteModel, fileSystemProvider *fs.Config) ([]eventstore.Command, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, addedWriteModel)
	if err != nil {
		return nil, err
	}
	if addedWriteModel.State.Exists() {
		return nil, caos_errs.ThrowAlreadyExists(nil, "INSTANCE-d93nfs", "Errors.IAM.DebugNotificationProvider.AlreadyExists")
	}

	events := []eventstore.Command{
		iam_repo.NewDebugNotificationProviderFileAddedEvent(ctx,
			instanceAgg,
			fileSystemProvider.Compact),
	}
	return events, nil
}

func (c *Commands) ChangeDefaultNotificationFile(ctx context.Context, fileSystemProvider *fs.Config) (*domain.ObjectDetails, error) {
	writeModel := NewInstanceDebugNotificationFileWriteModel(ctx)
	instanceAgg := InstanceAggregateFromWriteModel(&writeModel.WriteModel)
	events, err := c.changeDefaultDebugNotificationProviderFile(ctx, instanceAgg, writeModel, fileSystemProvider)
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

func (c *Commands) changeDefaultDebugNotificationProviderFile(ctx context.Context, instanceAgg *eventstore.Aggregate, existingProvider *InstanceDebugNotificationFileWriteModel, fileSystemProvider *fs.Config) ([]eventstore.Command, error) {
	err := c.defaultDebugNotificationProviderFileWriteModelByID(ctx, existingProvider)
	if err != nil {
		return nil, err
	}
	if !existingProvider.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-fm9wl", "Errors.IAM.DebugNotificationProvider.NotFound")
	}
	events := make([]eventstore.Command, 0)
	changedEvent, hasChanged := existingProvider.NewChangedEvent(ctx,
		instanceAgg,
		fileSystemProvider.Compact)
	if hasChanged {
		events = append(events, changedEvent)
	}
	if len(events) == 0 {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "INSTANCE-5M9vdd", "Errors.IAM.LoginPolicy.NotChanged")

	}
	return events, nil
}

func (c *Commands) RemoveDefaultNotificationFile(ctx context.Context) (*domain.ObjectDetails, error) {
	existingProvider := NewInstanceDebugNotificationFileWriteModel(ctx)
	instanceAgg := InstanceAggregateFromWriteModel(&existingProvider.WriteModel)
	err := c.defaultDebugNotificationProviderFileWriteModelByID(ctx, existingProvider)
	if err != nil {
		return nil, err
	}
	if !existingProvider.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-dj9ew", "Errors.IAM.DebugNotificationProvider.NotFound")
	}

	events, err := c.eventstore.Push(ctx, iam_repo.NewDebugNotificationProviderFileRemovedEvent(ctx, instanceAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingProvider, events...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingProvider.DebugNotificationWriteModel.WriteModel), nil
}

func (c *Commands) defaultDebugNotificationProviderFileWriteModelByID(ctx context.Context, writeModel *InstanceDebugNotificationFileWriteModel) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return err
	}
	return nil
}
