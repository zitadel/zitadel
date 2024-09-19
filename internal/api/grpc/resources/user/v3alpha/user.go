package user

import (
	"context"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	resource_object "github.com/zitadel/zitadel/v2/internal/api/grpc/resources/object/v3alpha"
	"github.com/zitadel/zitadel/v2/internal/command"
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/zerrors"
	object "github.com/zitadel/zitadel/v2/pkg/grpc/object/v3alpha"
	"github.com/zitadel/zitadel/v2/pkg/grpc/resources/user/v3alpha"
)

func (s *Server) CreateUser(ctx context.Context, req *user.CreateUserRequest) (_ *user.CreateUserResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	schemauser, err := createUserRequestToCreateSchemaUser(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := s.command.CreateSchemaUser(ctx, schemauser, s.userCodeAlg); err != nil {
		return nil, err
	}
	return &user.CreateUserResponse{
		Details:   resource_object.DomainToDetailsPb(schemauser.Details, object.OwnerType_OWNER_TYPE_ORG, schemauser.Details.ResourceOwner),
		EmailCode: gu.Ptr(schemauser.ReturnCodeEmail),
		PhoneCode: gu.Ptr(schemauser.ReturnCodePhone),
	}, nil
}

func createUserRequestToCreateSchemaUser(ctx context.Context, req *user.CreateUserRequest) (*command.CreateSchemaUser, error) {
	data, err := req.GetUser().GetData().MarshalJSON()
	if err != nil {
		return nil, err
	}

	return &command.CreateSchemaUser{
		ResourceOwner: organizationToCreateResourceOwner(ctx, req.Organization),
		SchemaID:      req.GetUser().GetSchemaId(),
		ID:            req.GetUser().GetUserId(),
		Data:          data,
	}, nil
}

func organizationToCreateResourceOwner(ctx context.Context, org *object.Organization) string {
	resourceOwner := authz.GetCtxData(ctx).OrgID
	if resourceOwnerReq := resource_object.ResourceOwnerFromOrganization(org); resourceOwnerReq != "" {
		return resourceOwnerReq
	}
	return resourceOwner
}

func organizationToUpdateResourceOwner(org *object.Organization) string {
	if resourceOwnerReq := resource_object.ResourceOwnerFromOrganization(org); resourceOwnerReq != "" {
		return resourceOwnerReq
	}
	return ""
}

func (s *Server) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (_ *user.DeleteUserResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.DeleteSchemaUser(ctx, organizationToUpdateResourceOwner(req.Organization), req.GetId())
	if err != nil {
		return nil, err
	}
	return &user.DeleteUserResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
	}, nil
}

func checkUserSchemaEnabled(ctx context.Context) error {
	if authz.GetInstance(ctx).Features().UserSchema {
		return nil
	}
	return zerrors.ThrowPreconditionFailed(nil, "TODO", "Errors.UserSchema.NotEnabled")
}

func (s *Server) PatchUser(ctx context.Context, req *user.PatchUserRequest) (_ *user.PatchUserResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	schemauser, err := patchUserRequestToChangeSchemaUser(req)
	if err != nil {
		return nil, err
	}

	if err := s.command.ChangeSchemaUser(ctx, schemauser, s.userCodeAlg); err != nil {
		return nil, err
	}
	return &user.PatchUserResponse{
		Details:   resource_object.DomainToDetailsPb(schemauser.Details, object.OwnerType_OWNER_TYPE_ORG, schemauser.Details.ResourceOwner),
		EmailCode: gu.Ptr(schemauser.ReturnCodeEmail),
		PhoneCode: gu.Ptr(schemauser.ReturnCodePhone),
	}, nil
}

func patchUserRequestToChangeSchemaUser(req *user.PatchUserRequest) (_ *command.ChangeSchemaUser, err error) {
	var data []byte
	if req.GetUser().Data != nil {
		data, err = req.GetUser().GetData().MarshalJSON()
		if err != nil {
			return nil, err
		}
	}

	var email *command.Email
	var phone *command.Phone
	if req.GetUser().GetContact() != nil {
		if req.GetUser().GetContact().GetEmail() != nil {
			email = &command.Email{
				Address: domain.EmailAddress(req.GetUser().GetContact().Email.Address),
			}
			if req.GetUser().GetContact().Email.GetIsVerified() {
				email.Verified = true
			}
			if req.GetUser().GetContact().Email.GetReturnCode() != nil {
				email.ReturnCode = true
			}
			if req.GetUser().GetContact().Email.GetSendCode() != nil {
				email.URLTemplate = req.GetUser().GetContact().Email.GetSendCode().GetUrlTemplate()
			}
		}
		if req.GetUser().GetContact().Phone != nil {
			phone = &command.Phone{
				Number: domain.PhoneNumber(req.GetUser().GetContact().Phone.Number),
			}
			if req.GetUser().GetContact().Phone.GetIsVerified() {
				phone.Verified = true
			}
			if req.GetUser().GetContact().Phone.GetReturnCode() != nil {
				phone.ReturnCode = true
			}
		}
	}
	return &command.ChangeSchemaUser{
		ResourceOwner: organizationToUpdateResourceOwner(req.Organization),
		ID:            req.GetId(),
		SchemaID:      req.GetUser().SchemaId,
		Data:          data,
		Email:         email,
		Phone:         phone,
	}, nil
}

func (s *Server) DeactivateUser(ctx context.Context, req *user.DeactivateUserRequest) (_ *user.DeactivateUserResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}

	details, err := s.command.DeactivateSchemaUser(ctx, organizationToUpdateResourceOwner(req.Organization), req.GetId())
	if err != nil {
		return nil, err
	}
	return &user.DeactivateUserResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
	}, nil
}

func (s *Server) ActivateUser(ctx context.Context, req *user.ActivateUserRequest) (_ *user.ActivateUserResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}

	details, err := s.command.ActivateSchemaUser(ctx, organizationToUpdateResourceOwner(req.Organization), req.GetId())
	if err != nil {
		return nil, err
	}
	return &user.ActivateUserResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
	}, nil
}

func (s *Server) LockUser(ctx context.Context, req *user.LockUserRequest) (_ *user.LockUserResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}

	details, err := s.command.LockSchemaUser(ctx, organizationToUpdateResourceOwner(req.Organization), req.GetId())
	if err != nil {
		return nil, err
	}
	return &user.LockUserResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
	}, nil
}

func (s *Server) UnlockUser(ctx context.Context, req *user.UnlockUserRequest) (_ *user.UnlockUserResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}

	details, err := s.command.UnlockSchemaUser(ctx, organizationToUpdateResourceOwner(req.Organization), req.GetId())
	if err != nil {
		return nil, err
	}
	return &user.UnlockUserResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
	}, nil
}
