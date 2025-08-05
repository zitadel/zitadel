package org

import (
	"context"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
)

func (s *Server) ListOrganizations(ctx context.Context, req *connect.Request[org.ListOrganizationsRequest]) (*connect.Response[org.ListOrganizationsResponse], error) {
	queries, err := listOrgRequestToModel(ctx, req)
	if err != nil {
		return nil, err
	}
	orgs, err := s.query.SearchOrgs(ctx, queries.Msg, s.checkPermission)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.ListOrganizationsResponse{
		Result:  organizationsToPb(orgs.Orgs),
		Details: object.ToListDetails(orgs.SearchResponse),
	}), nil
}

func listOrgRequestToModel(ctx context.Context, req *connect.Request[org.ListOrganizationsRequest]) (*connect.Response[query.OrgSearchQueries], error) {
	offset, limit, asc := object.ListQueryToQuery(req.Msg.Query)
	queries, err := orgQueriesToQuery(ctx, req.Msg.Queries)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&query.OrgSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			SortingColumn: fieldNameToOrganizationColumn(req.Msg.SortingColumn),
			Asc:           asc,
		},
		Queries: queries,
	}), nil
}

func orgQueriesToQuery(ctx context.Context, queries []*org.SearchQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = orgQueryToQuery(ctx, query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func orgQueryToQuery(ctx context.Context, orgQuery *org.SearchQuery) (query.SearchQuery, error) {
	switch q := orgQuery.Query.(type) {
	case *org.SearchQuery_DomainQuery:
		return query.NewOrgVerifiedDomainSearchQuery(object.TextMethodToQuery(q.DomainQuery.Method), q.DomainQuery.Domain)
	case *org.SearchQuery_NameQuery:
		return query.NewOrgNameSearchQuery(object.TextMethodToQuery(q.NameQuery.Method), q.NameQuery.Name)
	case *org.SearchQuery_StateQuery:
		return query.NewOrgStateSearchQuery(orgStateToDomain(q.StateQuery.State))
	case *org.SearchQuery_IdQuery:
		return query.NewOrgIDSearchQuery(q.IdQuery.Id)
	case *org.SearchQuery_DefaultQuery:
		return query.NewOrgIDSearchQuery(authz.GetInstance(ctx).DefaultOrganisationID())
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
			CreationDate:  organization.CreationDate,
		}),
		State: orgStateToPb(organization.State),
	}
}
