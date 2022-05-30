package system

import (
	"context"

	"github.com/zitadel/zitadel/internal/query"
	system_pb "github.com/zitadel/zitadel/pkg/grpc/system"
)

func (s *Server) ListViews(ctx context.Context, _ *system_pb.ListViewsRequest) (*system_pb.ListViewsResponse, error) {
	currentSequences, err := s.query.SearchCurrentSequences(ctx, new(query.CurrentSequencesSearchQueries))
	if err != nil {
		return nil, err
	}
	convertedCurrentSequences := CurrentSequencesToPb(currentSequences)
	views, err := s.administrator.GetViews()
	if err != nil {
		return nil, err
	}
	convertedViews := ViewsToPb(views)

	convertedCurrentSequences = append(convertedCurrentSequences, convertedViews...)
	return &system_pb.ListViewsResponse{Result: convertedCurrentSequences}, nil
}

func (s *Server) ClearView(ctx context.Context, req *system_pb.ClearViewRequest) (*system_pb.ClearViewResponse, error) {
	var err error
	if req.Database != s.database {
		err = s.administrator.ClearView(ctx, req.Database, req.ViewName)
	} else {
		err = s.query.ClearCurrentSequence(ctx, req.ViewName)
	}
	if err != nil {
		return nil, err
	}
	return &system_pb.ClearViewResponse{}, nil
}
