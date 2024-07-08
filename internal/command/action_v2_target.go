package command

import (
	"context"
	"net/url"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/target"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddTarget struct {
	models.ObjectRoot

	Name             string
	TargetType       domain.TargetType
	Endpoint         string
	Timeout          time.Duration
	InterruptOnError bool
}

func (a *AddTarget) IsValid() error {
	if a.Name == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-ddqbm9us5p", "Errors.Target.Invalid")
	}
	if a.Timeout == 0 {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-39f35d8uri", "Errors.Target.NoTimeout")
	}
	_, err := url.Parse(a.Endpoint)
	if err != nil || a.Endpoint == "" {
		return zerrors.ThrowInvalidArgument(err, "COMMAND-1r2k6qo6wg", "Errors.Target.InvalidURL")
	}

	return nil
}

func (c *Commands) AddTarget(ctx context.Context, add *AddTarget, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-brml926e2d", "Errors.IDMissing")
	}

	if err := add.IsValid(); err != nil {
		return nil, err
	}

	if add.AggregateID == "" {
		add.AggregateID, err = id_generator.Next()
		if err != nil {
			return nil, err
		}
	}
	wm, err := c.getTargetWriteModelByID(ctx, add.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if wm.State.Exists() {
		return nil, zerrors.ThrowAlreadyExists(nil, "INSTANCE-9axkz0jvzm", "Errors.Target.AlreadyExists")
	}

	pushedEvents, err := c.eventstore.Push(ctx, target.NewAddedEvent(
		ctx,
		TargetAggregateFromWriteModel(&wm.WriteModel),
		add.Name,
		add.TargetType,
		add.Endpoint,
		add.Timeout,
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
	TargetType       *domain.TargetType
	Endpoint         *string
	Timeout          *time.Duration
	InterruptOnError *bool
}

func (a *ChangeTarget) IsValid() error {
	if a.AggregateID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-1l6ympeagp", "Errors.IDMissing")
	}
	if a.Name != nil && *a.Name == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-d1wx4lm0zr", "Errors.Target.Invalid")
	}
	if a.Timeout != nil && *a.Timeout == 0 {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-08b39vdi57", "Errors.Target.NoTimeout")
	}
	if a.Endpoint != nil {
		_, err := url.Parse(*a.Endpoint)
		if err != nil || *a.Endpoint == "" {
			return zerrors.ThrowInvalidArgument(err, "COMMAND-jsbaera7b6", "Errors.Target.InvalidURL")
		}
	}
	return nil
}

func (c *Commands) ChangeTarget(ctx context.Context, change *ChangeTarget, resourceOwner string) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-zqibgg0wwh", "Errors.IDMissing")
	}
	if err := change.IsValid(); err != nil {
		return nil, err
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
		change.TargetType,
		change.Endpoint,
		change.Timeout,
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

func (c *Commands) existsTargetsByIDs(ctx context.Context, ids []string, resourceOwner string) bool {
	wm := NewTargetsExistsWriteModel(ids, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, wm)
	if err != nil {
		return false
	}
	return wm.AllExists()
}

func (c *Commands) getTargetWriteModelByID(ctx context.Context, id string, resourceOwner string) (*TargetWriteModel, error) {
	wm := NewTargetWriteModel(id, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, wm)
	if err != nil {
		return nil, err
	}
	return wm, nil
}
