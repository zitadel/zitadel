package group

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/repository/group"
	group_v2 "github.com/zitadel/zitadel/pkg/grpc/group/v2"
)

func protoAttrsToGroupUser(protoAttrs map[string]*structpb.Value) map[string]group.AttributeValue {
	if len(protoAttrs) == 0 {
		return nil
	}
	attrs := make(map[string]group.AttributeValue, len(protoAttrs))
	for k, v := range protoAttrs {
		if v == nil {
			continue
		}
		switch val := v.Kind.(type) {
		case *structpb.Value_StringValue:
			attrs[k] = group.AttributeValue{val.StringValue}
		case *structpb.Value_ListValue:
			if val.ListValue != nil && len(val.ListValue.Values) > 0 {
				strVals := make([]string, 0, len(val.ListValue.Values))
				for _, item := range val.ListValue.Values {
					if str, ok := item.Kind.(*structpb.Value_StringValue); ok {
						strVals = append(strVals, str.StringValue)
					}
				}
				if len(strVals) > 0 {
					attrs[k] = group.AttributeValue(strVals)
				}
			}
		}
	}
	return attrs
}

func (s *Server) AddUsersToGroup(ctx context.Context, c *connect.Request[group_v2.AddUsersToGroupRequest]) (*connect.Response[group_v2.AddUsersToGroupResponse], error) {
	reqUsers := c.Msg.GetUsers()
	users := make([]group.GroupUser, 0, len(reqUsers))
	for _, u := range reqUsers {
		users = append(users, group.GroupUser{
			UserID:     u.GetUserId(),
			Attributes: protoAttrsToGroupUser(u.GetAttributes()),
		})
	}
	addUsersResp, err := s.command.AddUsersToGroup(ctx, c.Msg.GetId(), users)
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
