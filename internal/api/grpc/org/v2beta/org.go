package org

import (
	"context"
	"errors"

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

func (s *Server) CreateOrganization(ctx context.Context, request *v2beta_org.CreateOrganizationRequest) (*v2beta_org.CreateOrganizationResponse, error) {
	orgSetup, err := createOrganizationRequestToCommand(request)
	if err != nil {
		return nil, err
	}
	createdOrg, err := s.command.SetUpOrg(ctx, orgSetup, false)
	if err != nil {
		return nil, err
	}
	return createdOrganizationToPb(createdOrg)
}

func (s *Server) UpdateOrganization(ctx context.Context, request *v2beta_org.UpdateOrganizationRequest) (*v2beta_org.UpdateOrganizationResponse, error) {
	org, err := s.command.ChangeOrg(ctx, request.Id, request.Name)
	if err != nil {
		return nil, err
	}

	return &v2beta_org.UpdateOrganizationResponse{
		ChangeDate: timestamppb.New(org.EventDate),
	}, nil
}

func (s *Server) ListOrganizations(ctx context.Context, request *v2beta_org.ListOrganizationsRequest) (*v2beta_org.ListOrganizationsResponse, error) {
	queries, err := listOrgRequestToModel(s.systemDefaults, request)
	if err != nil {
		return nil, err
	}
	orgs, err := s.query.SearchOrgs(ctx, queries, s.checkPermission)
	if err != nil {
		return nil, err
	}
	return &v2beta_org.ListOrganizationsResponse{
		Organizations: OrgViewsToPb(orgs.Orgs),
		Pagination: &filter.PaginationResponse{
			TotalResult:  orgs.Count,
			AppliedLimit: uint64(request.GetPagination().GetLimit()),
		},
	}, nil
}

func (s *Server) DeleteOrganization(ctx context.Context, request *v2beta_org.DeleteOrganizationRequest) (*v2beta_org.DeleteOrganizationResponse, error) {
	details, err := s.command.RemoveOrg(ctx, request.Id)
	if err != nil {
		var notFoundError *zerrors.NotFoundError
		if errors.As(err, &notFoundError) {
			return &v2beta_org.DeleteOrganizationResponse{}, nil
		}
		return nil, err
	}
	return &v2beta_org.DeleteOrganizationResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}, nil
}

func (s *Server) SetOrganizationMetadata(ctx context.Context, request *v2beta_org.SetOrganizationMetadataRequest) (*v2beta_org.SetOrganizationMetadataResponse, error) {
	result, err := s.command.BulkSetOrgMetadata(ctx, request.OrganizationId, BulkSetOrgMetadataToDomain(request)...)
	if err != nil {
		return nil, err
	}
	return &org.SetOrganizationMetadataResponse{
		SetDate: timestamppb.New(result.EventDate),
	}, nil
}

func (s *Server) ListOrganizationMetadata(ctx context.Context, request *v2beta_org.ListOrganizationMetadataRequest) (*v2beta_org.ListOrganizationMetadataResponse, error) {
	metadataQueries, err := ListOrgMetadataToDomain(s.systemDefaults, request)
	if err != nil {
		return nil, err
	}
	res, err := s.query.SearchOrgMetadata(ctx, true, request.OrganizationId, metadataQueries, false)
	if err != nil {
		return nil, err
	}
	return &v2beta_org.ListOrganizationMetadataResponse{
		Metadata: metadata.OrgMetadataListToPb(res.Metadata),
		Pagination: &filter.PaginationResponse{
			TotalResult:  res.Count,
			AppliedLimit: uint64(request.GetPagination().GetLimit()),
		},
	}, nil
}

func (s *Server) DeleteOrganizationMetadata(ctx context.Context, request *v2beta_org.DeleteOrganizationMetadataRequest) (*v2beta_org.DeleteOrganizationMetadataResponse, error) {
	result, err := s.command.BulkRemoveOrgMetadata(ctx, request.OrganizationId, request.Keys...)
	if err != nil {
		return nil, err
	}
	return &v2beta_org.DeleteOrganizationMetadataResponse{
		DeletionDate: timestamppb.New(result.EventDate),
	}, nil
}

func (s *Server) DeactivateOrganization(ctx context.Context, request *org.DeactivateOrganizationRequest) (*org.DeactivateOrganizationResponse, error) {
	objectDetails, err := s.command.DeactivateOrg(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &org.DeactivateOrganizationResponse{
		ChangeDate: timestamppb.New(objectDetails.EventDate),
	}, nil
}

func (s *Server) ActivateOrganization(ctx context.Context, request *org.ActivateOrganizationRequest) (*org.ActivateOrganizationResponse, error) {
	objectDetails, err := s.command.ReactivateOrg(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &org.ActivateOrganizationResponse{
		ChangeDate: timestamppb.New(objectDetails.EventDate),
	}, err
}

func (s *Server) AddOrganizationDomain(ctx context.Context, request *org.AddOrganizationDomainRequest) (*org.AddOrganizationDomainResponse, error) {
	userIDs, err := s.getClaimedUserIDsOfOrgDomain(ctx, request.Domain, request.OrganizationId)
	if err != nil {
		return nil, err
	}
	details, err := s.command.AddOrgDomain(ctx, request.OrganizationId, request.Domain, userIDs)
	if err != nil {
		return nil, err
	}
	return &org.AddOrganizationDomainResponse{
		CreationDate: timestamppb.New(details.EventDate),
	}, nil
}

func (s *Server) ListOrganizationDomains(ctx context.Context, req *org.ListOrganizationDomainsRequest) (*org.ListOrganizationDomainsResponse, error) {
	queries, err := ListOrgDomainsRequestToModel(s.systemDefaults, req)
	if err != nil {
		return nil, err
	}
	orgIDQuery, err := query.NewOrgDomainOrgIDSearchQuery(req.OrganizationId)
	if err != nil {
		return nil, err
	}
	queries.Queries = append(queries.Queries, orgIDQuery)

	domains, err := s.query.SearchOrgDomains(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return &org.ListOrganizationDomainsResponse{
		Domains: object.DomainsToPb(domains.Domains),
		Pagination: &filter.PaginationResponse{
			TotalResult:  domains.Count,
			AppliedLimit: uint64(req.GetPagination().GetLimit()),
		},
	}, nil
}

func (s *Server) DeleteOrganizationDomain(ctx context.Context, req *org.DeleteOrganizationDomainRequest) (*org.DeleteOrganizationDomainResponse, error) {
	details, err := s.command.RemoveOrgDomain(ctx, RemoveOrgDomainRequestToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &org.DeleteOrganizationDomainResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}, err
}

func (s *Server) GenerateOrganizationDomainValidation(ctx context.Context, req *org.GenerateOrganizationDomainValidationRequest) (*org.GenerateOrganizationDomainValidationResponse, error) {
	token, url, err := s.command.GenerateOrgDomainValidation(ctx, GenerateOrgDomainValidationRequestToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &org.GenerateOrganizationDomainValidationResponse{
		Token: token,
		Url:   url,
	}, nil
}

func (s *Server) VerifyOrganizationDomain(ctx context.Context, request *org.VerifyOrganizationDomainRequest) (*org.VerifyOrganizationDomainResponse, error) {
	userIDs, err := s.getClaimedUserIDsOfOrgDomain(ctx, request.Domain, request.OrganizationId)
	if err != nil {
		return nil, err
	}
	details, err := s.command.ValidateOrgDomain(ctx, ValidateOrgDomainRequestToDomain(ctx, request), userIDs)
	if err != nil {
		return nil, err
	}
	return &org.VerifyOrganizationDomainResponse{
		ChangeDate: timestamppb.New(details.EventDate),
	}, nil
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
	queries := make([]query.SearchQuery, 0, 2)
	loginName, err := query.NewUserPreferredLoginNameSearchQuery("@"+orgDomain, query.TextEndsWithIgnoreCase)
	if err != nil {
		return nil, err
	}
	queries = append(queries, loginName)
	if orgID != "" {
		owner, err := query.NewUserResourceOwnerSearchQuery(orgID, query.TextNotEquals)
		if err != nil {
			return nil, err
		}
		queries = append(queries, owner)
	}
	users, err := s.query.SearchUsers(ctx, &query.UserSearchQueries{Queries: queries}, nil)
	if err != nil {
		return nil, err
	}
	userIDs := make([]string, len(users.Users))
	for i, user := range users.Users {
		userIDs[i] = user.ID
	}
	return userIDs, nil
}
