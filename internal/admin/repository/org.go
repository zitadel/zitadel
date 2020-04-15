package repository

import (
	"context"

	admin_model "github.com/caos/zitadel/internal/admin/model"
)

type OrgRepository interface {
	SetUpOrg(context.Context, *admin_model.SetupOrg) (*admin_model.SetupOrg, error)
}
