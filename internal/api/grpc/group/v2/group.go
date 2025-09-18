package group

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	group_v2 "github.com/zitadel/zitadel/pkg/grpc/group/v2"
)

// CreateGroup creates a new user group in the specified organization.
func (s *Server) CreateGroup(ctx context.Context, req *connect.Request[group_v2.CreateGroupRequest]) (*connect.Response[group_v2.CreateGroupResponse], error) {
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
	return connect.NewResponse(&group_v2.CreateGroupResponse{
		Id:           groupDetails.ID,
		CreationDate: timestamppb.New(groupDetails.EventDate),
	}), nil
}

// UpdateGroup updates a user group.
func (s *Server) UpdateGroup(ctx context.Context, req *connect.Request[group_v2.UpdateGroupRequest]) (*connect.Response[group_v2.UpdateGroupResponse], error) {
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
	return connect.NewResponse(&group_v2.UpdateGroupResponse{
		ChangeDate: timestamppb.New(details.EventDate),
	}), nil
}

// DeleteGroup deletes a user group from an organization.
func (s *Server) DeleteGroup(ctx context.Context, req *connect.Request[group_v2.DeleteGroupRequest]) (*connect.Response[group_v2.DeleteGroupResponse], error) {
	details, err := s.command.DeleteGroup(ctx, req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	var deletionDate *timestamppb.Timestamp
	if !details.EventDate.IsZero() {
		deletionDate = timestamppb.New(details.EventDate)
	}
	return connect.NewResponse(&group_v2.DeleteGroupResponse{
		DeletionDate: deletionDate,
	}), nil
}
