package admin

import (
	"github.com/zitadel/zitadel/v2/internal/api/grpc/object"
	org_grpc "github.com/zitadel/zitadel/v2/internal/api/grpc/org"
	"github.com/zitadel/zitadel/v2/internal/query"
	"github.com/zitadel/zitadel/v2/pkg/grpc/admin"
	"github.com/zitadel/zitadel/v2/pkg/grpc/org"
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

func fieldNameToOrgColumn(fieldName org.OrgFieldName) query.Column {
	switch fieldName {
	case org.OrgFieldName_ORG_FIELD_NAME_NAME:
		return query.OrgColumnName
	default:
		return query.Column{}
	}
}
