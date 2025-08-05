package user

import (
	"context"

	resource_object "github.com/zitadel/zitadel/internal/api/grpc/resources/object/v3alpha"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/resources/user/v3alpha"
)

func (s *Server) SetContactPhone(ctx context.Context, req *user.SetContactPhoneRequest) (_ *user.SetContactPhoneResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	schemauser := setContactPhoneRequestToChangeSchemaUserPhone(req)
	details, err := s.command.ChangeSchemaUserPhone(ctx, schemauser)
	if err != nil {
		return nil, err
	}
	return &user.SetContactPhoneResponse{
		Details:          resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
		VerificationCode: schemauser.ReturnCode,
	}, nil
}

func setContactPhoneRequestToChangeSchemaUserPhone(req *user.SetContactPhoneRequest) *command.ChangeSchemaUserPhone {
	return &command.ChangeSchemaUserPhone{
		ResourceOwner: organizationToUpdateResourceOwner(req.Organization),
		ID:            req.GetId(),
		Phone:         setPhoneToPhone(req.Phone),
	}
}

func setPhoneToPhone(setPhone *user.SetPhone) *command.Phone {
	if setPhone == nil {
		return nil
	}
	return &command.Phone{
		Number:     domain.PhoneNumber(setPhone.Number),
		ReturnCode: setPhone.GetReturnCode() != nil,
		Verified:   setPhone.GetIsVerified(),
	}
}

func (s *Server) VerifyContactPhone(ctx context.Context, req *user.VerifyContactPhoneRequest) (_ *user.VerifyContactPhoneResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.VerifySchemaUserPhone(ctx, organizationToUpdateResourceOwner(req.Organization), req.GetId(), req.GetVerificationCode())
	if err != nil {
		return nil, err
	}
	return &user.VerifyContactPhoneResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
	}, nil
}

func (s *Server) ResendContactPhoneCode(ctx context.Context, req *user.ResendContactPhoneCodeRequest) (_ *user.ResendContactPhoneCodeResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	schemauser := resendContactPhoneCodeRequestToResendSchemaUserPhoneCode(req)
	details, err := s.command.ResendSchemaUserPhoneCode(ctx, schemauser)
	if err != nil {
		return nil, err
	}
	return &user.ResendContactPhoneCodeResponse{
		Details:          resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
		VerificationCode: schemauser.PlainCode,
	}, nil
}

func resendContactPhoneCodeRequestToResendSchemaUserPhoneCode(req *user.ResendContactPhoneCodeRequest) *command.ResendSchemaUserPhoneCode {
	return &command.ResendSchemaUserPhoneCode{
		ResourceOwner: organizationToUpdateResourceOwner(req.Organization),
		ID:            req.GetId(),
		ReturnCode:    req.GetReturnCode() != nil,
	}
}
