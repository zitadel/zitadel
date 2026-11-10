package group

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	group_v2 "github.com/zitadel/zitadel/pkg/grpc/group/v2"
)

func (s *Server) AddUsersToGroup(ctx context.Context, c *connect.Request[group_v2.AddUsersToGroupRequest]) (*connect.Response[group_v2.AddUsersToGroupResponse], error) {
	addUsersResp, err := s.command.AddUsersToGroup(ctx, c.Msg.GetId(), c.Msg.GetUserIds())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&group_v2.AddUsersToGroupResponse{
		ChangeDate: timestamppb.New(addUsersResp.EventDate),
	}), nil
}

func (s *Server) RemoveUsersFromGroup(ctx context.Context, c *connect.Request[group_v2.RemoveUsersFromGroupRequest]) (*connect.Response[group_v2.RemoveUsersFromGroupResponse], error) {
	details, err := s.command.RemoveUsersFromGroup(ctx, c.Msg.GetId(), c.Msg.GetUserIds())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&group_v2.RemoveUsersFromGroupResponse{
		ChangeDate: timestamppb.New(details.EventDate),
	}), nil
}
