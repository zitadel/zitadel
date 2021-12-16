package admin

import (
	member_grpc "github.com/caos/zitadel/internal/api/grpc/member"
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/query"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func AddIAMMemberToDomain(req *admin_pb.AddIAMMemberRequest) *domain.Member {
	return &domain.Member{
		UserID: req.UserId,
		Roles:  req.Roles,
	}
}

func UpdateIAMMemberToDomain(req *admin_pb.UpdateIAMMemberRequest) *domain.Member {
	return &domain.Member{
		UserID: req.UserId,
		Roles:  req.Roles,
	}
}

func ListIAMMembersRequestToQuery(req *admin_pb.ListIAMMembersRequest) (*query.IAMMembersQuery, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := member_grpc.MemberQueriesToQuery(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.IAMMembersQuery{
		MembersQuery: query.MembersQuery{
			SearchRequest: query.SearchRequest{
				Offset: offset,
				Limit:  limit,
				Asc:    asc,
				// SortingColumn: model.IAMMemberSearchKey, //TOOD: not implemented in proto
			},
			Queries: queries,
		},
	}, nil
}
