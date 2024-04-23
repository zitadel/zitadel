package actions

import (
	"context"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/resouces/action/v2"
)

func (s *Server) CreateTarget(ctx context.Context, request *action.CreateTargetRequest) (*action.CreateTargetResponse, error) {
	return nil, zerrors.ThrowUnimplementedf(nil, "ACT-NtK93", "Not implemented. Got: %v", request)
}

func (s *Server) UpdateTarget(ctx context.Context, request *action.UpdateTargetRequest) (*action.UpdateTargetResponse, error) {
	return nil, zerrors.ThrowUnimplementedf(nil, "ACT-NtK93", "Not implemented. Got: %v", request)
}

func (s *Server) DeleteTarget(ctx context.Context, request *action.DeleteTargetRequest) (*action.DeleteTargetResponse, error) {
	return nil, zerrors.ThrowUnimplementedf(nil, "ACT-NtK93", "Not implemented. Got: %v", request)
}

func (s *Server) ListTargets(ctx context.Context, request *action.ListTargetsRequest) (*action.ListTargetsResponse, error) {
	return nil, zerrors.ThrowUnimplementedf(nil, "ACT-NtK93", "Not implemented. Got: %v", request)
}

func (s *Server) GetTargetByID(ctx context.Context, request *action.GetTargetByIDRequest) (*action.GetTargetByIDResponse, error) {
	return nil, zerrors.ThrowUnimplementedf(nil, "ACT-NtK93", "Not implemented. Got: %v", request)
}

func (s *Server) SetExecution(ctx context.Context, request *action.SetExecutionRequest) (*action.SetExecutionResponse, error) {
	return nil, zerrors.ThrowUnimplementedf(nil, "ACT-NtK93", "Not implemented. Got: %v", request)
}

func (s *Server) DeleteExecution(ctx context.Context, request *action.DeleteExecutionRequest) (*action.DeleteExecutionResponse, error) {
	return nil, zerrors.ThrowUnimplementedf(nil, "ACT-NtK93", "Not implemented. Got: %v", request)
}

func (s *Server) ListExecutions(ctx context.Context, request *action.ListExecutionsRequest) (*action.ListExecutionsResponse, error) {
	return nil, zerrors.ThrowUnimplementedf(nil, "ACT-NtK93", "Not implemented. Got: %v", request)
}

func (s *Server) ListExecutionFunctions(ctx context.Context, request *action.ListExecutionFunctionsRequest) (*action.ListExecutionFunctionsResponse, error) {
	return nil, zerrors.ThrowUnimplementedf(nil, "ACT-NtK93", "Not implemented. Got: %v", request)
}

func (s *Server) ListExecutionMethods(ctx context.Context, request *action.ListExecutionMethodsRequest) (*action.ListExecutionMethodsResponse, error) {
	return nil, zerrors.ThrowUnimplementedf(nil, "ACT-NtK93", "Not implemented. Got: %v", request)
}

func (s *Server) ListExecutionServices(ctx context.Context, request *action.ListExecutionServicesRequest) (*action.ListExecutionServicesResponse, error) {
	return nil, zerrors.ThrowUnimplementedf(nil, "ACT-NtK93", "Not implemented. Got: %v", request)
}
