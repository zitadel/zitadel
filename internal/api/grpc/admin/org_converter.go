package admin

import (
	org_grpc "github.com/caos/zitadel/internal/api/grpc/org"
	"github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/admin"
)

func orgListRequestToModel(req *admin.ListOrgsRequest) (*model.OrgSearchRequest, error) {
	queries, err := org_grpc.OrgQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &model.OrgSearchRequest{
		Offset:  req.MetaData.Offset,
		Limit:   uint64(req.MetaData.Limit),
		Asc:     req.MetaData.Asc,
		Queries: queries,
	}, nil
}

func setUpOrgOrgToModel(req *admin.SetUpOrgRequest_Org) *domain.Org {
	org := &domain.Org{
		Name:    req.Name,
		Domains: []*domain.OrgDomain{},
	}
	if req.Domain != "" {
		org.Domains = append(org.Domains, req.Domain)
	}
	return org
}
