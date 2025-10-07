package instancev2

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/zerrors"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

func DeleteInstance(ctx context.Context, instanceID string) (*connect.Response[instance.DeleteInstanceResponse], error) {
	instanceDeleteCmd := domain.NewDeleteInstanceCommand(instanceID)

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

func GetInstance(ctx context.Context, instanceID string) (*connect.Response[instance.GetInstanceResponse], error) {
	instanceGetCmd := domain.NewGetInstanceCommand(instanceID)

	err := domain.Invoke(ctx, instanceGetCmd, domain.WithInstanceRepo(repository.InstanceRepository()))

	if err != nil {
		if errors.Is(err, &database.NoRowFoundError{}) {
			return nil, zerrors.ThrowNotFound(err, "INST-QVrUwc", "instance not found")
		}
		return nil, err
	}

	return &connect.Response[instance.GetInstanceResponse]{
		Msg: &instance.GetInstanceResponse{
			Instance: instanceGetCmd.ResultToGRPC(),
		},
	}, nil
}
