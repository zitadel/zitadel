package admin

import (
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
	return &model.IDPProviderSearchRequest{
		Offset: req.MetaData.Offset,
		Limit:  uint64(req.MetaData.Limit),
		Asc:    req.MetaData.Asc,
		// SortingColumn: model.IDPProviderSearchKey, //TODO: not in proto
		// Queries: []*model.IDPProviderSearchQuery, //TODO: not in proto
	}
}
