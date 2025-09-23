package orgv2

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/api/org/v2/convert"
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	v2_org "github.com/zitadel/zitadel/pkg/grpc/org/v2"
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

	// TODO(IAM-Marco): I noticed while debugging that this is calling twice the commands (I think?)
	// It's hard to debug, I haven't spent too much into it. Only drawback is pushing events twice.
	batchCmd := domain.BatchCommands(orgUpdtCmd, domainAddCmd, domainSetPrimaryCmd, domainRemoveCmd)

	err := domain.Invoke(ctx, batchCmd, domain.WithOrganizationRepo(repository.OrganizationRepository))
	if err != nil {
		return nil, err
	}

	return &connect.Response[v2beta_org.UpdateOrganizationResponse]{
		Msg: &v2beta_org.UpdateOrganizationResponse{
			// TODO(IAM-Marco): Change this with the real update date when OrganizationRepo.Update()
			// returns the timestamp
			ChangeDate: timestamppb.Now(),
		},
	}, nil
}

func ListOrganizations(ctx context.Context, request *connect.Request[v2_org.ListOrganizationsRequest]) (*connect.Response[v2_org.ListOrganizationsResponse], error) {
	orgListCmd := domain.NewListOrgsCommand(request.Msg)

	err := domain.Invoke(ctx, orgListCmd)
	if err != nil {
		return nil, err
	}

	orgs := orgListCmd.ResultToGRPC()
	return &connect.Response[v2_org.ListOrganizationsResponse]{
		Msg: &v2_org.ListOrganizationsResponse{
			Result: orgListCmd.ResultToGRPC(),
			Details: &object.ListDetails{
				// TODO(IAM-Marco): Return correct result once permissions are in place
				TotalResult: uint64(len(orgs)),
			},
			SortingColumn: request.Msg.GetSortingColumn(),
		},
	}, nil
}

// TODO(IAM-Marco): Remove in V5
func ListOrganizationsBeta(ctx context.Context, request *connect.Request[v2beta_org.ListOrganizationsRequest]) (*connect.Response[v2beta_org.ListOrganizationsResponse], error) {
	orgListCmd := domain.NewListOrgsCommand(convert.OrganizationBetaRequestToV2Request(request.Msg))

	err := domain.Invoke(ctx, orgListCmd)
	if err != nil {
		return nil, err
	}

	orgs := orgListCmd.ResultToGRPCBeta()
	return &connect.Response[v2beta_org.ListOrganizationsResponse]{
		Msg: &v2beta_org.ListOrganizationsResponse{
			Organizations: orgs,
			Pagination: &filter.PaginationResponse{
				TotalResult:  uint64(len(orgs)),
				AppliedLimit: uint64(request.Msg.GetPagination().GetLimit()),
			},
		},
	}, nil
}
