package orgv2

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
