package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

func (s *Server) RequestPasswordReset(ctx context.Context, req *user.RequestPasswordResetRequest) (_ *user.RequestPasswordResetResponse, err error) {
	var resourceOwner string // TODO: check if still needed
	var details *domain.ObjectDetails
	var code *string

	switch m := req.GetMedium().(type) {
	case *user.RequestPasswordResetRequest_SendLink:
		details, code, err = s.command.RequestPasswordResetURLTemplate(ctx, req.GetUserId(), resourceOwner, m.SendLink.GetUrlTemplate())
	case *user.RequestPasswordResetRequest_ReturnCode:
		details, code, err = s.command.RequestPasswordResetReturnCode(ctx, req.GetUserId(), resourceOwner)
	case nil:
		details, code, err = s.command.RequestPasswordReset(ctx, req.GetUserId(), resourceOwner)
	default:
		err = caos_errs.ThrowUnimplementedf(nil, "USERv2-SDeeg", "verification oneOf %T in method RequestPasswordReset not implemented", m)
	}
	if err != nil {
		return nil, err
	}

	return &user.RequestPasswordResetResponse{
		Details:          object.DomainToDetailsPb(details),
		VerificationCode: code,
	}, nil
}

func (s *Server) SetPassword(ctx context.Context, req *user.SetPasswordRequest) (_ *user.SetPasswordResponse, err error) {
	var resourceOwner = authz.GetCtxData(ctx).ResourceOwner
	var details *domain.ObjectDetails

	switch v := req.GetVerification().(type) {
	case *user.SetPasswordRequest_CurrentPassword:
		details, err = s.command.ChangePassword(ctx, resourceOwner, req.GetUserId(), v.CurrentPassword, req.GetNewPassword().GetPassword(), "")
	case *user.SetPasswordRequest_VerificationCode:
		details, err = s.command.SetPasswordWithVerifyCode(ctx, resourceOwner, req.GetUserId(), v.VerificationCode, req.GetNewPassword().GetPassword(), "")
	case nil:
		details, err = s.command.SetPassword(ctx, resourceOwner, req.GetUserId(), req.GetNewPassword().GetPassword(), req.GetNewPassword().GetChangeRequired())
	default:
		err = caos_errs.ThrowUnimplementedf(nil, "USERv2-SFdf2", "verification oneOf %T in method SetPasswordRequest not implemented", v)
	}
	if err != nil {
		return nil, err
	}

	return &user.SetPasswordResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}
