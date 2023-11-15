package admin

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	object_pb "github.com/zitadel/zitadel/v2/internal/api/grpc/object"
	"github.com/zitadel/zitadel/v2/pkg/grpc/admin"
)

func (s *Server) ListMilestones(ctx context.Context, req *admin.ListMilestonesRequest) (*admin.ListMilestonesResponse, error) {
	queries, err := listMilestonesToModel(authz.GetInstance(ctx).InstanceID(), req)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.SearchMilestones(ctx, []string{authz.GetInstance(ctx).InstanceID()}, queries)
	if err != nil {
		return nil, err
	}
	return &admin.ListMilestonesResponse{
		Result:  milestoneViewsToPb(resp.Milestones),
		Details: object_pb.ToListDetails(resp.Count, resp.Sequence, resp.LastRun),
	}, nil
}
