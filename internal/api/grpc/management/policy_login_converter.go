package management

import (
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/zitadel/zitadel/internal/api/grpc/policy"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func addLoginPolicyToDomain(p *mgmt_pb.AddCustomLoginPolicyRequest) *domain.LoginPolicy {
	return &domain.LoginPolicy{
		AllowUsernamePassword:      p.AllowUsernamePassword,
		AllowRegister:              p.AllowRegister,
		AllowExternalIDP:           p.AllowExternalIdp,
		ForceMFA:                   p.ForceMfa,
		PasswordlessType:           policy_grpc.PasswordlessTypeToDomain(p.PasswordlessType),
		HidePasswordReset:          p.HidePasswordReset,
		PasswordCheckLifetime:      p.PasswordCheckLifetime.AsDuration(),
		ExternalLoginCheckLifetime: p.ExternalLoginCheckLifetime.AsDuration(),
		MFAInitSkipLifetime:        p.MfaInitSkipLifetime.AsDuration(),
		SecondFactorCheckLifetime:  p.SecondFactorCheckLifetime.AsDuration(),
		MultiFactorCheckLifetime:   p.MultiFactorCheckLifetime.AsDuration(),
	}
}

func updateLoginPolicyToDomain(p *mgmt_pb.UpdateCustomLoginPolicyRequest) *domain.LoginPolicy {
	return &domain.LoginPolicy{
		AllowUsernamePassword:      p.AllowUsernamePassword,
		AllowRegister:              p.AllowRegister,
		AllowExternalIDP:           p.AllowExternalIdp,
		ForceMFA:                   p.ForceMfa,
		PasswordlessType:           policy_grpc.PasswordlessTypeToDomain(p.PasswordlessType),
		HidePasswordReset:          p.HidePasswordReset,
		PasswordCheckLifetime:      p.PasswordCheckLifetime.AsDuration(),
		ExternalLoginCheckLifetime: p.ExternalLoginCheckLifetime.AsDuration(),
		MFAInitSkipLifetime:        p.MfaInitSkipLifetime.AsDuration(),
		SecondFactorCheckLifetime:  p.SecondFactorCheckLifetime.AsDuration(),
		MultiFactorCheckLifetime:   p.MultiFactorCheckLifetime.AsDuration(),
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
