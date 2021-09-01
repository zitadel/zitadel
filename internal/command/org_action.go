package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/action"
)

func (c *Commands) AddAction(ctx context.Context, addAction *domain.Action, resourceOwner string) (_ string, _ *domain.ObjectDetails, err error) {
	if !addAction.IsValid() {
		return "", nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-eg2gf", "Errors.Action.Invalid") //TODO: i18n
	}
	addAction.AggregateID, err = c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	actionModel := NewActionWriteModel(addAction.AggregateID, resourceOwner)
	actionAgg := ActionAggregateFromWriteModel(&actionModel.WriteModel)

	pushedEvents, err := c.eventstore.PushEvents(ctx, action.NewAddedEvent(
		ctx,
		actionAgg,
		addAction.Name,
		addAction.Script,
		addAction.Timeout,
		addAction.AllowedToFail,
	))
	if err != nil {
		return "", nil, err
	}
	err = AppendAndReduce(actionModel, pushedEvents...)
	if err != nil {
		return "", nil, err
	}
	return actionModel.AggregateID, writeModelToObjectDetails(&actionModel.WriteModel), nil
}

func (c *Commands) ChangeAction(ctx context.Context, actionChange *domain.Action, resourceOwner string) (*domain.ObjectDetails, error) {
	if !actionChange.IsValid() || actionChange.AggregateID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-Df2f3", "Errors.Action.Invalid") //TODO: i18n
	}

	existingAction, err := c.getActionWriteModelByID(ctx, actionChange.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingAction.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-Sfg2t", "Errors.Action.NotFound") //TODO: i18n
	}

	actionAgg := ActionAggregateFromWriteModel(&existingAction.WriteModel)
	changedEvent, err := existingAction.NewChangedEvent(
		ctx,
		actionAgg,
		actionChange.Name,
		actionChange.Script,
		actionChange.Timeout,
		actionChange.AllowedToFail)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingAction, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingAction.WriteModel), nil
}

func (c *Commands) DeleteAction(ctx context.Context, actionID string, resourceOwner string) (*domain.ObjectDetails, error) {
	if actionID == "" || resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-Gfg3g", "Errors.Action.ActionIDMissing")
	}

	existingAction, err := c.getActionWriteModelByID(ctx, actionID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingAction.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-Dgh4h", "Errors.Action.NotFound")
	}
	actionAgg := ActionAggregateFromWriteModel(&existingAction.WriteModel)
	events := []eventstore.EventPusher{
		action.NewRemovedEvent(ctx, actionAgg, existingAction.Name),
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingAction, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingAction.WriteModel), nil
}

func (c *Commands) DeleteFlow(ctx context.Context, flowType domain.FlowType, resourceOwner string) (*domain.ObjectDetails, error) {
	if !flowType.Valid() || resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-Dfw2h", "Errors.Action.FlowTypeMissing")
	}
	//TODO: !
	//existingFlow, err := c.getFlowWriteModelByID(ctx, flowType, resourceOwner)
	//if err != nil {
	//	return nil, err
	//}
	//if !existingFlow.State.Exists() {
	//	return nil, caos_errs.ThrowNotFound(nil, "COMMAND-Dgh4h", "Errors.Flow.NotFound")
	//}
	//actionAgg := ActionAggregateFromWriteModel(&existingFlow.WriteModel)
	//events := []eventstore.EventPusher{
	//	action.NewFlowRemovedEvent(ctx, actionAgg, flowType),
	//}
	//pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	//if err != nil {
	//	return nil, err
	//}
	//err = AppendAndReduce(existingFlow, pushedEvents...)
	//if err != nil {
	//	return nil, err
	//}
	//return writeModelToObjectDetails(&existingFlow.WriteModel), nil
	return nil, nil
}

func (c *Commands) getActionWriteModelByID(ctx context.Context, actionID string, resourceOwner string) (*ActionWriteModel, error) {
	actionWriteModel := NewActionWriteModel(actionID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, actionWriteModel)
	if err != nil {
		return nil, err
	}
	return actionWriteModel, nil
}

//TODO: !
//func (c *Commands) getFlowWriteModelByID(ctx context.Context, flowType string, resourceOwner string) (*FlowWriteModel, error) {
//	flowWriteModel := NewFlowWriteModel(flowType, resourceOwner)
//	err := c.eventstore.FilterToQueryReducer(ctx, flowWriteModel)
//	if err != nil {
//		return nil, err
//	}
//	return flowWriteModel, nil
//}
