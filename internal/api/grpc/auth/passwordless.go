package auth

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/object"
	user_grpc "github.com/caos/zitadel/internal/api/grpc/user"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
)

func (s *Server) ListMyPasswordless(ctx context.Context, _ *auth_pb.ListMyPasswordlessRequest) (*auth_pb.ListMyPasswordlessResponse, error) {
	tokens, err := s.repo.GetMyPasswordless(ctx)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyPasswordlessResponse{
		Result: user_grpc.WebAuthNTokensViewToPb(tokens),
	}, nil
}

func (s *Server) AddMyPasswordless(ctx context.Context, _ *auth_pb.AddMyPasswordlessRequest) (*auth_pb.AddMyPasswordlessResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	u2f, err := s.command.HumanAddPasswordlessSetup(ctx, ctxData.UserID, ctxData.ResourceOwner, false)
	if err != nil {
		return nil, err
	}
	return &auth_pb.AddMyPasswordlessResponse{
		Key: user_grpc.WebAuthNTokenToWebAuthNKeyPb(u2f),
		Details: object.ToDetailsPb(
			u2f.Sequence,
			u2f.CreationDate,
			u2f.ChangeDate,
			u2f.ResourceOwner,
		),
	}, nil
}

func (s *Server) VerifyMyPasswordless(ctx context.Context, req *auth_pb.VerifyMyPasswordlessRequest) (*auth_pb.VerifyMyPasswordlessResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	err := s.command.HumanHumanPasswordlessSetup(ctx, ctxData.UserID, ctxData.OrgID, req.Verification.TokenName, "", req.Verification.PublicKeyCredential)
	if err != nil {
		return nil, err
	}
	//TODO: response from business
	return &auth_pb.VerifyMyPasswordlessResponse{}, nil
}

func (s *Server) RemoveMyPasswordless(ctx context.Context, req *auth_pb.RemoveMyPasswordlessRequest) (*auth_pb.RemoveMyPasswordlessResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	err := s.command.HumanRemovePasswordless(ctx, ctxData.UserID, req.TokenId, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	//TODO: response from business
	return &auth_pb.RemoveMyPasswordlessResponse{}, nil
}
