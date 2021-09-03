package management

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func createActionRequestToDomain(req *mgmt_pb.CreateActionRequest) *domain.Action {
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

func listActionsToQuery(req *mgmt_pb.ListActionsRequest) *interface{} {
	//offset, limit, asc := object.ListQueryToModel(req.Query)
	//return &iam_model.IDPConfigSearchRequest{
	//	Offset:        offset,
	//	Limit:         limit,
	//	Asc:           asc,
	//	SortingColumn: action_grpc.ActionFieldNameToDomain(req.SortingColumn),
	//	Queries:       actionQueries(req.Queries),
	//}
	return nil
}

//func actionQueries(queries []*mgmt_pb.ActionQuery) {
//	for i, query := range queries {
//		switch query.Query {
//		case *action_pb.ActionIDQuery:
//
//		}
//	}
//}
