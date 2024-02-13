package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Execution struct {
	models.ObjectRoot

	Name             string
	ExecutionType    domain.ExecutionType
	URL              string
	Timeout          time.Duration
	Async            bool
	InterruptOnError bool
}

func (a *Execution) IsValid() bool {
	return a.Name != ""
}

func (c *Commands) AddExecution(ctx context.Context, add *Execution, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	if !add.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-ddqbm9us5p", "Errors.Execution.Invalid")
	}

	if add.AggregateID == "" {
		add.AggregateID, err = c.idGenerator.Next()
		if err != nil {
			return nil, err
		}
	}

	executionModel := NewExecutionWriteModel(add.AggregateID, resourceOwner)
	pushedEvents, err := c.eventstore.Push(ctx, execution.NewAddedEvent(
		ctx,
		ExecutionAggregateFromWriteModel(&executionModel.WriteModel),
		add.Name,
		add.ExecutionType,
		add.URL,
		add.Timeout,
		add.Async,
		add.InterruptOnError,
	))
	if err != nil {
		return nil, err
	}
	if err := AppendAndReduce(executionModel, pushedEvents...); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&executionModel.WriteModel), nil
}

func (c *Commands) ChangeExecution(ctx context.Context, change *Execution, resourceOwner string) (*domain.ObjectDetails, error) {
	if change.AggregateID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-1l6ympeagp", "Errors.Execution.Invalid")
	}

	existing, err := c.getExecutionWriteModelByID(ctx, change.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !existing.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-xj14f2cccn", "Errors.Execution.NotFound")
	}

	changedEvent := existing.NewChangedEvent(
		ctx,
		ExecutionAggregateFromWriteModel(&existing.WriteModel),
		change.ExecutionType,
		change.URL,
		change.Timeout,
		change.Async,
		change.InterruptOnError)
	if changedEvent == nil {
		return writeModelToObjectDetails(&existing.WriteModel), nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existing, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existing.WriteModel), nil
}

func (c *Commands) DeleteExecution(ctx context.Context, id, resourceOwner string) (*domain.ObjectDetails, error) {
	if id == "" || resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-obqos2l3no", "Errors.IDMissing")
	}

	existing, err := c.getExecutionWriteModelByID(ctx, id, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !existing.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-k4s7ucu0ax", "Errors.Execution.NotFound")
	}

	if err := c.pushAppendAndReduce(ctx,
		existing,
		execution.NewRemovedEvent(ctx,
			ExecutionAggregateFromWriteModel(&existing.WriteModel),
			existing.Name,
		),
	); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existing.WriteModel), nil
}

func (c *Commands) getExecutionWriteModelByID(ctx context.Context, id string, resourceOwner string) (*ExecutionWriteModel, error) {
	wm := NewExecutionWriteModel(id, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, wm)
	if err != nil {
		return nil, err
	}
	return wm, nil
}
