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

func (c *Commands) DeactivateAction(ctx context.Context, actionID string, resourceOwner string) (*domain.ObjectDetails, error) {
	if actionID == "" || resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-DAhk5", "Errors.Action.ActionIDMissing")
	}

	existingAction, err := c.getActionWriteModelByID(ctx, actionID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingAction.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-NRmhu", "Errors.Action.NotFound")
	}
	if existingAction.State != domain.ActionStateActive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Dgj92", "Errors.Action.NotActive")
	}
	actionAgg := ActionAggregateFromWriteModel(&existingAction.WriteModel)
	events := []eventstore.EventPusher{
		action.NewDeactivatedEvent(ctx, actionAgg),
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

func (c *Commands) ReactivateAction(ctx context.Context, actionID string, resourceOwner string) (*domain.ObjectDetails, error) {
	if actionID == "" || resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-BNm56", "Errors.Action.ActionIDMissing")
	}

	existingAction, err := c.getActionWriteModelByID(ctx, actionID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingAction.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-Aa22g", "Errors.Action.NotFound")
	}
	if existingAction.State != domain.ActionStateInactive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-J53zh", "Errors.Action.NotInactive")
	}
	actionAgg := ActionAggregateFromWriteModel(&existingAction.WriteModel)
	events := []eventstore.EventPusher{
		action.NewReactivatedEvent(ctx, actionAgg),
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

func (c *Commands) removeActionsFromOrg(ctx context.Context, resourceOwner string) ([]eventstore.EventPusher, error) {
	existingActions, err := c.getActionsByOrgWriteModelByID(ctx, resourceOwner)
	if err != nil {
		return nil, err
	}
	if len(existingActions.Actions) == 0 {
		return nil, nil
	}
	events := make([]eventstore.EventPusher, 0, len(existingActions.Actions))
	for id, name := range existingActions.Actions {
		actionAgg := NewActionAggregate(id, resourceOwner)
		events = append(events, action.NewRemovedEvent(ctx, actionAgg, name))
	}
	return events, nil
}

func (c *Commands) getActionWriteModelByID(ctx context.Context, actionID string, resourceOwner string) (*ActionWriteModel, error) {
	actionWriteModel := NewActionWriteModel(actionID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, actionWriteModel)
	if err != nil {
		return nil, err
	}
	return actionWriteModel, nil
}

func (c *Commands) getActionsByOrgWriteModelByID(ctx context.Context, resourceOwner string) (*ActionsListByOrgModel, error) {
	actionWriteModel := NewActionsListByOrgModel(resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, actionWriteModel)
	if err != nil {
		return nil, err
	}
	return actionWriteModel, nil
}
