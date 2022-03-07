package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/notification/channels/fs"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddDebugNotificationProviderLog(ctx context.Context, fileSystemProvider *fs.FSConfig) (*domain.ObjectDetails, error) {
	writeModel := NewIAMDebugNotificationLogWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&writeModel.WriteModel)
	events, err := c.addDefaultDebugNotificationLog(ctx, iamAgg, writeModel, fileSystemProvider)
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

func (c *Commands) addDefaultDebugNotificationLog(ctx context.Context, iamAgg *eventstore.Aggregate, addedWriteModel *IAMDebugNotificationLogWriteModel, fileSystemProvider *fs.FSConfig) ([]eventstore.Command, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, addedWriteModel)
	if err != nil {
		return nil, err
	}
	if addedWriteModel.State.Exists() {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-3h0fs", "Errors.IAM.DebugNotificationProvider.AlreadyExists")
	}

	events := []eventstore.Command{
		iam_repo.NewDebugNotificationProviderLogAddedEvent(ctx,
			iamAgg,
			fileSystemProvider.Compact),
	}
	return events, nil
}

func (c *Commands) ChangeDefaultNotificationLog(ctx context.Context, fileSystemProvider *fs.FSConfig) (*domain.ObjectDetails, error) {
	writeModel := NewIAMDebugNotificationLogWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&writeModel.WriteModel)
	event, err := c.changeDefaultDebugNotificationProviderLog(ctx, iamAgg, writeModel, fileSystemProvider)
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

func (c *Commands) changeDefaultDebugNotificationProviderLog(ctx context.Context, iamAgg *eventstore.Aggregate, existingProvider *IAMDebugNotificationLogWriteModel, fileSystemProvider *fs.FSConfig) (eventstore.Command, error) {
	err := c.defaultDebugNotificationProviderLogWriteModelByID(ctx, existingProvider)
	if err != nil {
		return nil, err
	}
	if !existingProvider.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-2h0s3", "Errors.IAM.DebugNotificationProvider.NotFound")
	}
	changedEvent, hasChanged := existingProvider.NewChangedEvent(ctx,
		iamAgg,
		fileSystemProvider.Compact)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-fn9p3", "Errors.IAM.LoginPolicy.NotChanged")
	}
	return changedEvent, nil
}

func (c *Commands) RemoveDefaultNotificationLog(ctx context.Context) (*domain.ObjectDetails, error) {
	existingProvider := NewIAMDebugNotificationLogWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&existingProvider.WriteModel)
	err := c.defaultDebugNotificationProviderLogWriteModelByID(ctx, existingProvider)
	if err != nil {
		return nil, err
	}
	if !existingProvider.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-39lse", "Errors.IAM.DebugNotificationProvider.NotFound")
	}

	events, err := c.eventstore.Push(ctx, iam_repo.NewDebugNotificationProviderLogRemovedEvent(ctx, iamAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingProvider, events...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingProvider.DebugNotificationWriteModel.WriteModel), nil
}

func (c *Commands) defaultDebugNotificationProviderLogWriteModelByID(ctx context.Context, writeModel *IAMDebugNotificationLogWriteModel) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return err
	}
	return nil
}
