package auth

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/api/grpc/user"
	"github.com/zitadel/zitadel/internal/domain"
	auth_pb "github.com/zitadel/zitadel/pkg/grpc/auth"
)

func (s *Server) GetMyEmail(ctx context.Context, _ *auth_pb.GetMyEmailRequest) (*auth_pb.GetMyEmailResponse, error) {
	email, err := s.query.GetHumanEmail(ctx, authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	return &auth_pb.GetMyEmailResponse{
		Email: user.ModelEmailToPb(email),
		Details: object.ToViewDetailsPb(
			email.Sequence,
			email.CreationDate,
			email.ChangeDate,
			email.ResourceOwner,
		),
	}, nil
}

func (s *Server) SetMyEmail(ctx context.Context, req *auth_pb.SetMyEmailRequest) (*auth_pb.SetMyEmailResponse, error) {
	emailCodeGenerator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeVerifyEmailCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	email, err := s.command.ChangeHumanEmail(ctx, UpdateMyEmailToDomain(ctx, req), emailCodeGenerator)
	if err != nil {
		return nil, err
	}
	return &auth_pb.SetMyEmailResponse{
		Details: object.ChangeToDetailsPb(
			email.Sequence,
			email.ChangeDate,
			email.ResourceOwner,
		),
	}, nil
}

func (s *Server) VerifyMyEmail(ctx context.Context, req *auth_pb.VerifyMyEmailRequest) (*auth_pb.VerifyMyEmailResponse, error) {
	emailCodeGenerator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeVerifyEmailCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	ctxData := authz.GetCtxData(ctx)
	objectDetails, err := s.command.VerifyHumanEmail(ctx, ctxData.UserID, req.Code, ctxData.ResourceOwner, emailCodeGenerator)
	if err != nil {
		return nil, err
	}
	return &auth_pb.VerifyMyEmailResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ResendMyEmailVerification(ctx context.Context, _ *auth_pb.ResendMyEmailVerificationRequest) (*auth_pb.ResendMyEmailVerificationResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	emailCodeGenerator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeVerifyEmailCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	objectDetails, err := s.command.CreateHumanEmailVerificationCode(ctx, ctxData.UserID, ctxData.ResourceOwner, emailCodeGenerator, "")
	if err != nil {
		return nil, err
	}
	return &auth_pb.ResendMyEmailVerificationResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}
