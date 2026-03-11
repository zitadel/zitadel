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
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
)

func UpdateOrganization(ctx context.Context, request *connect.Request[org.UpdateOrganizationRequest]) (*connect.Response[org.UpdateOrganizationResponse], error) {
	orgUpdateCmd := domain.NewUpdateOrgCommand(request.Msg.GetOrganizationId(), request.Msg.GetName())

	// TODO(IAM-Marco) Finish implementation in https://github.com/zitadel/zitadel/issues/10447
	domainAddCmd := domain.NewAddOrgDomainCommand(request.Msg.GetOrganizationId(), request.Msg.GetName())
	domainSetPrimaryCmd := domain.NewSetPrimaryOrgDomainCommand(request.Msg.GetOrganizationId(), request.Msg.GetName())
	domainRemoveCmd := domain.NewRemoveOrgDomainCommand(request.Msg.GetOrganizationId(), orgUpdateCmd.OldDomainName, orgUpdateCmd.IsOldDomainVerified)

	batchCmd := domain.BatchExecutors(orgUpdateCmd, domainAddCmd, domainSetPrimaryCmd, domainRemoveCmd)

	err := domain.Invoke(ctx, batchCmd, domain.WithOrganizationRepo(repository.OrganizationRepository()))
	if err != nil {
		return nil, err
	}

	return &connect.Response[org.UpdateOrganizationResponse]{
		Msg: &org.UpdateOrganizationResponse{
			// TODO(IAM-Marco): Change this with the real update date when OrganizationRepo.Update()
			// returns the timestamp. See https://github.com/zitadel/zitadel/issues/10881
			ChangeDate: timestamppb.Now(),
		},
	}, nil
}

func DeleteOrganization(ctx context.Context, request *connect.Request[org.DeleteOrganizationRequest]) (*connect.Response[org.DeleteOrganizationResponse], error) {
	orgDeleteCmd := domain.NewDeleteOrgCommand(request.Msg.GetOrganizationId())

	err := domain.Invoke(ctx, orgDeleteCmd,
		domain.WithOrganizationRepo(repository.OrganizationRepository()),
		domain.WithProjectRepo(repository.ProjectRepository()),
	)
	if err != nil {
		var notFoundError *database.NoRowFoundError
		if errors.As(err, &notFoundError) {
			return connect.NewResponse(&org.DeleteOrganizationResponse{}), nil
		}
		return nil, err
	}

	return &connect.Response[org.DeleteOrganizationResponse]{
		Msg: &org.DeleteOrganizationResponse{
			// TODO(IAM-Marco): Change this with the real update date when OrganizationRepo.Delete()
			// returns the timestamp. See https://github.com/zitadel/zitadel/issues/10881
			DeletionDate: timestamppb.Now(),
		},
	}, nil
}

func DeactivateOrganization(ctx context.Context, request *connect.Request[org.DeactivateOrganizationRequest]) (*connect.Response[org.DeactivateOrganizationResponse], error) {
	orgDeactivateCmd := domain.NewDeactivateOrgCommand(request.Msg.GetOrganizationId())

	err := domain.Invoke(ctx, orgDeactivateCmd,
		domain.WithOrganizationRepo(repository.OrganizationRepository()),
	)
	if err != nil {
		return nil, err
	}

	return &connect.Response[org.DeactivateOrganizationResponse]{
		Msg: &org.DeactivateOrganizationResponse{
			// TODO(IAM-Marco): Change this with the real update date when OrganizationRepo.Update()
			// returns the timestamp. See https://github.com/zitadel/zitadel/issues/10881
			ChangeDate: timestamppb.Now(),
		},
	}, nil
}

func ActivateOrganization(ctx context.Context, request *connect.Request[org.ActivateOrganizationRequest]) (*connect.Response[org.ActivateOrganizationResponse], error) {
	orgActivateCmd := domain.NewActivateOrgCommand(request.Msg.GetOrganizationId())

	err := domain.Invoke(ctx, orgActivateCmd,
		domain.WithOrganizationRepo(repository.OrganizationRepository()),
	)
	if err != nil {
		return nil, err
	}

	return &connect.Response[org.ActivateOrganizationResponse]{
		Msg: &org.ActivateOrganizationResponse{
			// TODO(IAM-Marco): Change this with the real update date when OrganizationRepo.Update()
			// returns the timestamp. See https://github.com/zitadel/zitadel/issues/10881
			ChangeDate: timestamppb.Now(),
		},
	}, nil
}

func ListOrganizations(ctx context.Context, request *connect.Request[org.ListOrganizationsRequest]) (*connect.Response[org.ListOrganizationsResponse], error) {
	orgListQuery := domain.NewListOrgsQuery(request.Msg)

	err := domain.Invoke(ctx, orgListQuery,
		domain.WithOrganizationRepo(repository.OrganizationRepository()),
		domain.WithOrganizationDomainRepo(repository.OrganizationDomainRepository()),
	)
	if err != nil {
		return nil, err
	}

	orgs := convert.DomainOrganizationListModelToGRPCResponse(orgListQuery.Result())
	return &connect.Response[org.ListOrganizationsResponse]{
		Msg: &org.ListOrganizationsResponse{
			Result: orgs,
			Details: &object.ListDetails{
				// TODO(IAM-Marco): Return correct result once permissions are in place
				TotalResult: uint64(len(orgs)),
			},
			SortingColumn: request.Msg.GetSortingColumn(),
		},
	}, nil
}
