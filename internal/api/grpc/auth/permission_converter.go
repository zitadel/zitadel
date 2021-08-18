package auth

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	user_grpc "github.com/caos/zitadel/internal/api/grpc/user"
	user_model "github.com/caos/zitadel/internal/user/model"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
)

func ListMyMembershipsRequestToModel(req *auth_pb.ListMyMembershipsRequest) (*user_model.UserMembershipSearchRequest, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := user_grpc.MembershipQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &user_model.UserMembershipSearchRequest{
		Offset: offset,
		Limit:  limit,
		Asc:    asc,
		//SortingColumn: //TODO: sorting
		Queries: queries,
	}, nil
}
