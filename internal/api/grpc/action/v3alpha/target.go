package action

import (
	"context"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	action "github.com/zitadel/zitadel/pkg/grpc/action/v3alpha"
)

func (s *Server) CreateTarget(ctx context.Context, req *action.CreateTargetRequest) (*action.CreateTargetResponse, error) {
	if err := checkExecutionEnabled(ctx); err != nil {
		return nil, err
	}

	add := createTargetToCommand(req)
	details, err := s.command.AddTarget(ctx, add, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	return &action.CreateTargetResponse{
		Id:      add.AggregateID,
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) UpdateTarget(ctx context.Context, req *action.UpdateTargetRequest) (*action.UpdateTargetResponse, error) {
	if err := checkExecutionEnabled(ctx); err != nil {
		return nil, err
	}

	details, err := s.command.ChangeTarget(ctx, updateTargetToCommand(req), authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	return &action.UpdateTargetResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) DeleteTarget(ctx context.Context, req *action.DeleteTargetRequest) (*action.DeleteTargetResponse, error) {
	if err := checkExecutionEnabled(ctx); err != nil {
		return nil, err
	}

	details, err := s.command.DeleteTarget(ctx, req.GetTargetId(), authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	return &action.DeleteTargetResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func createTargetToCommand(req *action.CreateTargetRequest) *command.AddTarget {
	var targetType domain.TargetType
	var url string
	switch t := req.GetTargetType().(type) {
	case *action.CreateTargetRequest_RestWebhook:
		targetType = domain.TargetTypeWebhook
		url = t.RestWebhook.GetUrl()
	case *action.CreateTargetRequest_RestRequestResponse:
		targetType = domain.TargetTypeRequestResponse
		url = t.RestRequestResponse.GetUrl()
	}
	return &command.AddTarget{
		Name:             req.GetName(),
		TargetType:       targetType,
		URL:              url,
		Timeout:          req.GetTimeout().AsDuration(),
		Async:            req.GetIsAsync(),
		InterruptOnError: req.GetInterruptOnError(),
	}
}

func updateTargetToCommand(req *action.UpdateTargetRequest) *command.ChangeTarget {
	if req == nil {
		return nil
	}
	target := &command.ChangeTarget{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.GetTargetId(),
		},
		Name: req.Name,
	}
	switch t := req.GetTargetType().(type) {
	case *action.UpdateTargetRequest_RestWebhook:
		target.TargetType = gu.Ptr(domain.TargetTypeWebhook)
		target.URL = gu.Ptr(t.RestWebhook.GetUrl())
	case *action.UpdateTargetRequest_RestRequestResponse:
		target.TargetType = gu.Ptr(domain.TargetTypeRequestResponse)
		target.URL = gu.Ptr(t.RestRequestResponse.GetUrl())
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
