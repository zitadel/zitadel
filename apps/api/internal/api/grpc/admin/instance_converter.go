package admin

import (
	instance_grpc "github.com/zitadel/zitadel/internal/api/grpc/instance"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/instance"
)

func ListInstanceDomainsRequestToModel(req *admin_pb.ListInstanceDomainsRequest) (*query.InstanceDomainSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := instance_grpc.DomainQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.InstanceDomainSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: fieldNameToInstanceDomainColumn(req.SortingColumn),
		},
		Queries: queries,
	}, nil
}

func fieldNameToInstanceDomainColumn(fieldName instance.DomainFieldName) query.Column {
	switch fieldName {
	case instance.DomainFieldName_DOMAIN_FIELD_NAME_DOMAIN:
		return query.InstanceDomainDomainCol
	case instance.DomainFieldName_DOMAIN_FIELD_NAME_PRIMARY:
		return query.InstanceDomainIsPrimaryCol
	case instance.DomainFieldName_DOMAIN_FIELD_NAME_GENERATED:
		return query.InstanceDomainIsGeneratedCol
	case instance.DomainFieldName_DOMAIN_FIELD_NAME_CREATION_DATE:
		return query.InstanceDomainCreationDateCol
	default:
		return query.Column{}
	}
}

func ListInstanceTrustedDomainsRequestToModel(req *admin_pb.ListInstanceTrustedDomainsRequest) (*query.InstanceTrustedDomainSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := instance_grpc.TrustedDomainQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.InstanceTrustedDomainSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: fieldNameToInstanceTrustedDomainColumn(req.SortingColumn),
		},
		Queries: queries,
	}, nil
}

func fieldNameToInstanceTrustedDomainColumn(fieldName instance.DomainFieldName) query.Column {
	switch fieldName {
	case instance.DomainFieldName_DOMAIN_FIELD_NAME_DOMAIN:
		return query.InstanceTrustedDomainDomainCol
	case instance.DomainFieldName_DOMAIN_FIELD_NAME_CREATION_DATE:
		return query.InstanceTrustedDomainCreationDateCol
	case instance.DomainFieldName_DOMAIN_FIELD_NAME_UNSPECIFIED,
		instance.DomainFieldName_DOMAIN_FIELD_NAME_PRIMARY,
		instance.DomainFieldName_DOMAIN_FIELD_NAME_GENERATED:
		return query.InstanceTrustedDomainCreationDateCol
	default:
		return query.Column{}
	}
}
