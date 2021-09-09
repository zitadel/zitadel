package admin

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	org_grpc "github.com/caos/zitadel/internal/api/grpc/org"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/caos/zitadel/pkg/grpc/org"
)

func listOrgRequestToModel(req *admin.ListOrgsRequest) (*query.OrgSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := org_grpc.OrgQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.OrgSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			SortingColumn: fieldNameToOrgColumn(req.SortingColumn),
			Asc:           asc,
		},
		Queries: queries,
	}, nil
}

func fieldNameToOrgColumn(fieldName org.OrgFieldName) query.SortingColumn {
	switch fieldName {
	case org.OrgFieldName_ORG_FIELD_NAME_NAME:
		return query.OrgColumnName
	default:
		return nil
	}
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
