package admin

import (
	org_grpc "github.com/caos/zitadel/internal/api/grpc/org"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
)

func listOrgRequestToModel(req *admin.ListOrgsRequest) (*model.OrgSearchRequest, error) {
	queries, err := org_grpc.OrgQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &model.OrgSearchRequest{
		Offset:  req.Query.Offset,
		Limit:   uint64(req.Query.Limit),
		Asc:     req.Query.Asc,
		Queries: queries,
	}, nil
}

func setUpOrgOrgToDomain(req *admin.SetUpOrgRequest_Org) *domain.Org {
	org := &domain.Org{
		Name:    req.Name,
		Domains: []*domain.OrgDomain{},
	}
	if req.Domain != "" {
		org.Domains = append(org.Domains, &domain.OrgDomain{Domain: req.Domain})
	}
	return org
}
