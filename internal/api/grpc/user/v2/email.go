package user

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) SetEmail(ctx context.Context, req *connect.Request[user.SetEmailRequest]) (resp *connect.Response[user.SetEmailResponse], err error) {
	var email *domain.Email

	switch v := req.Msg.GetVerification().(type) {
	case *user.SetEmailRequest_SendCode:
		email, err = s.command.ChangeUserEmailURLTemplate(ctx, req.Msg.GetUserId(), req.Msg.GetEmail(), s.userCodeAlg, v.SendCode.GetUrlTemplate())
	case *user.SetEmailRequest_ReturnCode:
		email, err = s.command.ChangeUserEmailReturnCode(ctx, req.Msg.GetUserId(), req.Msg.GetEmail(), s.userCodeAlg)
	case *user.SetEmailRequest_IsVerified:
		email, err = s.command.ChangeUserEmailVerified(ctx, req.Msg.GetUserId(), req.Msg.GetEmail())
	case nil:
		email, err = s.command.ChangeUserEmail(ctx, req.Msg.GetUserId(), req.Msg.GetEmail(), s.userCodeAlg)
	default:
		err = zerrors.ThrowUnimplementedf(nil, "USERv2-Ahng0", "verification oneOf %T in method SetEmail not implemented", v)
	}
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&user.SetEmailResponse{
		Details: &object.Details{
			Sequence:      email.Sequence,
			ChangeDate:    timestamppb.New(email.ChangeDate),
			ResourceOwner: email.ResourceOwner,
		},
		VerificationCode: email.PlainCode,
	}), nil
}

func (s *Server) ResendEmailCode(ctx context.Context, req *connect.Request[user.ResendEmailCodeRequest]) (resp *connect.Response[user.ResendEmailCodeResponse], err error) {
	var email *domain.Email

	switch v := req.Msg.GetVerification().(type) {
	case *user.ResendEmailCodeRequest_SendCode:
		email, err = s.command.ResendUserEmailCodeURLTemplate(ctx, req.Msg.GetUserId(), s.userCodeAlg, v.SendCode.GetUrlTemplate())
	case *user.ResendEmailCodeRequest_ReturnCode:
		email, err = s.command.ResendUserEmailReturnCode(ctx, req.Msg.GetUserId(), s.userCodeAlg)
	case nil:
		email, err = s.command.ResendUserEmailCode(ctx, req.Msg.GetUserId(), s.userCodeAlg)
	default:
		err = zerrors.ThrowUnimplementedf(nil, "USERv2-faj0l0nj5x", "verification oneOf %T in method ResendEmailCode not implemented", v)
	}
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&user.ResendEmailCodeResponse{
		Details: &object.Details{
			Sequence:      email.Sequence,
			ChangeDate:    timestamppb.New(email.ChangeDate),
			ResourceOwner: email.ResourceOwner,
		},
		VerificationCode: email.PlainCode,
	}), nil
}

func (s *Server) SendEmailCode(ctx context.Context, req *connect.Request[user.SendEmailCodeRequest]) (resp *connect.Response[user.SendEmailCodeResponse], err error) {
	var email *domain.Email

	switch v := req.Msg.GetVerification().(type) {
	case *user.SendEmailCodeRequest_SendCode:
		email, err = s.command.SendUserEmailCodeURLTemplate(ctx, req.Msg.GetUserId(), s.userCodeAlg, v.SendCode.GetUrlTemplate())
	case *user.SendEmailCodeRequest_ReturnCode:
		email, err = s.command.SendUserEmailReturnCode(ctx, req.Msg.GetUserId(), s.userCodeAlg)
	case nil:
		email, err = s.command.SendUserEmailCode(ctx, req.Msg.GetUserId(), s.userCodeAlg)
	default:
		err = zerrors.ThrowUnimplementedf(nil, "USERv2-faj0l0nj5x", "verification oneOf %T in method SendEmailCode not implemented", v)
	}
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&user.SendEmailCodeResponse{
		Details: &object.Details{
			Sequence:      email.Sequence,
			ChangeDate:    timestamppb.New(email.ChangeDate),
			ResourceOwner: email.ResourceOwner,
		},
		VerificationCode: email.PlainCode,
	}), nil
}

func (s *Server) VerifyEmail(ctx context.Context, req *connect.Request[user.VerifyEmailRequest]) (*connect.Response[user.VerifyEmailResponse], error) {
	details, err := s.command.VerifyUserEmail(ctx,
		req.Msg.GetUserId(),
		req.Msg.GetVerificationCode(),
		s.userCodeAlg,
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.VerifyEmailResponse{
		Details: &object.Details{
			Sequence:      details.Sequence,
			ChangeDate:    timestamppb.New(details.EventDate),
			ResourceOwner: details.ResourceOwner,
		},
	}), nil
}
