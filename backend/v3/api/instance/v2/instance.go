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
	"github.com/zitadel/zitadel/internal/zerrors"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

func DeleteInstance(ctx context.Context, request *connect.Request[instance.DeleteInstanceRequest]) (*connect.Response[instance.DeleteInstanceResponse], error) {
	instanceDeleteCmd := domain.NewDeleteInstanceCommand(request.Msg.GetInstanceId())

	err := domain.Invoke(ctx, instanceDeleteCmd, domain.WithInstanceRepo(repository.InstanceRepository()))

	if err != nil {
		if errors.Is(err, &database.NoRowFoundError{}) {
			return nil, zerrors.ThrowNotFound(err, "INST-QVrUwc", "instance not found")
		}
		return nil, err
	}

	return &connect.Response[instance.DeleteInstanceResponse]{
		Msg: &instance.DeleteInstanceResponse{
			// TODO(IAM-Marco): Change this with the real update date when OrganizationRepo.Update()
			// returns the timestamp
			DeletionDate: timestamppb.Now(),
		},
	}, nil
}

func GetInstance(ctx context.Context, request *connect.Request[instance.GetInstanceRequest]) (*connect.Response[instance.GetInstanceResponse], error) {
	instanceGetCmd := domain.NewGetInstanceCommand(request.Msg.GetInstanceId())

	err := domain.Invoke(ctx, instanceGetCmd, domain.WithInstanceRepo(repository.InstanceRepository()))

	if err != nil {
		if errors.Is(err, &database.NoRowFoundError{}) {
			return nil, zerrors.ThrowNotFound(err, "INST-QVrUwc", "instance not found")
		}
		return nil, err
	}

	return &connect.Response[instance.GetInstanceResponse]{
		Msg: &instance.GetInstanceResponse{
			Instance: convert.DomainInstanceModelToGRPCResponse(instanceGetCmd.Result()),
		},
	}, nil
}

func UpdateInstance(ctx context.Context, request *connect.Request[instance.UpdateInstanceRequest]) (*connect.Response[instance.UpdateInstanceResponse], error) {
	instanceUpdateCmd := domain.NewUpdateInstanceCommand(request.Msg.GetInstanceId(), request.Msg.GetInstanceName())

	err := domain.Invoke(ctx, instanceUpdateCmd, domain.WithInstanceRepo(repository.InstanceRepository()))

	if err != nil {
		return nil, err
	}

	return &connect.Response[instance.UpdateInstanceResponse]{
		Msg: &instance.UpdateInstanceResponse{
			// TODO(IAM-Marco): Change this with the real update date when InstanceRepo.Update()
			// returns the timestamp
			ChangeDate: timestamppb.Now(),
		},
	}, nil
}

func ListInstances(ctx context.Context, request *connect.Request[instance.ListInstancesRequest]) (*connect.Response[instance.ListInstancesResponse], error) {
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
	return &connect.Response[instance.ListInstancesResponse]{
		Msg: &instance.ListInstancesResponse{
			Instances: convert.DomainInstanceListModelToGRPCResponse(instances),
			Pagination: &filter.PaginationResponse{
				TotalResult:  uint64(len(instances)),
				AppliedLimit: uint64(request.Msg.GetPagination().GetLimit()),
			},
		},
	}, nil
}
