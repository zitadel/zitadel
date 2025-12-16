package instance

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	instancev2 "github.com/zitadel/zitadel/backend/v3/api/instance/v2"
	"github.com/zitadel/zitadel/internal/api/authz"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

func (s *Server) DeleteInstance(ctx context.Context, request *connect.Request[instance.DeleteInstanceRequest]) (*connect.Response[instance.DeleteInstanceResponse], error) {
	if authz.GetFeatures(ctx).EnableRelationalTables {
		return instancev2.DeleteInstance(ctx, request)
	}

	obj, err := s.command.RemoveInstance(ctx, request.Msg.GetInstanceId())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&instance.DeleteInstanceResponse{
		DeletionDate: timestamppb.New(obj.EventDate),
	}), nil

}

func (s *Server) UpdateInstance(ctx context.Context, request *connect.Request[instance.UpdateInstanceRequest]) (*connect.Response[instance.UpdateInstanceResponse], error) {
	if authz.GetFeatures(ctx).EnableRelationalTables {
		return instancev2.UpdateInstance(ctx, request)
	}

	obj, err := s.command.UpdateInstance(ctx, request.Msg.GetInstanceName())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&instance.UpdateInstanceResponse{
		ChangeDate: timestamppb.New(obj.EventDate),
	}), nil
}
