package management

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/object"
	user_grpc "github.com/caos/zitadel/internal/api/grpc/user"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/usergrant/model"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func ListUserGrantsRequestToModel(ctx context.Context, req *mgmt_pb.ListUserGrantRequest) *model.UserGrantSearchRequest {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	request := &model.UserGrantSearchRequest{
		Offset:  offset,
		Limit:   limit,
		Asc:     asc,
		Queries: user_grpc.UserGrantQueriesToModel(req.Queries),
	}
	request.Queries = append(request.Queries, &model.UserGrantSearchQuery{
		Key:    model.UserGrantSearchKeyResourceOwner,
		Method: domain.SearchMethodEquals,
		Value:  authz.GetCtxData(ctx).OrgID,
	})
	return request
}

func AddUserGrantRequestToDomain(req *mgmt_pb.AddUserGrantRequest) *domain.UserGrant {
	return &domain.UserGrant{
		UserID:         req.UserId,
		ProjectID:      req.ProjectId,
		ProjectGrantID: req.ProjectGrantId,
		RoleKeys:       req.RoleKeys,
	}
}

func UpdateUserGrantRequestToDomain(req *mgmt_pb.UpdateUserGrantRequest) *domain.UserGrant {
	return &domain.UserGrant{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.GrantId,
		},
		UserID:   req.UserId,
		RoleKeys: req.RoleKeys,
	}

}
