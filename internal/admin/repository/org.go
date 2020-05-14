package repository

import (
	"context"

	admin_model "github.com/caos/zitadel/internal/admin/model"
	org_model "github.com/caos/zitadel/internal/org/model"
)

type OrgRepository interface {
	SetUpOrg(context.Context, *admin_model.SetupOrg) (*admin_model.SetupOrg, error)
	IsOrgUnique(ctx context.Context, name, domain string) (bool, error)
	OrgByID(ctx context.Context, id string) (*org_model.Org, error)
	SearchOrgs(ctx context.Context) ([]*org_model.Org, error)
}
