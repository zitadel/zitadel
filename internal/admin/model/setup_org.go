package model

import org_model "github.com/caos/zitadel/internal/org/model"

type SetupOrg struct {
	*org_model.Org
	User interface{}
}
