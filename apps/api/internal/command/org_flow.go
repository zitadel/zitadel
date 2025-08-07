package command

import (
	"context"
	"slices"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) ClearFlow(ctx context.Context, flowType domain.FlowType, resourceOwner string) (*domain.ObjectDetails, error) {
	if !flowType.Valid() || resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Dfw2h", "Errors.Flow.FlowTypeMissing")
	}
	existingFlow, err := c.getOrgFlowWriteModelByType(ctx, flowType, resourceOwner)
	if err != nil {
		return nil, err
	}
	if len(existingFlow.Triggers) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-DgGh3", "Errors.Flow.Empty")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingFlow.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewFlowClearedEvent(ctx, orgAgg, flowType))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingFlow, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingFlow.WriteModel), nil
}

func (c *Commands) SetTriggerActions(ctx context.Context, flowType domain.FlowType, triggerType domain.TriggerType, actionIDs []string, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if !flowType.Valid() || !triggerType.Valid() || resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Dfhj5", "Errors.Flow.FlowTypeMissing")
	}
	if !flowType.HasTrigger(triggerType) {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Dfgh6", "Errors.Flow.WrongTriggerType")
	}
	existingFlow, err := c.getOrgFlowWriteModelByType(ctx, flowType, resourceOwner)
	if err != nil {
		return nil, err
	}
	if slices.Equal(existingFlow.Triggers[triggerType], actionIDs) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Nfh52", "Errors.Flow.NoChanges")
	}
	if len(actionIDs) > 0 {
		exists, err := c.actionsIDsExist(ctx, actionIDs, resourceOwner)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-dg422", "Errors.Flow.ActionIDsNotExist")
		}
	}
	orgAgg := OrgAggregateFromWriteModel(&existingFlow.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewTriggerActionsSetEvent(ctx, orgAgg, flowType, triggerType, actionIDs))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingFlow, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingFlow.WriteModel), nil
}

func (c *Commands) getOrgFlowWriteModelByType(ctx context.Context, flowType domain.FlowType, resourceOwner string) (_ *OrgFlowWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	flowWriteModel := NewOrgFlowWriteModel(flowType, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, flowWriteModel)
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
