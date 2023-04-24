package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) ListViews(ctx context.Context, _ *admin_pb.ListViewsRequest) (*admin_pb.ListViewsResponse, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceIDQuery, err := query.NewCurrentStatesInstanceIDSearchQuery(instanceID)
	if err != nil {
		return nil, err
	}
	currentSequences, err := s.query.SearchCurrentStates(ctx, &query.CurrentStateSearchQueries{
		Queries: []query.SearchQuery{instanceIDQuery},
	})
	if err != nil {
		return nil, err
	}
	convertedCurrentSequences := CurrentSequencesToPb(s.database, currentSequences)
	views, err := s.administrator.GetViews(instanceID)
	if err != nil {
		return nil, err
	}
	convertedViews := ViewsToPb(views)

	convertedCurrentSequences = append(convertedCurrentSequences, convertedViews...)
	return &admin_pb.ListViewsResponse{Result: convertedCurrentSequences}, nil
}
