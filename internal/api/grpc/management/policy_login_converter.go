package management

import (
	idp_grpc "github.com/zitadel/zitadel/internal/api/grpc/idp"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/zitadel/zitadel/internal/api/grpc/policy"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func AddLoginPolicyToCommand(p *mgmt_pb.AddCustomLoginPolicyRequest) *command.AddLoginPolicy {
	return &command.AddLoginPolicy{
		AllowUsernamePassword:      p.AllowUsernamePassword,
		AllowRegister:              p.AllowRegister,
		AllowExternalIDP:           p.AllowExternalIdp,
		ForceMFA:                   p.ForceMfa,
		ForceMFALocalOnly:          p.ForceMfaLocalOnly,
		PasswordlessType:           policy_grpc.PasswordlessTypeToDomain(p.PasswordlessType),
		HidePasswordReset:          p.HidePasswordReset,
		IgnoreUnknownUsernames:     p.IgnoreUnknownUsernames,
		AllowDomainDiscovery:       p.AllowDomainDiscovery,
		DefaultRedirectURI:         p.DefaultRedirectUri,
		PasswordCheckLifetime:      p.PasswordCheckLifetime.AsDuration(),
		ExternalLoginCheckLifetime: p.ExternalLoginCheckLifetime.AsDuration(),
		MFAInitSkipLifetime:        p.MfaInitSkipLifetime.AsDuration(),
		SecondFactorCheckLifetime:  p.SecondFactorCheckLifetime.AsDuration(),
		MultiFactorCheckLifetime:   p.MultiFactorCheckLifetime.AsDuration(),
		SecondFactors:              policy_grpc.SecondFactorsTypesToDomain(p.SecondFactors),
		MultiFactors:               policy_grpc.MultiFactorsTypesToDomain(p.MultiFactors),
		IDPProviders:               addLoginPolicyIDPsToCommand(p.Idps),
		DisableLoginWithEmail:      p.DisableLoginWithEmail,
		DisableLoginWithPhone:      p.DisableLoginWithPhone,
	}
}
func addLoginPolicyIDPsToCommand(idps []*mgmt_pb.AddCustomLoginPolicyRequest_IDP) []*command.AddLoginPolicyIDP {
	providers := make([]*command.AddLoginPolicyIDP, len(idps))
	for i, idp := range idps {
		providers[i] = &command.AddLoginPolicyIDP{
			Type:     idp_grpc.IDPProviderTypeFromPb(idp.OwnerType),
			ConfigID: idp.IdpId,
		}
	}
	return providers
}

func updateLoginPolicyToCommand(p *mgmt_pb.UpdateCustomLoginPolicyRequest) *command.ChangeLoginPolicy {
	return &command.ChangeLoginPolicy{
		AllowUsernamePassword:      p.AllowUsernamePassword,
		AllowRegister:              p.AllowRegister,
		AllowExternalIDP:           p.AllowExternalIdp,
		ForceMFA:                   p.ForceMfa,
		ForceMFALocalOnly:          p.ForceMfaLocalOnly,
		PasswordlessType:           policy_grpc.PasswordlessTypeToDomain(p.PasswordlessType),
		HidePasswordReset:          p.HidePasswordReset,
		IgnoreUnknownUsernames:     p.IgnoreUnknownUsernames,
		AllowDomainDiscovery:       p.AllowDomainDiscovery,
		DisableLoginWithEmail:      p.DisableLoginWithEmail,
		DisableLoginWithPhone:      p.DisableLoginWithPhone,
		DefaultRedirectURI:         p.DefaultRedirectUri,
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
