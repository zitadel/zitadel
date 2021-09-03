package command

import (
	"context"
	"reflect"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/org"
)

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

func (c *Commands) SetTriggerActions(ctx context.Context, flowType domain.FlowType, triggerType domain.TriggerType, actionIDs []string, resourceOwner string) (*domain.ObjectDetails, error) {
	if !flowType.Valid() || !triggerType.Valid() || resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-Dfhj5", "Errors.Flow.FlowTypeMissing")
	}
	if !flowType.HasTrigger(triggerType) {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-Dfgh6", "Errors.Flow.WrongTriggerType")
	}
	existingFlow, err := c.getOrgFlowWriteModelByType(ctx, flowType, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingFlow.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-Dgh4h", "Errors.Flow.NotFound")
	}
	if reflect.DeepEqual(existingFlow.Triggers[triggerType], actionIDs) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Nfh52", "Errors.Flow.NoChanges")
	}
	exists, err := c.actionsIDsExist(ctx, actionIDs, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-dg422", "Errors.Flow.ActionIDsNotExist")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingFlow.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org.NewTriggerActionsSetEvent(ctx, orgAgg, flowType, triggerType, actionIDs))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingFlow, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingFlow.WriteModel), nil
}

func (c *Commands) getOrgFlowWriteModelByType(ctx context.Context, flowType domain.FlowType, resourceOwner string) (*OrgFlowWriteModel, error) {
	flowWriteModel := NewOrgFlowWriteModel(flowType, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, flowWriteModel)
	if err != nil {
		return nil, err
	}
	return flowWriteModel, nil
}

func (c *Commands) actionsIDsExist(ctx context.Context, ids []string, resourceOwner string) (bool, error) {
	actionIDsModel := NewActionsExistModel(ids, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, actionIDsModel)
	return len(actionIDsModel.actionIDs) == len(actionIDsModel.checkedIDs), err
}
