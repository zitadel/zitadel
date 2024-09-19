package action

import (
	"context"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	resource_object "github.com/zitadel/zitadel/v2/internal/api/grpc/resources/object/v3alpha"
	"github.com/zitadel/zitadel/v2/internal/command"
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/eventstore/v1/models"
	object "github.com/zitadel/zitadel/v2/pkg/grpc/object/v3alpha"
	action "github.com/zitadel/zitadel/v2/pkg/grpc/resources/action/v3alpha"
)

func (s *Server) CreateTarget(ctx context.Context, req *action.CreateTargetRequest) (*action.CreateTargetResponse, error) {
	if err := checkActionsEnabled(ctx); err != nil {
		return nil, err
	}
	add := createTargetToCommand(req)
	instanceID := authz.GetInstance(ctx).InstanceID()
	details, err := s.command.AddTarget(ctx, add, instanceID)
	if err != nil {
		return nil, err
	}
	return &action.CreateTargetResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_INSTANCE, instanceID),
	}, nil
}

func (s *Server) PatchTarget(ctx context.Context, req *action.PatchTargetRequest) (*action.PatchTargetResponse, error) {
	if err := checkActionsEnabled(ctx); err != nil {
		return nil, err
	}
	instanceID := authz.GetInstance(ctx).InstanceID()
	details, err := s.command.ChangeTarget(ctx, patchTargetToCommand(req), instanceID)
	if err != nil {
		return nil, err
	}
	return &action.PatchTargetResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_INSTANCE, instanceID),
	}, nil
}

func (s *Server) DeleteTarget(ctx context.Context, req *action.DeleteTargetRequest) (*action.DeleteTargetResponse, error) {
	if err := checkActionsEnabled(ctx); err != nil {
		return nil, err
	}
	instanceID := authz.GetInstance(ctx).InstanceID()
	details, err := s.command.DeleteTarget(ctx, req.GetId(), instanceID)
	if err != nil {
		return nil, err
	}
	return &action.DeleteTargetResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_INSTANCE, instanceID),
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
