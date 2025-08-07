package system

import (
	"context"

	"github.com/zitadel/zitadel/internal/query"
	system_pb "github.com/zitadel/zitadel/pkg/grpc/system"
)

func (s *Server) ListViews(ctx context.Context, _ *system_pb.ListViewsRequest) (*system_pb.ListViewsResponse, error) {
	currentSequences, err := s.query.SearchCurrentStates(ctx, new(query.CurrentStateSearchQueries))
	if err != nil {
		return nil, err
	}
	return &system_pb.ListViewsResponse{Result: CurrentSequencesToPb(s.database, currentSequences)}, nil
}

func (s *Server) ClearView(ctx context.Context, req *system_pb.ClearViewRequest) (*system_pb.ClearViewResponse, error) {
	err := s.query.ClearCurrentSequence(ctx, req.ViewName)
	if err != nil {
		return nil, err
	}
	return &system_pb.ClearViewResponse{}, nil
}
