package instance

import (
	"context"
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"

	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

func (s *Server) DeleteInstance(ctx context.Context, request *instance.DeleteInstanceRequest) (*instance.DeleteInstanceResponse, error) {
	instanceID := strings.TrimSpace(request.GetInstanceId())
	if err := validateParam(instanceID, "instance_id"); err != nil {
		return nil, err
	}

	obj, err := s.command.RemoveInstance(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	return &instance.DeleteInstanceResponse{
		DeletionDate: timestamppb.New(obj.EventDate),
	}, nil

}

func (s *Server) UpdateInstance(ctx context.Context, request *instance.UpdateInstanceRequest) (*instance.UpdateInstanceResponse, error) {
	instanceName := strings.TrimSpace(request.GetInstanceName())
	if err := validateParam(instanceName, "instance_name"); err != nil {
		return nil, err
	}

	obj, err := s.command.UpdateInstance(ctx, instanceName)
	if err != nil {
		return nil, err
	}

	return &instance.UpdateInstanceResponse{
		ChangeDate: timestamppb.New(obj.EventDate),
	}, nil
}
