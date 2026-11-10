package admin

import (
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	org_grpc "github.com/zitadel/zitadel/internal/api/grpc/org"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
)

func listOrgRequestToModel(defaults systemdefaults.SystemDefaults, req *admin.ListOrgsRequest) (*query.OrgSearchQueries, error) {
	offset, limit, asc, err := object.ListQueryToModel(defaults, req.Query)
	if err != nil {
		return nil, err
	}
	queries, err := org_grpc.OrgQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.OrgSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			SortingColumn: org_grpc.FieldNameToOrgColumn(req.SortingColumn),
			Asc:           asc,
		},
		Queries: queries,
	}, nil
}
