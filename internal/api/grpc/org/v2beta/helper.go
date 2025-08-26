package org

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	// TODO fix below
	filter "github.com/zitadel/zitadel/internal/api/grpc/filter/v2beta"
	metadata "github.com/zitadel/zitadel/internal/api/grpc/metadata/v2beta"
	v2beta_object "github.com/zitadel/zitadel/internal/api/grpc/object/v2beta"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	v2beta "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	v2beta_org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

// NOTE: most of this code is copied from `internal/api/grpc/admin/*`, as we will eventually axe the previous versons of the API,
// we will have code duplication until then

func listOrgRequestToModel(systemDefaults systemdefaults.SystemDefaults, request *v2beta_org.ListOrganizationsRequest) (*query.OrgSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(systemDefaults, request.Pagination)
	if err != nil {
		return nil, err
	}
	queries, err := OrgQueriesToModel(request.Filter)
	if err != nil {
		return nil, err
	}
	return &query.OrgSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			SortingColumn: FieldNameToOrgColumn(request.SortingColumn),
			Asc:           asc,
		},
		Queries: queries,
	}, nil
}

func OrganizationViewToPb(org *query.Org) *v2beta_org.Organization {
	return &v2beta_org.Organization{
		Id:            org.ID,
		State:         OrgStateToPb(org.State),
		Name:          org.Name,
		PrimaryDomain: org.Domain,
		CreationDate:  timestamppb.New(org.CreationDate),
		ChangedDate:   timestamppb.New(org.ChangeDate),
	}
}

func OrgStateToPb(state domain.OrgState) v2beta_org.OrgState {
	switch state {
	case domain.OrgStateActive:
		return v2beta_org.OrgState_ORG_STATE_ACTIVE
	case domain.OrgStateInactive:
		return v2beta_org.OrgState_ORG_STATE_INACTIVE
	case domain.OrgStateRemoved:
		// added to please golangci-lint
		return v2beta_org.OrgState_ORG_STATE_REMOVED
	case domain.OrgStateUnspecified:
		// added to please golangci-lint
		return v2beta_org.OrgState_ORG_STATE_UNSPECIFIED
	default:
		return v2beta_org.OrgState_ORG_STATE_UNSPECIFIED
	}
}

func createdOrganizationToPb(createdOrg *command.CreatedOrg) (_ *connect.Response[org.CreateOrganizationResponse], err error) {
	admins := make([]*org.OrganizationAdmin, len(createdOrg.OrgAdmins))
	for i, admin := range createdOrg.OrgAdmins {
		switch admin := admin.(type) {
		case *command.CreatedOrgAdmin:
			admins[i] = &org.OrganizationAdmin{
				OrganizationAdmin: &org.OrganizationAdmin_CreatedAdmin{
					CreatedAdmin: &org.CreatedAdmin{
						UserId:    admin.ID,
						EmailCode: admin.EmailCode,
						PhoneCode: admin.PhoneCode,
					},
				},
			}
		case *command.AssignedOrgAdmin:
			admins[i] = &org.OrganizationAdmin{
				OrganizationAdmin: &org.OrganizationAdmin_AssignedAdmin{
					AssignedAdmin: &org.AssignedAdmin{
						UserId: admin.ID,
					},
				},
			}
		}
	}
	return connect.NewResponse(&org.CreateOrganizationResponse{
		CreationDate:       timestamppb.New(createdOrg.ObjectDetails.EventDate),
		Id:                 createdOrg.ObjectDetails.ResourceOwner,
		OrganizationAdmins: admins,
	}), nil
}

func OrgViewsToPb(orgs []*query.Org) []*v2beta_org.Organization {
	o := make([]*v2beta_org.Organization, len(orgs))
	for i, org := range orgs {
		o[i] = OrganizationViewToPb(org)
	}
	return o
}

func OrgQueriesToModel(queries []*v2beta_org.OrganizationSearchFilter) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = OrgQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func OrgQueryToModel(apiQuery *v2beta_org.OrganizationSearchFilter) (query.SearchQuery, error) {
	switch q := apiQuery.Filter.(type) {
	case *v2beta_org.OrganizationSearchFilter_DomainFilter:
		return query.NewOrgVerifiedDomainSearchQuery(v2beta_object.TextMethodToQuery(q.DomainFilter.Method), q.DomainFilter.Domain)
	case *v2beta_org.OrganizationSearchFilter_NameFilter:
		return query.NewOrgNameSearchQuery(v2beta_object.TextMethodToQuery(q.NameFilter.Method), q.NameFilter.Name)
	case *v2beta_org.OrganizationSearchFilter_StateFilter:
		return query.NewOrgStateSearchQuery(OrgStateToDomain(q.StateFilter.State))
	case *v2beta_org.OrganizationSearchFilter_IdFilter:
		return query.NewOrgIDSearchQuery(q.IdFilter.Id)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-vR9nC", "List.Query.Invalid")
	}
}

func OrgStateToDomain(state v2beta_org.OrgState) domain.OrgState {
	switch state {
	case v2beta_org.OrgState_ORG_STATE_ACTIVE:
		return domain.OrgStateActive
	case v2beta_org.OrgState_ORG_STATE_INACTIVE:
		return domain.OrgStateInactive
	case v2beta_org.OrgState_ORG_STATE_REMOVED:
		// added to please golangci-lint
		return domain.OrgStateRemoved
	case v2beta_org.OrgState_ORG_STATE_UNSPECIFIED:
		fallthrough
	default:
		return domain.OrgStateUnspecified
	}
}

func FieldNameToOrgColumn(fieldName v2beta_org.OrgFieldName) query.Column {
	switch fieldName {
	case v2beta_org.OrgFieldName_ORG_FIELD_NAME_NAME:
		return query.OrgColumnName
	case v2beta_org.OrgFieldName_ORG_FIELD_NAME_CREATION_DATE:
		return query.OrgColumnCreationDate
	case v2beta_org.OrgFieldName_ORG_FIELD_NAME_UNSPECIFIED:
		return query.Column{}
	default:
		return query.Column{}
	}
}

func ListOrgDomainsRequestToModel(systemDefaults systemdefaults.SystemDefaults, request *org.ListOrganizationDomainsRequest) (*query.OrgDomainSearchQueries, error) {
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
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		// SortingColumn: //TODO: sorting
		Queries: queries,
	}, nil
}

func ListQueryToModel(query *v2beta.ListQuery) (offset, limit uint64, asc bool) {
	if query == nil {
		return 0, 0, false
	}
	return query.Offset, uint64(query.Limit), query.Asc
}

func DomainQueriesToModel(queries []*v2beta_org.DomainSearchFilter) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = DomainQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func DomainQueryToModel(searchQuery *v2beta_org.DomainSearchFilter) (query.SearchQuery, error) {
	switch q := searchQuery.Filter.(type) {
	case *v2beta_org.DomainSearchFilter_DomainNameFilter:
		return query.NewOrgDomainDomainSearchQuery(v2beta_object.TextMethodToQuery(q.DomainNameFilter.Method), q.DomainNameFilter.Name)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Ags89", "List.Query.Invalid")
	}
}

func RemoveOrgDomainRequestToDomain(ctx context.Context, req *v2beta_org.DeleteOrganizationDomainRequest) *domain.OrgDomain {
	return &domain.OrgDomain{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.OrganizationId,
		},
		Domain: req.Domain,
	}
}

func GenerateOrgDomainValidationRequestToDomain(ctx context.Context, req *v2beta_org.GenerateOrganizationDomainValidationRequest) *domain.OrgDomain {
	return &domain.OrgDomain{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.OrganizationId,
		},
		Domain:         req.Domain,
		ValidationType: v2beta_object.DomainValidationTypeToDomain(req.Type),
	}
}

func ValidateOrgDomainRequestToDomain(ctx context.Context, req *v2beta_org.VerifyOrganizationDomainRequest) *domain.OrgDomain {
	return &domain.OrgDomain{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.OrganizationId,
		},
		Domain: req.Domain,
	}
}

func BulkSetOrgMetadataToDomain(req *v2beta_org.SetOrganizationMetadataRequest) []*domain.Metadata {
	metadata := make([]*domain.Metadata, len(req.Metadata))
	for i, data := range req.Metadata {
		metadata[i] = &domain.Metadata{
			Key:   data.Key,
			Value: data.Value,
		}
	}
	return metadata
}

func ListOrgMetadataToDomain(systemDefaults systemdefaults.SystemDefaults, request *v2beta_org.ListOrganizationMetadataRequest) (*query.OrgMetadataSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(systemDefaults, request.Pagination)
	if err != nil {
		return nil, err
	}
	queries, err := metadata.OrgMetadataQueriesToQuery(request.Filter)
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
