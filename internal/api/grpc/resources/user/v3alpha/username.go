package user

import (
	"context"

	resource_object "github.com/zitadel/zitadel/internal/api/grpc/resources/object/v3alpha"
	"github.com/zitadel/zitadel/internal/command"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/resources/user/v3alpha"
)

func (s *Server) AddUsername(ctx context.Context, req *user.AddUsernameRequest) (_ *user.AddUsernameResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.AddUsername(ctx, addUsernameRequestToAddUsername(req))
	if err != nil {
		return nil, err
	}
	return &user.AddUsernameResponse{
		Details:    resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
		UsernameId: details.ID,
	}, nil
}

func addUsernameRequestToAddUsername(req *user.AddUsernameRequest) *command.AddUsername {
	if req == nil {
		return nil
	}
	return &command.AddUsername{
		ResourceOwner: organizationToUpdateResourceOwner(req.Organization),
		UserID:        req.GetId(),
		Username:      setUsernameToAddUsername(req.GetUsername()),
	}
}

func setUsernameToAddUsername(req *user.SetUsername) *command.Username {
	if req == nil {
		return nil
	}
	return &command.Username{
		Username:      req.GetUsername(),
		IsOrgSpecific: req.GetIsOrganizationSpecific(),
	}
}

func (s *Server) RemoveUsername(ctx context.Context, req *user.RemoveUsernameRequest) (_ *user.RemoveUsernameResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.DeleteUsername(ctx, organizationToUpdateResourceOwner(req.Organization), req.GetId(), req.GetUsernameId())
	if err != nil {
		return nil, err
	}
	return &user.RemoveUsernameResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
	}, nil
}
