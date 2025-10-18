package instance

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/cmd/build"
	"github.com/zitadel/zitadel/internal/api/grpc/filter/v2"
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
		Id:            inst.ID,
		Name:          inst.Name,
		CustomDomains: DomainsToPb(inst.Domains),
		Version:       build.Version(),
		ChangeDate:    timestamppb.New(inst.ChangeDate),
		CreationDate:  timestamppb.New(inst.CreationDate),
	}
}

func DomainsToPb(domains []*query.InstanceDomain) []*instance.CustomDomain {
	d := []*instance.CustomDomain{}
	for _, dm := range domains {
		pbDomain := DomainToPb(dm)
		d = append(d, pbDomain)
	}
	return d
}

func DomainToPb(d *query.InstanceDomain) *instance.CustomDomain {
	return &instance.CustomDomain{
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

	queries, err := filtersToQueries(req.GetFilters())
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

func filtersToQueries(filters []*instance.Filter) (_ []query.SearchQuery, err error) {
	q := []query.SearchQuery{}
	for _, filter := range filters {
		model, err := instanceFilterToQuery(filter)
		if err != nil {
			return nil, err
		}
		q = append(q, model)
	}
	return q, nil
}

func instanceFilterToQuery(filter *instance.Filter) (query.SearchQuery, error) {
	switch q := filter.GetFilter().(type) {
	case *instance.Filter_InIdsFilter:
		return query.NewInstanceIDsListSearchQuery(q.InIdsFilter.GetIds()...)
	case *instance.Filter_CustomDomainsFilter:
		return query.NewInstanceDomainsListSearchQuery(q.CustomDomainsFilter.GetDomains()...)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "INST-3m0se", "List.Query.Invalid")
	}
}

func ListCustomDomainsRequestToModel(req *instance.ListCustomDomainsRequest, defaults systemdefaults.SystemDefaults) (*query.InstanceDomainSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(defaults, req.GetPagination())
	if err != nil {
		return nil, err
	}

	queries, err := customDomainFiltersToQueries(req.GetFilters())
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

func customDomainFiltersToQueries(filters []*instance.CustomDomainFilter) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(filters))
	for i, filter := range filters {
		q[i], err = customDomainFilterToQuery(filter)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func customDomainFilterToQuery(filter *instance.CustomDomainFilter) (query.SearchQuery, error) {
	switch q := filter.GetFilter().(type) {
	case *instance.CustomDomainFilter_DomainFilter:
		return query.NewInstanceDomainDomainSearchQuery(object.TextMethodToQuery(q.DomainFilter.GetMethod()), q.DomainFilter.GetDomain())
	case *instance.CustomDomainFilter_GeneratedFilter:
		return query.NewInstanceDomainGeneratedSearchQuery(q.GeneratedFilter)
	case *instance.CustomDomainFilter_PrimaryFilter:
		return query.NewInstanceDomainPrimarySearchQuery(q.PrimaryFilter)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "INST-Ags42", "List.Query.Invalid")
	}
}

func ListTrustedDomainsRequestToModel(req *instance.ListTrustedDomainsRequest, defaults systemdefaults.SystemDefaults) (*query.InstanceTrustedDomainSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(defaults, req.GetPagination())
	if err != nil {
		return nil, err
	}

	queries, err := trustedDomainFiltersToQueries(req.GetFilters())
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

func trustedDomainFiltersToQueries(filters []*instance.TrustedDomainFilter) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(filters))
	for i, filter := range filters {
		q[i], err = trustedDomainQueryToModel(filter)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func trustedDomainQueryToModel(filter *instance.TrustedDomainFilter) (query.SearchQuery, error) {
	switch q := filter.GetFilter().(type) {
	case *instance.TrustedDomainFilter_DomainFilter:
		return query.NewInstanceTrustedDomainDomainSearchQuery(object.TextMethodToQuery(q.DomainFilter.GetMethod()), q.DomainFilter.GetDomain())
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
