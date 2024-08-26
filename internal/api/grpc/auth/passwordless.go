package auth

import (
	"context"

	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	user_grpc "github.com/zitadel/zitadel/internal/api/grpc/user"
	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	auth_pb "github.com/zitadel/zitadel/pkg/grpc/auth"
	user_pb "github.com/zitadel/zitadel/pkg/grpc/user"
)

func (s *Server) ListMyPasswordless(ctx context.Context, _ *auth_pb.ListMyPasswordlessRequest) (*auth_pb.ListMyPasswordlessResponse, error) {
	query := new(query.UserAuthMethodSearchQueries)
	err := query.AppendUserIDQuery(authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	err = query.AppendAuthMethodQuery(domain.UserAuthMethodTypePasswordless)
	if err != nil {
		return nil, err
	}
	err = query.AppendStateQuery(domain.MFAStateReady)
	if err != nil {
		return nil, err
	}
	authMethods, err := s.query.SearchUserAuthMethods(ctx, query, nil)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyPasswordlessResponse{
		Result: user_grpc.UserAuthMethodsToWebAuthNTokenPb(authMethods),
	}, nil
}

func (s *Server) AddMyPasswordless(ctx context.Context, _ *auth_pb.AddMyPasswordlessRequest) (*auth_pb.AddMyPasswordlessResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	token, err := s.command.HumanAddPasswordlessSetup(ctx, ctxData.UserID, ctxData.ResourceOwner, domain.AuthenticatorAttachmentUnspecified)
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

func (s *Server) AddMyPasswordlessLink(ctx context.Context, _ *auth_pb.AddMyPasswordlessLinkRequest) (*auth_pb.AddMyPasswordlessLinkResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	passwordlessInitCode, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypePasswordlessInitCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	initCode, err := s.command.HumanAddPasswordlessInitCode(ctx, ctxData.UserID, ctxData.ResourceOwner, passwordlessInitCode)
	if err != nil {
		return nil, err
	}
	return &auth_pb.AddMyPasswordlessLinkResponse{
		Details:    object.AddToDetailsPb(initCode.Sequence, initCode.ChangeDate, initCode.ResourceOwner),
		Link:       initCode.Link(http.DomainContext(ctx).Origin() + login.HandlerPrefix + login.EndpointPasswordlessRegistration),
		Expiration: durationpb.New(initCode.Expiration),
	}, nil
}

func (s *Server) SendMyPasswordlessLink(ctx context.Context, _ *auth_pb.SendMyPasswordlessLinkRequest) (*auth_pb.SendMyPasswordlessLinkResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	passwordlessInitCode, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypePasswordlessInitCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	initCode, err := s.command.HumanSendPasswordlessInitCode(ctx, ctxData.UserID, ctxData.ResourceOwner, passwordlessInitCode)
	if err != nil {
		return nil, err
	}
	return &auth_pb.SendMyPasswordlessLinkResponse{
		Details: object.AddToDetailsPb(initCode.Sequence, initCode.ChangeDate, initCode.ResourceOwner),
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
