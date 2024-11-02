package org

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/api/grpc/user/v2"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
)

func (s *Server) AddOrganization(ctx context.Context, request *org.AddOrganizationRequest) (*org.AddOrganizationResponse, error) {
	orgSetup, err := addOrganizationRequestToCommand(request)
	if err != nil {
		return nil, err
	}
	createdOrg, err := s.command.SetUpOrg(ctx, orgSetup, false)
	if err != nil {
		return nil, err
	}
	return createdOrganizationToPb(createdOrg)
}

func addOrganizationRequestToCommand(request *org.AddOrganizationRequest) (*command.OrgSetup, error) {
	admins, err := addOrganizationRequestAdminsToCommand(request.GetAdmins())
	if err != nil {
		return nil, err
	}
	return &command.OrgSetup{
		Name:         request.GetName(),
		CustomDomain: "",
		Admins:       admins,
	}, nil
}

func addOrganizationRequestAdminsToCommand(requestAdmins []*org.AddOrganizationRequest_Admin) (admins []*command.OrgSetupAdmin, err error) {
	admins = make([]*command.OrgSetupAdmin, len(requestAdmins))
	for i, admin := range requestAdmins {
		admins[i], err = addOrganizationRequestAdminToCommand(admin)
		if err != nil {
			return nil, err
		}
	}
	return admins, nil
}

func addOrganizationRequestAdminToCommand(admin *org.AddOrganizationRequest_Admin) (*command.OrgSetupAdmin, error) {
	switch a := admin.GetUserType().(type) {
	case *org.AddOrganizationRequest_Admin_UserId:
		return &command.OrgSetupAdmin{
			ID:    a.UserId,
			Roles: admin.GetRoles(),
		}, nil
	case *org.AddOrganizationRequest_Admin_Human:
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

func createdOrganizationToPb(createdOrg *command.CreatedOrg) (_ *org.AddOrganizationResponse, err error) {
	admins := make([]*org.AddOrganizationResponse_CreatedAdmin, len(createdOrg.CreatedAdmins))
	for i, admin := range createdOrg.CreatedAdmins {
		admins[i] = &org.AddOrganizationResponse_CreatedAdmin{
			UserId:    admin.ID,
			EmailCode: admin.EmailCode,
			PhoneCode: admin.PhoneCode,
		}
	}
	return &org.AddOrganizationResponse{
		Details:        object.DomainToDetailsPb(createdOrg.ObjectDetails),
		OrganizationId: createdOrg.ObjectDetails.ResourceOwner,
		CreatedAdmins:  admins,
	}, nil
}
