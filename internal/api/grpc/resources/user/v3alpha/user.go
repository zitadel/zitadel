package user

import (
	"context"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/api/authz"
	resource_object "github.com/zitadel/zitadel/internal/api/grpc/resources/object/v3alpha"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	"github.com/zitadel/zitadel/pkg/grpc/resources/user/v3alpha"
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
		Details:   resource_object.DomainToDetailsPb(schemauser.Details, object.OwnerType_OWNER_TYPE_ORG, schemauser.ResourceOwner),
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
		ResourceOwner: authz.GetCtxData(ctx).OrgID,
		SchemaID:      req.GetUser().GetSchemaId(),
		ID:            req.GetUser().GetUserId(),
		Data:          data,
	}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (_ *user.DeleteUserResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.DeleteSchemaUser(ctx, req.GetId())
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
	schemauser, err := patchUserRequestToChangeSchemaUser(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := s.command.ChangeSchemaUser(ctx, schemauser, s.userCodeAlg); err != nil {
		return nil, err
	}
	return &user.PatchUserResponse{
		Details:   resource_object.DomainToDetailsPb(schemauser.Details, object.OwnerType_OWNER_TYPE_ORG, schemauser.ResourceOwner),
		EmailCode: gu.Ptr(schemauser.ReturnCodeEmail),
		PhoneCode: gu.Ptr(schemauser.ReturnCodePhone),
	}, nil
}

func patchUserRequestToChangeSchemaUser(ctx context.Context, req *user.PatchUserRequest) (*command.ChangeSchemaUser, error) {
	data, err := req.GetUser().GetData().MarshalJSON()
	if err != nil {
		return nil, err
	}

	var email *command.Email
	if req.GetUser().GetContact().Email != nil {
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
	var phone *command.Phone
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
	return &command.ChangeSchemaUser{
		ResourceOwner: authz.GetCtxData(ctx).OrgID,
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

	details, err := s.command.DeactivateSchemaUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &user.DeactivateUserResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
	}, nil
}

func (s *Server) ReactivateUser(ctx context.Context, req *user.ReactivateUserRequest) (_ *user.ReactivateUserResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}

	details, err := s.command.ReactivateSchemaUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &user.ReactivateUserResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
	}, nil
}

func (s *Server) LockUser(ctx context.Context, req *user.LockUserRequest) (_ *user.LockUserResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}

	details, err := s.command.LockSchemaUser(ctx, req.GetId())
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

	details, err := s.command.UnlockSchemaUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &user.UnlockUserResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
	}, nil
}
