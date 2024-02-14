package execution

import (
	"context"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	execution "github.com/zitadel/zitadel/pkg/grpc/execution/v3alpha"
)

func (s *Server) CreateTarget(ctx context.Context, req *execution.CreateTargetRequest) (*execution.CreateTargetResponse, error) {
	add := createTargetToCommand(req)
	details, err := s.command.AddTarget(ctx, add, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	return &execution.CreateTargetResponse{
		Id:      add.AggregateID,
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) UpdateTarget(ctx context.Context, req *execution.UpdateTargetRequest) (*execution.UpdateTargetResponse, error) {
	details, err := s.command.ChangeTarget(ctx, updateTargetToCommand(req), authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	return &execution.UpdateTargetResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) DeleteTarget(ctx context.Context, req *execution.DeleteTargetRequest) (*execution.DeleteTargetResponse, error) {
	details, err := s.command.DeleteTarget(ctx, req.GetTargetId(), authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	return &execution.DeleteTargetResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func createTargetToCommand(req *execution.CreateTargetRequest) *command.AddTarget {
	return &command.AddTarget{
		Name:             req.GetName(),
		TargetType:       targetTypeToDomain(req.GetType()),
		URL:              req.GetUrl(),
		Timeout:          req.GetTimeout().AsDuration(),
		Async:            req.GetIsAsync(),
		InterruptOnError: req.GetInterruptOnError(),
	}
}

func updateTargetToCommand(req *execution.UpdateTargetRequest) *command.ChangeTarget {
	if req == nil {
		return nil
	}
	target := &command.ChangeTarget{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.GetTargetId(),
		},
		Name: req.Name,
		URL:  req.Url,
	}
	if req.Type != nil {
		target.TargetType = gu.Ptr(targetTypeToDomain(req.GetType()))
	}
	if req.Timeout != nil {
		target.Timeout = gu.Ptr(req.GetTimeout().AsDuration())
	}
	if req.ExecutionType != nil {
		target.Async = gu.Ptr(req.GetIsAsync())
		target.InterruptOnError = gu.Ptr(req.GetInterruptOnError())
	}
	return target
}

func targetTypeToDomain(executionType execution.TargetType) domain.TargetType {
	switch executionType {
	case execution.TargetType_TARGET_TYPE_UNSPECIFIED:
		return domain.TargetTypeUnspecified
	case execution.TargetType_TARGET_TYPE_REST_WEBHOOK:
		return domain.TargetTypeWebhook
	case execution.TargetType_TARGET_TYPE_REST_REQUEST_RESPONSE:
		return domain.TargetTypeRequestResponse
	default:
		return domain.TargetTypeUnspecified
	}
}
