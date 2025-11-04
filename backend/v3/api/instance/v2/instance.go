package instancev2

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/api/instance/v2/convert"
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
	filter_v2 "github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	filter_v2beta "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	instance_v2 "github.com/zitadel/zitadel/pkg/grpc/instance/v2"
	instance_v2beta "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

func DeleteInstanceBeta(ctx context.Context, request *connect.Request[instance_v2beta.DeleteInstanceRequest]) (*connect.Response[instance_v2beta.DeleteInstanceResponse], error) {
	instanceDeleteCmd := domain.NewDeleteInstanceCommand(request.Msg.GetInstanceId())

	err := domain.Invoke(ctx, instanceDeleteCmd, domain.WithInstanceRepo(repository.InstanceRepository()))

	if err != nil {
		if errors.Is(err, &database.NoRowFoundError{}) {
			return &connect.Response[instance_v2beta.DeleteInstanceResponse]{}, nil
		}
		return nil, err
	}

	return &connect.Response[instance_v2beta.DeleteInstanceResponse]{
		Msg: &instance_v2beta.DeleteInstanceResponse{
			// TODO(IAM-Marco): Change this with the real update date when OrganizationRepo.Update()
			// returns the timestamp
			DeletionDate: timestamppb.Now(),
		},
	}, nil
}

func GetInstanceBeta(ctx context.Context, request *connect.Request[instance_v2beta.GetInstanceRequest]) (*connect.Response[instance_v2beta.GetInstanceResponse], error) {
	instanceGetCmd := domain.NewGetInstanceCommand(request.Msg.GetInstanceId())

	err := domain.Invoke(ctx, instanceGetCmd, domain.WithInstanceRepo(repository.InstanceRepository()))

	if err != nil {
		if errors.Is(err, &database.NoRowFoundError{}) {
			return nil, zerrors.ThrowNotFound(err, "INST-QVrUwc", "instance not found")
		}
		return nil, err
	}

	return &connect.Response[instance_v2beta.GetInstanceResponse]{
		Msg: &instance_v2beta.GetInstanceResponse{
			Instance: convert.DomainInstanceModelToGRPCBetaResponse(instanceGetCmd.Result()),
		},
	}, nil
}

func UpdateInstanceBeta(ctx context.Context, request *connect.Request[instance_v2beta.UpdateInstanceRequest]) (*connect.Response[instance_v2beta.UpdateInstanceResponse], error) {
	instanceUpdateCmd := domain.NewUpdateInstanceCommand(request.Msg.GetInstanceId(), request.Msg.GetInstanceName())

	err := domain.Invoke(ctx, instanceUpdateCmd, domain.WithInstanceRepo(repository.InstanceRepository()))

	if err != nil {
		return nil, err
	}

	return &connect.Response[instance_v2beta.UpdateInstanceResponse]{
		Msg: &instance_v2beta.UpdateInstanceResponse{
			// TODO(IAM-Marco): Change this with the real update date when InstanceRepo.Update()
			// returns the timestamp
			ChangeDate: timestamppb.Now(),
		},
	}, nil
}

func ListInstancesBeta(ctx context.Context, request *connect.Request[instance_v2beta.ListInstancesRequest]) (*connect.Response[instance_v2beta.ListInstancesResponse], error) {
	instancesListCmd := domain.NewListInstancesCommand(convert.ListInstancesBetaRequestToV2Request(request.Msg))

	err := domain.Invoke(
		ctx,
		instancesListCmd,
		domain.WithInstanceRepo(repository.InstanceRepository()),
		domain.WithInstanceDomainRepo(repository.InstanceDomainRepository()),
	)
	if err != nil {
		return nil, err
	}

	instances := instancesListCmd.Result()
	return &connect.Response[instance_v2beta.ListInstancesResponse]{
		Msg: &instance_v2beta.ListInstancesResponse{
			Instances: convert.DomainInstanceListModelToGRPCBetaResponse(instances),
			Pagination: &filter_v2beta.PaginationResponse{
				// TODO(IAM-Marco): return correct value. Tracked in https://github.com/zitadel/zitadel/issues/10955
				TotalResult:  uint64(len(instances)),
				AppliedLimit: uint64(request.Msg.GetPagination().GetLimit()),
			},
		},
	}, nil
}

func DeleteInstance(ctx context.Context, request *connect.Request[instance_v2.DeleteInstanceRequest]) (*connect.Response[instance_v2.DeleteInstanceResponse], error) {
	instanceDeleteCmd := domain.NewDeleteInstanceCommand(request.Msg.GetInstanceId())

	err := domain.Invoke(ctx, instanceDeleteCmd, domain.WithInstanceRepo(repository.InstanceRepository()))

	if err != nil {
		if errors.Is(err, &database.NoRowFoundError{}) {
			return nil, nil
		}
		return nil, err
	}

	return &connect.Response[instance_v2.DeleteInstanceResponse]{
		Msg: &instance_v2.DeleteInstanceResponse{
			// TODO(IAM-Marco): Change this with the real update date when OrganizationRepo.Update()
			// returns the timestamp
			DeletionDate: timestamppb.Now(),
		},
	}, nil
}

func GetInstance(ctx context.Context, request *connect.Request[instance_v2.GetInstanceRequest]) (*connect.Response[instance_v2.GetInstanceResponse], error) {
	instanceGetCmd := domain.NewGetInstanceCommand(request.Msg.GetInstanceId())

	err := domain.Invoke(ctx, instanceGetCmd, domain.WithInstanceRepo(repository.InstanceRepository()))

	if err != nil {
		if errors.Is(err, &database.NoRowFoundError{}) {
			return nil, zerrors.ThrowNotFound(err, "INST-QVrUwc", "instance not found")
		}
		return nil, err
	}

	return &connect.Response[instance_v2.GetInstanceResponse]{
		Msg: &instance_v2.GetInstanceResponse{
			Instance: convert.DomainInstanceModelToGRPCResponse(instanceGetCmd.Result()),
		},
	}, nil
}

func UpdateInstance(ctx context.Context, request *connect.Request[instance_v2.UpdateInstanceRequest]) (*connect.Response[instance_v2.UpdateInstanceResponse], error) {
	instanceUpdateCmd := domain.NewUpdateInstanceCommand(request.Msg.GetInstanceId(), request.Msg.GetInstanceName())

	err := domain.Invoke(ctx, instanceUpdateCmd, domain.WithInstanceRepo(repository.InstanceRepository()))

	if err != nil {
		return nil, err
	}

	return &connect.Response[instance_v2.UpdateInstanceResponse]{
		Msg: &instance_v2.UpdateInstanceResponse{
			// TODO(IAM-Marco): Change this with the real update date when InstanceRepo.Update()
			// returns the timestamp
			ChangeDate: timestamppb.Now(),
		},
	}, nil
}

func ListInstances(ctx context.Context, request *connect.Request[instance_v2.ListInstancesRequest]) (*connect.Response[instance_v2.ListInstancesResponse], error) {
	instancesListCmd := domain.NewListInstancesCommand(request.Msg)

	err := domain.Invoke(
		ctx,
		instancesListCmd,
		domain.WithInstanceRepo(repository.InstanceRepository()),
		domain.WithInstanceDomainRepo(repository.InstanceDomainRepository()),
	)
	if err != nil {
		return nil, err
	}

	instances := instancesListCmd.Result()
	return &connect.Response[instance_v2.ListInstancesResponse]{
		Msg: &instance_v2.ListInstancesResponse{
			Instances: convert.DomainInstanceListModelToGRPCResponse(instances),
			Pagination: &filter_v2.PaginationResponse{
				// TODO(IAM-Marco): return correct value. Tracked in https://github.com/zitadel/zitadel/issues/10955
				TotalResult:  uint64(len(instances)),
				AppliedLimit: uint64(request.Msg.GetPagination().GetLimit()),
			},
		},
	}, nil
}

func AddCustomDomain(ctx context.Context, request *connect.Request[instance_v2beta.AddCustomDomainRequest]) (*connect.Response[instance_v2beta.AddCustomDomainResponse], error) {
	addCustomDomainCmd := domain.NewAddInstanceDomainCommand(request.Msg.GetInstanceId(), request.Msg.GetDomain())
	oidcConfigUpdateCmd := domain.NewOIDCConfigurationUpdate(request.Msg.GetDomain(), authz.GetInstance(ctx).ProjectID(), authz.GetInstance(ctx).ConsoleApplicationID())

	batchExec := domain.BatchExecutors(
		addCustomDomainCmd,
		oidcConfigUpdateCmd,
	)

	err := domain.Invoke(
		ctx,
		batchExec,
		domain.WithInstanceDomainRepo(repository.InstanceDomainRepository()),
		// domain.WithOIDCConfigurationRepo(repository.OIDCConfigurationRepository()),
	)

	if err != nil {
		return nil, err
	}

	return &connect.Response[instance_v2beta.AddCustomDomainResponse]{
		Msg: &instance_v2beta.AddCustomDomainResponse{
			// TODO(IAM-Marco): Return correct value. Tracked in https://github.com/zitadel/zitadel/issues/10881
			CreationDate: timestamppb.Now(),
		},
	}, nil
}
