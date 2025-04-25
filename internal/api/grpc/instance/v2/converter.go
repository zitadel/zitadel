package instance

import (
	"github.com/zitadel/zitadel/cmd/build"
	filter "github.com/zitadel/zitadel/internal/api/grpc/filter/v2beta"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/instance/v2"
)

func InstancesToPb(instances []*query.Instance) []*instance.Instance {
	list := []*instance.Instance{}
	for _, instance := range instances {
		list = append(list, ToProtoObject(instance))
	}
	return list
}

func ToProtoObject(inst *query.Instance) *instance.Instance {
	return &instance.Instance{
		Id:      inst.ID,
		Name:    inst.Name,
		Domains: DomainsToPb(inst.Domains),
		Version: build.Version(),
		Details: object.ToViewDetailsPb(inst.Sequence, inst.CreationDate, inst.ChangeDate, inst.ID),
	}
}

func DomainsToPb(domains []*query.InstanceDomain) []*instance.Domain {
	d := []*instance.Domain{}
	for _, dm := range domains {
		pbDomain := DomainToPb(dm)
		d = append(d, pbDomain)
	}
	return d
}

func DomainToPb(d *query.InstanceDomain) *instance.Domain {
	return &instance.Domain{
		Domain:    d.Domain,
		Primary:   d.IsPrimary,
		Generated: d.IsGenerated,
		Details: object.ToViewDetailsPb(
			d.Sequence,
			d.CreationDate,
			d.ChangeDate,
			d.InstanceID,
		),
	}
}

func ListInstancesRequestToModel(req *instance.ListInstancesRequest, sysDefaults systemdefaults.SystemDefaults) (*query.InstanceSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(sysDefaults, req.GetPagination())
	if err != nil {
		return nil, err
	}

	queries, err := instanceQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}

	return &query.InstanceSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: fieldNameToInstanceColumn(req.SortingColumn),
		},
		Queries: queries,
	}, nil

}

func fieldNameToInstanceColumn(fieldName instance.FieldName) query.Column {
	switch fieldName {
	case instance.FieldName_FIELD_NAME_ID:
		return query.InstanceColumnID
	case instance.FieldName_FIELD_NAME_NAME:
		return query.InstanceColumnName
	case instance.FieldName_FIELD_NAME_CREATION_DATE:
		return query.InstanceColumnCreationDate
	default:
		return query.Column{}
	}
}

func instanceQueriesToModel(queries []*instance.Query) (_ []query.SearchQuery, err error) {
	q := []query.SearchQuery{}
	for _, query := range queries {
		model, err := instanceQueryToModel(query)
		if err != nil {
			return nil, err
		}
		q = append(q, model)
	}
	return q, nil
}

func instanceQueryToModel(searchQuery *instance.Query) (query.SearchQuery, error) {
	switch q := searchQuery.Query.(type) {
	case *instance.Query_IdQuery:
		return query.NewInstanceIDsListSearchQuery(q.IdQuery.Ids...)
	case *instance.Query_DomainQuery:
		return query.NewInstanceDomainsListSearchQuery(q.DomainQuery.Domains...)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "INST-3m0se", "List.Query.Invalid")
	}
}
