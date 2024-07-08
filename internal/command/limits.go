package command

import (
	"context"
	"errors"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/limits"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SetLimits struct {
	AuditLogRetention *time.Duration
	Block             *bool
}

// SetLimits creates new limits or updates existing limits.
func (c *Commands) SetLimits(
	ctx context.Context,
	setLimits *SetLimits,
) (*domain.ObjectDetails, error) {
	instanceId := authz.GetInstance(ctx).InstanceID()
	wm, err := c.getLimitsWriteModel(ctx, instanceId)
	if err != nil {
		return nil, err
	}
	cmds, err := c.setLimitsCommands(ctx, wm, setLimits)
	if err != nil {
		return nil, err
	}
	if len(cmds) > 0 {
		events, err := c.eventstore.Push(ctx, cmds...)
		if err != nil {
			return nil, err
		}
		err = AppendAndReduce(wm, events...)
		if err != nil {
			return nil, err
		}
	}
	return writeModelToObjectDetails(&wm.WriteModel), err
}

type SetInstanceLimitsBulk struct {
	InstanceID string
	SetLimits
}

func (c *Commands) SetInstanceLimitsBulk(
	ctx context.Context,
	bulk []*SetInstanceLimitsBulk,
) (bulkDetails *domain.ObjectDetails, targetsDetails []*domain.ObjectDetails, err error) {
	bulkWm, err := c.getBulkInstanceLimitsWriteModel(ctx, bulk)
	if err != nil {
		return nil, nil, err
	}
	cmds := make([]eventstore.Command, 0)
	for _, t := range bulk {
		targetWM, ok := bulkWm.writeModels[t.InstanceID]
		if !ok {
			return nil, nil, zerrors.ThrowInternal(nil, "COMMAND-5HWA9", "Errors.Limits.NotFound")
		}
		targetCMDs, setErr := c.setLimitsCommands(ctx, targetWM, &t.SetLimits)
		err = errors.Join(err, setErr)
		cmds = append(cmds, targetCMDs...)
	}
	if err != nil {
		return nil, nil, err
	}
	if len(cmds) > 0 {
		events, err := c.eventstore.Push(ctx, cmds...)
		if err != nil {
			return nil, nil, err
		}
		err = AppendAndReduce(bulkWm, events...)
		if err != nil {
			return nil, nil, err
		}
	}
	targetDetails := make([]*domain.ObjectDetails, len(bulk))
	for i, t := range bulk {
		targetDetails[i] = writeModelToObjectDetails(&bulkWm.writeModels[t.InstanceID].WriteModel)
	}
	details := writeModelToObjectDetails(&bulkWm.WriteModel)
	details.ResourceOwner = ""
	return details, targetDetails, err
}

func (c *Commands) setLimitsCommands(ctx context.Context, wm *limitsWriteModel, setLimits *SetLimits) (cmds []eventstore.Command, err error) {
	aggregateId := wm.AggregateID
	if aggregateId == "" {
		aggregateId, err = id_generator.Next()
		if err != nil {
			return nil, err
		}
	}
	aggregate := limits.NewAggregate(aggregateId, wm.InstanceID)
	createCmds, err := c.SetLimitsCommand(aggregate, wm, setLimits)()
	if err != nil {
		return nil, err
	}
	cmds, err = createCmds(ctx, nil)
	return cmds, err
}

func (c *Commands) ResetLimits(ctx context.Context) (*domain.ObjectDetails, error) {
	instanceId := authz.GetInstance(ctx).InstanceID()
	wm, err := c.getLimitsWriteModel(ctx, instanceId)
	if err != nil {
		return nil, err
	}
	if wm.AggregateID == "" {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-9JToT", "Errors.Limits.NotFound")
	}
	aggregate := limits.NewAggregate(wm.AggregateID, instanceId)
	events := []eventstore.Command{limits.NewResetEvent(ctx, &aggregate.Aggregate)}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(wm, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func (c *Commands) getLimitsWriteModel(ctx context.Context, instanceId string) (*limitsWriteModel, error) {
	wm := newLimitsWriteModel(instanceId)
	return wm, c.eventstore.FilterToQueryReducer(ctx, wm)
}

func (c *Commands) getBulkInstanceLimitsWriteModel(ctx context.Context, target []*SetInstanceLimitsBulk) (*limitsBulkWriteModel, error) {
	wm := newLimitsBulkWriteModel()
	for _, t := range target {
		wm.addWriteModel(t.InstanceID)
	}
	return wm, c.eventstore.FilterToQueryReducer(ctx, wm)
}

func (c *Commands) SetLimitsCommand(a *limits.Aggregate, wm *limitsWriteModel, setLimits *SetLimits) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if setLimits == nil || (setLimits.AuditLogRetention == nil && setLimits.Block == nil) {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-4M9vs", "Errors.Limits.NoneSpecified")
		}
		return func(ctx context.Context, _ preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			changes := wm.NewChanges(setLimits)
			if len(changes) == 0 {
				return nil, nil
			}
			return []eventstore.Command{limits.NewSetEvent(
				eventstore.NewBaseEventForPush(
					ctx,
					&a.Aggregate,
					limits.SetEventType,
				),
				changes...,
			)}, nil
		}, nil
	}
}
