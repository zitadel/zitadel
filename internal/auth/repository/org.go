package repository

import (
	"context"
	auth_model "github.com/caos/zitadel/internal/auth/model"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	org_model "github.com/caos/zitadel/internal/org/model"
)

type OrgRepository interface {
	RegisterOrg(context.Context, *auth_model.RegisterOrg) (*auth_model.RegisterOrg, error)
	OrgByPrimaryDomain(primaryDomain string) (*org_model.OrgView, error)
	GetOrgIAMPolicy(ctx context.Context, orgID string) (*iam_model.OrgIAMPolicyView, error)
	GetDefaultOrgIAMPolicy(ctx context.Context) (*iam_model.OrgIAMPolicyView, error)
	GetIDPConfigByID(ctx context.Context, idpConfigID string) (*iam_model.IDPConfigView, error)
	GetMyPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error)
}
