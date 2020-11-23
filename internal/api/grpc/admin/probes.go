package admin

import (
	"context"
	"github.com/caos/zitadel/internal/telemetry/metrics"

	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) Healthz(_ context.Context, e *empty.Empty) (*empty.Empty, error) {
	metrics.M.PrintCounters()
	return &empty.Empty{}, nil
}
