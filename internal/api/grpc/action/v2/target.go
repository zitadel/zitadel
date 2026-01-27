package action

import (
	"context"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/muhlemmer/gu"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	target_domain "github.com/zitadel/zitadel/internal/execution/target"
	"github.com/zitadel/zitadel/pkg/grpc/action/v2"
)

func (s *Server) CreateTarget(ctx context.Context, req *connect.Request[action.CreateTargetRequest]) (*connect.Response[action.CreateTargetResponse], error) {
	add := createTargetToCommand(req.Msg)
	instanceID := authz.GetInstance(ctx).InstanceID()
	createdAt, err := s.command.AddTarget(ctx, add, instanceID)
	if err != nil {
		return nil, err
	}
	var creationDate *timestamppb.Timestamp
	if !createdAt.IsZero() {
		creationDate = timestamppb.New(createdAt)
	}
	return connect.NewResponse(&action.CreateTargetResponse{
		Id:           add.AggregateID,
		CreationDate: creationDate,
		SigningKey:   add.SigningKey,
	}), nil
}

func (s *Server) UpdateTarget(ctx context.Context, req *connect.Request[action.UpdateTargetRequest]) (*connect.Response[action.UpdateTargetResponse], error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	update := updateTargetToCommand(req.Msg)
	changedAt, err := s.command.ChangeTarget(ctx, update, instanceID)
	if err != nil {
		return nil, err
	}
	var changeDate *timestamppb.Timestamp
	if !changedAt.IsZero() {
		changeDate = timestamppb.New(changedAt)
	}
	return connect.NewResponse(&action.UpdateTargetResponse{
		ChangeDate: changeDate,
		SigningKey: update.SigningKey,
	}), nil
}

func (s *Server) DeleteTarget(ctx context.Context, req *connect.Request[action.DeleteTargetRequest]) (*connect.Response[action.DeleteTargetResponse], error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	deletedAt, err := s.command.DeleteTarget(ctx, req.Msg.GetId(), instanceID)
	if err != nil {
		return nil, err
	}
	var deletionDate *timestamppb.Timestamp
	if !deletedAt.IsZero() {
		deletionDate = timestamppb.New(deletedAt)
	}
	return connect.NewResponse(&action.DeleteTargetResponse{
		DeletionDate: deletionDate,
	}), nil
}

func (s *Server) AddPublicKey(ctx context.Context, req *connect.Request[action.AddPublicKeyRequest]) (*connect.Response[action.AddPublicKeyResponse], error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	var expirationDate time.Time
	if req.Msg.GetExpirationDate() != nil {
		expirationDate = req.Msg.GetExpirationDate().AsTime()
	}
	key := &command.TargetPublicKey{
		TargetID:   strings.TrimSpace(req.Msg.GetTargetId()),
		PublicKey:  req.Msg.GetPublicKey(),
		Expiration: expirationDate,
	}
	creationDate, err := s.command.AddTargetPublicKey(ctx, key, instanceID)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&action.AddPublicKeyResponse{
		KeyId:        key.KeyID,
		CreationDate: timestamppb.New(creationDate),
	}), nil
}

func (s *Server) ActivatePublicKey(ctx context.Context, req *connect.Request[action.ActivatePublicKeyRequest]) (*connect.Response[action.ActivatePublicKeyResponse], error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	changeDate, err := s.command.ActivateTargetPublicKey(ctx, strings.TrimSpace(req.Msg.GetTargetId()), strings.TrimSpace(req.Msg.GetKeyId()), instanceID)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&action.ActivatePublicKeyResponse{
		ChangeDate: timestamppb.New(changeDate),
	}), nil
}

func (s *Server) DeactivatePublicKey(ctx context.Context, req *connect.Request[action.DeactivatePublicKeyRequest]) (*connect.Response[action.DeactivatePublicKeyResponse], error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	changeDate, err := s.command.DeactivateTargetPublicKey(ctx, strings.TrimSpace(req.Msg.GetTargetId()), strings.TrimSpace(req.Msg.GetKeyId()), instanceID)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&action.DeactivatePublicKeyResponse{
		ChangeDate: timestamppb.New(changeDate),
	}), nil
}

func (s *Server) RemovePublicKey(ctx context.Context, req *connect.Request[action.RemovePublicKeyRequest]) (*connect.Response[action.RemovePublicKeyResponse], error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	deletionDate, err := s.command.RemoveTargetPublicKey(ctx, strings.TrimSpace(req.Msg.GetTargetId()), strings.TrimSpace(req.Msg.GetKeyId()), instanceID)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&action.RemovePublicKeyResponse{
		DeletionDate: timestamppb.New(deletionDate),
	}), nil
}

func createTargetToCommand(req *action.CreateTargetRequest) *command.AddTarget {
	var (
		targetType       target_domain.TargetType
		interruptOnError bool
	)
	switch t := req.GetTargetType().(type) {
	case *action.CreateTargetRequest_RestWebhook:
		targetType = target_domain.TargetTypeWebhook
		interruptOnError = t.RestWebhook.InterruptOnError
	case *action.CreateTargetRequest_RestCall:
		targetType = target_domain.TargetTypeCall
		interruptOnError = t.RestCall.InterruptOnError
	case *action.CreateTargetRequest_RestAsync:
		targetType = target_domain.TargetTypeAsync
	}
	return &command.AddTarget{
		Name:             req.GetName(),
		TargetType:       targetType,
		Endpoint:         req.GetEndpoint(),
		Timeout:          req.GetTimeout().AsDuration(),
		InterruptOnError: interruptOnError,
		PayloadType:      payloadTypeToDomain(req.GetPayloadType()),
	}
}

func payloadTypeToDomain(payloadType action.PayloadType) target_domain.PayloadType {
	switch payloadType {
	case action.PayloadType_PAYLOAD_TYPE_UNSPECIFIED:
		return target_domain.PayloadTypeUnspecified
	case action.PayloadType_PAYLOAD_TYPE_JSON:
		return target_domain.PayloadTypeJSON
	case action.PayloadType_PAYLOAD_TYPE_JWT:
		return target_domain.PayloadTypeJWT
	case action.PayloadType_PAYLOAD_TYPE_JWE:
		return target_domain.PayloadTypeJWE
	default:
		return target_domain.PayloadTypeUnspecified
	}
}

func updateTargetToCommand(req *action.UpdateTargetRequest) *command.ChangeTarget {
	// TODO handle expiration, currently only immediate expiration is supported
	expirationSigningKey := req.GetExpirationSigningKey() != nil

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
		PayloadType:          payloadTypeToDomain(req.GetPayloadType()),
	}
	if req.TargetType != nil {
		switch t := req.GetTargetType().(type) {
		case *action.UpdateTargetRequest_RestWebhook:
			target.TargetType = gu.Ptr(target_domain.TargetTypeWebhook)
			target.InterruptOnError = gu.Ptr(t.RestWebhook.InterruptOnError)
		case *action.UpdateTargetRequest_RestCall:
			target.TargetType = gu.Ptr(target_domain.TargetTypeCall)
			target.InterruptOnError = gu.Ptr(t.RestCall.InterruptOnError)
		case *action.UpdateTargetRequest_RestAsync:
			target.TargetType = gu.Ptr(target_domain.TargetTypeAsync)
			target.InterruptOnError = gu.Ptr(false)
		}
	}
	if req.Timeout != nil {
		target.Timeout = gu.Ptr(req.GetTimeout().AsDuration())
	}
	return target
}
