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

func (e *ExecutionAPICondition) ID(executionType domain.ExecutionType) string {
	if e.Method != "" {
		return execution.ID(executionType, e.Method)
	}
	if e.Service != "" {
		return execution.ID(executionType, e.Service)
	}
	if e.All {
		return execution.IDAll(executionType)
	}
	return ""
}

func (e *ExecutionAPICondition) Existing(c *Commands) error {
	if e.Method != "" && !c.GrpcMethodExisting(e.Method) {
		return zerrors.ThrowNotFound(nil, "COMMAND-vysplsevt8", "Errors.Execution.ConditionInvalid")
	}
	if e.Service != "" && !c.GrpcServiceExisting(e.Service) {
		return zerrors.ThrowNotFound(nil, "COMMAND-qu6dfhiioq", "Errors.Execution.ConditionInvalid")
	}
	return nil
}

func (c *Commands) SetExecutionRequest(ctx context.Context, cond *ExecutionAPICondition, set *SetExecution, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	if err := cond.IsValid(); err != nil {
		return nil, err
	}
	if err := cond.Existing(c); err != nil {
		return nil, err
	}
	if set.AggregateID == "" {
		set.AggregateID = cond.ID(domain.ExecutionTypeRequest)
	}
	return c.setExecution(ctx, set, resourceOwner)
}

func (c *Commands) SetExecutionResponse(ctx context.Context, cond *ExecutionAPICondition, set *SetExecution, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	if err := cond.IsValid(); err != nil {
		return nil, err
	}
	if err := cond.Existing(c); err != nil {
		return nil, err
	}
	if set.AggregateID == "" {
		set.AggregateID = cond.ID(domain.ExecutionTypeResponse)
	}
	return c.setExecution(ctx, set, resourceOwner)
}

type ExecutionFunctionCondition string

func (e ExecutionFunctionCondition) IsValid() error {
	if e == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-5folwn5jws", "Errors.Execution.Invalid")
	}
	return nil
}

func (e ExecutionFunctionCondition) ID() string {
	return execution.ID(domain.ExecutionTypeFunction, string(e))
}

func (e ExecutionFunctionCondition) Existing(c *Commands) error {
	if !c.ActionFunctionExisting(string(e)) {
		return zerrors.ThrowNotFound(nil, "COMMAND-cdy39t0ksr", "Errors.Execution.ConditionInvalid")
	}
	return nil
}

func (c *Commands) SetExecutionFunction(ctx context.Context, cond ExecutionFunctionCondition, set *SetExecution, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	if err := cond.IsValid(); err != nil {
		return nil, err
	}
	if err := cond.Existing(c); err != nil {
		return nil, err
	}
	if set.AggregateID == "" {
		set.AggregateID = cond.ID()
	}
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
		return execution.ID(domain.ExecutionTypeEvent, e.Event)
	}
	if e.Group != "" {
		return execution.ID(domain.ExecutionTypeEvent, e.Group)
	}
	if e.All {
		return execution.IDAll(domain.ExecutionTypeEvent)
	}
	return ""
}

func (e *ExecutionEventCondition) Existing(c *Commands) error {
	if e.Event != "" && !c.EventExisting(e.Event) {
		return zerrors.ThrowNotFound(nil, "COMMAND-74aaqj8fv9", "Errors.Execution.ConditionInvalid")
	}
	if e.Group != "" && !c.EventGroupExisting(e.Group) {
		return zerrors.ThrowNotFound(nil, "COMMAND-er5oneb5lz", "Errors.Execution.ConditionInvalid")
	}
	return nil
}

func (c *Commands) SetExecutionEvent(ctx context.Context, cond *ExecutionEventCondition, set *SetExecution, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	if err := cond.IsValid(); err != nil {
		return nil, err
	}
	if err := cond.Existing(c); err != nil {
		return nil, err
	}
	if set.AggregateID == "" {
		set.AggregateID = cond.ID()
	}
	return c.setExecution(ctx, set, resourceOwner)
}

type SetExecution struct {
	models.ObjectRoot

	Targets  []string
	Includes []string
}

func (e *SetExecution) IsValid() error {
	if len(e.Targets) == 0 && len(e.Includes) == 0 {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-56bteot2uj", "Errors.Execution.NoTargets")
	}
	if len(e.Targets) > 0 && len(e.Includes) > 0 {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-5zleae34r1", "Errors.Execution.Invalid")
	}
	return nil
}

func (e *SetExecution) Existing(c *Commands, ctx context.Context, resourceOwner string) error {
	if len(e.Targets) > 0 && !c.existsTargetsByIDs(ctx, e.Targets, resourceOwner) {
		return zerrors.ThrowNotFound(nil, "COMMAND-17e8fq1ggk", "Errors.Target.NotFound")
	}
	if len(e.Includes) > 0 && !c.existsExecutionsByIDs(ctx, e.Includes, resourceOwner) {
		return zerrors.ThrowNotFound(nil, "COMMAND-slgj0l4cdz", "Errors.Execution.IncludeNotFound")
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

	wm := NewExecutionWriteModel(set.AggregateID, resourceOwner)
	// Check if targets and includes for execution are existing
	if err := set.Existing(c, ctx, resourceOwner); err != nil {
		return nil, err
	}

	if err := c.pushAppendAndReduce(ctx, wm, execution.NewSetEvent(
		ctx,
		ExecutionAggregateFromWriteModel(&wm.WriteModel),
		set.Targets,
		set.Includes,
	)); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func (c *Commands) DeleteExecutionRequest(ctx context.Context, cond *ExecutionAPICondition, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	if err := cond.IsValid(); err != nil {
		return nil, err
	}
	return c.deleteExecution(ctx, cond.ID(domain.ExecutionTypeRequest), resourceOwner)
}

func (c *Commands) DeleteExecutionResponse(ctx context.Context, cond *ExecutionAPICondition, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	if err := cond.IsValid(); err != nil {
		return nil, err
	}
	return c.deleteExecution(ctx, cond.ID(domain.ExecutionTypeResponse), resourceOwner)
}

func (c *Commands) DeleteExecutionFunction(ctx context.Context, cond ExecutionFunctionCondition, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	if err := cond.IsValid(); err != nil {
		return nil, err
	}
	return c.deleteExecution(ctx, cond.ID(), resourceOwner)
}

func (c *Commands) DeleteExecutionEvent(ctx context.Context, cond *ExecutionEventCondition, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	if err := cond.IsValid(); err != nil {
		return nil, err
	}
	return c.deleteExecution(ctx, cond.ID(), resourceOwner)
}

func (c *Commands) deleteExecution(ctx context.Context, aggID string, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	if resourceOwner == "" || aggID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-cnic97c0g3", "Errors.IDMissing")
	}

	wm, err := c.getExecutionWriteModelByID(ctx, aggID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !wm.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-suq2upd3rt", "Errors.Execution.NotFound")
	}
	if err := c.pushAppendAndReduce(ctx, wm, execution.NewRemovedEvent(
		ctx,
		ExecutionAggregateFromWriteModel(&wm.WriteModel),
	)); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func (c *Commands) existsExecutionsByIDs(ctx context.Context, ids []string, resourceOwner string) bool {
	wm := NewExecutionsExistWriteModel(ids, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, wm)
	if err != nil {
		return false
	}
	return wm.AllExists()
}

func (c *Commands) getExecutionWriteModelByID(ctx context.Context, id string, resourceOwner string) (*ExecutionWriteModel, error) {
	wm := NewExecutionWriteModel(id, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, wm)
	if err != nil {
		return nil, err
	}
	return wm, nil
}
