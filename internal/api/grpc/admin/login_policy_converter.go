package admin

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/iam/model"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func updateLoginPolicyToDomain(p *admin_pb.UpdateLoginPolicyRequest) *domain.LoginPolicy {
	return &domain.LoginPolicy{
		AllowUsernamePassword: p.AllowUsernamePassword,
		AllowRegister:         p.AllowRegister,
		AllowExternalIDP:      p.AllowExternalIdp,
		ForceMFA:              p.ForceMfa,
		PasswordlessType:      policy_grpc.PasswordlessTypeToDomain(p.PasswordlessType),
	}
}

func ListLoginPolicyIDPsRequestToModel(req *admin_pb.ListLoginPolicyIDPsRequest) *model.IDPProviderSearchRequest {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	return &model.IDPProviderSearchRequest{
		Offset: offset,
		Limit:  limit,
		Asc:    asc,
		// SortingColumn: model.IDPProviderSearchKey, //TODO: not in proto
		// Queries: []*model.IDPProviderSearchQuery, //TODO: not in proto
	}
}
