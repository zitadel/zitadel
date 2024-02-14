package user

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

func (s *Server) SetEmail(ctx context.Context, req *user.SetEmailRequest) (resp *user.SetEmailResponse, err error) {
	var resourceOwner string // TODO: check if still needed
	var email *domain.Email

	switch v := req.GetVerification().(type) {
	case *user.SetEmailRequest_SendCode:
		email, err = s.command.ChangeUserEmailURLTemplate(ctx, req.GetUserId(), resourceOwner, req.GetEmail(), s.userCodeAlg, v.SendCode.GetUrlTemplate())
	case *user.SetEmailRequest_ReturnCode:
		email, err = s.command.ChangeUserEmailReturnCode(ctx, req.GetUserId(), resourceOwner, req.GetEmail(), s.userCodeAlg)
	case *user.SetEmailRequest_IsVerified:
		email, err = s.command.ChangeUserEmailVerified(ctx, req.GetUserId(), resourceOwner, req.GetEmail())
	case nil:
		email, err = s.command.ChangeUserEmail(ctx, req.GetUserId(), resourceOwner, req.GetEmail(), s.userCodeAlg)
	default:
		err = zerrors.ThrowUnimplementedf(nil, "USERv2-Ahng0", "verification oneOf %T in method SetEmail not implemented", v)
	}
	if err != nil {
		return nil, err
	}

	return &user.SetEmailResponse{
		Details: &object.Details{
			Sequence:      email.Sequence,
			ChangeDate:    timestamppb.New(email.ChangeDate),
			ResourceOwner: email.ResourceOwner,
		},
		VerificationCode: email.PlainCode,
	}, nil
}

func (s *Server) VerifyEmail(ctx context.Context, req *user.VerifyEmailRequest) (*user.VerifyEmailResponse, error) {
	details, err := s.command.VerifyUserEmail(ctx,
		req.GetUserId(),
		"", // TODO: check if still needed
		req.GetVerificationCode(),
		s.userCodeAlg,
	)
	if err != nil {
		return nil, err
	}
	return &user.VerifyEmailResponse{
		Details: &object.Details{
			Sequence:      details.Sequence,
			ChangeDate:    timestamppb.New(details.EventDate),
			ResourceOwner: details.ResourceOwner,
		},
	}, nil
}
