package admin

import (
	"context"
	"slices"

	object_pb "github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) ListMilestones(ctx context.Context, req *admin.ListMilestonesRequest) (*admin.ListMilestonesResponse, error) {
	reached, err := milestoneQueriesReached(req.GetQueries())
	if err != nil {
		return nil, err
	}
	model, err := s.query.GetMilestones(ctx)
	if err != nil {
		return nil, err
	}
	if reached {
		model.Milestones = slices.DeleteFunc(model.Milestones, func(milestone *query.Milestone) bool {
			return milestone.ReachedDate.IsZero()
		})
	}
	return &admin.ListMilestonesResponse{
		Result:  milestoneViewsToPb(model.Milestones),
		Details: object_pb.ToListDetails(uint64(len(model.Milestones)), model.ProcessedSequence, model.ChangeDate),
	}, nil
}
