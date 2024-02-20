package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type ExecutionAPICondition struct {
	Method  string
	Service string
	All     bool
}

func (e *ExecutionAPICondition) IsValid() error {
	if e.Method == "" && e.Service == "" && !e.All {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-3tkej630e6", "Errors.Execution.Invalid")
	}
	// never set two conditions
	if e.Method != "" && (e.Service != "" || e.All) ||
		e.Service != "" && (e.Method != "" || e.All) ||
		e.All && (e.Method != "" || e.Service != "") {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-nee5q8aszq", "Errors.Execution.Invalid")
	}
	return nil
}

func (e *ExecutionAPICondition) ID() string {
	if e.Method != "" {
		return execution.IDFromGRPC(e.Method)
	}
	if e.Service != "" {
		return execution.IDFromGRPC(e.Service)
	}
	if e.All {
		return execution.IDFromGRPCAll()
	}
	return ""
}

func (c *Commands) SetExecutionRequest(ctx context.Context, cond *ExecutionAPICondition, set *SetExecution, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	if err := cond.IsValid(); err != nil {
		return nil, err
	}
	if cond.Method != "" && !c.grpcMethodExisting(cond.Method) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-vysplsevt8", "Errors.Execution.ConditionInvalid")
	}
	if cond.Service != "" && !c.grpcServiceExisting(cond.Service) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-qu6dfhiioq", "Errors.Execution.ConditionInvalid")
	}

	if set.AggregateID == "" {
		set.AggregateID = cond.ID()
	}
	set.ExecutionType = domain.ExecutionTypeRequest
	return c.setExecution(ctx, set, resourceOwner)
}

func (c *Commands) SetExecutionResponse(ctx context.Context, cond *ExecutionAPICondition, set *SetExecution, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	if set.AggregateID == "" {
		set.AggregateID = cond.ID()
	}
	set.ExecutionType = domain.ExecutionTypeResponse
	return c.setExecution(ctx, set, resourceOwner)
}

func (c *Commands) SetExecutionFunction(ctx context.Context, function string, set *SetExecution, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	if set.AggregateID == "" {
		set.AggregateID = execution.IDFromFunction(function)
	}
	set.ExecutionType = domain.ExecutionTypeFunction
	return c.setExecution(ctx, set, resourceOwner)
}

type ExecutionEventCondition struct {
	Event string
	Group string
	All   bool
}

func (e *ExecutionEventCondition) IsValid() error {
	if e.Event == "" && e.Group == "" && !e.All {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-w5smb6v7qu", "Errors.Execution.Invalid")
	}
	// never set two conditions
	if e.Event != "" && (e.Group != "" || e.All) ||
		e.Group != "" && (e.Event != "" || e.All) ||
		e.All && (e.Event != "" || e.Group != "") {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-hdm4zl1hmd", "Errors.Execution.Invalid")
	}
	return nil
}

func (e *ExecutionEventCondition) ID() string {
	if e.Event != "" {
		return execution.IDFromEvent(e.Event)
	}
	if e.Group != "" {
		return execution.IDFromEvent(e.Group)
	}
	if e.All {
		return execution.IDFromEventAll()
	}
	return ""
}

func (c *Commands) SetExecutionEvent(ctx context.Context, cond *ExecutionEventCondition, set *SetExecution, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	if set.AggregateID == "" {
		set.AggregateID = cond.ID()
	}
	set.ExecutionType = domain.ExecutionTypeEvent
	return c.setExecution(ctx, set, resourceOwner)
}

type SetExecution struct {
	models.ObjectRoot

	ExecutionType domain.ExecutionType
	Targets       []string
	Includes      []string
}

func (e *SetExecution) IsValid() error {
	if !e.ExecutionType.Valid() {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-wf8juv9lut", "Errors.Execution.Invalid")
	}
	if len(e.Targets) == 0 && len(e.Includes) == 0 {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-56bteot2uj", "Errors.Execution.NoTargets")
	}
	if len(e.Targets) > 0 && len(e.Includes) > 0 {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-5zleae34r1", "Errors.Execution.Invalid")
	}
	return nil
}

func (c *Commands) setExecution(ctx context.Context, set *SetExecution, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	if resourceOwner == "" || set.AggregateID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-gg3a6ol4om", "Errors.IDMissing")
	}
	if err := set.IsValid(); err != nil {
		return nil, err
	}

	wm, err := c.getExecutionWriteModelByIDAndType(ctx, set.AggregateID, resourceOwner, set.ExecutionType)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, execution.NewSetEvent(
		ctx,
		ExecutionAggregateFromWriteModel(&wm.WriteModel),
		wm.ExecutionType,
		set.Targets,
		set.Includes,
	))
	if err != nil {
		return nil, err
	}
	if err := AppendAndReduce(wm, pushedEvents...); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func (c *Commands) RemoveExecutionRequest(ctx context.Context, cond *ExecutionAPICondition, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	return c.removeExecution(ctx, cond.ID(), resourceOwner, domain.ExecutionTypeRequest)
}

func (c *Commands) RemoveExecutionResponse(ctx context.Context, cond *ExecutionAPICondition, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	return c.removeExecution(ctx, cond.ID(), resourceOwner, domain.ExecutionTypeResponse)
}

func (c *Commands) RemoveExecutionFunction(ctx context.Context, function string, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	return c.removeExecution(ctx, function, resourceOwner, domain.ExecutionTypeFunction)
}

func (c *Commands) RemoveExecutionEvent(ctx context.Context, cond *ExecutionEventCondition, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	return c.removeExecution(ctx, cond.ID(), resourceOwner, domain.ExecutionTypeEvent)
}

func (c *Commands) removeExecution(ctx context.Context, aggID string, resourceOwner string, executionType domain.ExecutionType) (_ *domain.ObjectDetails, err error) {
	if resourceOwner == "" || aggID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-cnic97c0g3", "Errors.IDMissing")
	}

	wm, err := c.getExecutionWriteModelByIDAndType(ctx, aggID, resourceOwner, executionType)
	if err != nil {
		return nil, err
	}
	if !wm.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-suq2upd3rt", "Errors.Execution.NotFound")
	}
	pushedEvents, err := c.eventstore.Push(ctx, execution.NewRemovedEvent(
		ctx,
		ExecutionAggregateFromWriteModel(&wm.WriteModel),
	))
	if err != nil {
		return nil, err
	}
	if err := AppendAndReduce(wm, pushedEvents...); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func (c *Commands) getExecutionWriteModelByIDAndType(ctx context.Context, id string, resourceOwner string, executionType domain.ExecutionType) (*ExecutionWriteModel, error) {
	wm := NewExecutionWriteModel(id, resourceOwner, executionType)
	err := c.eventstore.FilterToQueryReducer(ctx, wm)
	if err != nil {
		return nil, err
	}
	return wm, nil
}
