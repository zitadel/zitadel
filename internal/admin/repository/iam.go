package repository

import (
	"context"

	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type IAMRepository interface {
	SearchIAMMembers(ctx context.Context, request *iam_model.IAMMemberSearchRequest) (*iam_model.IAMMemberSearchResponse, error)

	GetIAMMemberRoles() []string

	SearchIDPConfigs(ctx context.Context, request *iam_model.IDPConfigSearchRequest) (*iam_model.IDPConfigSearchResponse, error)

	GetDefaultLoginPolicy(ctx context.Context) (*iam_model.LoginPolicyView, error)
	SearchDefaultIDPProviders(ctx context.Context, request *iam_model.IDPProviderSearchRequest) (*iam_model.IDPProviderSearchResponse, error)
	SearchDefaultSecondFactors(ctx context.Context) (*iam_model.SecondFactorsSearchResponse, error)
	SearchDefaultMultiFactors(ctx context.Context) (*iam_model.MultiFactorsSearchResponse, error)

	GetDefaultLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error)

	GetDefaultPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error)

	GetDefaultPasswordAgePolicy(ctx context.Context) (*iam_model.PasswordAgePolicyView, error)

	GetDefaultPasswordLockoutPolicy(ctx context.Context) (*iam_model.PasswordLockoutPolicyView, error)

	GetDefaultOrgIAMPolicy(ctx context.Context) (*iam_model.OrgIAMPolicyView, error)
}
