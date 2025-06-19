package action

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/zerrors"
	action "github.com/zitadel/zitadel/pkg/grpc/action/v2beta"
)

func (s *Server) SetExecution(ctx context.Context, req *action.SetExecutionRequest) (*action.SetExecutionResponse, error) {
	reqTargets := req.GetTargets()
	targets := make([]*execution.Target, len(reqTargets))
	for i, target := range reqTargets {
		targets[i] = &execution.Target{Type: domain.ExecutionTargetTypeTarget, Target: target}
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
		SetDate: timestamppb.New(details.EventDate),
	}, nil
}

func (s *Server) ListExecutionFunctions(ctx context.Context, _ *action.ListExecutionFunctionsRequest) (*action.ListExecutionFunctionsResponse, error) {
	return &action.ListExecutionFunctionsResponse{
		Functions: s.ListActionFunctions(),
	}, nil
}

func (s *Server) ListExecutionMethods(ctx context.Context, _ *action.ListExecutionMethodsRequest) (*action.ListExecutionMethodsResponse, error) {
	return &action.ListExecutionMethodsResponse{
		Methods: s.ListGRPCMethods(),
	}, nil
}

func (s *Server) ListExecutionServices(ctx context.Context, _ *action.ListExecutionServicesRequest) (*action.ListExecutionServicesResponse, error) {
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
