package auth

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
)

func (s *Server) UpdateMyPassword(ctx context.Context, req *auth_pb.UpdateMyPasswordRequest) (*auth_pb.UpdateMyPasswordResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	err := s.command.ChangePassword(ctx, ctxData.OrgID, ctxData.UserID, req.OldPassword, req.NewPassword, "")
	if err != nil {
		return nil, err
	}
	//TODO: returns values
	return &auth_pb.UpdateMyPasswordResponse{}, nil
}
