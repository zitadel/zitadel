package instance

import (
	"context"

	"github.com/zitadel/zitadel/pkg/grpc/instance/v2"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) GetInstance(ctx context.Context, _ *emptypb.Empty) (*instance.GetInstanceResponse, error) {
	inst, err := s.query.Instance(ctx, true)
	if err != nil {
		return nil, err
	}

	return &instance.GetInstanceResponse{
		Instance: ToProtoObject(inst),
	}, nil
}
