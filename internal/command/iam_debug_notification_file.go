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

func (c *Commands) AddDebugNotificationProviderFile(ctx context.Context, fileSystemProvider *fs.FSConfig) (*domain.ObjectDetails, error) {
	writeModel := NewIAMDebugNotificationFileWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&writeModel.WriteModel)
	events, err := c.addDefaultDebugNotificationFile(ctx, iamAgg, writeModel, fileSystemProvider)
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

func (c *Commands) addDefaultDebugNotificationFile(ctx context.Context, iamAgg *eventstore.Aggregate, addedWriteModel *IAMDebugNotificationFileWriteModel, fileSystemProvider *fs.FSConfig) ([]eventstore.Command, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, addedWriteModel)
	if err != nil {
		return nil, err
	}
	if addedWriteModel.State.Exists() {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-d93nfs", "Errors.IAM.DebugNotificationProvider.AlreadyExists")
	}

	events := []eventstore.Command{
		iam_repo.NewDebugNotificationProviderFileAddedEvent(ctx,
			iamAgg,
			fileSystemProvider.Compact),
	}
	return events, nil
}

func (c *Commands) ChangeDefaultNotificationFile(ctx context.Context, fileSystemProvider *fs.FSConfig) (*domain.ObjectDetails, error) {
	writeModel := NewIAMDebugNotificationFileWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&writeModel.WriteModel)
	events, err := c.changeDefaultDebugNotificationProviderFile(ctx, iamAgg, writeModel, fileSystemProvider)
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

func (c *Commands) changeDefaultDebugNotificationProviderFile(ctx context.Context, iamAgg *eventstore.Aggregate, existingProvider *IAMDebugNotificationFileWriteModel, fileSystemProvider *fs.FSConfig) ([]eventstore.Command, error) {
	err := c.defaultDebugNotificationProviderFileWriteModelByID(ctx, existingProvider)
	if err != nil {
		return nil, err
	}
	if !existingProvider.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-fm9wl", "Errors.IAM.DebugNotificationProvider.NotFound")
	}
	events := make([]eventstore.Command, 0)
	changedEvent, hasChanged := existingProvider.NewChangedEvent(ctx,
		iamAgg,
		fileSystemProvider.Compact)
	if hasChanged {
		events = append(events, changedEvent)
	}
	if len(events) == 0 {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-5M9vdd", "Errors.IAM.LoginPolicy.NotChanged")

	}
	return events, nil
}

func (c *Commands) RemoveDefaultNotificationFile(ctx context.Context) (*domain.ObjectDetails, error) {
	existingProvider := NewIAMDebugNotificationFileWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&existingProvider.WriteModel)
	err := c.defaultDebugNotificationProviderFileWriteModelByID(ctx, existingProvider)
	if err != nil {
		return nil, err
	}
	if !existingProvider.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-dj9ew", "Errors.IAM.DebugNotificationProvider.NotFound")
	}

	events, err := c.eventstore.Push(ctx, iam_repo.NewDebugNotificationProviderFileRemovedEvent(ctx, iamAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingProvider, events...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingProvider.DebugNotificationWriteModel.WriteModel), nil
}

func (c *Commands) defaultDebugNotificationProviderFileWriteModelByID(ctx context.Context, writeModel *IAMDebugNotificationFileWriteModel) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return err
	}
	return nil
}
