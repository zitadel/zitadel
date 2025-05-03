package org

import (
	"context"

	metadata "github.com/zitadel/zitadel/internal/api/grpc/metadata/v2beta"
	v2beta_object "github.com/zitadel/zitadel/internal/api/grpc/object/v2beta"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"

	// TODO fix below
	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	v2beta_org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"

	v2beta "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
)

// NOTE: most of this code is copied from `internal/api/grpc/admin/*`, as we will eventually axe the previous versons of the API,
// we will have code duplication until then

func listOrgRequestToModel(request *v2beta_org.ListOrganizationsRequest) (*query.OrgSearchQueries, error) {
	offset, limit, asc := v2beta_object.ListQueryToModel(request.Query)
	queries, err := OrgQueriesToModel(request.Queries)
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
		Details: v2beta_object.ToViewDetailsPb(
			org.Sequence,
			org.CreationDate,
			org.ChangeDate,
			org.ResourceOwner,
		),
	}
}

func OrgStateToPb(state domain.OrgState) v2beta_org.OrgState {
	switch state {
	case domain.OrgStateActive:
		return v2beta_org.OrgState_ORG_STATE_ACTIVE
	case domain.OrgStateInactive:
		return v2beta_org.OrgState_ORG_STATE_INACTIVE
	default:
		return v2beta_org.OrgState_ORG_STATE_UNSPECIFIED
	}
}

func createdOrganizationToPb(createdOrg *command.CreatedOrg) (_ *org.CreateOrganizationResponse, err error) {
	admins := make([]*org.CreateOrganizationResponse_CreatedAdmin, len(createdOrg.CreatedAdmins))
	for i, admin := range createdOrg.CreatedAdmins {
		admins[i] = &org.CreateOrganizationResponse_CreatedAdmin{
			UserId:    admin.ID,
			EmailCode: admin.EmailCode,
			PhoneCode: admin.PhoneCode,
		}
	}
	return &org.CreateOrganizationResponse{
		Details:        v2beta_object.DomainToDetailsPb(createdOrg.ObjectDetails),
		OrganizationId: createdOrg.ObjectDetails.ResourceOwner,
		CreatedAdmins:  admins,
	}, nil
}

func OrgViewsToPb(orgs []*query.Org) []*v2beta_org.Organization {
	o := make([]*v2beta_org.Organization, len(orgs))
	for i, org := range orgs {
		o[i] = OrgViewToPb(org)
	}
	return o
}

func OrgQueriesToModel(queries []*v2beta_org.OrgQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = OrgQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func OrgQueryToModel(apiQuery *v2beta_org.OrgQuery) (query.SearchQuery, error) {
	switch q := apiQuery.Query.(type) {
	case *v2beta_org.OrgQuery_DomainQuery:
		return query.NewOrgVerifiedDomainSearchQuery(v2beta_object.TextMethodToQuery(q.DomainQuery.Method), q.DomainQuery.Domain)
	case *v2beta_org.OrgQuery_NameQuery:
		return query.NewOrgNameSearchQuery(v2beta_object.TextMethodToQuery(q.NameQuery.Method), q.NameQuery.Name)
	case *v2beta_org.OrgQuery_StateQuery:
		return query.NewOrgStateSearchQuery(OrgStateToDomain(q.StateQuery.State))
	case *v2beta_org.OrgQuery_IdQuery:
		return query.NewOrgIDSearchQuery(q.IdQuery.Id)
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
	case v2beta_org.OrgFieldName_ORG_FIELD_NAME_UNSPECIFIED:
		return query.Column{}
	default:
		return query.Column{}
	}
}

func OrgViewToPb(org *query.Org) *v2beta_org.Organization {
	return &v2beta_org.Organization{
		Id:            org.ID,
		State:         OrgStateToPb(org.State),
		Name:          org.Name,
		PrimaryDomain: org.Domain,
		Details: v2beta_object.ToViewDetailsPb(
			org.Sequence,
			org.CreationDate,
			org.ChangeDate,
			org.ResourceOwner,
		),
	}
}

func ListOrgDomainsRequestToModel(req *org.ListOrganizationDomainsRequest) (*query.OrgDomainSearchQueries, error) {
	offset, limit, asc := ListQueryToModel(req.Query)
	// queries, err := org_grpc.DomainQueriesToModel(req.Queries)
	queries, err := DomainQueriesToModel(req.Queries)
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

func DomainQueriesToModel(queries []*v2beta_org.DomainSearchQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = DomainQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func DomainQueryToModel(searchQuery *v2beta_org.DomainSearchQuery) (query.SearchQuery, error) {
	switch q := searchQuery.Query.(type) {
	case *v2beta_org.DomainSearchQuery_DomainNameQuery:
		// return query.NewOrgDomainDomainSearchQuery(object.TextMethodToQuery(q.DomainNameQuery.Method), q.DomainNameQuery.Name)
		return query.NewOrgDomainDomainSearchQuery(v2beta_object.TextMethodToQuery(q.DomainNameQuery.Method), q.DomainNameQuery.Name)
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

func ListOrgMetadataToDomain(request *v2beta_org.ListOrganizationMetadataRequest) (*query.OrgMetadataSearchQueries, error) {
	offset, limit, asc := v2beta_object.ListQueryToModel(request.Query)
	queries, err := metadata.OrgMetadataQueriesToQuery(request.Queries)
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
