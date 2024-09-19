package org

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/query"
	"github.com/zitadel/zitadel/v2/internal/zerrors"
	"github.com/zitadel/zitadel/v2/pkg/grpc/org/v2"
)

func (s *Server) ListOrganizations(ctx context.Context, req *org.ListOrganizationsRequest) (*org.ListOrganizationsResponse, error) {
	queries, err := listOrgRequestToModel(req)
	if err != nil {
		return nil, err
	}
	orgs, err := s.query.SearchOrgs(ctx, queries, s.checkPermission)
	if err != nil {
		return nil, err
	}
	return &org.ListOrganizationsResponse{
		Result:  organizationsToPb(orgs.Orgs),
		Details: object.ToListDetails(orgs.SearchResponse),
	}, nil
}

func listOrgRequestToModel(req *org.ListOrganizationsRequest) (*query.OrgSearchQueries, error) {
	offset, limit, asc := object.ListQueryToQuery(req.Query)
	queries, err := orgQueriesToQuery(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.OrgSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			SortingColumn: fieldNameToOrganizationColumn(req.SortingColumn),
			Asc:           asc,
		},
		Queries: queries,
	}, nil
}

func orgQueriesToQuery(queries []*org.SearchQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = orgQueryToQuery(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func orgQueryToQuery(orgQuery *org.SearchQuery) (query.SearchQuery, error) {
	switch q := orgQuery.Query.(type) {
	case *org.SearchQuery_DomainQuery:
		return query.NewOrgDomainSearchQuery(object.TextMethodToQuery(q.DomainQuery.Method), q.DomainQuery.Domain)
	case *org.SearchQuery_NameQuery:
		return query.NewOrgNameSearchQuery(object.TextMethodToQuery(q.NameQuery.Method), q.NameQuery.Name)
	case *org.SearchQuery_StateQuery:
		return query.NewOrgStateSearchQuery(orgStateToDomain(q.StateQuery.State))
	case *org.SearchQuery_IdQuery:
		return query.NewOrgIDSearchQuery(q.IdQuery.Id)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-vR9nC", "List.Query.Invalid")
	}
}

func orgStateToPb(state domain.OrgState) org.OrganizationState {
	switch state {
	case domain.OrgStateActive:
		return org.OrganizationState_ORGANIZATION_STATE_ACTIVE
	case domain.OrgStateInactive:
		return org.OrganizationState_ORGANIZATION_STATE_INACTIVE
	case domain.OrgStateRemoved:
		return org.OrganizationState_ORGANIZATION_STATE_REMOVED
	case domain.OrgStateUnspecified:
		fallthrough
	default:
		return org.OrganizationState_ORGANIZATION_STATE_UNSPECIFIED
	}
}

func orgStateToDomain(state org.OrganizationState) domain.OrgState {
	switch state {
	case org.OrganizationState_ORGANIZATION_STATE_ACTIVE:
		return domain.OrgStateActive
	case org.OrganizationState_ORGANIZATION_STATE_INACTIVE:
		return domain.OrgStateInactive
	case org.OrganizationState_ORGANIZATION_STATE_REMOVED:
		return domain.OrgStateRemoved
	case org.OrganizationState_ORGANIZATION_STATE_UNSPECIFIED:
		fallthrough
	default:
		return domain.OrgStateUnspecified
	}
}

func fieldNameToOrganizationColumn(fieldName org.OrganizationFieldName) query.Column {
	switch fieldName {
	case org.OrganizationFieldName_ORGANIZATION_FIELD_NAME_NAME:
		return query.OrgColumnName
	case org.OrganizationFieldName_ORGANIZATION_FIELD_NAME_UNSPECIFIED:
		return query.Column{}
	default:
		return query.Column{}
	}
}

func organizationsToPb(orgs []*query.Org) []*org.Organization {
	o := make([]*org.Organization, len(orgs))
	for i, org := range orgs {
		o[i] = organizationToPb(org)
	}
	return o
}

func organizationToPb(organization *query.Org) *org.Organization {
	return &org.Organization{
		Id:            organization.ID,
		Name:          organization.Name,
		PrimaryDomain: organization.Domain,
		Details: object.DomainToDetailsPb(&domain.ObjectDetails{
			Sequence:      organization.Sequence,
			EventDate:     organization.ChangeDate,
			ResourceOwner: organization.ResourceOwner,
		}),
		State: orgStateToPb(organization.State),
	}
}
