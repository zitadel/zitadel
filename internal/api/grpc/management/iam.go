package management

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetIam(ctx context.Context, _ *empty.Empty) (*management.Iam, error) {
	iam, err := s.iam.IAMByID(ctx, s.systemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	return iamFromModel(iam), nil
}
