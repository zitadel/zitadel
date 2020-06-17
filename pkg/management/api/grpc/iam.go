package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetIam(ctx context.Context, _ *empty.Empty) (*Iam, error) {
	iam, err := s.iam.IamByID(ctx, s.systemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	return iamFromModel(iam), nil
}
