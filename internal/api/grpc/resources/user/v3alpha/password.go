package user

import (
	"context"

	resource_object "github.com/zitadel/zitadel/internal/api/grpc/resources/object/v3alpha"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
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
	pw, verification := setPasswordToSetSchemaUserPassword(req.GetNewPassword())
	return &command.SetSchemaUserPassword{
		ResourceOwner: organizationToUpdateResourceOwner(req.Organization),
		UserID:        req.GetId(),
		Password:      pw,
		Verification:  verification,
	}
}

func setPasswordToSetSchemaUserPassword(req *user.SetPassword) (*command.SchemaUserPassword, *command.SchemaUserPasswordVerification) {
	return setPasswordToSchemaUserPassword(req.GetPassword(), req.GetHash(), req.GetChangeRequired()),
		setPasswordToSchemaUserPasswordVerification(req.GetCurrentPassword(), req.GetVerificationCode())
}

func setPasswordToSchemaUserPassword(pw string, hash string, changeRequired bool) *command.SchemaUserPassword {
	if pw == "" && hash == "" {
		return nil
	}
	return &command.SchemaUserPassword{
		Password:            pw,
		EncodedPasswordHash: hash,
		ChangeRequired:      changeRequired,
	}
}

func setPasswordToSchemaUserPasswordVerification(pw string, code string) *command.SchemaUserPasswordVerification {
	if pw == "" && code == "" {
		return nil
	}
	return &command.SchemaUserPasswordVerification{
		CurrentPassword: pw,
		Code:            code,
	}
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

func (s *Server) RequestPasswordReset(ctx context.Context, req *user.RequestPasswordResetRequest) (_ *user.RequestPasswordResetResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	schemauser := requestPasswordResetRequestToRequestSchemaUserPasswordReset(req)
	details, err := s.command.RequestSchemaUserPasswordReset(ctx, schemauser)
	if err != nil {
		return nil, err
	}
	return &user.RequestPasswordResetResponse{
		Details:          resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
		VerificationCode: schemauser.PlainCode,
	}, nil
}

func requestPasswordResetRequestToRequestSchemaUserPasswordReset(req *user.RequestPasswordResetRequest) *command.RequestSchemaUserPasswordReset {
	var notificationType domain.NotificationType
	if req.GetSendEmail() != nil {
		notificationType = domain.NotificationTypeEmail
	}
	if req.GetSendSms() != nil {
		notificationType = domain.NotificationTypeSms
	}
	return &command.RequestSchemaUserPasswordReset{
		ResourceOwner:    organizationToUpdateResourceOwner(req.Organization),
		UserID:           req.GetId(),
		URLTemplate:      req.GetSendEmail().GetUrlTemplate(),
		ReturnCode:       req.GetReturnCode() != nil,
		NotificationType: notificationType,
	}
}
