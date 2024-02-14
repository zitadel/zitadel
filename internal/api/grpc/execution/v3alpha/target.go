package execution

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	execution "github.com/zitadel/zitadel/pkg/grpc/execution/v3alpha"
)

func (s *Server) CreateTarget(ctx context.Context, req *execution.CreateTargetRequest) (*execution.CreateTargetResponse, error) {
	details, err := s.command.AddExecution(ctx, createTargetToCommand(req), "") // TODO: RO?
	if err != nil {
		return nil, err
	}
	return &execution.CreateTargetResponse{
		Id:      "", // TODO: id?
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) UpdateTarget(ctx context.Context, req *execution.UpdateTargetRequest) (*execution.UpdateTargetResponse, error) {
	details, err := s.command.ChangeExecution(ctx, updateTargetToCommand(req), "") // TODO: RO?
	if err != nil {
		return nil, err
	}
	return &execution.UpdateTargetResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) DeleteTarget(ctx context.Context, req *execution.DeleteTargetRequest) (*execution.DeleteTargetResponse, error) {
	details, err := s.command.DeleteExecution(ctx, req.GetTargetId(), "") // TODO: RO?
	if err != nil {
		return nil, err
	}
	return &execution.DeleteTargetResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func createTargetToCommand(req *execution.CreateTargetRequest) *command.Execution {
	return &command.Execution{
		Name:             req.GetName(),
		ExecutionType:    executionTypeToDomain(req.GetType()),
		URL:              req.GetUrl(),
		Timeout:          req.GetTimeout().AsDuration(),
		Async:            req.GetIsAsync(),
		InterruptOnError: req.GetInterruptOnError(),
	}
}

func updateTargetToCommand(req *execution.UpdateTargetRequest) *command.Execution {
	return &command.Execution{
		Name:             req.GetName(),
		ExecutionType:    executionTypeToDomain(req.GetType()),
		URL:              req.GetUrl(),
		Timeout:          req.GetTimeout().AsDuration(),
		Async:            req.GetIsAsync(),
		InterruptOnError: req.GetInterruptOnError(),
	}
}

func executionTypeToDomain(executionType execution.TargetType) domain.ExecutionType {
	switch executionType {
	case execution.TargetType_TARGET_TYPE_UNSPECIFIED:
		return domain.ExecutionTypeUndefined
	case execution.TargetType_TARGET_TYPE_REST_WEBHOOK:
		return domain.ExecutionTypeWebhook
	case execution.TargetType_TARGET_TYPE_REST_REQUEST_RESPONSE:
		return domain.ExecutionTypeRequestResponse
	default:
		return domain.ExecutionTypeUndefined
	}
}
