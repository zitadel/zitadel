package admin

import (
	member_grpc "github.com/caos/zitadel/internal/api/grpc/member"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/iam/model"
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

func ListIAMMemberRequestToModel(req *admin_pb.ListIAMMembersRequest) *model.IAMMemberSearchRequest {
	return &model.IAMMemberSearchRequest{
		Offset: req.Query.Offset,
		Limit:  uint64(req.Query.Limit),
		Asc:    req.Query.Asc,
		// SortingColumn: model.IAMMemberSearchKey, //TOOD: not implemented in proto
		Queries: member_grpc.MemberQueriesToIAMMember(req.Queries),
	}
}
