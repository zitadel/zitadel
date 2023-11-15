package admin

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	"github.com/zitadel/zitadel/v2/internal/query"
	admin_pb "github.com/zitadel/zitadel/v2/pkg/grpc/admin"
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
	return &admin_pb.ListViewsResponse{Result: CurrentSequencesToPb(s.database, currentSequences)}, nil
}
