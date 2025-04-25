package org

import (
	"context"

	object "github.com/zitadel/zitadel/internal/api/grpc/object/v2beta"
	user "github.com/zitadel/zitadel/internal/api/grpc/user/v2beta"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/zerrors"
	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

func (s *Server) CreateOrganization(ctx context.Context, request *org.CreateOrganizationRequest) (*org.CreateOrganizationResponse, error) {
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

func (s *Server) UpdateOrganization(ctx context.Context, request *org.UpdateOrganizationRequest) (*org.UpdateOrganizationResponse, error) {
	updated_org, err := s.command.UpdateOrg(ctx, request.Id, request.Name)
	if err != nil {
		return nil, err
	}

	return &org.UpdateOrganizationResponse{
		Details: object.DomainToDetailsPb(updated_org),
	}, nil
}

func createOrganizationRequestToCommand(request *org.CreateOrganizationRequest) (*command.OrgSetup, error) {
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

func createOrganizationRequestAdminsToCommand(requestAdmins []*org.CreateOrganizationRequest_Admin) (admins []*command.OrgSetupAdmin, err error) {
	admins = make([]*command.OrgSetupAdmin, len(requestAdmins))
	for i, admin := range requestAdmins {
		admins[i], err = createOrganizationRequestAdminToCommand(admin)
		if err != nil {
			return nil, err
		}
	}
	return admins, nil
}

func createOrganizationRequestAdminToCommand(admin *org.CreateOrganizationRequest_Admin) (*command.OrgSetupAdmin, error) {
	switch a := admin.GetUserType().(type) {
	case *org.CreateOrganizationRequest_Admin_UserId:
		return &command.OrgSetupAdmin{
			ID:    a.UserId,
			Roles: admin.GetRoles(),
		}, nil
	case *org.CreateOrganizationRequest_Admin_Human:
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
		Details:        object.DomainToDetailsPb(createdOrg.ObjectDetails),
		OrganizationId: createdOrg.ObjectDetails.ResourceOwner,
		CreatedAdmins:  admins,
	}, nil
}
