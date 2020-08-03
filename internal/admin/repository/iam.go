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

	AddOidcIdpConfig(ctx context.Context, idp *iam_model.IdpConfig) (*iam_model.IdpConfig, error)
	ChangeIdpConfig(ctx context.Context, idp *iam_model.IdpConfig) (*iam_model.IdpConfig, error)
	DeactivateIdpConfig(ctx context.Context, idpConfigID string) (*iam_model.IdpConfig, error)
	ReactivateIdpConfig(ctx context.Context, idpConfigID string) (*iam_model.IdpConfig, error)
	RemoveIdpConfig(ctx context.Context, idpConfigID string) error
	ChangeOidcIdpConfig(ctx context.Context, oidcConfig *iam_model.OidcIdpConfig) (*iam_model.OidcIdpConfig, error)
}
