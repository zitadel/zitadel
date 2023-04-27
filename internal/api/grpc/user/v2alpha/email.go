package user

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2alpha"
)

func (s *Server) SetEmail(ctx context.Context, req *SetEmailRequest) (resp *SetEmailResponse, err error) {
	var resourceOwner string // TODO: check if still needed
	var email *domain.Email

	switch v := req.GetVerification().(type) {
	case *SetEmailRequest_SendCode:
		email, err = s.command.ChangeUserEmailURLTemplate(ctx, req.GetUserId(), resourceOwner, req.GetEmail(), s.userCodeAlg, v.SendCode.GetUrlTemplate())
	case *SetEmailRequest_ReturnCode:
		email, err = s.command.ChangeUserEmailReturnCode(ctx, req.GetUserId(), resourceOwner, req.GetEmail(), s.userCodeAlg)
	case *SetEmailRequest_IsVerified:
		if v.IsVerified {
			email, err = s.command.ChangeUserEmailVerified(ctx, req.GetUserId(), resourceOwner, req.GetEmail())
		} else {
			email, err = s.command.ChangeUserEmail(ctx, req.GetUserId(), resourceOwner, req.GetEmail(), s.userCodeAlg)
		}
	case nil:
		email, err = s.command.ChangeUserEmail(ctx, req.GetUserId(), resourceOwner, req.GetEmail(), s.userCodeAlg)
	default:
		err = caos_errs.ThrowUnimplementedf(nil, "USERv2-Ahng0", "verification oneOf %T in method SetEmail not implemented", v)
	}
	if err != nil {
		return nil, err
	}

	return &SetEmailResponse{
		Details: &object.Details{
			Sequence:      email.Sequence,
			ChangeDate:    timestamppb.New(email.ChangeDate),
			ResourceOwner: email.ResourceOwner,
		},
		VerificationCode: email.PlainCode,
	}, nil
}

func (s *Server) VerifyEmail(ctx context.Context, req *VerifyEmailRequest) (*VerifyEmailResponse, error) {
	details, err := s.command.VerifyUserEmail(ctx,
		req.GetUserId(),
		"", // TODO: check if still needed
		req.GetVerificationCode(),
		s.userCodeAlg,
	)
	if err != nil {
		return nil, err
	}
	return &VerifyEmailResponse{
		Details: &object.Details{
			Sequence:      details.Sequence,
			ChangeDate:    timestamppb.New(details.EventDate),
			ResourceOwner: details.ResourceOwner,
		},
	}, nil
}
