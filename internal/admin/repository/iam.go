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

	SearchIdpConfigs(ctx context.Context, request *iam_model.IDPConfigSearchRequest) (*iam_model.IDPConfigSearchResponse, error)
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
}
