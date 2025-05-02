package org

import (
	"time"

	v2beta_object "github.com/zitadel/zitadel/internal/api/grpc/object/v2beta"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	v2beta_org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"

	v2beta "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"

	org_pb "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NOTE: most of this code is copied from `internal/api/grpc/admin/*`, as we will eventually axe the previous versons of the API,
// we will have code duplication until then

func listOrgRequestToModel(request *v2beta_org.ListOrganizationsRequest) (*query.OrgSearchQueries, error) {
	offset, limit, asc := v2beta_object.ListQueryToModel(request.Query)
	// queries, err := org_pb.OrgQueriesToModel(request.Queries)
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

func OrganizationViewToPb(org *query.Org) *org_pb.Organization {
	return &org_pb.Organization{
		Id:            org.ID,
		State:         OrgStateToPb(org.State),
		Name:          org.Name,
		PrimaryDomain: org.Domain,
		Details: ToViewDetailsPb(
			org.Sequence,
			org.CreationDate,
			org.ChangeDate,
			org.ResourceOwner,
		),
	}
}

func OrgStateToPb(state domain.OrgState) org_pb.OrgState {
	switch state {
	case domain.OrgStateActive:
		return org_pb.OrgState_ORG_STATE_ACTIVE
	case domain.OrgStateInactive:
		return org_pb.OrgState_ORG_STATE_INACTIVE
	default:
		return org_pb.OrgState_ORG_STATE_UNSPECIFIED
	}
}

func ToViewDetailsPb(
	sequence uint64,
	creationDate,
	changeDate time.Time,
	resourceOwner string,
) *v2beta.Details {
	details := &v2beta.Details{
		Sequence:      sequence,
		ResourceOwner: resourceOwner,
	}
	if !creationDate.IsZero() {
		details.CreationDate = timestamppb.New(creationDate)
	}
	if !changeDate.IsZero() {
		details.ChangeDate = timestamppb.New(changeDate)
	}
	return details
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

func OrgViewsToPb(orgs []*query.Org) []*org_pb.Organization {
	o := make([]*org_pb.Organization, len(orgs))
	for i, org := range orgs {
		o[i] = OrgViewToPb(org)
	}
	return o
}

func OrgQueriesToModel(queries []*org_pb.OrgQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = OrgQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func OrgQueryToModel(apiQuery *org_pb.OrgQuery) (query.SearchQuery, error) {
	switch q := apiQuery.Query.(type) {
	case *org_pb.OrgQuery_DomainQuery:
		return query.NewOrgVerifiedDomainSearchQuery(v2beta_object.TextMethodToQuery(q.DomainQuery.Method), q.DomainQuery.Domain)
	case *org_pb.OrgQuery_NameQuery:
		return query.NewOrgNameSearchQuery(v2beta_object.TextMethodToQuery(q.NameQuery.Method), q.NameQuery.Name)
	case *org_pb.OrgQuery_StateQuery:
		return query.NewOrgStateSearchQuery(OrgStateToDomain(q.StateQuery.State))
	case *org_pb.OrgQuery_IdQuery:
		return query.NewOrgIDSearchQuery(q.IdQuery.Id)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-vR9nC", "List.Query.Invalid")
	}
}

func OrgStateToDomain(state org_pb.OrgState) domain.OrgState {
	switch state {
	case org_pb.OrgState_ORG_STATE_ACTIVE:
		return domain.OrgStateActive
	case org_pb.OrgState_ORG_STATE_INACTIVE:
		return domain.OrgStateInactive
	case org_pb.OrgState_ORG_STATE_UNSPECIFIED:
		fallthrough
	default:
		return domain.OrgStateUnspecified
	}
}

func FieldNameToOrgColumn(fieldName org_pb.OrgFieldName) query.Column {
	switch fieldName {
	case org_pb.OrgFieldName_ORG_FIELD_NAME_NAME:
		return query.OrgColumnName
	case org_pb.OrgFieldName_ORG_FIELD_NAME_UNSPECIFIED:
		return query.Column{}
	default:
		return query.Column{}
	}
}

func OrgViewToPb(org *query.Org) *org_pb.Organization {
	return &org_pb.Organization{
		Id:            org.ID,
		State:         OrgStateToPb(org.State),
		Name:          org.Name,
		PrimaryDomain: org.Domain,
		Details: ToViewDetailsPb(
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

func DomainQueriesToModel(queries []*org_pb.DomainSearchQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = DomainQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func DomainQueryToModel(searchQuery *org_pb.DomainSearchQuery) (query.SearchQuery, error) {
	switch q := searchQuery.Query.(type) {
	case *org_pb.DomainSearchQuery_DomainNameQuery:
		// return query.NewOrgDomainDomainSearchQuery(object.TextMethodToQuery(q.DomainNameQuery.Method), q.DomainNameQuery.Name)
		return query.NewOrgDomainDomainSearchQuery(v2beta_object.TextMethodToQuery(q.DomainNameQuery.Method), q.DomainNameQuery.Name)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Ags89", "List.Query.Invalid")
	}
}
