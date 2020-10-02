package repository

import (
	"context"
	auth_model "github.com/caos/zitadel/internal/auth/model"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type OrgRepository interface {
	RegisterOrg(context.Context, *auth_model.RegisterOrg) (*auth_model.RegisterOrg, error)
	GetOrgIamPolicy(ctx context.Context, orgID string) (*iam_model.OrgIAMPolicy, error)
	GetDefaultOrgIamPolicy(ctx context.Context) *iam_model.OrgIAMPolicy
	GetIDPConfigByID(ctx context.Context, idpConfigID string) (*iam_model.IDPConfigView, error)
	GetMyPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error)
}
