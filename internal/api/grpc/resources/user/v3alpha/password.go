package user

import (
	"context"

	resource_object "github.com/zitadel/zitadel/internal/api/grpc/resources/object/v3alpha"
	"github.com/zitadel/zitadel/internal/command"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/resources/user/v3alpha"
)

func (s *Server) SetPassword(ctx context.Context, req *user.SetPasswordRequest) (_ *user.SetPasswordResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.SetSchemaUserPassword(ctx, setPasswordRequestToSetSchemaUserPassword(req))
	if err != nil {
		return nil, err
	}
	return &user.SetPasswordResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
	}, nil
}

func setPasswordRequestToSetSchemaUserPassword(req *user.SetPasswordRequest) *command.SetSchemaUserPassword {
	return &command.SetSchemaUserPassword{
		ResourceOwner:       organizationToUpdateResourceOwner(req.Organization),
		UserID:              req.GetId(),
		Password:            req.GetNewPassword().GetPassword(),
		EncodedPasswordHash: req.GetNewPassword().GetHash(),
		ChangeRequired:      req.GetNewPassword().GetChangeRequired(),
	}
}

/*
func (s *Server) RemovePassword(ctx context.Context, req *user.RemovePasswordRequest) (_ *user.RemovePasswordResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.DeleteSchemaUserPassword(ctx, organizationToUpdateResourceOwner(req.Organization), req.GetId())
	if err != nil {
		return nil, err
	}
	return &user.RemovePasswordResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
	}, nil
}

func (s *Server) RemovePassword(ctx context.Context, req *user.RemovePasswordRequest) (_ *user.RemovePasswordResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.DeleteSchemaUserPassword(ctx, organizationToUpdateResourceOwner(req.Organization), req.GetId())
	if err != nil {
		return nil, err
	}
	return &user.RemovePasswordResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
	}, nil
}
*/
