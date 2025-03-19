package action

import (
	"context"

	"github.com/muhlemmer/gu"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	action "github.com/zitadel/zitadel/pkg/grpc/action/v2beta"
)

func (s *Server) CreateTarget(ctx context.Context, req *action.CreateTargetRequest) (*action.CreateTargetResponse, error) {
	if err := checkActionsEnabled(ctx); err != nil {
		return nil, err
	}
	add := createTargetToCommand(req)
	instanceID := authz.GetInstance(ctx).InstanceID()
	createdAt, err := s.command.AddTarget(ctx, add, instanceID)
	if err != nil {
		return nil, err
	}
	var creationDate *timestamppb.Timestamp
	if !createdAt.IsZero() {
		creationDate = timestamppb.New(createdAt)
	}
	return &action.CreateTargetResponse{
		Id:           add.AggregateID,
		CreationDate: creationDate,
		SigningKey:   add.SigningKey,
	}, nil
}

func (s *Server) UpdateTarget(ctx context.Context, req *action.UpdateTargetRequest) (*action.UpdateTargetResponse, error) {
	if err := checkActionsEnabled(ctx); err != nil {
		return nil, err
	}
	instanceID := authz.GetInstance(ctx).InstanceID()
	update := updateTargetToCommand(req)
	changedAt, err := s.command.ChangeTarget(ctx, update, instanceID)
	if err != nil {
		return nil, err
	}
	var changeDate *timestamppb.Timestamp
	if !changedAt.IsZero() {
		changeDate = timestamppb.New(changedAt)
	}
	return &action.UpdateTargetResponse{
		ChangeDate: changeDate,
		SigningKey: update.SigningKey,
	}, nil
}

func (s *Server) DeleteTarget(ctx context.Context, req *action.DeleteTargetRequest) (*action.DeleteTargetResponse, error) {
	if err := checkActionsEnabled(ctx); err != nil {
		return nil, err
	}
	instanceID := authz.GetInstance(ctx).InstanceID()
	deletedAt, err := s.command.DeleteTarget(ctx, req.GetId(), instanceID)
	if err != nil {
		return nil, err
	}
	var deletionDate *timestamppb.Timestamp
	if !deletedAt.IsZero() {
		deletionDate = timestamppb.New(deletedAt)
	}
	return &action.DeleteTargetResponse{
		DeletionDate: deletionDate,
	}, nil
}

func createTargetToCommand(req *action.CreateTargetRequest) *command.AddTarget {
	var (
		targetType       domain.TargetType
		interruptOnError bool
	)
	switch t := req.GetTargetType().(type) {
	case *action.CreateTargetRequest_RestWebhook:
		targetType = domain.TargetTypeWebhook
		interruptOnError = t.RestWebhook.InterruptOnError
	case *action.CreateTargetRequest_RestCall:
		targetType = domain.TargetTypeCall
		interruptOnError = t.RestCall.InterruptOnError
	case *action.CreateTargetRequest_RestAsync:
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

func updateTargetToCommand(req *action.UpdateTargetRequest) *command.ChangeTarget {
	expirationSigningKey := false
	// TODO handle expiration, currently only immediate expiration is supported
	if req.GetExpirationSigningKey() != nil {
		expirationSigningKey = true
	}

	if req == nil {
		return nil
	}
	target := &command.ChangeTarget{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.GetId(),
		},
		Name:                 req.Name,
		Endpoint:             req.Endpoint,
		ExpirationSigningKey: expirationSigningKey,
	}
	if req.TargetType != nil {
		switch t := req.GetTargetType().(type) {
		case *action.UpdateTargetRequest_RestWebhook:
			target.TargetType = gu.Ptr(domain.TargetTypeWebhook)
			target.InterruptOnError = gu.Ptr(t.RestWebhook.InterruptOnError)
		case *action.UpdateTargetRequest_RestCall:
			target.TargetType = gu.Ptr(domain.TargetTypeCall)
			target.InterruptOnError = gu.Ptr(t.RestCall.InterruptOnError)
		case *action.UpdateTargetRequest_RestAsync:
			target.TargetType = gu.Ptr(domain.TargetTypeAsync)
			target.InterruptOnError = gu.Ptr(false)
		}
	}
	if req.Timeout != nil {
		target.Timeout = gu.Ptr(req.GetTimeout().AsDuration())
	}
	return target
}
