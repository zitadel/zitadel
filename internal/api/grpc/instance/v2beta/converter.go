package instance

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/cmd/build"
	filter "github.com/zitadel/zitadel/internal/api/grpc/filter/v2beta"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
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
		Id:           inst.ID,
		Name:         inst.Name,
		Domains:      DomainsToPb(inst.Domains),
		Version:      build.Version(),
		ChangeDate:   timestamppb.New(inst.ChangeDate),
		CreationDate: timestamppb.New(inst.CreationDate),
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
		Domain:       d.Domain,
		Primary:      d.IsPrimary,
		Generated:    d.IsGenerated,
		InstanceId:   d.InstanceID,
		CreationDate: timestamppb.New(d.CreationDate),
	}
}

func ListInstancesRequestToModel(req *instance.ListInstancesRequest, sysDefaults systemdefaults.SystemDefaults) (*query.InstanceSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(sysDefaults, req.GetPagination())
	if err != nil {
		return nil, err
	}

	queries, err := instanceQueriesToModel(req.GetQueries())
	if err != nil {
		return nil, err
	}

	return &query.InstanceSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: fieldNameToInstanceColumn(req.GetSortingColumn()),
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
	case instance.FieldName_FIELD_NAME_UNSPECIFIED:
		fallthrough
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
	switch q := searchQuery.GetQuery().(type) {
	case *instance.Query_IdQuery:
		return query.NewInstanceIDsListSearchQuery(q.IdQuery.GetIds()...)
	case *instance.Query_DomainQuery:
		return query.NewInstanceDomainsListSearchQuery(q.DomainQuery.GetDomains()...)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "INST-3m0se", "List.Query.Invalid")
	}
}

func ListCustomDomainsRequestToModel(req *instance.ListCustomDomainsRequest, defaults systemdefaults.SystemDefaults) (*query.InstanceDomainSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(defaults, req.GetPagination())
	if err != nil {
		return nil, err
	}

	queries, err := domainQueriesToModel(req.GetQueries())
	if err != nil {
		return nil, err
	}

	return &query.InstanceDomainSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: fieldNameToInstanceDomainColumn(req.GetSortingColumn()),
		},
		Queries: queries,
	}, nil
}

func fieldNameToInstanceDomainColumn(fieldName instance.DomainFieldName) query.Column {
	switch fieldName {
	case instance.DomainFieldName_DOMAIN_FIELD_NAME_DOMAIN:
		return query.InstanceDomainDomainCol
	case instance.DomainFieldName_DOMAIN_FIELD_NAME_GENERATED:
		return query.InstanceDomainIsGeneratedCol
	case instance.DomainFieldName_DOMAIN_FIELD_NAME_PRIMARY:
		return query.InstanceDomainIsPrimaryCol
	case instance.DomainFieldName_DOMAIN_FIELD_NAME_CREATION_DATE:
		return query.InstanceDomainCreationDateCol
	case instance.DomainFieldName_DOMAIN_FIELD_NAME_UNSPECIFIED:
		fallthrough
	default:
		return query.Column{}
	}
}

func domainQueriesToModel(queries []*instance.DomainSearchQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = domainQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func domainQueryToModel(searchQuery *instance.DomainSearchQuery) (query.SearchQuery, error) {
	switch q := searchQuery.GetQuery().(type) {
	case *instance.DomainSearchQuery_DomainQuery:
		return query.NewInstanceDomainDomainSearchQuery(object.TextMethodToQuery(q.DomainQuery.GetMethod()), q.DomainQuery.GetDomain())
	case *instance.DomainSearchQuery_GeneratedQuery:
		return query.NewInstanceDomainGeneratedSearchQuery(q.GeneratedQuery.GetGenerated())
	case *instance.DomainSearchQuery_PrimaryQuery:
		return query.NewInstanceDomainPrimarySearchQuery(q.PrimaryQuery.GetPrimary())
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "INST-Ags42", "List.Query.Invalid")
	}
}

func ListTrustedDomainsRequestToModel(req *instance.ListTrustedDomainsRequest, defaults systemdefaults.SystemDefaults) (*query.InstanceTrustedDomainSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(defaults, req.GetPagination())
	if err != nil {
		return nil, err
	}

	queries, err := trustedDomainQueriesToModel(req.GetQueries())
	if err != nil {
		return nil, err
	}

	return &query.InstanceTrustedDomainSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: fieldNameToInstanceTrustedDomainColumn(req.GetSortingColumn()),
		},
		Queries: queries,
	}, nil
}

func trustedDomainQueriesToModel(queries []*instance.TrustedDomainSearchQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = trustedDomainQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func trustedDomainQueryToModel(searchQuery *instance.TrustedDomainSearchQuery) (query.SearchQuery, error) {
	switch q := searchQuery.GetQuery().(type) {
	case *instance.TrustedDomainSearchQuery_DomainQuery:
		return query.NewInstanceTrustedDomainDomainSearchQuery(object.TextMethodToQuery(q.DomainQuery.GetMethod()), q.DomainQuery.GetDomain())
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "INST-Ags42", "List.Query.Invalid")
	}
}

func trustedDomainsToPb(domains []*query.InstanceTrustedDomain) []*instance.TrustedDomain {
	d := make([]*instance.TrustedDomain, len(domains))
	for i, domain := range domains {
		d[i] = trustedDomainToPb(domain)
	}
	return d
}

func trustedDomainToPb(d *query.InstanceTrustedDomain) *instance.TrustedDomain {
	return &instance.TrustedDomain{
		Domain:       d.Domain,
		InstanceId:   d.InstanceID,
		CreationDate: timestamppb.New(d.CreationDate),
	}
}

func fieldNameToInstanceTrustedDomainColumn(fieldName instance.TrustedDomainFieldName) query.Column {
	switch fieldName {
	case instance.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_DOMAIN:
		return query.InstanceTrustedDomainDomainCol
	case instance.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_CREATION_DATE:
		return query.InstanceTrustedDomainCreationDateCol
	case instance.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_UNSPECIFIED:
		fallthrough
	default:
		return query.Column{}
	}
}
