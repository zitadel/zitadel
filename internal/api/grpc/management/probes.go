package management

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	"github.com/golang/protobuf/ptypes/empty"
	pb_struct "github.com/golang/protobuf/ptypes/struct"
)

func (s *Server) Healthz(_ context.Context, e *empty.Empty) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-mkd3y", "Not implemented")
}

func (s *Server) Ready(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-pl6BM", "Not implemented")
}

func (s *Server) Validate(ctx context.Context, _ *empty.Empty) (*pb_struct.Struct, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-2wxF", "Not implemented")
}
