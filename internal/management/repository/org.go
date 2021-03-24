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

	GetMyOrgIamPolicy(ctx context.Context) (*iam_model.OrgIAMPolicyView, error)

	SearchIDPConfigs(ctx context.Context, request *iam_model.IDPConfigSearchRequest) (*iam_model.IDPConfigSearchResponse, error)
	IDPConfigByID(ctx context.Context, id string) (*iam_model.IDPConfigView, error)
	AddOIDCIDPConfig(ctx context.Context, idp *iam_model.IDPConfig) (*iam_model.IDPConfig, error)
	ChangeIDPConfig(ctx context.Context, idp *iam_model.IDPConfig) (*iam_model.IDPConfig, error)
	DeactivateIDPConfig(ctx context.Context, idpConfigID string) (*iam_model.IDPConfig, error)
	ReactivateIDPConfig(ctx context.Context, idpConfigID string) (*iam_model.IDPConfig, error)
	RemoveIDPConfig(ctx context.Context, idpConfigID string) error
	ChangeOIDCIDPConfig(ctx context.Context, oidcConfig *iam_model.OIDCIDPConfig) (*iam_model.OIDCIDPConfig, error)

	GetLoginPolicy(ctx context.Context) (*iam_model.LoginPolicyView, error)
	GetDefaultLoginPolicy(ctx context.Context) (*iam_model.LoginPolicyView, error)
	AddLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error)
	ChangeLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error)
	RemoveLoginPolicy(ctx context.Context) error
	SearchIDPProviders(ctx context.Context, request *iam_model.IDPProviderSearchRequest) (*iam_model.IDPProviderSearchResponse, error)
	AddIDPProviderToLoginPolicy(ctx context.Context, provider *iam_model.IDPProvider) (*iam_model.IDPProvider, error)
	RemoveIDPProviderFromLoginPolicy(ctx context.Context, provider *iam_model.IDPProvider) error
	SearchSecondFactors(ctx context.Context) (*iam_model.SecondFactorsSearchResponse, error)
	AddSecondFactorToLoginPolicy(ctx context.Context, mfa iam_model.SecondFactorType) (iam_model.SecondFactorType, error)
	RemoveSecondFactorFromLoginPolicy(ctx context.Context, mfa iam_model.SecondFactorType) error
	SearchMultiFactors(ctx context.Context) (*iam_model.MultiFactorsSearchResponse, error)
	AddMultiFactorToLoginPolicy(ctx context.Context, mfa iam_model.MultiFactorType) (iam_model.MultiFactorType, error)
	RemoveMultiFactorFromLoginPolicy(ctx context.Context, mfa iam_model.MultiFactorType) error

	GetPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error)
	GetDefaultPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error)
	AddPasswordComplexityPolicy(ctx context.Context, policy *iam_model.PasswordComplexityPolicy) (*iam_model.PasswordComplexityPolicy, error)
	ChangePasswordComplexityPolicy(ctx context.Context, policy *iam_model.PasswordComplexityPolicy) (*iam_model.PasswordComplexityPolicy, error)
	RemovePasswordComplexityPolicy(ctx context.Context) error

	GetPasswordAgePolicy(ctx context.Context) (*iam_model.PasswordAgePolicyView, error)
	GetDefaultPasswordAgePolicy(ctx context.Context) (*iam_model.PasswordAgePolicyView, error)
	AddPasswordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) (*iam_model.PasswordAgePolicy, error)
	ChangePasswordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) (*iam_model.PasswordAgePolicy, error)
	RemovePasswordAgePolicy(ctx context.Context) error

	GetPasswordLockoutPolicy(ctx context.Context) (*iam_model.PasswordLockoutPolicyView, error)
	GetDefaultPasswordLockoutPolicy(ctx context.Context) (*iam_model.PasswordLockoutPolicyView, error)
	AddPasswordLockoutPolicy(ctx context.Context, policy *iam_model.PasswordLockoutPolicy) (*iam_model.PasswordLockoutPolicy, error)
	ChangePasswordLockoutPolicy(ctx context.Context, policy *iam_model.PasswordLockoutPolicy) (*iam_model.PasswordLockoutPolicy, error)
	RemovePasswordLockoutPolicy(ctx context.Context) error

	GetDefaultMailTemplate(ctx context.Context) (*iam_model.MailTemplateView, error)
	GetMailTemplate(ctx context.Context) (*iam_model.MailTemplateView, error)
	AddMailTemplate(ctx context.Context, template *iam_model.MailTemplate) (*iam_model.MailTemplate, error)
	ChangeMailTemplate(ctx context.Context, template *iam_model.MailTemplate) (*iam_model.MailTemplate, error)
	RemoveMailTemplate(ctx context.Context) error

	GetDefaultMailTexts(ctx context.Context) (*iam_model.MailTextsView, error)
	GetMailTexts(ctx context.Context) (*iam_model.MailTextsView, error)
	AddMailText(ctx context.Context, mailText *iam_model.MailText) (*iam_model.MailText, error)
	ChangeMailText(ctx context.Context, mailText *iam_model.MailText) (*iam_model.MailText, error)
	RemoveMailText(ctx context.Context, mailText *iam_model.MailText) error

	GetLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error)
	GetDefaultLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error)
	AddLabelPolicy(ctx context.Context, policy *iam_model.LabelPolicy) (*iam_model.LabelPolicy, error)
	ChangeLabelPolicy(ctx context.Context, policy *iam_model.LabelPolicy) (*iam_model.LabelPolicy, error)
	RemoveLabelPolicy(ctx context.Context) error
}
