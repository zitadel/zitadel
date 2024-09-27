package user

import (
	"context"
	"time"

	resource_object "github.com/zitadel/zitadel/internal/api/grpc/resources/object/v3alpha"
	"github.com/zitadel/zitadel/internal/command"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/resources/user/v3alpha"
)

func (s *Server) AddPublicKey(ctx context.Context, req *user.AddPublicKeyRequest) (_ *user.AddPublicKeyResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	pk := addPublicKeyRequestToAddPublicKey(req)
	details, err := s.command.AddPublicKey(ctx, pk)
	if err != nil {
		return nil, err
	}
	return &user.AddPublicKeyResponse{
		Details:     resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
		PublicKeyId: details.ID,
		PrivateKey:  pk.PrivateKey,
	}, nil
}

func addPublicKeyRequestToAddPublicKey(req *user.AddPublicKeyRequest) *command.AddPublicKey {
	expDate := time.Time{}
	if req.GetPublicKey().GetExpirationDate() != nil {
		expDate = req.GetPublicKey().GetExpirationDate().AsTime()
	}
	return &command.AddPublicKey{
		ResourceOwner:  organizationToUpdateResourceOwner(req.Organization),
		UserID:         req.GetId(),
		PublicKey:      req.GetPublicKey().GetPublicKey().GetPublicKey(),
		ExpirationDate: expDate,
	}
}

func (s *Server) RemovePublicKey(ctx context.Context, req *user.RemovePublicKeyRequest) (_ *user.RemovePublicKeyResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.DeletePublicKey(ctx, organizationToUpdateResourceOwner(req.Organization), req.GetId(), req.GetPublicKeyId())
	if err != nil {
		return nil, err
	}
	return &user.RemovePublicKeyResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
	}, nil
}
