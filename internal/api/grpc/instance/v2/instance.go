package instance

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/instance/v2"
)

func (s *Server) DeleteInstance(ctx context.Context, request *instance.DeleteInstanceRequest) (*instance.DeleteInstanceResponse, error) {
	instanceID := strings.TrimSpace(request.GetInstanceId())
	if instanceID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "instance_id", "instance id must not be empty")
	}

	obj, err := s.command.RemoveInstance(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	return &instance.DeleteInstanceResponse{
		Details: object.DomainToDetailsPb(obj),
	}, nil

}

func (s *Server) UpdateInstance(ctx context.Context, request *instance.UpdateInstanceRequest) (*instance.UpdateInstanceResponse, error) {
	instanceName := strings.TrimSpace(request.GetInstanceName())
	if instanceName == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "instance_name", "instance name must not be empty")
	}

	obj, err := s.command.UpdateInstance(ctx, instanceName)
	if err != nil {
		return nil, err
	}

	return &instance.UpdateInstanceResponse{
		Details: object.DomainToDetailsPb(obj),
	}, nil
}
