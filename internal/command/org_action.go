package command

import (
	"context"
	"sort"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/action"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddActionWithID(ctx context.Context, addAction *domain.Action, resourceOwner, actionID string) (_ string, _ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	existingAction, err := c.getActionWriteModelByID(ctx, actionID, resourceOwner)
	if err != nil {
		return "", nil, err
	}
	if existingAction.State != domain.ActionStateUnspecified {
		return "", nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-nau2k", "Errors.Action.AlreadyExisting")
	}

	return c.addActionWithID(ctx, addAction, resourceOwner, actionID)
}

func (c *Commands) AddAction(ctx context.Context, addAction *domain.Action, resourceOwner string) (_ string, _ *domain.ObjectDetails, err error) {
	if !addAction.IsValid() {
		return "", nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-eg2gf", "Errors.Action.Invalid")
	}

	actionID, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}

	return c.addActionWithID(ctx, addAction, resourceOwner, actionID)
}

func (c *Commands) addActionWithID(ctx context.Context, addAction *domain.Action, resourceOwner, actionID string) (_ string, _ *domain.ObjectDetails, err error) {
	addAction.AggregateID = actionID
	actionModel := NewActionWriteModel(addAction.AggregateID, resourceOwner)
	actionAgg := ActionAggregateFromWriteModel(&actionModel.WriteModel)

	pushedEvents, err := c.eventstore.Push(ctx, action.NewAddedEvent(
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
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Df2f3", "Errors.Action.Invalid")
	}

	existingAction, err := c.getActionWriteModelByID(ctx, actionChange.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingAction.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-Sfg2t", "Errors.Action.NotFound")
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
	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
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
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-DAhk5", "Errors.IDMissing")
	}

	existingAction, err := c.getActionWriteModelByID(ctx, actionID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingAction.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-NRmhu", "Errors.Action.NotFound")
	}
	if existingAction.State != domain.ActionStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Dgj92", "Errors.Action.NotActive")
	}
	actionAgg := ActionAggregateFromWriteModel(&existingAction.WriteModel)
	events := []eventstore.Command{
		action.NewDeactivatedEvent(ctx, actionAgg),
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
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
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-BNm56", "Errors.IDMissing")
	}

	existingAction, err := c.getActionWriteModelByID(ctx, actionID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingAction.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-Aa22g", "Errors.Action.NotFound")
	}
	if existingAction.State != domain.ActionStateInactive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-J53zh", "Errors.Action.NotInactive")
	}

	actionAgg := ActionAggregateFromWriteModel(&existingAction.WriteModel)
	events := []eventstore.Command{
		action.NewReactivatedEvent(ctx, actionAgg),
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingAction, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingAction.WriteModel), nil
}

func (c *Commands) DeleteAction(ctx context.Context, actionID, resourceOwner string, flowTypes ...domain.FlowType) (*domain.ObjectDetails, error) {
	if actionID == "" || resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Gfg3g", "Errors.IDMissing")
	}

	existingAction, err := c.getActionWriteModelByID(ctx, actionID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingAction.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-Dgh4h", "Errors.Action.NotFound")
	}
	actionAgg := ActionAggregateFromWriteModel(&existingAction.WriteModel)
	events := []eventstore.Command{
		action.NewRemovedEvent(ctx, actionAgg, existingAction.Name),
	}
	orgAgg := org.NewAggregate(resourceOwner).Aggregate
	for _, flowType := range flowTypes {
		events = append(events, org.NewTriggerActionsCascadeRemovedEvent(ctx, &orgAgg, flowType, actionID))
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingAction, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingAction.WriteModel), nil
}

func (c *Commands) removeActionsFromOrg(ctx context.Context, resourceOwner string) ([]eventstore.Command, error) {
	existingActions, err := c.getActionsByOrgWriteModelByID(ctx, resourceOwner)
	if err != nil {
		return nil, err
	}
	if len(existingActions.Actions) == 0 {
		return nil, nil
	}
	events := make([]eventstore.Command, 0, len(existingActions.Actions))
	for id, existingAction := range existingActions.Actions {
		actionAgg := NewActionAggregate(id, resourceOwner)
		events = append(events, action.NewRemovedEvent(ctx, actionAgg, existingAction.Name))
	}
	return events, nil
}

func (c *Commands) deactivateNotAllowedActionsFromOrg(ctx context.Context, resourceOwner string, maxAllowed int) ([]eventstore.Command, error) {
	existingActions, err := c.getActionsByOrgWriteModelByID(ctx, resourceOwner)
	if err != nil {
		return nil, err
	}
	activeActions := make([]*ActionWriteModel, 0, len(existingActions.Actions))
	for _, existingAction := range existingActions.Actions {
		if existingAction.State == domain.ActionStateActive {
			activeActions = append(activeActions, existingAction)
		}
	}
	if len(activeActions) <= maxAllowed {
		return nil, nil
	}
	sort.Slice(activeActions, func(i, j int) bool {
		return activeActions[i].WriteModel.ChangeDate.Before(activeActions[j].WriteModel.ChangeDate)
	})
	events := make([]eventstore.Command, 0, len(existingActions.Actions))
	for i := maxAllowed; i < len(activeActions); i++ {
		actionAgg := NewActionAggregate(activeActions[i].AggregateID, resourceOwner)
		events = append(events, action.NewDeactivatedEvent(ctx, actionAgg))
	}
	return events, nil
}

func (c *Commands) getActionWriteModelByID(ctx context.Context, actionID string, resourceOwner string) (_ *ActionWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	actionWriteModel := NewActionWriteModel(actionID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, actionWriteModel)
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
