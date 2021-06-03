package management

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/iam/model"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func addLoginPolicyToDomain(p *mgmt_pb.AddCustomLoginPolicyRequest) *domain.LoginPolicy {
	return &domain.LoginPolicy{
		AllowUsernamePassword: p.AllowUsernamePassword,
		AllowRegister:         p.AllowRegister,
		AllowExternalIDP:      p.AllowExternalIdp,
		ForceMFA:              p.ForceMfa,
		PasswordlessType:      policy_grpc.PasswordlessTypeToDomain(p.PasswordlessType),
		HidePasswordReset:     p.HidePasswordReset,
	}
}

func updateLoginPolicyToDomain(p *mgmt_pb.UpdateCustomLoginPolicyRequest) *domain.LoginPolicy {
	return &domain.LoginPolicy{
		AllowUsernamePassword: p.AllowUsernamePassword,
		AllowRegister:         p.AllowRegister,
		AllowExternalIDP:      p.AllowExternalIdp,
		ForceMFA:              p.ForceMfa,
		PasswordlessType:      policy_grpc.PasswordlessTypeToDomain(p.PasswordlessType),
		HidePasswordReset:     p.HidePasswordReset,
	}
}

func ListLoginPolicyIDPsRequestToModel(req *mgmt_pb.ListLoginPolicyIDPsRequest) *model.IDPProviderSearchRequest {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	return &model.IDPProviderSearchRequest{
		Offset: offset,
		Limit:  limit,
		Asc:    asc,
		// SortingColumn: model.IDPProviderSearchKey, //TODO: not in proto
		// Queries: []*model.IDPProviderSearchQuery, //TODO: not in proto
	}
}
