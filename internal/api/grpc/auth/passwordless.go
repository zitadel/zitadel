package auth

import (
	"context"

	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/object"
	user_grpc "github.com/caos/zitadel/internal/api/grpc/user"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
	user_pb "github.com/caos/zitadel/pkg/grpc/user"
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
	token, err := s.command.HumanAddPasswordlessSetup(ctx, ctxData.UserID, ctxData.ResourceOwner, false)
	if err != nil {
		return nil, err
	}
	return &auth_pb.AddMyPasswordlessResponse{
		Key: &user_pb.WebAuthNKey{
			PublicKey: token.CredentialCreationData,
		},
		Details: object.AddToDetailsPb(
			token.Sequence,
			token.ChangeDate,
			token.ResourceOwner,
		),
	}, nil
}

func (s *Server) AddMyPasswordlessLink(ctx context.Context, req *auth_pb.AddMyPasswordlessLinkRequest) (*auth_pb.AddMyPasswordlessLinkResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	initCode, err := s.command.HumanAddPasswordlessInitCode(ctx, ctxData.UserID, ctxData.ResourceOwner, req.Send)
	if err != nil {
		return nil, err
	}
	var linkAdded auth_pb.LinkAdded
	if initCode.Active {
		linkAdded = &auth_pb.AddMyPasswordlessLinkResponse_Added{
			Added: &auth_pb.AddMyPasswordlessLinkResponse_Link{
				CodeId:     initCode.CodeID,
				Code:       initCode.Code,
				Link:       initCode.Code,
				Expiration: durationpb.New(initCode.Expiration),
			},
		}
	} else {
		linkAdded = &auth_pb.AddMyPasswordlessLinkResponse_Send{
			Send: true,
		}
	}

	return &auth_pb.AddMyPasswordlessLinkResponse{
		Details:   object.AddToDetailsPb(initCode.Sequence, initCode.ChangeDate, initCode.ResourceOwner),
		LinkAdded: linkAdded,
	}, nil
}

func (s *Server) VerifyMyPasswordless(ctx context.Context, req *auth_pb.VerifyMyPasswordlessRequest) (*auth_pb.VerifyMyPasswordlessResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	objectDetails, err := s.command.HumanHumanPasswordlessSetup(ctx, ctxData.UserID, ctxData.ResourceOwner, req.Verification.TokenName, "", req.Verification.PublicKeyCredential)
	if err != nil {
		return nil, err
	}
	return &auth_pb.VerifyMyPasswordlessResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) RemoveMyPasswordless(ctx context.Context, req *auth_pb.RemoveMyPasswordlessRequest) (*auth_pb.RemoveMyPasswordlessResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	objectDetails, err := s.command.HumanRemovePasswordless(ctx, ctxData.UserID, req.TokenId, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.RemoveMyPasswordlessResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}
