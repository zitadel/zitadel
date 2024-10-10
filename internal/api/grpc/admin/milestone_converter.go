package admin

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/milestone"
	"github.com/zitadel/zitadel/internal/zerrors"
	milestone_pb "github.com/zitadel/zitadel/pkg/grpc/milestone"
)

func milestoneQueriesReached(queries []*milestone_pb.MilestoneQuery) (reached bool, err error) {
	for _, query := range queries {
		switch q := query.Query.(type) {
		case *milestone_pb.MilestoneQuery_IsReachedQuery:
			reached = q.IsReachedQuery.GetReached()
		default:
			return false, zerrors.ThrowInvalidArgument(nil, "ADMIN-sE7pc", "List.Query.Invalid")
		}
	}
	return reached, nil
}

func milestoneViewsToPb(milestones []*query.Milestone) []*milestone_pb.Milestone {
	resp := make([]*milestone_pb.Milestone, len(milestones))
	for i, idp := range milestones {
		resp[i] = modelMilestoneViewToPb(idp)
	}
	return resp
}

func modelMilestoneViewToPb(m *query.Milestone) *milestone_pb.Milestone {
	mspb := &milestone_pb.Milestone{
		Type: modelMilestoneTypeToPb(m.Type),
	}
	if !m.ReachedDate.IsZero() {
		mspb.ReachedDate = timestamppb.New(m.ReachedDate)
	}
	return mspb
}

func modelMilestoneTypeToPb(t milestone.Type) milestone_pb.MilestoneType {
	switch t {
	case milestone.InstanceCreated:
		return milestone_pb.MilestoneType_MILESTONE_TYPE_INSTANCE_CREATED
	case milestone.AuthenticationSucceededOnInstance:
		return milestone_pb.MilestoneType_MILESTONE_TYPE_AUTHENTICATION_SUCCEEDED_ON_INSTANCE
	case milestone.ProjectCreated:
		return milestone_pb.MilestoneType_MILESTONE_TYPE_PROJECT_CREATED
	case milestone.ApplicationCreated:
		return milestone_pb.MilestoneType_MILESTONE_TYPE_APPLICATION_CREATED
	case milestone.AuthenticationSucceededOnApplication:
		return milestone_pb.MilestoneType_MILESTONE_TYPE_AUTHENTICATION_SUCCEEDED_ON_APPLICATION
	case milestone.InstanceDeleted:
		return milestone_pb.MilestoneType_MILESTONE_TYPE_INSTANCE_DELETED
	default:
		return milestone_pb.MilestoneType_MILESTONE_TYPE_UNSPECIFIED
	}
}
