package management

import (
	"context"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetZitadelDocs(ctx context.Context, _ *empty.Empty) (*management.ZitadelDocs, error) {
	user, err := s.usergrant.UserGrantByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return userGrantViewFromModel(user), nil
}
