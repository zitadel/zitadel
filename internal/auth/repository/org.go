package repository

import (
	"context"
	auth_model "github.com/caos/zitadel/internal/auth/model"
	org_model "github.com/caos/zitadel/internal/org/model"
)

type OrgRepository interface {
	RegisterOrg(context.Context, *auth_model.RegisterOrg) (*auth_model.RegisterOrg, error)
	GetOrgIamPolicy(ctx context.Context, orgID string) (*org_model.OrgIAMPolicy, error)
}
