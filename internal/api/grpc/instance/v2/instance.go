package instance

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/pkg/grpc/instance/v2"
)

func (s *Server) DeleteInstance(ctx context.Context, request *connect.Request[instance.DeleteInstanceRequest]) (*connect.Response[instance.DeleteInstanceResponse], error) {
	// Deleting an instance is currently only allowed with system permissions,
	// so we directly check for them in the auth interceptor and do not check here again.
	obj, err := s.command.RemoveInstance(ctx, request.Msg.GetInstanceId())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&instance.DeleteInstanceResponse{
		DeletionDate: timestamppb.New(obj.EventDate),
	}), nil

}

func (s *Server) UpdateInstance(ctx context.Context, request *connect.Request[instance.UpdateInstanceRequest]) (*connect.Response[instance.UpdateInstanceResponse], error) {
	if err := s.checkPermission(ctx, domain.PermissionSystemInstanceWrite, domain.PermissionInstanceWrite); err != nil {
		return nil, err
	}
	obj, err := s.command.UpdateInstance(ctx, request.Msg.GetInstanceName())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&instance.UpdateInstanceResponse{
		ChangeDate: timestamppb.New(obj.EventDate),
	}), nil
}
