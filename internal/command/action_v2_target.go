package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/target"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddTarget struct {
	models.ObjectRoot

	Name             string
	ExecutionType    domain.TargetType
	URL              string
	Timeout          time.Duration
	Async            bool
	InterruptOnError bool
}

func (a *AddTarget) IsValid() bool {
	return a.Name != ""
}

func (c *Commands) AddTarget(ctx context.Context, add *AddTarget, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	if !add.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-ddqbm9us5p", "Errors.Target.Invalid")
	}

	if add.AggregateID == "" {
		add.AggregateID, err = c.idGenerator.Next()
		if err != nil {
			return nil, err
		}
	}

	wm := NewTargetWriteModel(add.AggregateID, resourceOwner)
	pushedEvents, err := c.eventstore.Push(ctx, target.NewAddedEvent(
		ctx,
		TargetAggregateFromWriteModel(&wm.WriteModel),
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
	if err := AppendAndReduce(wm, pushedEvents...); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

type ChangeTarget struct {
	models.ObjectRoot

	Name             *string
	ExecutionType    *domain.TargetType
	URL              *string
	Timeout          *time.Duration
	Async            *bool
	InterruptOnError *bool
}

func (c *Commands) ChangeTarget(ctx context.Context, change *ChangeTarget, resourceOwner string) (*domain.ObjectDetails, error) {
	if change.AggregateID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-1l6ympeagp", "Errors.Target.Invalid")
	}

	existing, err := c.getTargetWriteModelByID(ctx, change.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !existing.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-xj14f2cccn", "Errors.Target.NotFound")
	}

	changedEvent := existing.NewChangedEvent(
		ctx,
		TargetAggregateFromWriteModel(&existing.WriteModel),
		change.Name,
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

func (c *Commands) DeleteTarget(ctx context.Context, id, resourceOwner string) (*domain.ObjectDetails, error) {
	if id == "" || resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-obqos2l3no", "Errors.IDMissing")
	}

	existing, err := c.getTargetWriteModelByID(ctx, id, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !existing.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-k4s7ucu0ax", "Errors.Target.NotFound")
	}

	if err := c.pushAppendAndReduce(ctx,
		existing,
		target.NewRemovedEvent(ctx,
			TargetAggregateFromWriteModel(&existing.WriteModel),
			existing.Name,
		),
	); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existing.WriteModel), nil
}

func (c *Commands) getTargetWriteModelByID(ctx context.Context, id string, resourceOwner string) (*TargetWriteModel, error) {
	wm := NewTargetWriteModel(id, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, wm)
	if err != nil {
		return nil, err
	}
	return wm, nil
}
