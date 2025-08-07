package auth

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	auth_pb "github.com/zitadel/zitadel/pkg/grpc/auth"
)

func (s *Server) RemoveMyAvatar(ctx context.Context, req *auth_pb.RemoveMyAvatarRequest) (*auth_pb.RemoveMyAvatarResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	objectDetails, err := s.command.RemoveHumanAvatar(ctx, ctxData.ResourceOwner, ctxData.UserID)
	if err != nil {
		return nil, err
	}
	return &auth_pb.RemoveMyAvatarResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}
