package management

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/query"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func addLoginPolicyToDomain(p *mgmt_pb.AddCustomLoginPolicyRequest) *domain.LoginPolicy {
	return &domain.LoginPolicy{
		AllowUsernamePassword:  p.AllowUsernamePassword,
		AllowRegister:          p.AllowRegister,
		AllowExternalIDP:       p.AllowExternalIdp,
		ForceMFA:               p.ForceMfa,
		PasswordlessType:       policy_grpc.PasswordlessTypeToDomain(p.PasswordlessType),
		HidePasswordReset:      p.HidePasswordReset,
		IgnoreUnknownUsernames: p.IgnoreUnknownUsernames,
	}
}

func updateLoginPolicyToDomain(p *mgmt_pb.UpdateCustomLoginPolicyRequest) *domain.LoginPolicy {
	return &domain.LoginPolicy{
		AllowUsernamePassword:  p.AllowUsernamePassword,
		AllowRegister:          p.AllowRegister,
		AllowExternalIDP:       p.AllowExternalIdp,
		ForceMFA:               p.ForceMfa,
		PasswordlessType:       policy_grpc.PasswordlessTypeToDomain(p.PasswordlessType),
		HidePasswordReset:      p.HidePasswordReset,
		IgnoreUnknownUsernames: p.IgnoreUnknownUsernames,
	}
}

func ListLoginPolicyIDPsRequestToQuery(req *mgmt_pb.ListLoginPolicyIDPsRequest) *query.IDPLoginPolicyLinksSearchQuery {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	return &query.IDPLoginPolicyLinksSearchQuery{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
	}
}
