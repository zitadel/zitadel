package grpc

import (
	admin_model "github.com/caos/zitadel/internal/admin/model"
	org_model "github.com/caos/zitadel/internal/org/model"
)

func setUpRequestToModel(setUp *OrgSetUpRequest) *admin_model.SetupOrg {
	return &admin_model.SetupOrg{
		Org: orgCreateRequestToModel(setUp.Org),
	}
}

func orgCreateRequestToModel(org *CreateOrgRequest) *org_model.Org {
	return &org_model.Org{
		Domain: org.Domain,
		Name:   org.Name,
	}
}
