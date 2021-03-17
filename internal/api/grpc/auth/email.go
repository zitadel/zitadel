package auth

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/api/grpc/user"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
)

func (s *Server) GetMyEmail(ctx context.Context, _ *auth_pb.GetMyEmailRequest) (*auth_pb.GetMyEmailResponse, error) {
	email, err := s.repo.MyEmail(ctx)
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
	email, err := s.command.ChangeHumanEmail(ctx, UpdateMyEmailToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &auth_pb.SetMyEmailResponse{
		Details: object.ToDetailsPb(
			email.Sequence,
			email.ChangeDate,
			email.ResourceOwner,
		),
	}, nil
}

func (s *Server) VerifyMyEmail(ctx context.Context, req *auth_pb.VerifyMyEmailRequest) (*auth_pb.VerifyMyEmailResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	objectDetails, err := s.command.VerifyHumanEmail(ctx, ctxData.UserID, req.Code, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.VerifyMyEmailResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ResendMyEmailVerification(ctx context.Context, _ *auth_pb.ResendMyEmailVerificationRequest) (*auth_pb.ResendMyEmailVerificationResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	objectDetails, err := s.command.CreateHumanEmailVerificationCode(ctx, ctxData.UserID, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ResendMyEmailVerificationResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}
