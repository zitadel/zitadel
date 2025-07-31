package org

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	metadata "github.com/zitadel/zitadel/internal/api/grpc/metadata/v2beta"
	object "github.com/zitadel/zitadel/internal/api/grpc/object/v2beta"
	user "github.com/zitadel/zitadel/internal/api/grpc/user/v2beta"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	v2beta_org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

func (s *Server) CreateOrganization(ctx context.Context, request *connect.Request[v2beta_org.CreateOrganizationRequest]) (*connect.Response[v2beta_org.CreateOrganizationResponse], error) {
	orgSetup, err := createOrganizationRequestToCommand(request.Msg)
	if err != nil {
		return nil, err
	}
	createdOrg, err := s.command.SetUpOrg(ctx, orgSetup, false)
	if err != nil {
		return nil, err
	}
	return createdOrganizationToPb(createdOrg)
}

func (s *Server) UpdateOrganization(ctx context.Context, request *connect.Request[v2beta_org.UpdateOrganizationRequest]) (*connect.Response[v2beta_org.UpdateOrganizationResponse], error) {
	org, err := s.command.ChangeOrg(ctx, request.Msg.GetId(), request.Msg.GetName())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v2beta_org.UpdateOrganizationResponse{
		ChangeDate: timestamppb.New(org.EventDate),
	}), nil
}

func (s *Server) ListOrganizations(ctx context.Context, request *connect.Request[v2beta_org.ListOrganizationsRequest]) (*connect.Response[v2beta_org.ListOrganizationsResponse], error) {
	queries, err := listOrgRequestToModel(s.systemDefaults, request.Msg)
	if err != nil {
		return nil, err
	}
	orgs, err := s.query.SearchOrgs(ctx, queries, s.checkPermission)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&v2beta_org.ListOrganizationsResponse{
		Organizations: OrgViewsToPb(orgs.Orgs),
		Pagination: &filter.PaginationResponse{
			TotalResult:  orgs.Count,
			AppliedLimit: uint64(request.Msg.GetPagination().GetLimit()),
		},
	}), nil
}

func (s *Server) DeleteOrganization(ctx context.Context, request *connect.Request[v2beta_org.DeleteOrganizationRequest]) (*connect.Response[v2beta_org.DeleteOrganizationResponse], error) {
	details, err := s.command.RemoveOrg(ctx, request.Msg.GetId())
	if err != nil {
		var notFoundError *zerrors.NotFoundError
		if errors.As(err, &notFoundError) {
			return connect.NewResponse(&v2beta_org.DeleteOrganizationResponse{}), nil
		}
		return nil, err
	}
	return connect.NewResponse(&v2beta_org.DeleteOrganizationResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}), nil
}

func (s *Server) SetOrganizationMetadata(ctx context.Context, request *connect.Request[v2beta_org.SetOrganizationMetadataRequest]) (*connect.Response[v2beta_org.SetOrganizationMetadataResponse], error) {
	result, err := s.command.BulkSetOrgMetadata(ctx, request.Msg.GetOrganizationId(), BulkSetOrgMetadataToDomain(request.Msg)...)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.SetOrganizationMetadataResponse{
		SetDate: timestamppb.New(result.EventDate),
	}), nil
}

func (s *Server) ListOrganizationMetadata(ctx context.Context, request *connect.Request[v2beta_org.ListOrganizationMetadataRequest]) (*connect.Response[v2beta_org.ListOrganizationMetadataResponse], error) {
	metadataQueries, err := ListOrgMetadataToDomain(s.systemDefaults, request.Msg)
	if err != nil {
		return nil, err
	}
	res, err := s.query.SearchOrgMetadata(ctx, true, request.Msg.GetOrganizationId(), metadataQueries, false)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&v2beta_org.ListOrganizationMetadataResponse{
		Metadata: metadata.OrgMetadataListToPb(res.Metadata),
		Pagination: &filter.PaginationResponse{
			TotalResult:  res.Count,
			AppliedLimit: uint64(request.Msg.GetPagination().GetLimit()),
		},
	}), nil
}

func (s *Server) DeleteOrganizationMetadata(ctx context.Context, request *connect.Request[v2beta_org.DeleteOrganizationMetadataRequest]) (*connect.Response[v2beta_org.DeleteOrganizationMetadataResponse], error) {
	result, err := s.command.BulkRemoveOrgMetadata(ctx, request.Msg.GetOrganizationId(), request.Msg.Keys...)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&v2beta_org.DeleteOrganizationMetadataResponse{
		DeletionDate: timestamppb.New(result.EventDate),
	}), nil
}

func (s *Server) DeactivateOrganization(ctx context.Context, request *connect.Request[org.DeactivateOrganizationRequest]) (*connect.Response[org.DeactivateOrganizationResponse], error) {
	objectDetails, err := s.command.DeactivateOrg(ctx, request.Msg.GetId())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.DeactivateOrganizationResponse{
		ChangeDate: timestamppb.New(objectDetails.EventDate),
	}), nil
}

func (s *Server) ActivateOrganization(ctx context.Context, request *connect.Request[org.ActivateOrganizationRequest]) (*connect.Response[org.ActivateOrganizationResponse], error) {
	objectDetails, err := s.command.ReactivateOrg(ctx, request.Msg.GetId())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.ActivateOrganizationResponse{
		ChangeDate: timestamppb.New(objectDetails.EventDate),
	}), err
}

func (s *Server) AddOrganizationDomain(ctx context.Context, request *connect.Request[org.AddOrganizationDomainRequest]) (*connect.Response[org.AddOrganizationDomainResponse], error) {
	userIDs, err := s.getClaimedUserIDsOfOrgDomain(ctx, request.Msg.GetDomain(), request.Msg.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	details, err := s.command.AddOrgDomain(ctx, request.Msg.GetOrganizationId(), request.Msg.GetDomain(), userIDs)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.AddOrganizationDomainResponse{
		CreationDate: timestamppb.New(details.EventDate),
	}), nil
}

func (s *Server) ListOrganizationDomains(ctx context.Context, req *connect.Request[org.ListOrganizationDomainsRequest]) (*connect.Response[org.ListOrganizationDomainsResponse], error) {
	queries, err := ListOrgDomainsRequestToModel(s.systemDefaults, req.Msg)
	if err != nil {
		return nil, err
	}
	orgIDQuery, err := query.NewOrgDomainOrgIDSearchQuery(req.Msg.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	queries.Queries = append(queries.Queries, orgIDQuery)

	domains, err := s.query.SearchOrgDomains(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.ListOrganizationDomainsResponse{
		Domains: object.DomainsToPb(domains.Domains),
		Pagination: &filter.PaginationResponse{
			TotalResult:  domains.Count,
			AppliedLimit: uint64(req.Msg.GetPagination().GetLimit()),
		},
	}), nil
}

func (s *Server) DeleteOrganizationDomain(ctx context.Context, req *connect.Request[org.DeleteOrganizationDomainRequest]) (*connect.Response[org.DeleteOrganizationDomainResponse], error) {
	details, err := s.command.RemoveOrgDomain(ctx, RemoveOrgDomainRequestToDomain(ctx, req.Msg))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.DeleteOrganizationDomainResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}), err
}

func (s *Server) GenerateOrganizationDomainValidation(ctx context.Context, req *connect.Request[org.GenerateOrganizationDomainValidationRequest]) (*connect.Response[org.GenerateOrganizationDomainValidationResponse], error) {
	token, url, err := s.command.GenerateOrgDomainValidation(ctx, GenerateOrgDomainValidationRequestToDomain(ctx, req.Msg))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.GenerateOrganizationDomainValidationResponse{
		Token: token,
		Url:   url,
	}), nil
}

func (s *Server) VerifyOrganizationDomain(ctx context.Context, request *connect.Request[org.VerifyOrganizationDomainRequest]) (*connect.Response[org.VerifyOrganizationDomainResponse], error) {
	userIDs, err := s.getClaimedUserIDsOfOrgDomain(ctx, request.Msg.GetDomain(), request.Msg.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	details, err := s.command.ValidateOrgDomain(ctx, ValidateOrgDomainRequestToDomain(ctx, request.Msg), userIDs)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.VerifyOrganizationDomainResponse{
		ChangeDate: timestamppb.New(details.EventDate),
	}), nil
}

func createOrganizationRequestToCommand(request *v2beta_org.CreateOrganizationRequest) (*command.OrgSetup, error) {
	admins, err := createOrganizationRequestAdminsToCommand(request.GetAdmins())
	if err != nil {
		return nil, err
	}
	return &command.OrgSetup{
		Name:         request.GetName(),
		CustomDomain: "",
		Admins:       admins,
		OrgID:        request.GetId(),
	}, nil
}

func createOrganizationRequestAdminsToCommand(requestAdmins []*v2beta_org.CreateOrganizationRequest_Admin) (admins []*command.OrgSetupAdmin, err error) {
	admins = make([]*command.OrgSetupAdmin, len(requestAdmins))
	for i, admin := range requestAdmins {
		admins[i], err = createOrganizationRequestAdminToCommand(admin)
		if err != nil {
			return nil, err
		}
	}
	return admins, nil
}

func createOrganizationRequestAdminToCommand(admin *v2beta_org.CreateOrganizationRequest_Admin) (*command.OrgSetupAdmin, error) {
	switch a := admin.GetUserType().(type) {
	case *v2beta_org.CreateOrganizationRequest_Admin_UserId:
		return &command.OrgSetupAdmin{
			ID:    a.UserId,
			Roles: admin.GetRoles(),
		}, nil
	case *v2beta_org.CreateOrganizationRequest_Admin_Human:
		human, err := user.AddUserRequestToAddHuman(a.Human)
		if err != nil {
			return nil, err
		}
		return &command.OrgSetupAdmin{
			Human: human,
			Roles: admin.GetRoles(),
		}, nil
	default:
		return nil, zerrors.ThrowUnimplementedf(nil, "ORGv2-SL2r8", "userType oneOf %T in method AddOrganization not implemented", a)
	}
}

func (s *Server) getClaimedUserIDsOfOrgDomain(ctx context.Context, orgDomain, orgID string) ([]string, error) {
	return s.query.SearchClaimedUserIDsOfOrgDomain(ctx, orgDomain, orgID)
}
