package user

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

func (s *Server) SetEmail(ctx context.Context, req *user.SetEmailRequest) (*user.SetEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetEmail not implemented")
}
