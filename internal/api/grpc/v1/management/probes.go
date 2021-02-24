package management

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) Healthz(_ context.Context, e *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
