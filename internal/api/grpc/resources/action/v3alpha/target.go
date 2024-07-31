package action

import (
	"context"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/api/authz"
	resource_object "github.com/zitadel/zitadel/internal/api/grpc/resources/object/v3alpha"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	action "github.com/zitadel/zitadel/pkg/grpc/resources/action/v3alpha"
)

func (s *Server) CreateTarget(ctx context.Context, req *action.CreateTargetRequest) (*action.CreateTargetResponse, error) {
	if err := checkExecutionEnabled(ctx); err != nil {
		return nil, err
	}
	add := createTargetToCommand(req)
	instance := targetOwnerInstance(ctx)
	details, err := s.command.AddTarget(ctx, add, instance.Id)
	if err != nil {
		return nil, err
	}
	return &action.CreateTargetResponse{
		Details: resource_object.DomainToDetailsPb(details, instance, add.AggregateID),
	}, nil
}

func (s *Server) PatchTarget(ctx context.Context, req *action.PatchTargetRequest) (*action.PatchTargetResponse, error) {
	if err := checkExecutionEnabled(ctx); err != nil {
		return nil, err
	}
	instance := targetOwnerInstance(ctx)
	details, err := s.command.ChangeTarget(ctx, patchTargetToCommand(req), instance.Id)
	if err != nil {
		return nil, err
	}
	return &action.PatchTargetResponse{
		Details: resource_object.DomainToDetailsPb(details, instance, req.GetId()),
	}, nil
}

func (s *Server) DeleteTarget(ctx context.Context, req *action.DeleteTargetRequest) (*action.DeleteTargetResponse, error) {
	if err := checkExecutionEnabled(ctx); err != nil {
		return nil, err
	}
	instance := targetOwnerInstance(ctx)
	details, err := s.command.DeleteTarget(ctx, req.GetId(), instance.Id)
	if err != nil {
		return nil, err
	}
	return &action.DeleteTargetResponse{
		Details: resource_object.DomainToDetailsPb(details, instance, req.GetId()),
	}, nil
}

func createTargetToCommand(req *action.CreateTargetRequest) *command.AddTarget {
	reqTarget := req.GetTarget()
	var (
		targetType       domain.TargetType
		interruptOnError bool
	)
	switch t := reqTarget.GetTargetType().(type) {
	case *action.Target_RestWebhook:
		targetType = domain.TargetTypeWebhook
		interruptOnError = t.RestWebhook.InterruptOnError
	case *action.Target_RestCall:
		targetType = domain.TargetTypeCall
		interruptOnError = t.RestCall.InterruptOnError
	case *action.Target_RestAsync:
		targetType = domain.TargetTypeAsync
	}
	return &command.AddTarget{
		Name:             reqTarget.GetName(),
		TargetType:       targetType,
		Endpoint:         reqTarget.GetEndpoint(),
		Timeout:          reqTarget.GetTimeout().AsDuration(),
		InterruptOnError: interruptOnError,
	}
}

func patchTargetToCommand(req *action.PatchTargetRequest) *command.ChangeTarget {
	reqTarget := req.GetTarget()
	if reqTarget == nil {
		return nil
	}
	target := &command.ChangeTarget{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.GetId(),
		},
		Name:     reqTarget.Name,
		Endpoint: reqTarget.Endpoint,
	}
	if reqTarget.TargetType != nil {
		switch t := reqTarget.GetTargetType().(type) {
		case *action.PatchTarget_RestWebhook:
			target.TargetType = gu.Ptr(domain.TargetTypeWebhook)
			target.InterruptOnError = gu.Ptr(t.RestWebhook.InterruptOnError)
		case *action.PatchTarget_RestCall:
			target.TargetType = gu.Ptr(domain.TargetTypeCall)
			target.InterruptOnError = gu.Ptr(t.RestCall.InterruptOnError)
		case *action.PatchTarget_RestAsync:
			target.TargetType = gu.Ptr(domain.TargetTypeAsync)
			target.InterruptOnError = gu.Ptr(false)
		}
	}
	if reqTarget.Timeout != nil {
		target.Timeout = gu.Ptr(reqTarget.GetTimeout().AsDuration())
	}
	return target
}

func targetOwnerInstance(ctx context.Context) *object.Owner {
	return &object.Owner{
		Type: object.OwnerType_OWNER_TYPE_INSTANCE,
		Id:   authz.GetInstance(ctx).InstanceID(),
	}
}
