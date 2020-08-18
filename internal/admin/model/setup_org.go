package model

import (
	org_model "github.com/caos/zitadel/internal/org/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
)

type SetupOrg struct {
	*org_model.Org
	*usr_model.Human
}
