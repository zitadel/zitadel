package org

import (
	"context"

	"connectrpc.com/connect"

	orgv2 "github.com/zitadel/zitadel/backend/v3/api/org/v2"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/filter/v2"
	"github.com/zitadel/zitadel/internal/api/grpc/metadata/v2"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	filter_pb "github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
)

func (s *Server) ListOrganizations(ctx context.Context, req *connect.Request[org.ListOrganizationsRequest]) (*connect.Response[org.ListOrganizationsResponse], error) {
	if authz.GetFeatures(ctx).EnableRelationalTables {
		return orgv2.ListOrganizations(ctx, req)
	}

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

func (s *Server) ListOrganizationMetadata(ctx context.Context, request *connect.Request[org.ListOrganizationMetadataRequest]) (*connect.Response[org.ListOrganizationMetadataResponse], error) {
	metadataQueries, err := listOrgMetadataToDomain(s.systemDefaults, request.Msg)
	if err != nil {
		return nil, err
	}
	res, err := s.query.SearchOrgMetadata(ctx, true, request.Msg.GetOrganizationId(), metadataQueries, false, true)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.ListOrganizationMetadataResponse{
		Metadata: metadata.OrgMetadataListToPb(res.Metadata),
		Pagination: &filter_pb.PaginationResponse{
			TotalResult:  res.Count,
			AppliedLimit: uint64(request.Msg.GetPagination().GetLimit()),
		},
	}), nil
}

func (s *Server) ListOrganizationDomains(ctx context.Context, req *connect.Request[org.ListOrganizationDomainsRequest]) (*connect.Response[org.ListOrganizationDomainsResponse], error) {
	queries, err := listOrgDomainsRequestToDomain(s.systemDefaults, req.Msg)
	if err != nil {
		return nil, err
	}
	orgIDQuery, err := query.NewOrgDomainOrgIDSearchQuery(req.Msg.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	queries.Queries = append(queries.Queries, orgIDQuery)

	domains, err := s.query.SearchOrgDomains(ctx, queries, false, true)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.ListOrganizationDomainsResponse{
		Domains: domainsToPb(domains.Domains),
		Pagination: &filter_pb.PaginationResponse{
			TotalResult:  domains.Count,
			AppliedLimit: uint64(req.Msg.GetPagination().GetLimit()),
		},
	}), nil
}

func listOrgDomainsRequestToDomain(systemDefaults systemdefaults.SystemDefaults, request *org.ListOrganizationDomainsRequest) (*query.OrgDomainSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(systemDefaults, request.Pagination)
	if err != nil {
		return nil, err
	}
	queries, err := DomainQueriesToModel(request.Filters)
	if err != nil {
		return nil, err
	}
	return &query.OrgDomainSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: fieldNameToOrganizationDomainColumn(request.GetSortingColumn()),
		},
		Queries: queries,
	}, nil
}

func fieldNameToOrganizationDomainColumn(column org.DomainFieldName) query.Column {
	switch column {
	case org.DomainFieldName_DOMAIN_FIELD_NAME_NAME:
		return query.OrgDomainDomainCol
	case org.DomainFieldName_DOMAIN_FIELD_NAME_CREATION_DATE:
		return query.OrgDomainCreationDateCol
	case org.DomainFieldName_DOMAIN_FIELD_NAME_UNSPECIFIED:
		return query.Column{}
	default:
		return query.Column{}
	}
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
	case org.OrganizationFieldName_ORGANIZATION_FIELD_NAME_CREATION_DATE:
		return query.OrgColumnCreationDate
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

func listOrgMetadataToDomain(systemDefaults systemdefaults.SystemDefaults, request *org.ListOrganizationMetadataRequest) (*query.OrgMetadataSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(systemDefaults, request.Pagination)
	if err != nil {
		return nil, err
	}
	queries, err := metadata.OrgMetadataQueriesToQuery(request.GetFilters())
	if err != nil {
		return nil, err
	}
	return &query.OrgMetadataSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: queries,
	}, nil
}

func DomainQueriesToModel(queries []*org.DomainSearchFilter) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = DomainQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func DomainQueryToModel(searchQuery *org.DomainSearchFilter) (query.SearchQuery, error) {
	switch q := searchQuery.Filter.(type) {
	case *org.DomainSearchFilter_DomainFilter:
		return query.NewOrgDomainDomainSearchQuery(object.TextMethodToQuery(q.DomainFilter.Method), q.DomainFilter.GetDomain())
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Ags89", "List.Query.Invalid")
	}
}

func domainsToPb(domains []*query.Domain) []*org.Domain {
	d := make([]*org.Domain, len(domains))
	for i, domain := range domains {
		d[i] = domainToPb(domain)
	}
	return d
}

func domainToPb(d *query.Domain) *org.Domain {
	return &org.Domain{
		OrganizationId: d.OrgID,
		Domain:         d.Domain,
		IsVerified:     d.IsVerified,
		IsPrimary:      d.IsPrimary,
		ValidationType: domainValidationTypeToPb(d.ValidationType),
	}
}

func domainValidationTypeToPb(validationType domain.OrgDomainValidationType) org.DomainValidationType {
	switch validationType {
	case domain.OrgDomainValidationTypeDNS:
		return org.DomainValidationType_DOMAIN_VALIDATION_TYPE_DNS
	case domain.OrgDomainValidationTypeHTTP:
		return org.DomainValidationType_DOMAIN_VALIDATION_TYPE_HTTP
	case domain.OrgDomainValidationTypeUnspecified:
		return org.DomainValidationType_DOMAIN_VALIDATION_TYPE_UNSPECIFIED
	default:
		return org.DomainValidationType_DOMAIN_VALIDATION_TYPE_UNSPECIFIED
	}
}
