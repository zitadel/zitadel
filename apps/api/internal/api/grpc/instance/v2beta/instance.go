package instance

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

func (s *Server) DeleteInstance(ctx context.Context, request *connect.Request[instance.DeleteInstanceRequest]) (*connect.Response[instance.DeleteInstanceResponse], error) {
	obj, err := s.command.RemoveInstance(ctx, request.Msg.GetInstanceId())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&instance.DeleteInstanceResponse{
		DeletionDate: timestamppb.New(obj.EventDate),
	}), nil

}

func (s *Server) UpdateInstance(ctx context.Context, request *connect.Request[instance.UpdateInstanceRequest]) (*connect.Response[instance.UpdateInstanceResponse], error) {
	obj, err := s.command.UpdateInstance(ctx, request.Msg.GetInstanceName())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&instance.UpdateInstanceResponse{
		ChangeDate: timestamppb.New(obj.EventDate),
	}), nil
}
