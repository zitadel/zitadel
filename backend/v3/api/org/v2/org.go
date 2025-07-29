package org

import (
	"context"
	"slices"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/api"
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

// ActivateOrganization implements [orgconnect.OrganizationServiceHandler].
func (s *Server) ActivateOrganization(ctx context.Context, req *connect.Request[org.ActivateOrganizationRequest]) (*connect.Response[org.ActivateOrganizationResponse], error) {
	err := domain.Invoke(ctx, domain.NewActivateOrganizationCommand(
		authz.GetInstance(ctx).InstanceID(),
		req.Msg.GetId(),
	))
	if err != nil {
		return nil, err
	}
	// DISCUSS(adlerhurst): does returning the ChangeDate bring any value?
	return connect.NewResponse(&org.ActivateOrganizationResponse{
		ChangeDate: timestamppb.Now(),
	}), nil
}

// CreateOrganization implements [orgconnect.OrganizationServiceHandler].
func (s *Server) CreateOrganization(ctx context.Context, req *connect.Request[org.CreateOrganizationRequest]) (*connect.Response[org.CreateOrganizationResponse], error) {
	// TODO: Implement admins
	opts := make([]domain.CreateOrganizationCommandOpts, 0, 1+len(req.Msg.Admins))
	if req.Msg.Id != nil {
		opts = append(opts, domain.WithOrganizationID(req.Msg.GetId()))
	}
	cmd := domain.NewCreateOrganizationCommand(
		authz.GetInstance(ctx).InstanceID(),
		req.Msg.GetName(),
		opts...,
	)
	err := domain.Invoke(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.CreateOrganizationResponse{
		Id:           cmd.ID,
		CreationDate: timestamppb.New(cmd.CreatedAt),
	}), nil
}

// DeactivateOrganization implements [orgconnect.OrganizationServiceHandler].
func (s *Server) DeactivateOrganization(ctx context.Context, req *connect.Request[org.DeactivateOrganizationRequest]) (*connect.Response[org.DeactivateOrganizationResponse], error) {
	err := domain.Invoke(ctx, domain.NewDeactivateOrganizationCommand(
		authz.GetInstance(ctx).InstanceID(),
		req.Msg.GetId(),
	))
	if err != nil {
		return nil, err
	}
	// DISCUSS(adlerhurst): does returning the ChangeDate bring any value?
	return connect.NewResponse(&org.DeactivateOrganizationResponse{
		ChangeDate: timestamppb.Now(),
	}), nil
}

// DeleteOrganization implements [orgconnect.OrganizationServiceHandler].
func (s *Server) DeleteOrganization(ctx context.Context, req *connect.Request[org.DeleteOrganizationRequest]) (*connect.Response[org.DeleteOrganizationResponse], error) {
	err := domain.Invoke(ctx, domain.NewDeleteOrganizationCommand(
		authz.GetInstance(ctx).InstanceID(),
		req.Msg.GetId(),
	))
	if err != nil {
		return nil, err
	}
	// DISCUSS(adlerhurst): does returning the DeletionDate bring any value?
	return connect.NewResponse(&org.DeleteOrganizationResponse{
		DeletionDate: timestamppb.Now(),
	}), nil
}

// ListOrganizations implements [orgconnect.OrganizationServiceHandler].
func (s *Server) ListOrganizations(ctx context.Context, req *connect.Request[org.ListOrganizationsRequest]) (*connect.Response[org.ListOrganizationsResponse], error) {
	opts := orgFiltersToDomain(req.Msg.GetFilter())
	opts = slices.Grow(opts, 2)
	opts = append(opts, api.V2BetaPaginationToDomain(req.Msg.GetPagination()))
	if req.Msg.SortingColumn != org.OrgFieldName_ORG_FIELD_NAME_UNSPECIFIED {
		opts = append(opts, domain.WithOrgQuerySortingColumn(orgFieldNameToDatabase(req.Msg.SortingColumn)))
	}

	query := domain.NewOrgsQuery(
		authz.GetInstance(ctx).InstanceID(),
		opts...,
	)
	err := domain.Invoke(ctx, query)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.ListOrganizationsResponse{
		Organizations: orgsFromDomain(query.Result),
		Pagination: &filter.PaginationResponse{
			AppliedLimit: uint64(req.Msg.Pagination.Limit),
			// TotalResult: TODO(adlerhurst): needs implementation in lower layers
		},
	}), nil
}

// UpdateOrganization implements [orgconnect.OrganizationServiceHandler].
func (s *Server) UpdateOrganization(ctx context.Context, req *connect.Request[org.UpdateOrganizationRequest]) (*connect.Response[org.UpdateOrganizationResponse], error) {
	opts := make([]domain.UpdateOrganizationCommandOpts, 0, 1)
	if req.Msg.Name != "" {
		opts = append(opts, domain.WithOrganizationName(req.Msg.GetName()))
	}

	err := domain.Invoke(ctx, domain.NewUpdateOrganizationCommand(
		authz.GetInstance(ctx).InstanceID(),
		req.Msg.GetId(),
		opts...,
	))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.UpdateOrganizationResponse{
		ChangeDate: timestamppb.Now(),
	}), nil
}

func orgFiltersToDomain(filters []*org.OrganizationSearchFilter) []domain.OrgsQueryOpts {
	opts := make([]domain.OrgsQueryOpts, len(filters))
	for i, filter := range filters {
		opts[i] = orgFilterToDomain(filter)
	}
	return opts
}

func orgFilterToDomain(filter *org.OrganizationSearchFilter) domain.OrgsQueryOpts {
	switch f := filter.Filter.(type) {
	case *org.OrganizationSearchFilter_NameFilter:
		return domain.WithOrgByNameQuery(api.V2BetaTextFilterToDatabase(f.NameFilter.Method), f.NameFilter.Name)
	case *org.OrganizationSearchFilter_DomainFilter:
		return domain.WithOrgByDomainQuery(api.V2BetaTextFilterToDatabase(f.DomainFilter.Method), f.DomainFilter.Domain)
	case *org.OrganizationSearchFilter_IdFilter:
		return domain.WithOrgByIDQuery(f.IdFilter.Id)
	case *org.OrganizationSearchFilter_StateFilter:
		return domain.WithOrgByStateQuery(api.V2BetaOrgStateToDomain(f.StateFilter.State))
	default:
		panic("unknown organization search filter: " + filter.String())
	}
}

func orgsFromDomain(orgs []*domain.Organization) []*org.Organization {
	result := make([]*org.Organization, len(orgs))
	for i, o := range orgs {
		result[i] = &org.Organization{
			Id:            o.ID,
			Name:          o.Name,
			State:         orgStateFromDomain(o.State),
			CreationDate:  timestamppb.New(o.CreatedAt),
			ChangedDate:   timestamppb.New(o.UpdatedAt),
			PrimaryDomain: orgPrimaryDomainFromDomain(o.Domains),
		}
	}
	return result
}

func orgStateFromDomain(state domain.OrgState) org.OrgState {
	switch state {
	case domain.OrgStateActive:
		return org.OrgState_ORG_STATE_ACTIVE
	case domain.OrgStateInactive:
		return org.OrgState_ORG_STATE_INACTIVE
	default:
		return org.OrgState_ORG_STATE_UNSPECIFIED
	}
}

func orgPrimaryDomainFromDomain(domains []*domain.OrganizationDomain) string {
	for _, d := range domains {
		if d.IsPrimary {
			return d.Domain
		}
	}
	return ""
}

func orgFieldNameToDatabase(fieldName org.OrgFieldName) func(query *domain.OrgsQuery) database.Column {
	switch fieldName {
	case org.OrgFieldName_ORG_FIELD_NAME_NAME:
		return domain.OrderOrgsByName
	case org.OrgFieldName_ORG_FIELD_NAME_CREATION_DATE:
		return domain.OrderOrgsByCreationDate
	default:
		return nil
	}
}
