package convert

import (
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	v2_org "github.com/zitadel/zitadel/pkg/grpc/org/v2"
	v2beta_org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

func OrganizationBetaRequestToV2Request(in *v2beta_org.ListOrganizationsRequest) *v2_org.ListOrganizationsRequest {
	return &v2_org.ListOrganizationsRequest{
		Query: &object.ListQuery{
			Offset: in.GetPagination().GetOffset(),
			Limit:  in.GetPagination().GetLimit(),
			Asc:    in.GetPagination().GetAsc(),
		},
		SortingColumn: organizationSortingColumnBetaToV2(in.GetSortingColumn()),
		Queries:       organizationQueriesBetaToV2(in.GetFilter()),
	}
}

func organizationSortingColumnBetaToV2(sc v2beta_org.OrgFieldName) v2_org.OrganizationFieldName {
	switch sc {
	case v2beta_org.OrgFieldName_ORG_FIELD_NAME_NAME:
		return v2_org.OrganizationFieldName_ORGANIZATION_FIELD_NAME_NAME
	case v2beta_org.OrgFieldName_ORG_FIELD_NAME_CREATION_DATE, v2beta_org.OrgFieldName_ORG_FIELD_NAME_UNSPECIFIED:
		fallthrough
	default:
		return v2_org.OrganizationFieldName_ORGANIZATION_FIELD_NAME_UNSPECIFIED
	}
}

func organizationQueriesBetaToV2(queries []*v2beta_org.OrganizationSearchFilter) []*v2_org.SearchQuery {
	toReturn := make([]*v2_org.SearchQuery, len(queries))

	for i, query := range queries {
		toReturn[i] = organizationQueryBetaToV2(query)
	}

	return toReturn
}

func organizationQueryBetaToV2(query *v2beta_org.OrganizationSearchFilter) *v2_org.SearchQuery {
	toReturn := &v2_org.SearchQuery{}

	switch assertedType := query.GetFilter().(type) {
	case *v2beta_org.OrganizationSearchFilter_DomainFilter:
		toReturn.Query = &v2_org.SearchQuery_DomainQuery{
			DomainQuery: &v2_org.OrganizationDomainQuery{
				Domain: assertedType.DomainFilter.GetDomain(),
				Method: TextQueryMethodBetaToV2(assertedType.DomainFilter.GetMethod()),
			},
		}

	case *v2beta_org.OrganizationSearchFilter_IdFilter:
		toReturn.Query = &v2_org.SearchQuery_IdQuery{
			IdQuery: &v2_org.OrganizationIDQuery{
				Id: assertedType.IdFilter.GetId(),
			},
		}
	case *v2beta_org.OrganizationSearchFilter_NameFilter:
		toReturn.Query = &v2_org.SearchQuery_NameQuery{
			NameQuery: &v2_org.OrganizationNameQuery{
				Name:   assertedType.NameFilter.GetName(),
				Method: TextQueryMethodBetaToV2(assertedType.NameFilter.GetMethod()),
			},
		}
	case *v2beta_org.OrganizationSearchFilter_StateFilter:
		toReturn.Query = &v2_org.SearchQuery_StateQuery{
			StateQuery: &v2_org.OrganizationStateQuery{
				State: organizationStateBetaToV2(assertedType.StateFilter.GetState()),
			},
		}
	default:
		return toReturn
	}
	return toReturn
}

func organizationStateBetaToV2(in v2beta_org.OrgState) v2_org.OrganizationState {
	switch in {
	case v2beta_org.OrgState_ORG_STATE_ACTIVE:
		return v2_org.OrganizationState_ORGANIZATION_STATE_ACTIVE
	case v2beta_org.OrgState_ORG_STATE_INACTIVE:
		return v2_org.OrganizationState_ORGANIZATION_STATE_INACTIVE
	case v2beta_org.OrgState_ORG_STATE_UNSPECIFIED, v2beta_org.OrgState_ORG_STATE_REMOVED:
		fallthrough
	default:
		return v2_org.OrganizationState_ORGANIZATION_STATE_UNSPECIFIED
	}
}
