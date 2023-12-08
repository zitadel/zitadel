package management

import (
	action_grpc "github.com/zitadel/zitadel/internal/api/grpc/action"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func CreateActionRequestToDomain(req *mgmt_pb.CreateActionRequest) *domain.Action {
	return &domain.Action{
		Name:          req.Name,
		Script:        req.Script,
		Timeout:       req.Timeout.AsDuration(),
		AllowedToFail: req.AllowedToFail,
	}
}

func updateActionRequestToDomain(req *mgmt_pb.UpdateActionRequest) *domain.Action {
	return &domain.Action{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.Id,
		},
		Name:          req.Name,
		Script:        req.Script,
		Timeout:       req.Timeout.AsDuration(),
		AllowedToFail: req.AllowedToFail,
	}
}

func listActionsToQuery(orgID string, req *mgmt_pb.ListActionsRequest) (_ *query.ActionSearchQueries, err error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries := make([]query.SearchQuery, len(req.Queries)+1)
	queries[0], err = query.NewActionResourceOwnerQuery(orgID)
	if err != nil {
		return nil, err
	}
	for i, actionQuery := range req.Queries {
		queries[i+1], err = ActionQueryToQuery(actionQuery.Query)
		if err != nil {
			return nil, err
		}
	}
	return &query.ActionSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: queries,
	}, nil
}

func ActionQueryToQuery(query interface{}) (query.SearchQuery, error) {
	switch q := query.(type) {
	case *mgmt_pb.ActionQuery_ActionNameQuery:
		return action_grpc.ActionNameQuery(q.ActionNameQuery)
	case *mgmt_pb.ActionQuery_ActionStateQuery:
		return action_grpc.ActionStateQuery(q.ActionStateQuery)
	case *mgmt_pb.ActionQuery_ActionIdQuery:
		return action_grpc.ActionIDQuery(q.ActionIdQuery)
	}
	return nil, zerrors.ThrowInvalidArgument(nil, "MGMT-dsg3z", "Errors.Query.InvalidRequest")
}
