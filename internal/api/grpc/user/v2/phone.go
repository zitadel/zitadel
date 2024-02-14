package user

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

func (s *Server) SetPhone(ctx context.Context, req *user.SetPhoneRequest) (resp *user.SetPhoneResponse, err error) {
	var resourceOwner string // TODO: check if still needed
	var phone *domain.Phone

	switch v := req.GetVerification().(type) {
	case *user.SetPhoneRequest_SendCode:
		phone, err = s.command.ChangeUserPhone(ctx, req.GetUserId(), resourceOwner, req.GetPhone(), s.userCodeAlg)
	case *user.SetPhoneRequest_ReturnCode:
		phone, err = s.command.ChangeUserPhoneReturnCode(ctx, req.GetUserId(), resourceOwner, req.GetPhone(), s.userCodeAlg)
	case *user.SetPhoneRequest_IsVerified:
		phone, err = s.command.ChangeUserPhoneVerified(ctx, req.GetUserId(), resourceOwner, req.GetPhone())
	case nil:
		phone, err = s.command.ChangeUserPhone(ctx, req.GetUserId(), resourceOwner, req.GetPhone(), s.userCodeAlg)
	default:
		err = zerrors.ThrowUnimplementedf(nil, "USERv2-Ahng0", "verification oneOf %T in method SetPhone not implemented", v)
	}
	if err != nil {
		return nil, err
	}

	return &user.SetPhoneResponse{
		Details: &object.Details{
			Sequence:      phone.Sequence,
			ChangeDate:    timestamppb.New(phone.ChangeDate),
			ResourceOwner: phone.ResourceOwner,
		},
		VerificationCode: phone.PlainCode,
	}, nil
}

func (s *Server) VerifyPhone(ctx context.Context, req *user.VerifyPhoneRequest) (*user.VerifyPhoneResponse, error) {
	details, err := s.command.VerifyUserPhone(ctx,
		req.GetUserId(),
		"", // TODO: check if still needed
		req.GetVerificationCode(),
		s.userCodeAlg,
	)
	if err != nil {
		return nil, err
	}
	return &user.VerifyPhoneResponse{
		Details: &object.Details{
			Sequence:      details.Sequence,
			ChangeDate:    timestamppb.New(details.EventDate),
			ResourceOwner: details.ResourceOwner,
		},
	}, nil
}
