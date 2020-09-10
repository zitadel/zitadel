package repository

import (
	"context"
	auth_model "github.com/caos/zitadel/internal/auth/model"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	org_model "github.com/caos/zitadel/internal/org/model"
)

type OrgRepository interface {
	RegisterOrg(context.Context, *auth_model.RegisterOrg) (*auth_model.RegisterOrg, error)
	GetOrgIamPolicy(ctx context.Context, orgID string) (*org_model.OrgIAMPolicy, error)
	GetDefaultOrgIamPolicy(ctx context.Context) (*org_model.OrgIAMPolicy, error)
	GetIDPConfigByID(ctx context.Context, idpConfigID string) (*iam_model.IDPConfigView, error)
}
