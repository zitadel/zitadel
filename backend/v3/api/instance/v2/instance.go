package instancev2

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/api/instance/v2/convert"
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	filter_v2 "github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	instance_v2 "github.com/zitadel/zitadel/pkg/grpc/instance/v2"
)

func DeleteInstance(ctx context.Context, request *connect.Request[instance_v2.DeleteInstanceRequest]) (*connect.Response[instance_v2.DeleteInstanceResponse], error) {
	instanceDeleteCmd := domain.NewDeleteInstanceCommand(request.Msg.GetInstanceId())

	err := domain.Invoke(ctx, instanceDeleteCmd, domain.WithInstanceRepo(repository.InstanceRepository()))

	if err != nil {
		return nil, err
	}

	var deletionDate *timestamppb.Timestamp
	if instanceDeleteCmd.DeleteTime != nil {
		deletionDate = timestamppb.New(*instanceDeleteCmd.DeleteTime)
	}
	return &connect.Response[instance_v2.DeleteInstanceResponse]{
		Msg: &instance_v2.DeleteInstanceResponse{
			// TODO(IAM-Marco): Change this with the real update date when OrganizationRepo.Update()
			// returns the timestamp
			DeletionDate: deletionDate,
		},
	}, nil
}

func GetInstance(ctx context.Context, request *connect.Request[instance_v2.GetInstanceRequest]) (*connect.Response[instance_v2.GetInstanceResponse], error) {
	instanceGetCmd := domain.NewGetInstanceCommand(request.Msg.GetInstanceId())

	err := domain.Invoke(ctx, instanceGetCmd, domain.WithInstanceRepo(repository.InstanceRepository()))

	if err != nil {
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
