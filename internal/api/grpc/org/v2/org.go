package org

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/api/grpc/user/v2"
	"github.com/zitadel/zitadel/internal/command"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	org "github.com/zitadel/zitadel/pkg/grpc/organisation/v2beta"
)

func (s *Server) AddOrganisation(ctx context.Context, request *org.AddOrganisationRequest) (*org.AddOrganisationResponse, error) {
	orgSetup, err := addOrganisationRequestToCommand(request)
	if err != nil {
		return nil, err
	}
	createdOrg, err := s.command.SetUpOrg(ctx, orgSetup, false)
	if err != nil {
		return nil, err
	}
	return createdOrganisationToPb(createdOrg)
}

func addOrganisationRequestToCommand(request *org.AddOrganisationRequest) (*command.OrgSetup, error) {
	admins, err := addOrganisationRequestAdminsToCommand(request.GetAdmins())
	if err != nil {
		return nil, err
	}
	return &command.OrgSetup{
		Name:         request.GetName(),
		CustomDomain: "",
		Admins:       admins,
	}, nil
}

func addOrganisationRequestAdminsToCommand(requestAdmins []*org.AddOrganisationRequest_Admin) (admins []*command.OrgSetupAdmin, err error) {
	admins = make([]*command.OrgSetupAdmin, len(requestAdmins))
	for i, admin := range requestAdmins {
		admins[i], err = addOrganisationRequestAdminToCommand(admin)
		if err != nil {
			return nil, err
		}
	}
	return admins, nil
}

func addOrganisationRequestAdminToCommand(admin *org.AddOrganisationRequest_Admin) (*command.OrgSetupAdmin, error) {
	switch a := admin.GetUserType().(type) {
	case *org.AddOrganisationRequest_Admin_UserId:
		return &command.OrgSetupAdmin{
			ID:    a.UserId,
			Roles: admin.GetRoles(),
		}, nil
	case *org.AddOrganisationRequest_Admin_Human:
		human, err := user.AddUserRequestToAddHuman(a.Human)
		if err != nil {
			return nil, err
		}
		return &command.OrgSetupAdmin{
			Human: human,
			Roles: admin.GetRoles(),
		}, nil
	default:
		return nil, caos_errs.ThrowUnimplementedf(nil, "ORGv2-SD2r1", "userType oneOf %T in method AddOrganisation not implemented", a)
	}
}

func createdOrganisationToPb(createdOrg *command.CreatedOrg) (_ *org.AddOrganisationResponse, err error) {
	admins := make([]*org.AddOrganisationResponse_CreatedAdmin, len(createdOrg.CreatedAdmins))
	for i, admin := range createdOrg.CreatedAdmins {
		admins[i] = &org.AddOrganisationResponse_CreatedAdmin{
			UserId:    admin.ID,
			EmailCode: admin.EmailCode,
			PhoneCode: admin.PhoneCode,
		}
	}
	return &org.AddOrganisationResponse{
		Details:        object.DomainToDetailsPb(createdOrg.ObjectDetails),
		OrganisationId: createdOrg.ObjectDetails.ResourceOwner,
		CreatedAdmins:  admins,
	}, nil
}
