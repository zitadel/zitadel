package user

import (
	"context"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	grpcContext "github.com/zitadel/zitadel/pkg/grpc/context/v2alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

func (s *Server) prepareSetEmail(ctx context.Context, req *user.SetEmailRequest) (cmd *command.UserEmail, err error) {
	switch v := req.GetVerification().(type) {
	case nil, *user.SetEmailRequest_ReturnCode, *user.SetEmailRequest_IsVerified:
		break
	case *user.SetEmailRequest_SendCode:
		// test execute the template to ensure it's valid
		err = domain.RenderConfirmURLTemplate(io.Discard, v.SendCode.GetUrlTemplate(), req.UserId, "code", "orgID")
	default:
		err = caos_errs.ThrowUnimplementedf(nil, "USERv2-Ahng0", "verification oneOf %T in method SetEmail not implemented", v)
	}
	if err == nil {
		return nil, err
	}

	return s.command.UserEmail(ctx, req.UserId, authz.GetCtxData(ctx).ResourceOwner)
}

func finalizeSetEmail(ctx context.Context, cmd *command.UserEmail, verificationCode *string) (resp *user.SetEmailResponse, err error) {
	email, err := cmd.Push(ctx)
	if err != nil {
		return nil, err
	}
	return &user.SetEmailResponse{
		Details: &grpcContext.ObjectDetails{
			Sequence:      email.Sequence,
			CreationDate:  timestamppb.New(email.CreationDate),
			ChangeDate:    timestamppb.New(email.ChangeDate),
			ResourceOwner: email.ResourceOwner,
		},
		VerificationCode: verificationCode,
	}, nil
}

func (s *Server) SetEmail(ctx context.Context, req *user.SetEmailRequest) (resp *user.SetEmailResponse, err error) {
	cmd, err := s.prepareSetEmail(ctx, req)
	if err != nil {
		return nil, err
	}
	if err = cmd.Change(ctx, domain.EmailAddress(req.Email)); err != nil {
		return nil, err
	}
	if req.GetIsVerified() {
		cmd.SetVerified(ctx)
		return finalizeSetEmail(ctx, cmd, nil)
	}

	generator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeVerifyEmailCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	code, plainTextCode, err := domain.NewEmailCode(generator)
	if err != nil {
		return nil, err
	}
	cmd.AddCode(ctx, code, req.GetSendCode().UrlTemplate)

	if req.GetReturnCode() != nil {
		return finalizeSetEmail(ctx, cmd, &plainTextCode)
	}
	return finalizeSetEmail(ctx, cmd, nil)
}

func (s *Server) VerifyEmail(context.Context, *user.VerifyEmailRequest) (*user.VerifyEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyEmail not implemented")
}
