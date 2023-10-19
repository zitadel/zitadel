package admin

import (
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/milestone"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	milestone_pb "github.com/zitadel/zitadel/pkg/grpc/milestone"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func listMilestonesToModel(instanceID string, req *admin_pb.ListMilestonesRequest) (*query.MilestonesSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := milestoneQueriesToModel(req.GetQueries())
	instanceIDQuery, err := query.NewTextQuery(query.MilestoneInstanceIDColID, instanceID, query.TextEquals)
	if err != nil {
		return nil, err
	}
	queries = append(queries, instanceIDQuery)
	return &query.MilestonesSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: milestoneFieldNameToSortingColumn(req.SortingColumn),
		},
		Queries: queries,
	}, nil
}

func milestoneQueriesToModel(queries []*milestone_pb.MilestoneQuery) (q []query.SearchQuery, err error) {
	q = make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = milestoneQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func milestoneQueryToModel(milestoneQuery *milestone_pb.MilestoneQuery) (query.SearchQuery, error) {
	switch q := milestoneQuery.Query.(type) {
	case *milestone_pb.MilestoneQuery_IsReachedQuery:
		if q.IsReachedQuery.GetReached() {
			return query.NewNotNullQuery(query.MilestoneReachedDateColID)
		}
		return query.NewIsNullQuery(query.MilestoneReachedDateColID)
	default:
		return nil, errors.ThrowInvalidArgument(nil, "ADMIN-sE7pc", "List.Query.Invalid")
	}
}

func milestoneFieldNameToSortingColumn(field milestone_pb.MilestoneFieldName) query.Column {
	switch field {
	case milestone_pb.MilestoneFieldName_MILESTONE_FIELD_NAME_REACHED_DATE:
		return query.MilestoneReachedDateColID
	default:
		return query.MilestoneTypeColID
	}
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
