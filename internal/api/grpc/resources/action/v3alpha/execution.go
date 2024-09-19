package action

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	resource_object "github.com/zitadel/zitadel/v2/internal/api/grpc/resources/object/v3alpha"
	"github.com/zitadel/zitadel/v2/internal/command"
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/repository/execution"
	"github.com/zitadel/zitadel/v2/internal/zerrors"
	object "github.com/zitadel/zitadel/v2/pkg/grpc/object/v3alpha"
	action "github.com/zitadel/zitadel/v2/pkg/grpc/resources/action/v3alpha"
)

func (s *Server) SetExecution(ctx context.Context, req *action.SetExecutionRequest) (*action.SetExecutionResponse, error) {
	if err := checkActionsEnabled(ctx); err != nil {
		return nil, err
	}
	reqTargets := req.GetExecution().GetTargets()
	targets := make([]*execution.Target, len(reqTargets))
	for i, target := range reqTargets {
		switch t := target.GetType().(type) {
		case *action.ExecutionTargetType_Include:
			include, err := conditionToInclude(t.Include)
			if err != nil {
				return nil, err
			}
			targets[i] = &execution.Target{Type: domain.ExecutionTargetTypeInclude, Target: include}
		case *action.ExecutionTargetType_Target:
			targets[i] = &execution.Target{Type: domain.ExecutionTargetTypeTarget, Target: t.Target}
		}
	}
	set := &command.SetExecution{
		Targets: targets,
	}
	var err error
	var details *domain.ObjectDetails
	instanceID := authz.GetInstance(ctx).InstanceID()
	switch t := req.GetCondition().GetConditionType().(type) {
	case *action.Condition_Request:
		cond := executionConditionFromRequest(t.Request)
		details, err = s.command.SetExecutionRequest(ctx, cond, set, instanceID)
	case *action.Condition_Response:
		cond := executionConditionFromResponse(t.Response)
		details, err = s.command.SetExecutionResponse(ctx, cond, set, instanceID)
	case *action.Condition_Event:
		cond := executionConditionFromEvent(t.Event)
		details, err = s.command.SetExecutionEvent(ctx, cond, set, instanceID)
	case *action.Condition_Function:
		details, err = s.command.SetExecutionFunction(ctx, command.ExecutionFunctionCondition(t.Function.GetName()), set, instanceID)
	default:
		err = zerrors.ThrowInvalidArgument(nil, "ACTION-5r5Ju", "Errors.Execution.ConditionInvalid")
	}
	if err != nil {
		return nil, err
	}
	return &action.SetExecutionResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_INSTANCE, instanceID),
	}, nil
}

func conditionToInclude(cond *action.Condition) (string, error) {
	switch t := cond.GetConditionType().(type) {
	case *action.Condition_Request:
		cond := executionConditionFromRequest(t.Request)
		if err := cond.IsValid(); err != nil {
			return "", err
		}
		return cond.ID(domain.ExecutionTypeRequest), nil
	case *action.Condition_Response:
		cond := executionConditionFromResponse(t.Response)
		if err := cond.IsValid(); err != nil {
			return "", err
		}
		return cond.ID(domain.ExecutionTypeRequest), nil
	case *action.Condition_Event:
		cond := executionConditionFromEvent(t.Event)
		if err := cond.IsValid(); err != nil {
			return "", err
		}
		return cond.ID(), nil
	case *action.Condition_Function:
		cond := command.ExecutionFunctionCondition(t.Function.GetName())
		if err := cond.IsValid(); err != nil {
			return "", err
		}
		return cond.ID(), nil
	default:
		return "", zerrors.ThrowInvalidArgument(nil, "ACTION-9BBob", "Errors.Execution.ConditionInvalid")
	}
}

func (s *Server) ListExecutionFunctions(ctx context.Context, _ *action.ListExecutionFunctionsRequest) (*action.ListExecutionFunctionsResponse, error) {
	if err := checkActionsEnabled(ctx); err != nil {
		return nil, err
	}
	return &action.ListExecutionFunctionsResponse{
		Functions: s.ListActionFunctions(),
	}, nil
}

func (s *Server) ListExecutionMethods(ctx context.Context, _ *action.ListExecutionMethodsRequest) (*action.ListExecutionMethodsResponse, error) {
	if err := checkActionsEnabled(ctx); err != nil {
		return nil, err
	}
	return &action.ListExecutionMethodsResponse{
		Methods: s.ListGRPCMethods(),
	}, nil
}

func (s *Server) ListExecutionServices(ctx context.Context, _ *action.ListExecutionServicesRequest) (*action.ListExecutionServicesResponse, error) {
	if err := checkActionsEnabled(ctx); err != nil {
		return nil, err
	}
	return &action.ListExecutionServicesResponse{
		Services: s.ListGRPCServices(),
	}, nil
}

func executionConditionFromRequest(request *action.RequestExecution) *command.ExecutionAPICondition {
	return &command.ExecutionAPICondition{
		Method:  request.GetMethod(),
		Service: request.GetService(),
		All:     request.GetAll(),
	}
}

func executionConditionFromResponse(response *action.ResponseExecution) *command.ExecutionAPICondition {
	return &command.ExecutionAPICondition{
		Method:  response.GetMethod(),
		Service: response.GetService(),
		All:     response.GetAll(),
	}
}

func executionConditionFromEvent(event *action.EventExecution) *command.ExecutionEventCondition {
	return &command.ExecutionEventCondition{
		Event: event.GetEvent(),
		Group: event.GetGroup(),
		All:   event.GetAll(),
	}
}
