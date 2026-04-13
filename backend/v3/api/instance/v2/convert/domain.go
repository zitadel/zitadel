package convert

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	filter_v2 "github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	instance_v2 "github.com/zitadel/zitadel/pkg/grpc/instance/v2"
	instance_v2beta "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

/*
 * Domain Model to GRPC v2Beta
 */

func DomainInstanceDomainListModelToGRPCBetaResponse(dms []*domain.InstanceDomain) []*instance_v2beta.Domain {
	toReturn := make([]*instance_v2beta.Domain, len(dms))
	for i, domain := range dms {
		isGenerated := domain.IsGenerated != nil && *domain.IsGenerated
		isPrimary := domain.IsPrimary != nil && *domain.IsPrimary
		toReturn[i] = &instance_v2beta.Domain{
			InstanceId:   domain.InstanceID,
			CreationDate: timestamppb.New(domain.CreatedAt),
			Domain:       domain.Domain,
			Primary:      isPrimary,
			Generated:    isGenerated,
		}
	}

	return toReturn
}

func TrustedDomainInstanceDomainListModelToGRPCBetaResponse(dms []*domain.InstanceDomain) []*instance_v2beta.TrustedDomain {
	toReturn := make([]*instance_v2beta.TrustedDomain, len(dms))
	for i, domain := range dms {
		toReturn[i] = &instance_v2beta.TrustedDomain{
			InstanceId:   domain.InstanceID,
			CreationDate: timestamppb.New(domain.CreatedAt),
			Domain:       domain.Domain,
		}
	}

	return toReturn
}

/*
 * Domain Model to GRPC v2
 */

func DomainInstanceDomainListModelToGRPCResponse(dms []*domain.InstanceDomain) []*instance_v2.CustomDomain {
	toReturn := make([]*instance_v2.CustomDomain, len(dms))
	for i, domain := range dms {
		isGenerated := domain.IsGenerated != nil && *domain.IsGenerated
		isPrimary := domain.IsPrimary != nil && *domain.IsPrimary
		toReturn[i] = &instance_v2.CustomDomain{
			InstanceId:   domain.InstanceID,
			CreationDate: timestamppb.New(domain.CreatedAt),
			Domain:       domain.Domain,
			Primary:      isPrimary,
			Generated:    isGenerated,
		}
	}

	return toReturn
}

func TrustedDomainInstanceDomainListModelToGRPCResponse(dms []*domain.InstanceDomain) []*instance_v2.TrustedDomain {
	toReturn := make([]*instance_v2.TrustedDomain, len(dms))
	for i, domain := range dms {
		toReturn[i] = &instance_v2.TrustedDomain{
			InstanceId:   domain.InstanceID,
			CreationDate: timestamppb.New(domain.CreatedAt),
			Domain:       domain.Domain,
		}
	}

	return toReturn
}

/*
 * GRPC v2Beta Requests to GRPC v2
 */

func ListCustomDomainsBetaRequestToV2Request(in *instance_v2beta.ListCustomDomainsRequest) *instance_v2.ListCustomDomainsRequest {
	return &instance_v2.ListCustomDomainsRequest{
		InstanceId: in.GetInstanceId(),
		Pagination: &filter_v2.PaginationRequest{
			Offset: in.GetPagination().GetOffset(),
			Limit:  in.GetPagination().GetLimit(),
			Asc:    in.GetPagination().GetAsc(),
		},
		SortingColumn: listCustomDomainsBetaSortingColToV2Request(in.GetSortingColumn()),
		Filters:       listCustomDomainsQueriesToV2Request(in.GetQueries()),
	}
}

func listCustomDomainsBetaSortingColToV2Request(domainFieldName instance_v2beta.DomainFieldName) instance_v2.DomainFieldName {
	switch domainFieldName {
	case instance_v2beta.DomainFieldName_DOMAIN_FIELD_NAME_CREATION_DATE:
		return instance_v2.DomainFieldName_DOMAIN_FIELD_NAME_CREATION_DATE
	case instance_v2beta.DomainFieldName_DOMAIN_FIELD_NAME_DOMAIN:
		return instance_v2.DomainFieldName_DOMAIN_FIELD_NAME_DOMAIN
	case instance_v2beta.DomainFieldName_DOMAIN_FIELD_NAME_GENERATED:
		return instance_v2.DomainFieldName_DOMAIN_FIELD_NAME_GENERATED
	case instance_v2beta.DomainFieldName_DOMAIN_FIELD_NAME_PRIMARY:
		return instance_v2.DomainFieldName_DOMAIN_FIELD_NAME_PRIMARY
	case instance_v2beta.DomainFieldName_DOMAIN_FIELD_NAME_UNSPECIFIED:
		fallthrough
	default:
		return instance_v2.DomainFieldName_DOMAIN_FIELD_NAME_UNSPECIFIED
	}
}

func listCustomDomainsQueriesToV2Request(domainSearchQuery []*instance_v2beta.DomainSearchQuery) []*instance_v2.CustomDomainFilter {
	toReturn := make([]*instance_v2.CustomDomainFilter, len(domainSearchQuery))

	for i, query := range domainSearchQuery {
		switch assertedType := query.GetQuery().(type) {
		case *instance_v2beta.DomainSearchQuery_DomainQuery:
			filter := &instance_v2.CustomDomainFilter_DomainFilter{
				DomainFilter: &instance_v2.DomainFilter{
					Domain: assertedType.DomainQuery.GetDomain(),
					Method: assertedType.DomainQuery.GetMethod(),
				},
			}
			toReturn[i] = &instance_v2.CustomDomainFilter{
				Filter: filter,
			}
		case *instance_v2beta.DomainSearchQuery_GeneratedQuery:
			filter := &instance_v2.CustomDomainFilter_GeneratedFilter{
				GeneratedFilter: assertedType.GeneratedQuery.GetGenerated(),
			}
			toReturn[i] = &instance_v2.CustomDomainFilter{
				Filter: filter,
			}
		case *instance_v2beta.DomainSearchQuery_PrimaryQuery:
			filter := &instance_v2.CustomDomainFilter_PrimaryFilter{
				PrimaryFilter: assertedType.PrimaryQuery.GetPrimary(),
			}
			toReturn[i] = &instance_v2.CustomDomainFilter{
				Filter: filter,
			}
		}
	}
	return toReturn
}

func ListTrustedDomainsBetaRequestToV2Request(in *instance_v2beta.ListTrustedDomainsRequest) *instance_v2.ListTrustedDomainsRequest {
	return &instance_v2.ListTrustedDomainsRequest{
		InstanceId: in.GetInstanceId(),
		Pagination: &filter_v2.PaginationRequest{
			Offset: in.GetPagination().GetOffset(),
			Limit:  in.GetPagination().GetLimit(),
			Asc:    in.GetPagination().GetAsc(),
		},
		SortingColumn: listTrustedDomainsBetaSortingColToV2Request(in.GetSortingColumn()),
		Filters:       listTrustedDomainsQueriesToV2Request(in.GetQueries()),
	}
}

func listTrustedDomainsBetaSortingColToV2Request(domainFieldName instance_v2beta.TrustedDomainFieldName) instance_v2.TrustedDomainFieldName {
	switch domainFieldName {
	case instance_v2beta.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_CREATION_DATE:
		return instance_v2.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_CREATION_DATE
	case instance_v2beta.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_DOMAIN:
		return instance_v2.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_DOMAIN
	case instance_v2beta.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_UNSPECIFIED:
		fallthrough
	default:
		return instance_v2.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_UNSPECIFIED
	}
}

func listTrustedDomainsQueriesToV2Request(domainSearchQuery []*instance_v2beta.TrustedDomainSearchQuery) []*instance_v2.TrustedDomainFilter {
	toReturn := make([]*instance_v2.TrustedDomainFilter, len(domainSearchQuery))

	for i, query := range domainSearchQuery {
		if assertedType, ok := query.GetQuery().(*instance_v2beta.TrustedDomainSearchQuery_DomainQuery); ok {
			filter := &instance_v2.TrustedDomainFilter_DomainFilter{
				DomainFilter: &instance_v2.DomainFilter{
					Domain: assertedType.DomainQuery.GetDomain(),
					Method: assertedType.DomainQuery.GetMethod(),
				},
			}
			toReturn[i] = &instance_v2.TrustedDomainFilter{
				Filter: filter,
			}
		}
	}
	return toReturn
}
