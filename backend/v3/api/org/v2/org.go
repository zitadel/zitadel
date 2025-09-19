package orgv2

import (
	"context"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	v2beta_org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

// import (
// 	"context"

// 	"github.com/zitadel/zitadel/backend/v3/domain"
// 	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
// )

// func CreateOrg(ctx context.Context, req *org.AddOrganizationRequest) (resp *org.AddOrganizationResponse, err error) {
// 	cmd := domain.NewAddOrgCommand(
// 		req.GetName(),
// 		addOrgAdminToCommand(req.GetAdmins()...)...,
// 	)
// 	err = domain.Invoke(ctx, cmd)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &org.AddOrganizationResponse{
// 		OrganizationId: cmd.ID,
// 	}, nil
// }

// func addOrgAdminToCommand(admins ...*org.AddOrganizationRequest_Admin) []*domain.AddMemberCommand {
// 	cmds := make([]*domain.AddMemberCommand, len(admins))
// 	for i, admin := range admins {
// 		cmds[i] = &domain.AddMemberCommand{
// 			UserID: admin.GetUserId(),
// 			Roles:  admin.GetRoles(),
// 		}
// 	}
// 	return cmds
// }

func UpdateOrganization(ctx context.Context, request *connect.Request[v2beta_org.UpdateOrganizationRequest]) (*connect.Response[v2beta_org.UpdateOrganizationResponse], error) {
	orgUpdtCmd := domain.NewUpdateOrgCommand(request.Msg.GetId(), request.Msg.GetName())

	// TODO(IAM-Marco) Finish implementation in https://github.com/zitadel/zitadel/issues/10447
	domainAddCmd := domain.NewAddOrgDomainCommand(request.Msg.GetId(), request.Msg.GetName())
	domainSetPrimaryCmd := domain.NewSetPrimaryOrgDomainCommand(request.Msg.GetId(), request.Msg.GetName())
	// TODO(IAM-Marco) Check if passing the pointer is actually working to retrieve the domain name and the DomainVerified
	domainRemoveCmd := domain.NewRemoveOrgDomainCommand(request.Msg.GetId(), orgUpdtCmd.OldDomainName, orgUpdtCmd.IsOldDomainVerified)

	batchCmd := domain.BatchCommands(orgUpdtCmd, domainAddCmd, domainSetPrimaryCmd, domainRemoveCmd)

	err := domain.Invoke(ctx, batchCmd, domain.WithOrganizationRepo(repository.OrganizationRepository(nil)))
	if err != nil {
		return nil, err
	}

	return &connect.Response[v2beta_org.UpdateOrganizationResponse]{
		Msg: &v2beta_org.UpdateOrganizationResponse{},
	}, nil
}
