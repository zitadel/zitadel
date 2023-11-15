package auth

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/api/grpc/object"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	auth_pb "github.com/zitadel/zitadel/v2/pkg/grpc/auth"
)

func (s *Server) UpdateMyPassword(ctx context.Context, req *auth_pb.UpdateMyPasswordRequest) (*auth_pb.UpdateMyPasswordResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	objectDetails, err := s.command.ChangePassword(ctx, ctxData.ResourceOwner, ctxData.UserID, req.OldPassword, req.NewPassword, "")
	if err != nil {
		return nil, err
	}
	return &auth_pb.UpdateMyPasswordResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}
