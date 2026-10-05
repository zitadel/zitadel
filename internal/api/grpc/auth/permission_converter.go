package auth

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	user_grpc "github.com/zitadel/zitadel/internal/api/grpc/user"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/query"
	auth_pb "github.com/zitadel/zitadel/pkg/grpc/auth"
)

func ListMyMembershipsRequestToModel(defaults systemdefaults.SystemDefaults, ctx context.Context, req *auth_pb.ListMyMembershipsRequest) (*query.MembershipSearchQuery, error) {
	offset, limit, asc, err := object.ListQueryToModel(defaults, req.Query)
	if err != nil {
		return nil, err
	}
	queries, err := user_grpc.MembershipQueriesToQuery(req.Queries)
	if err != nil {
		return nil, err
	}
	userQuery, err := query.NewMembershipUserIDQuery(authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	queries = append(queries, userQuery)
	return &query.MembershipSearchQuery{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
			//SortingColumn: //TODO: sorting
		},
		Queries: queries,
	}, nil
}
