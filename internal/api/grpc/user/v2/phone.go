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

func (s *Server) SetPhone(ctx context.Context, req *connect.Request[user.SetPhoneRequest]) (resp *connect.Response[user.SetPhoneResponse], err error) {
	var phone *domain.Phone

	switch v := req.Msg.GetVerification().(type) {
	case *user.SetPhoneRequest_SendCode:
		phone, err = s.command.ChangeUserPhone(ctx, req.Msg.GetUserId(), req.Msg.GetPhone(), s.userCodeAlg)
	case *user.SetPhoneRequest_ReturnCode:
		phone, err = s.command.ChangeUserPhoneReturnCode(ctx, req.Msg.GetUserId(), req.Msg.GetPhone(), s.userCodeAlg)
	case *user.SetPhoneRequest_IsVerified:
		phone, err = s.command.ChangeUserPhoneVerified(ctx, req.Msg.GetUserId(), req.Msg.GetPhone())
	case nil:
		phone, err = s.command.ChangeUserPhone(ctx, req.Msg.GetUserId(), req.Msg.GetPhone(), s.userCodeAlg)
	default:
		err = zerrors.ThrowUnimplementedf(nil, "USERv2-Ahng0", "verification oneOf %T in method SetPhone not implemented", v)
	}
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&user.SetPhoneResponse{
		Details: &object.Details{
			Sequence:      phone.Sequence,
			ChangeDate:    timestamppb.New(phone.ChangeDate),
			ResourceOwner: phone.ResourceOwner,
		},
		VerificationCode: phone.PlainCode,
	}), nil
}

func (s *Server) RemovePhone(ctx context.Context, req *connect.Request[user.RemovePhoneRequest]) (resp *connect.Response[user.RemovePhoneResponse], err error) {
	details, err := s.command.RemoveUserPhone(ctx,
		req.Msg.GetUserId(),
	)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&user.RemovePhoneResponse{
		Details: &object.Details{
			Sequence:      details.Sequence,
			ChangeDate:    timestamppb.New(details.EventDate),
			ResourceOwner: details.ResourceOwner,
		},
	}), nil
}

func (s *Server) ResendPhoneCode(ctx context.Context, req *connect.Request[user.ResendPhoneCodeRequest]) (resp *connect.Response[user.ResendPhoneCodeResponse], err error) {
	var phone *domain.Phone
	switch v := req.Msg.GetVerification().(type) {
	case *user.ResendPhoneCodeRequest_SendCode:
		phone, err = s.command.ResendUserPhoneCode(ctx, req.Msg.GetUserId(), s.userCodeAlg)
	case *user.ResendPhoneCodeRequest_ReturnCode:
		phone, err = s.command.ResendUserPhoneCodeReturnCode(ctx, req.Msg.GetUserId(), s.userCodeAlg)
	case nil:
		phone, err = s.command.ResendUserPhoneCode(ctx, req.Msg.GetUserId(), s.userCodeAlg)
	default:
		err = zerrors.ThrowUnimplementedf(nil, "USERv2-ResendUserPhoneCode", "verification oneOf %T in method SetPhone not implemented", v)
	}
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&user.ResendPhoneCodeResponse{
		Details: &object.Details{
			Sequence:      phone.Sequence,
			ChangeDate:    timestamppb.New(phone.ChangeDate),
			ResourceOwner: phone.ResourceOwner,
		},
		VerificationCode: phone.PlainCode,
	}), nil
}

func (s *Server) VerifyPhone(ctx context.Context, req *connect.Request[user.VerifyPhoneRequest]) (*connect.Response[user.VerifyPhoneResponse], error) {
	details, err := s.command.VerifyUserPhone(ctx,
		req.Msg.GetUserId(),
		req.Msg.GetVerificationCode(),
		s.userCodeAlg,
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.VerifyPhoneResponse{
		Details: &object.Details{
			Sequence:      details.Sequence,
			ChangeDate:    timestamppb.New(details.EventDate),
			ResourceOwner: details.ResourceOwner,
		},
	}), nil
}
