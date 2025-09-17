package group

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/pkg/grpc/group/v2"
)

// CreateGroup creates a new user group in the specified organization.
func (s *Server) CreateGroup(ctx context.Context, req *connect.Request[group.CreateGroupRequest]) (*connect.Response[group.CreateGroupResponse], error) {
	userGroup := &domain.Group{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   req.Msg.GetId(),
			ResourceOwner: authz.GetCtxData(ctx).ResourceOwner,
		},
		Name:           req.Msg.GetName(),
		Description:    req.Msg.GetDescription(),
		OrganizationID: req.Msg.GetOrganizationId(),
	}

	groupDetails, err := s.command.CreateGroup(ctx, userGroup)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&group.CreateGroupResponse{
		Id:           groupDetails.ID,
		CreationDate: timestamppb.New(groupDetails.CreationDate),
	}), nil
}

// UpdateGroup updates a user group.
func (s *Server) UpdateGroup(ctx context.Context, req *connect.Request[group.UpdateGroupRequest]) (*connect.Response[group.UpdateGroupResponse], error) {
	userGroup := &domain.Group{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   req.Msg.GetId(),
			ResourceOwner: authz.GetCtxData(ctx).ResourceOwner,
		},
		Name:        req.Msg.GetName(),
		Description: req.Msg.GetDescription(),
	}

	details, err := s.command.UpdateGroup(ctx, userGroup)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&group.UpdateGroupResponse{
		ChangeDate: timestamppb.New(details.EventDate),
	}), nil
}

// DeleteGroup deletes a user group from an organization.
func (s *Server) DeleteGroup(ctx context.Context, req *connect.Request[group.DeleteGroupRequest]) (*connect.Response[group.DeleteGroupResponse], error) {
	details, err := s.command.DeleteGroup(ctx, req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&group.DeleteGroupResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}), nil
}
