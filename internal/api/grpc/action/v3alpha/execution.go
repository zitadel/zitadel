package action

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/execution"
	action "github.com/zitadel/zitadel/pkg/grpc/action/v3alpha"
)

func (s *Server) ListExecutionFunctions(_ context.Context, _ *action.ListExecutionFunctionsRequest) (*action.ListExecutionFunctionsResponse, error) {
	return &action.ListExecutionFunctionsResponse{
		Functions: s.ListActionFunctions(),
	}, nil
}

func (s *Server) ListExecutionMethods(_ context.Context, _ *action.ListExecutionMethodsRequest) (*action.ListExecutionMethodsResponse, error) {
	return &action.ListExecutionMethodsResponse{
		Methods: s.ListGRPCMethods(),
	}, nil
}

func (s *Server) ListExecutionServices(_ context.Context, _ *action.ListExecutionServicesRequest) (*action.ListExecutionServicesResponse, error) {
	return &action.ListExecutionServicesResponse{
		Services: s.ListGRPCServices(),
	}, nil
}

func (s *Server) SetExecution(ctx context.Context, req *action.SetExecutionRequest) (*action.SetExecutionResponse, error) {
	if err := checkExecutionEnabled(ctx); err != nil {
		return nil, err
	}

	targets := make([]*execution.Target, len(req.Targets))
	for i, target := range req.Targets {
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
	switch t := req.GetCondition().GetConditionType().(type) {
	case *action.Condition_Request:
		cond := executionConditionFromRequest(t.Request)
		details, err = s.command.SetExecutionRequest(ctx, cond, set, authz.GetInstance(ctx).InstanceID())
		if err != nil {
			return nil, err
		}
	case *action.Condition_Response:
		cond := executionConditionFromResponse(t.Response)
		details, err = s.command.SetExecutionResponse(ctx, cond, set, authz.GetInstance(ctx).InstanceID())
		if err != nil {
			return nil, err
		}
	case *action.Condition_Event:
		cond := executionConditionFromEvent(t.Event)
		details, err = s.command.SetExecutionEvent(ctx, cond, set, authz.GetInstance(ctx).InstanceID())
		if err != nil {
			return nil, err
		}
	case *action.Condition_Function:
		details, err = s.command.SetExecutionFunction(ctx, command.ExecutionFunctionCondition(t.Function.GetName()), set, authz.GetInstance(ctx).InstanceID())
		if err != nil {
			return nil, err
		}
	}
	return &action.SetExecutionResponse{
		Details: object.DomainToDetailsPb(details),
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
	}
	return "", nil
}

func (s *Server) DeleteExecution(ctx context.Context, req *action.DeleteExecutionRequest) (*action.DeleteExecutionResponse, error) {
	if err := checkExecutionEnabled(ctx); err != nil {
		return nil, err
	}

	var err error
	var details *domain.ObjectDetails
	switch t := req.GetCondition().GetConditionType().(type) {
	case *action.Condition_Request:
		cond := executionConditionFromRequest(t.Request)
		details, err = s.command.DeleteExecutionRequest(ctx, cond, authz.GetInstance(ctx).InstanceID())
		if err != nil {
			return nil, err
		}
	case *action.Condition_Response:
		cond := executionConditionFromResponse(t.Response)
		details, err = s.command.DeleteExecutionResponse(ctx, cond, authz.GetInstance(ctx).InstanceID())
		if err != nil {
			return nil, err
		}
	case *action.Condition_Event:
		cond := executionConditionFromEvent(t.Event)
		details, err = s.command.DeleteExecutionEvent(ctx, cond, authz.GetInstance(ctx).InstanceID())
		if err != nil {
			return nil, err
		}
	case *action.Condition_Function:
		details, err = s.command.DeleteExecutionFunction(ctx, command.ExecutionFunctionCondition(t.Function.GetName()), authz.GetInstance(ctx).InstanceID())
		if err != nil {
			return nil, err
		}
	}
	return &action.DeleteExecutionResponse{
		Details: object.DomainToDetailsPb(details),
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
