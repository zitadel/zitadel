package repository

import (
	"context"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type IamRepository interface {
	SearchIamMembers(ctx context.Context, request *iam_model.IamMemberSearchRequest) (*iam_model.IamMemberSearchResponse, error)
	AddIamMember(ctx context.Context, member *iam_model.IamMember) (*iam_model.IamMember, error)
	ChangeIamMember(ctx context.Context, member *iam_model.IamMember) (*iam_model.IamMember, error)
	RemoveIamMember(ctx context.Context, userID string) error

	GetIamMemberRoles() []string

	SearchIdpConfigs(ctx context.Context, request *iam_model.IdpConfigSearchRequest) (*iam_model.IdpConfigSearchResponse, error)
	IdpConfigByID(ctx context.Context, id string) (*iam_model.IdpConfigView, error)
	AddOidcIdpConfig(ctx context.Context, idp *iam_model.IdpConfig) (*iam_model.IdpConfig, error)
	ChangeIdpConfig(ctx context.Context, idp *iam_model.IdpConfig) (*iam_model.IdpConfig, error)
	DeactivateIdpConfig(ctx context.Context, idpConfigID string) (*iam_model.IdpConfig, error)
	ReactivateIdpConfig(ctx context.Context, idpConfigID string) (*iam_model.IdpConfig, error)
	RemoveIdpConfig(ctx context.Context, idpConfigID string) error
	ChangeOidcIdpConfig(ctx context.Context, oidcConfig *iam_model.OidcIdpConfig) (*iam_model.OidcIdpConfig, error)

	GetDefaultLoginPolicy(ctx context.Context) (*iam_model.LoginPolicyView, error)
	AddDefaultLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error)
	ChangeDefaultLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error)
	SearchDefaultIdpProviders(ctx context.Context, request *iam_model.IdpProviderSearchRequest) (*iam_model.IdpProviderSearchResponse, error)
	AddIdpProviderToLoginPolicy(ctx context.Context, provider *iam_model.IdpProvider) (*iam_model.IdpProvider, error)
	RemoveIdpProviderFromLoginPolicy(ctx context.Context, provider *iam_model.IdpProvider) error
}
