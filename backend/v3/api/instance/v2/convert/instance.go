package convert

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/cmd/build"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	instance_v2 "github.com/zitadel/zitadel/pkg/grpc/instance/v2"
	instance_v2beta "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

/*
 * Domain Model to GRPC v2Beta
 */

func DomainInstanceListModelToGRPCBetaResponse(instances []*domain.Instance) []*instance_v2beta.Instance {
	toReturn := make([]*instance_v2beta.Instance, len(instances))

	for i, inst := range instances {
		toReturn[i] = DomainInstanceModelToGRPCBetaResponse(inst)
	}

	return toReturn
}

func DomainInstanceModelToGRPCBetaResponse(inst *domain.Instance) *instance_v2beta.Instance {
	return &instance_v2beta.Instance{
		Id:           inst.ID,
		ChangeDate:   timestamppb.New(inst.UpdatedAt),
		CreationDate: timestamppb.New(inst.CreatedAt),
		State:        instance_v2beta.State_STATE_RUNNING, // TODO(IAM-Marco): Not sure what to put here
		Name:         inst.Name,
		Version:      build.Version(),
		Domains:      domainInstanceDomainListModelToGRPCBetaResponse(inst.Domains),
	}
}

func domainInstanceDomainListModelToGRPCBetaResponse(dms []*domain.InstanceDomain) []*instance_v2beta.Domain {
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

/*
 * Domain Model to GRPC v2
 */

func DomainInstanceListModelToGRPCResponse(instances []*domain.Instance) []*instance_v2.Instance {
	toReturn := make([]*instance_v2.Instance, len(instances))

	for i, inst := range instances {
		toReturn[i] = DomainInstanceModelToGRPCResponse(inst)
	}

	return toReturn
}

func DomainInstanceModelToGRPCResponse(inst *domain.Instance) *instance_v2.Instance {
	return &instance_v2.Instance{
		Id:            inst.ID,
		ChangeDate:    timestamppb.New(inst.UpdatedAt),
		CreationDate:  timestamppb.New(inst.CreatedAt),
		State:         instance_v2.State_STATE_RUNNING, // TODO(IAM-Marco): Not sure what to put here
		Name:          inst.Name,
		Version:       build.Version(),
		CustomDomains: domainInstanceDomainListModelToGRPCResponse(inst.Domains),
	}
}

func domainInstanceDomainListModelToGRPCResponse(dms []*domain.InstanceDomain) []*instance_v2.CustomDomain {
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

/*
 * GRPC v2Beta Requests to GRPC v2
 */

func ListInstancesBetaRequestToV2Request(in *instance_v2beta.ListInstancesRequest) *instance_v2.ListInstancesRequest {
	return &instance_v2.ListInstancesRequest{
		Pagination: &filter.PaginationRequest{
			Offset: in.GetPagination().GetOffset(),
			Limit:  in.GetPagination().GetLimit(),
			Asc:    in.GetPagination().GetAsc(),
		},
		SortingColumn: listInstancesBetaSortingColToV2Request(in.SortingColumn),
		Filters:       listInstancesQueriesToV2Request(in.GetQueries()),
	}
}

func listInstancesQueriesToV2Request(queries []*instance_v2beta.Query) []*instance_v2.Filter {
	toReturn := make([]*instance_v2.Filter, len(queries))
	for i, query := range queries {
		switch assertedQuery := query.GetQuery().(type) {
		case *instance_v2beta.Query_DomainQuery:
			filter := &instance_v2.Filter_CustomDomainsFilter{
				CustomDomainsFilter: &instance_v2.CustomDomainsFilter{
					Domains: assertedQuery.DomainQuery.GetDomains(),
				},
			}
			toReturn[i] = &instance_v2.Filter{Filter: filter}
		case *instance_v2beta.Query_IdQuery:
			filter := &instance_v2.Filter_InIdsFilter{
				InIdsFilter: &filter.InIDsFilter{
					Ids: assertedQuery.IdQuery.GetIds(),
				},
			}
			toReturn[i] = &instance_v2.Filter{Filter: filter}
		}
	}
	return toReturn
}

func listInstancesBetaSortingColToV2Request(fieldName *instance_v2beta.FieldName) instance_v2.FieldName {
	if fieldName == nil {
		return instance_v2.FieldName_FIELD_NAME_UNSPECIFIED
	}

	switch *fieldName {
	case instance_v2beta.FieldName_FIELD_NAME_CREATION_DATE:
		return instance_v2.FieldName_FIELD_NAME_CREATION_DATE
	case instance_v2beta.FieldName_FIELD_NAME_ID:
		return instance_v2.FieldName_FIELD_NAME_ID
	case instance_v2beta.FieldName_FIELD_NAME_NAME:
		return instance_v2.FieldName_FIELD_NAME_NAME
	case instance_v2beta.FieldName_FIELD_NAME_UNSPECIFIED:
		fallthrough
	default:
		return instance_v2.FieldName_FIELD_NAME_UNSPECIFIED
	}
}
