package admin

import (
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/zitadel/zitadel/internal/api/grpc/policy"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func updateLoginPolicyToCommand(p *admin_pb.UpdateLoginPolicyRequest) *command.ChangeLoginPolicy {
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

func ListLoginPolicyIDPsRequestToQuery(defaults systemdefaults.SystemDefaults, req *admin_pb.ListLoginPolicyIDPsRequest) (*query.IDPLoginPolicyLinksSearchQuery, error) {
	offset, limit, asc, err := object.ListQueryToModel(defaults, req.Query)
	if err != nil {
		return nil, err
	}
	return &query.IDPLoginPolicyLinksSearchQuery{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
	}, nil
}
