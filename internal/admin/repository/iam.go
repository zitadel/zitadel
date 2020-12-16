package repository

import (
	"context"

	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type IAMRepository interface {
	SearchIAMMembers(ctx context.Context, request *iam_model.IAMMemberSearchRequest) (*iam_model.IAMMemberSearchResponse, error)
	AddIAMMember(ctx context.Context, member *iam_model.IAMMember) (*iam_model.IAMMember, error)
	ChangeIAMMember(ctx context.Context, member *iam_model.IAMMember) (*iam_model.IAMMember, error)
	RemoveIAMMember(ctx context.Context, userID string) error

	GetIAMMemberRoles() []string

	SearchIDPConfigs(ctx context.Context, request *iam_model.IDPConfigSearchRequest) (*iam_model.IDPConfigSearchResponse, error)
	IDPConfigByID(ctx context.Context, id string) (*iam_model.IDPConfigView, error)
	AddOIDCIDPConfig(ctx context.Context, idp *iam_model.IDPConfig) (*iam_model.IDPConfig, error)
	ChangeIDPConfig(ctx context.Context, idp *iam_model.IDPConfig) (*iam_model.IDPConfig, error)
	DeactivateIDPConfig(ctx context.Context, idpConfigID string) (*iam_model.IDPConfig, error)
	ReactivateIDPConfig(ctx context.Context, idpConfigID string) (*iam_model.IDPConfig, error)
	RemoveIDPConfig(ctx context.Context, idpConfigID string) error
	ChangeOidcIDPConfig(ctx context.Context, oidcConfig *iam_model.OIDCIDPConfig) (*iam_model.OIDCIDPConfig, error)

	GetDefaultLoginPolicy(ctx context.Context) (*iam_model.LoginPolicyView, error)
	AddDefaultLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error)
	ChangeDefaultLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error)
	SearchDefaultIDPProviders(ctx context.Context, request *iam_model.IDPProviderSearchRequest) (*iam_model.IDPProviderSearchResponse, error)
	AddIDPProviderToLoginPolicy(ctx context.Context, provider *iam_model.IDPProvider) (*iam_model.IDPProvider, error)
	RemoveIDPProviderFromLoginPolicy(ctx context.Context, provider *iam_model.IDPProvider) error
	SearchDefaultSecondFactors(ctx context.Context) (*iam_model.SecondFactorsSearchResponse, error)
	AddSecondFactorToLoginPolicy(ctx context.Context, mfa iam_model.SecondFactorType) (iam_model.SecondFactorType, error)
	RemoveSecondFactorFromLoginPolicy(ctx context.Context, mfa iam_model.SecondFactorType) error
	SearchDefaultMultiFactors(ctx context.Context) (*iam_model.MultiFactorsSearchResponse, error)
	AddMultiFactorToLoginPolicy(ctx context.Context, mfa iam_model.MultiFactorType) (iam_model.MultiFactorType, error)
	RemoveMultiFactorFromLoginPolicy(ctx context.Context, mfa iam_model.MultiFactorType) error

	GetDefaultLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error)
	AddDefaultLabelPolicy(ctx context.Context, policy *iam_model.LabelPolicy) (*iam_model.LabelPolicy, error)
	ChangeDefaultLabelPolicy(ctx context.Context, policy *iam_model.LabelPolicy) (*iam_model.LabelPolicy, error)

	GetDefaultPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error)
	AddDefaultPasswordComplexityPolicy(ctx context.Context, policy *iam_model.PasswordComplexityPolicy) (*iam_model.PasswordComplexityPolicy, error)
	ChangeDefaultPasswordComplexityPolicy(ctx context.Context, policy *iam_model.PasswordComplexityPolicy) (*iam_model.PasswordComplexityPolicy, error)

	GetDefaultPasswordAgePolicy(ctx context.Context) (*iam_model.PasswordAgePolicyView, error)
	AddDefaultPasswordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) (*iam_model.PasswordAgePolicy, error)
	ChangeDefaultPasswordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) (*iam_model.PasswordAgePolicy, error)

	GetDefaultPasswordLockoutPolicy(ctx context.Context) (*iam_model.PasswordLockoutPolicyView, error)
	AddDefaultPasswordLockoutPolicy(ctx context.Context, policy *iam_model.PasswordLockoutPolicy) (*iam_model.PasswordLockoutPolicy, error)
	ChangeDefaultPasswordLockoutPolicy(ctx context.Context, policy *iam_model.PasswordLockoutPolicy) (*iam_model.PasswordLockoutPolicy, error)

	GetDefaultOrgIAMPolicy(ctx context.Context) (*iam_model.OrgIAMPolicyView, error)
	AddDefaultOrgIAMPolicy(ctx context.Context, policy *iam_model.OrgIAMPolicy) (*iam_model.OrgIAMPolicy, error)
	ChangeDefaultOrgIAMPolicy(ctx context.Context, policy *iam_model.OrgIAMPolicy) (*iam_model.OrgIAMPolicy, error)
}
