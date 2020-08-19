package repository

import (
	"context"
	iam_model "github.com/caos/zitadel/internal/iam/model"

	org_model "github.com/caos/zitadel/internal/org/model"
)

type OrgRepository interface {
	OrgByID(ctx context.Context, id string) (*org_model.OrgView, error)
	OrgByDomainGlobal(ctx context.Context, domain string) (*org_model.OrgView, error)
	CreateOrg(ctx context.Context, name string) (*org_model.Org, error)
	UpdateOrg(ctx context.Context, org *org_model.Org) (*org_model.Org, error)
	DeactivateOrg(ctx context.Context, id string) (*org_model.Org, error)
	ReactivateOrg(ctx context.Context, id string) (*org_model.Org, error)
	OrgChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool) (*org_model.OrgChanges, error)

	SearchMyOrgDomains(ctx context.Context, request *org_model.OrgDomainSearchRequest) (*org_model.OrgDomainSearchResponse, error)
	AddMyOrgDomain(ctx context.Context, domain *org_model.OrgDomain) (*org_model.OrgDomain, error)
	GenerateMyOrgDomainValidation(ctx context.Context, domain *org_model.OrgDomain) (string, string, error)
	ValidateMyOrgDomain(ctx context.Context, domain *org_model.OrgDomain) error
	SetMyPrimaryOrgDomain(ctx context.Context, domain *org_model.OrgDomain) error
	RemoveMyOrgDomain(ctx context.Context, domain string) error

	SearchMyOrgMembers(ctx context.Context, request *org_model.OrgMemberSearchRequest) (*org_model.OrgMemberSearchResponse, error)
	AddMyOrgMember(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error)
	ChangeMyOrgMember(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error)
	RemoveMyOrgMember(ctx context.Context, userID string) error

	GetOrgMemberRoles() []string

	GetMyOrgIamPolicy(ctx context.Context) (*org_model.OrgIamPolicy, error)

	SearchIdpConfigs(ctx context.Context, request *iam_model.IdpConfigSearchRequest) (*iam_model.IdpConfigSearchResponse, error)
	IdpConfigByID(ctx context.Context, id string) (*iam_model.IdpConfigView, error)
	AddOidcIdpConfig(ctx context.Context, idp *iam_model.IdpConfig) (*iam_model.IdpConfig, error)
	ChangeIdpConfig(ctx context.Context, idp *iam_model.IdpConfig) (*iam_model.IdpConfig, error)
	DeactivateIdpConfig(ctx context.Context, idpConfigID string) (*iam_model.IdpConfig, error)
	ReactivateIdpConfig(ctx context.Context, idpConfigID string) (*iam_model.IdpConfig, error)
	RemoveIdpConfig(ctx context.Context, idpConfigID string) error
	ChangeOidcIdpConfig(ctx context.Context, oidcConfig *iam_model.OidcIdpConfig) (*iam_model.OidcIdpConfig, error)

	GetLoginPolicy(ctx context.Context) (*iam_model.LoginPolicyView, error)
	AddLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error)
	ChangeLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error)
	RemoveLoginPolicy(ctx context.Context) error
	SearchIdpProviders(ctx context.Context, request *iam_model.IdpProviderSearchRequest) (*iam_model.IdpProviderSearchResponse, error)
	AddIdpProviderToLoginPolicy(ctx context.Context, provider *iam_model.IdpProvider) (*iam_model.IdpProvider, error)
	RemoveIdpProviderFromLoginPolicy(ctx context.Context, provider *iam_model.IdpProvider) error
}
