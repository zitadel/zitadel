package management

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetIam(ctx context.Context, _ *empty.Empty) (*management.Iam, error) {
	iam, err := s.project.GetIAMByID(ctx)
	if err != nil {
		return nil, err
	}
	return iamFromModel(iam), nil
}
