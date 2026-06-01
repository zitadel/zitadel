package orgv2

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/api/org/v2/convert"
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	v2_org "github.com/zitadel/zitadel/pkg/grpc/org/v2"
	v2beta_org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

// =================
// v2Beta endpoints
// =================

func UpdateOrganizationBeta(ctx context.Context, request *connect.Request[v2beta_org.UpdateOrganizationRequest]) (*connect.Response[v2beta_org.UpdateOrganizationResponse], error) {
	orgUpdateCmd := domain.NewUpdateOrgCommand(request.Msg.GetId(), request.Msg.GetName())

	// TODO(IAM-Marco) Finish implementation in https://github.com/zitadel/zitadel/issues/10447
	domainAddCmd := domain.NewAddOrgDomainCommand(request.Msg.GetId(), request.Msg.GetName())
	domainSetPrimaryCmd := domain.NewSetPrimaryOrgDomainCommand(request.Msg.GetId(), request.Msg.GetName())
	domainRemoveCmd := domain.NewRemoveOrgDomainCommand(request.Msg.GetId(), orgUpdateCmd.OldDomainName, orgUpdateCmd.IsOldDomainVerified)

	batchCmd := domain.BatchExecutors(orgUpdateCmd, domainAddCmd, domainSetPrimaryCmd, domainRemoveCmd)

	err := domain.Invoke(ctx, batchCmd, domain.WithOrganizationRepo(repository.OrganizationRepository()))
	if err != nil {
		return nil, err
	}

	return &connect.Response[v2beta_org.UpdateOrganizationResponse]{
		Msg: &v2beta_org.UpdateOrganizationResponse{
			// TODO(IAM-Marco): Change this with the real update date when OrganizationRepo.Update()
			// returns the timestamp. See https://github.com/zitadel/zitadel/issues/10881
			ChangeDate: timestamppb.Now(),
		},
	}, nil
}

// TODO(IAM-Marco): Remove in V5 (see https://github.com/zitadel/zitadel/issues/10877)
func ListOrganizationsBeta(ctx context.Context, request *connect.Request[v2beta_org.ListOrganizationsRequest]) (*connect.Response[v2beta_org.ListOrganizationsResponse], error) {
	orgListQuery := domain.NewListOrgsQuery(convert.OrganizationBetaRequestToV2Request(request.Msg))

	err := domain.Invoke(ctx, orgListQuery,
		domain.WithOrganizationRepo(repository.OrganizationRepository()),
		domain.WithOrganizationDomainRepo(repository.OrganizationDomainRepository()),
	)
	if err != nil {
		return nil, err
	}

	orgs := convert.DomainOrganizationListModelToGRPCBetaResponse(orgListQuery.Result())
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

func DeleteOrganizationBeta(ctx context.Context, request *connect.Request[v2beta_org.DeleteOrganizationRequest]) (*connect.Response[v2beta_org.DeleteOrganizationResponse], error) {
	orgDeleteCmd := domain.NewDeleteOrgCommand(request.Msg.GetId())

	err := domain.Invoke(ctx, orgDeleteCmd,
		domain.WithOrganizationRepo(repository.OrganizationRepository()),
		domain.WithProjectRepo(repository.ProjectRepository()),
	)
	if err != nil {
		var notFoundError *database.NoRowFoundError
		if errors.As(err, &notFoundError) {
			return connect.NewResponse(&v2beta_org.DeleteOrganizationResponse{}), nil
		}
		return nil, err
	}

	return &connect.Response[v2beta_org.DeleteOrganizationResponse]{
		Msg: &v2beta_org.DeleteOrganizationResponse{
			// TODO(IAM-Marco): Change this with the real update date when OrganizationRepo.Delete()
			// returns the timestamp. See https://github.com/zitadel/zitadel/issues/10881
			DeletionDate: timestamppb.Now(),
		},
	}, nil
}

func DeactivateOrganizationBeta(ctx context.Context, request *connect.Request[v2beta_org.DeactivateOrganizationRequest]) (*connect.Response[v2beta_org.DeactivateOrganizationResponse], error) {
	orgDeactivateCmd := domain.NewDeactivateOrgCommand(request.Msg.GetId())

	err := domain.Invoke(ctx, orgDeactivateCmd,
		domain.WithOrganizationRepo(repository.OrganizationRepository()),
	)
	if err != nil {
		return nil, err
	}

	return &connect.Response[v2beta_org.DeactivateOrganizationResponse]{
		Msg: &v2beta_org.DeactivateOrganizationResponse{
			// TODO(IAM-Marco): Change this with the real update date when OrganizationRepo.Update()
			// returns the timestamp. See https://github.com/zitadel/zitadel/issues/10881
			ChangeDate: timestamppb.Now(),
		},
	}, nil
}

func ActivateOrganizationBeta(ctx context.Context, request *connect.Request[v2beta_org.ActivateOrganizationRequest]) (*connect.Response[v2beta_org.ActivateOrganizationResponse], error) {
	orgActivateCmd := domain.NewActivateOrgCommand(request.Msg.GetId())

	err := domain.Invoke(ctx, orgActivateCmd,
		domain.WithOrganizationRepo(repository.OrganizationRepository()),
	)
	if err != nil {
		return nil, err
	}

	return &connect.Response[v2beta_org.ActivateOrganizationResponse]{
		Msg: &v2beta_org.ActivateOrganizationResponse{
			// TODO(IAM-Marco): Change this with the real update date when OrganizationRepo.Update()
			// returns the timestamp. See https://github.com/zitadel/zitadel/issues/10881
			ChangeDate: timestamppb.Now(),
		},
	}, nil
}

// =================
// v2 endpoints
// =================

func ListOrganizations(ctx context.Context, request *connect.Request[v2_org.ListOrganizationsRequest]) (*connect.Response[v2_org.ListOrganizationsResponse], error) {
	orgListQuery := domain.NewListOrgsQuery(request.Msg)

	err := domain.Invoke(ctx, orgListQuery,
		domain.WithOrganizationRepo(repository.OrganizationRepository()),
		domain.WithOrganizationDomainRepo(repository.OrganizationDomainRepository()),
	)
	if err != nil {
		return nil, err
	}

	orgs := convert.DomainOrganizationListModelToGRPCResponse(orgListQuery.Result())
	return &connect.Response[v2_org.ListOrganizationsResponse]{
		Msg: &v2_org.ListOrganizationsResponse{
			Result: orgs,
			Details: &object.ListDetails{
				// TODO(IAM-Marco): Return correct result once permissions are in place
				TotalResult: uint64(len(orgs)),
			},
			SortingColumn: request.Msg.GetSortingColumn(),
		},
	}, nil
}
