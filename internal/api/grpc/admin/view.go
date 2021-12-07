package admin

import (
	"context"

	"github.com/caos/zitadel/internal/query"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) ListViews(ctx context.Context, _ *admin_pb.ListViewsRequest) (*admin_pb.ListViewsResponse, error) {
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
	return &admin_pb.ListViewsResponse{Result: convertedCurrentSequences}, nil
}

func (s *Server) ClearView(ctx context.Context, req *admin_pb.ClearViewRequest) (*admin_pb.ClearViewResponse, error) {
	var err error
	if req.Database != "zitadel" {
		err = s.administrator.ClearView(ctx, req.Database, req.ViewName)
	} else {
		err = s.query.ClearCurrentSequence(ctx, req.ViewName)
	}
	if err != nil {
		return nil, err
	}
	return &admin_pb.ClearViewResponse{}, nil
}
