package group

import (
	"connectrpc.com/connect"
	"context"
	"github.com/zitadel/zitadel/internal/zerrors"
	group "github.com/zitadel/zitadel/pkg/grpc/group/v2beta"
)

func (s *Server) CreateGroup(ctx context.Context, req *connect.Request[group.CreateGroupRequest]) (*connect.Response[group.CreateGroupResponse], error) {
	return nil, zerrors.ThrowUnimplemented(nil, "GRP-qazd", "Not implemented")
}

func (s *Server) ListGroups(ctx context.Context, req *connect.Request[group.ListGroupsRequest]) (*connect.Response[group.ListGroupsResponse], error) {
	return nil, zerrors.ThrowUnimplemented(nil, "GRP-qazd", "Not implemented")
}

func (s *Server) UpdateGroup(ctx context.Context, req *connect.Request[group.UpdateGroupRequest]) (*connect.Response[group.UpdateGroupResponse], error) {
	return nil, zerrors.ThrowUnimplemented(nil, "GRP-qazd", "Not implemented")
}

func (s *Server) DeleteGroup(ctx context.Context, req *connect.Request[group.DeleteGroupRequest]) (*connect.Response[group.DeleteGroupResponse], error) {
	return nil, zerrors.ThrowUnimplemented(nil, "GRP-qazd", "Not implemented")
}
