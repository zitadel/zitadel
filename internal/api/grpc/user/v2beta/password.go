package user

import (
	"context"

	"connectrpc.com/connect"

	object "github.com/zitadel/zitadel/internal/api/grpc/object/v2beta"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

func (s *Server) PasswordReset(ctx context.Context, req *connect.Request[user.PasswordResetRequest]) (_ *connect.Response[user.PasswordResetResponse], err error) {
	var details *domain.ObjectDetails
	var code *string

	switch m := req.Msg.GetMedium().(type) {
	case *user.PasswordResetRequest_SendLink:
		details, code, err = s.command.RequestPasswordResetURLTemplate(ctx, req.Msg.GetUserId(), m.SendLink.GetUrlTemplate(), notificationTypeToDomain(m.SendLink.GetNotificationType()))
	case *user.PasswordResetRequest_ReturnCode:
		details, code, err = s.command.RequestPasswordResetReturnCode(ctx, req.Msg.GetUserId())
	case nil:
		details, code, err = s.command.RequestPasswordReset(ctx, req.Msg.GetUserId())
	default:
		err = zerrors.ThrowUnimplementedf(nil, "USERv2-SDeeg", "verification oneOf %T in method RequestPasswordReset not implemented", m)
	}
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&user.PasswordResetResponse{
		Details:          object.DomainToDetailsPb(details),
		VerificationCode: code,
	}), nil
}

func notificationTypeToDomain(notificationType user.NotificationType) domain.NotificationType {
	switch notificationType {
	case user.NotificationType_NOTIFICATION_TYPE_Email:
		return domain.NotificationTypeEmail
	case user.NotificationType_NOTIFICATION_TYPE_SMS:
		return domain.NotificationTypeSms
	case user.NotificationType_NOTIFICATION_TYPE_Unspecified:
		return domain.NotificationTypeEmail
	default:
		return domain.NotificationTypeEmail
	}
}

func (s *Server) SetPassword(ctx context.Context, req *connect.Request[user.SetPasswordRequest]) (_ *connect.Response[user.SetPasswordResponse], err error) {
	var details *domain.ObjectDetails

	switch v := req.Msg.GetVerification().(type) {
	case *user.SetPasswordRequest_CurrentPassword:
		details, err = s.command.ChangePassword(ctx, "", req.Msg.GetUserId(), v.CurrentPassword, req.Msg.GetNewPassword().GetPassword(), "", req.Msg.GetNewPassword().GetChangeRequired())
	case *user.SetPasswordRequest_VerificationCode:
		details, err = s.command.SetPasswordWithVerifyCode(ctx, "", req.Msg.GetUserId(), v.VerificationCode, req.Msg.GetNewPassword().GetPassword(), "", req.Msg.GetNewPassword().GetChangeRequired())
	case nil:
		details, err = s.command.SetPassword(ctx, "", req.Msg.GetUserId(), req.Msg.GetNewPassword().GetPassword(), req.Msg.GetNewPassword().GetChangeRequired())
	default:
		err = zerrors.ThrowUnimplementedf(nil, "USERv2-SFdf2", "verification oneOf %T in method SetPasswordRequest not implemented", v)
	}
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&user.SetPasswordResponse{
		Details: object.DomainToDetailsPb(details),
	}), nil
}
