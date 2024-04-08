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
	var targetType domain.TargetType
	var interruptOnError bool
	switch t := req.GetTargetType().(type) {
	case *execution.CreateTargetRequest_RestWebhook:
		targetType = domain.TargetTypeWebhook
		interruptOnError = t.RestWebhook.InterruptOnError
	case *execution.CreateTargetRequest_RestCall:
		targetType = domain.TargetTypeCall
		interruptOnError = t.RestCall.InterruptOnError
	case *execution.CreateTargetRequest_RestAsync:
		targetType = domain.TargetTypeAsync
	}
	return &command.AddTarget{
		Name:             req.GetName(),
		TargetType:       targetType,
		Endpoint:         req.GetEndpoint(),
		Timeout:          req.GetTimeout().AsDuration(),
		InterruptOnError: interruptOnError,
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
		Name:     req.Name,
		Endpoint: req.Endpoint,
	}
	if req.TargetType != nil {
		switch t := req.GetTargetType().(type) {
		case *execution.UpdateTargetRequest_RestWebhook:
			target.TargetType = gu.Ptr(domain.TargetTypeWebhook)
			target.InterruptOnError = gu.Ptr(t.RestWebhook.InterruptOnError)
		case *execution.UpdateTargetRequest_RestCall:
			target.TargetType = gu.Ptr(domain.TargetTypeCall)
			target.InterruptOnError = gu.Ptr(t.RestCall.InterruptOnError)
		case *execution.UpdateTargetRequest_RestAsync:
			target.TargetType = gu.Ptr(domain.TargetTypeAsync)
			target.InterruptOnError = gu.Ptr(false)
		}
	}
	if req.Timeout != nil {
		target.Timeout = gu.Ptr(req.GetTimeout().AsDuration())
	}
	return target
}
