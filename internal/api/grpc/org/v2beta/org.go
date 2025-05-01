package org

import (
	"context"

	object "github.com/zitadel/zitadel/internal/api/grpc/object/v2beta"
	user "github.com/zitadel/zitadel/internal/api/grpc/user/v2beta"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	org, err := s.command.UpdateOrg(ctx, request.Id, request.Name)
	if err != nil {
		return nil, err
	}

	return &v2beta_org.UpdateOrganizationResponse{
		Details: object.DomainToDetailsPb(org),
	}, nil
}

func (s *Server) GetOrganizationByID(ctx context.Context, request *v2beta_org.GetOrganizationByIDRequest) (*v2beta_org.GetOrganizationByIDResponse, error) {
	org, err := s.query.OrgByID(ctx, true, request.Id)
	if err != nil {
		return nil, err
	}
	return &v2beta_org.GetOrganizationByIDResponse{
		Organization: OrganizationViewToPb(org),
	}, nil
}

func (s *Server) ListOrganizations(ctx context.Context, request *v2beta_org.ListOrganizationsRequest) (*v2beta_org.ListOrganizationsResponse, error) {
	queries, err := listOrgRequestToModel(request)
	if err != nil {
		return nil, err
	}
	orgs, err := s.query.SearchOrgs(ctx, queries, nil)
	if err != nil {
		return nil, err
	}
	return &v2beta_org.ListOrganizationsResponse{
		Result:  OrgViewsToPb(orgs.Orgs),
		Details: object.ToListDetails(orgs.SearchResponse),
	}, nil
}

func (s *Server) DeleteOrganization(ctx context.Context, request *v2beta_org.DeleteOrganizationRequest) (*v2beta_org.DeleteOrganizationResponse, error) {
	details, err := s.command.RemoveOrg(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &v2beta_org.DeleteOrganizationResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) DeactivateOrganization(ctx context.Context, request *org.DeactivateOrganizationRequest) (*org.DeactivateOrganizationResponse, error) {
	objectDetails, err := s.command.DeactivateOrg(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &org.DeactivateOrganizationResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ReactivateOrganization(ctx context.Context, request *org.ReactivateOrganizationRequest) (*org.ReactivateOrganizationResponse, error) {
	objectDetails, err := s.command.ReactivateOrg(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &org.ReactivateOrganizationResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, err
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
		return nil, zerrors.ThrowUnimplementedf(nil, "ORGv2-SD2r1", "userType oneOf %T in method AddOrganization not implemented", a)
	}
}
