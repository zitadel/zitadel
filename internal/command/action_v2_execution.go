package command

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	EventGroupSuffix = ".*"
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
	for _, target := range set.Targets {
		if err = target.Validate(); err != nil {
			return nil, err
		}
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
	for _, target := range set.Targets {
		if err = target.Validate(); err != nil {
			return nil, err
		}
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
	for _, target := range set.Targets {
		if err = target.Validate(); err != nil {
			return nil, err
		}
	}
	if err := cond.Existing(c); err != nil {
		return nil, err
	}
	for _, target := range set.Targets {
		if err = target.Validate(); err != nil {
			return nil, err
		}
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
		group := e.Group
		if !strings.HasSuffix(e.Group, EventGroupSuffix) {
			group += EventGroupSuffix
		}
		return execution.ID(domain.ExecutionTypeEvent, group)
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
	for _, target := range set.Targets {
		if err = target.Validate(); err != nil {
			return nil, err
		}
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

	Targets []*execution.Target
}

func (t SetExecution) GetIncludes() []string {
	includes := make([]string, 0)
	for i := range t.Targets {
		if t.Targets[i].Type == domain.ExecutionTargetTypeInclude {
			includes = append(includes, t.Targets[i].Target)
		}
	}
	return includes
}

func (t SetExecution) GetTargets() []string {
	targets := make([]string, 0)
	for i := range t.Targets {
		if t.Targets[i].Type == domain.ExecutionTargetTypeTarget {
			targets = append(targets, t.Targets[i].Target)
		}
	}
	return targets
}

func (e *SetExecution) Existing(c *Commands, ctx context.Context, resourceOwner string) error {
	targets := e.GetTargets()
	if len(targets) > 0 && !c.existsTargetsByIDs(ctx, targets, resourceOwner) {
		return zerrors.ThrowNotFound(nil, "COMMAND-17e8fq1ggk", "Errors.Target.NotFound")
	}
	includes := e.GetIncludes()
	if len(includes) > 0 && !c.existsExecutionsByIDs(ctx, includes, resourceOwner) {
		return zerrors.ThrowNotFound(nil, "COMMAND-slgj0l4cdz", "Errors.Execution.IncludeNotFound")
	}
	get, set := createIncludeCacheFunctions()
	// maxLevels could be configurable, but set as 3 for now
	return checkForIncludeCircular(ctx, e.AggregateID, resourceOwner, includes, c.getExecutionIncludes(get, set), 3)
}

func (c *Commands) setExecution(ctx context.Context, set *SetExecution, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	if resourceOwner == "" || set.AggregateID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-gg3a6ol4om", "Errors.IDMissing")
	}
	wm, err := c.getExecutionWriteModelByID(ctx, set.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	// Check if targets and includes for execution are existing
	if wm.ExecutionTargetsEqual(set.Targets) {
		return writeModelToObjectDetails(&wm.WriteModel), err
	}
	if err := set.Existing(c, ctx, resourceOwner); err != nil {
		return nil, err
	}
	if err := c.pushAppendAndReduce(ctx, wm, execution.NewSetEventV2(
		ctx,
		ExecutionAggregateFromWriteModel(&wm.WriteModel),
		set.Targets,
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

func createIncludeCacheFunctions() (func(s string) ([]string, bool), func(s string, strings []string)) {
	tempCache := make(map[string][]string)
	return func(s string) ([]string, bool) {
			include, ok := tempCache[s]
			return include, ok
		}, func(s string, strings []string) {
			tempCache[s] = strings
		}
}

type includeCacheFunc func(ctx context.Context, id string, resourceOwner string) ([]string, error)

func checkForIncludeCircular(ctx context.Context, id string, resourceOwner string, includes []string, cache includeCacheFunc, maxLevels int) error {
	if len(includes) == 0 {
		return nil
	}
	level := 0
	for _, include := range includes {
		if id == include {
			return zerrors.ThrowPreconditionFailed(nil, "COMMAND-mo1cmjp5k7", "Errors.Execution.CircularInclude")
		}
		if err := checkForIncludeCircularRecur(ctx, []string{id}, resourceOwner, include, cache, maxLevels, level); err != nil {
			return err
		}
	}
	return nil
}

func (c *Commands) getExecutionIncludes(
	getCache func(string) ([]string, bool),
	setCache func(string, []string),
) includeCacheFunc {
	return func(ctx context.Context, id string, resourceOwner string) ([]string, error) {
		included, ok := getCache(id)
		if !ok {
			included, err := c.getExecutionWriteModelByID(ctx, id, resourceOwner)
			if err != nil {
				return nil, err
			}
			includes := included.IncludeList()
			setCache(id, includes)
			return includes, nil
		}
		return included, nil
	}
}

func checkForIncludeCircularRecur(ctx context.Context, ids []string, resourceOwner string, include string, cache includeCacheFunc, maxLevels, level int) error {
	included, err := cache(ctx, include, resourceOwner)
	if err != nil {
		return err
	}
	currentLevel := level + 1
	if currentLevel >= maxLevels {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-gbhd3g57oo", "Errors.Execution.MaxLevelsInclude")
	}
	for _, includedInclude := range included {
		if include == includedInclude {
			return zerrors.ThrowPreconditionFailed(nil, "COMMAND-iuch02i656", "Errors.Execution.CircularInclude")
		}
		for _, id := range ids {
			if includedInclude == id {
				return zerrors.ThrowPreconditionFailed(nil, "COMMAND-819opvhgjv", "Errors.Execution.CircularInclude")
			}
		}
		if err := checkForIncludeCircularRecur(ctx, append(ids, include), resourceOwner, includedInclude, cache, maxLevels, currentLevel); err != nil {
			return err
		}
	}
	return nil
}
