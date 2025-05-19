package instance

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

func (s *Server) DeleteInstance(ctx context.Context, request *instance.DeleteInstanceRequest) (*instance.DeleteInstanceResponse, error) {
	obj, err := s.command.RemoveInstance(ctx, request.GetInstanceId())
	if err != nil {
		return nil, err
	}

	return &instance.DeleteInstanceResponse{
		DeletionDate: timestamppb.New(obj.EventDate),
	}, nil

}

func (s *Server) UpdateInstance(ctx context.Context, request *instance.UpdateInstanceRequest) (*instance.UpdateInstanceResponse, error) {
	obj, err := s.command.UpdateInstance(ctx, request.GetInstanceName())
	if err != nil {
		return nil, err
	}

	return &instance.UpdateInstanceResponse{
		ChangeDate: timestamppb.New(obj.EventDate),
	}, nil
}
