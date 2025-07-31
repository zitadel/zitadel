package internal_permission

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/zerrors"
	internal_permission "github.com/zitadel/zitadel/pkg/grpc/internal_permission/v2beta"
)

func (s *Server) CreateAdministrator(ctx context.Context, req *connect.Request[internal_permission.CreateAdministratorRequest]) (*connect.Response[internal_permission.CreateAdministratorResponse], error) {
	var creationDate *timestamppb.Timestamp

	switch resource := req.Msg.GetResource().GetResource().(type) {
	case *internal_permission.ResourceType_Instance:
		if resource.Instance {
			member, err := s.command.AddInstanceMember(ctx, createAdministratorInstanceToCommand(authz.GetInstance(ctx).InstanceID(), req.Msg.UserId, req.Msg.Roles))
			if err != nil {
				return nil, err
			}
			if !member.EventDate.IsZero() {
				creationDate = timestamppb.New(member.EventDate)
			}
		}
	case *internal_permission.ResourceType_OrganizationId:
		member, err := s.command.AddOrgMember(ctx, createAdministratorOrganizationToCommand(resource, req.Msg.UserId, req.Msg.Roles))
		if err != nil {
			return nil, err
		}
		if !member.EventDate.IsZero() {
			creationDate = timestamppb.New(member.EventDate)
		}
	case *internal_permission.ResourceType_ProjectId:
		member, err := s.command.AddProjectMember(ctx, createAdministratorProjectToCommand(resource, req.Msg.UserId, req.Msg.Roles))
		if err != nil {
			return nil, err
		}
		if !member.EventDate.IsZero() {
			creationDate = timestamppb.New(member.EventDate)
		}
	case *internal_permission.ResourceType_ProjectGrant_:
		member, err := s.command.AddProjectGrantMember(ctx, createAdministratorProjectGrantToCommand(resource, req.Msg.UserId, req.Msg.Roles))
		if err != nil {
			return nil, err
		}
		if !member.EventDate.IsZero() {
			creationDate = timestamppb.New(member.EventDate)
		}
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "ADMIN-IbPp47HDP5", "Errors.Invalid.Argument")
	}

	return connect.NewResponse(&internal_permission.CreateAdministratorResponse{
		CreationDate: creationDate,
	}), nil
}

func createAdministratorInstanceToCommand(instanceID, userID string, roles []string) *command.AddInstanceMember {
	return &command.AddInstanceMember{
		InstanceID: instanceID,
		UserID:     userID,
		Roles:      roles,
	}
}

func createAdministratorOrganizationToCommand(req *internal_permission.ResourceType_OrganizationId, userID string, roles []string) *command.AddOrgMember {
	return &command.AddOrgMember{
		OrgID:  req.OrganizationId,
		UserID: userID,
		Roles:  roles,
	}
}

func createAdministratorProjectToCommand(req *internal_permission.ResourceType_ProjectId, userID string, roles []string) *command.AddProjectMember {
	return &command.AddProjectMember{
		ProjectID: req.ProjectId,
		UserID:    userID,
		Roles:     roles,
	}
}

func createAdministratorProjectGrantToCommand(req *internal_permission.ResourceType_ProjectGrant_, userID string, roles []string) *command.AddProjectGrantMember {
	return &command.AddProjectGrantMember{
		GrantID:   req.ProjectGrant.ProjectGrantId,
		ProjectID: req.ProjectGrant.ProjectId,
		UserID:    userID,
		Roles:     roles,
	}
}

func (s *Server) UpdateAdministrator(ctx context.Context, req *connect.Request[internal_permission.UpdateAdministratorRequest]) (*connect.Response[internal_permission.UpdateAdministratorResponse], error) {
	var changeDate *timestamppb.Timestamp

	switch resource := req.Msg.GetResource().GetResource().(type) {
	case *internal_permission.ResourceType_Instance:
		if resource.Instance {
			member, err := s.command.ChangeInstanceMember(ctx, updateAdministratorInstanceToCommand(authz.GetInstance(ctx).InstanceID(), req.Msg.UserId, req.Msg.Roles))
			if err != nil {
				return nil, err
			}
			if !member.EventDate.IsZero() {
				changeDate = timestamppb.New(member.EventDate)
			}
		}
	case *internal_permission.ResourceType_OrganizationId:
		member, err := s.command.ChangeOrgMember(ctx, updateAdministratorOrganizationToCommand(resource, req.Msg.UserId, req.Msg.Roles))
		if err != nil {
			return nil, err
		}
		if !member.EventDate.IsZero() {
			changeDate = timestamppb.New(member.EventDate)
		}
	case *internal_permission.ResourceType_ProjectId:
		member, err := s.command.ChangeProjectMember(ctx, updateAdministratorProjectToCommand(resource, req.Msg.UserId, req.Msg.Roles))
		if err != nil {
			return nil, err
		}
		if !member.EventDate.IsZero() {
			changeDate = timestamppb.New(member.EventDate)
		}
	case *internal_permission.ResourceType_ProjectGrant_:
		member, err := s.command.ChangeProjectGrantMember(ctx, updateAdministratorProjectGrantToCommand(resource, req.Msg.UserId, req.Msg.Roles))
		if err != nil {
			return nil, err
		}
		if !member.EventDate.IsZero() {
			changeDate = timestamppb.New(member.EventDate)
		}
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "ADMIN-i0V2IbdloZ", "Errors.Invalid.Argument")
	}

	return connect.NewResponse(&internal_permission.UpdateAdministratorResponse{
		ChangeDate: changeDate,
	}), nil
}

func updateAdministratorInstanceToCommand(instanceID, userID string, roles []string) *command.ChangeInstanceMember {
	return &command.ChangeInstanceMember{
		InstanceID: instanceID,
		UserID:     userID,
		Roles:      roles,
	}
}

func updateAdministratorOrganizationToCommand(req *internal_permission.ResourceType_OrganizationId, userID string, roles []string) *command.ChangeOrgMember {
	return &command.ChangeOrgMember{
		OrgID:  req.OrganizationId,
		UserID: userID,
		Roles:  roles,
	}
}

func updateAdministratorProjectToCommand(req *internal_permission.ResourceType_ProjectId, userID string, roles []string) *command.ChangeProjectMember {
	return &command.ChangeProjectMember{
		ProjectID: req.ProjectId,
		UserID:    userID,
		Roles:     roles,
	}
}

func updateAdministratorProjectGrantToCommand(req *internal_permission.ResourceType_ProjectGrant_, userID string, roles []string) *command.ChangeProjectGrantMember {
	return &command.ChangeProjectGrantMember{
		GrantID:   req.ProjectGrant.ProjectGrantId,
		ProjectID: req.ProjectGrant.ProjectId,
		UserID:    userID,
		Roles:     roles,
	}
}

func (s *Server) DeleteAdministrator(ctx context.Context, req *connect.Request[internal_permission.DeleteAdministratorRequest]) (*connect.Response[internal_permission.DeleteAdministratorResponse], error) {
	var deletionDate *timestamppb.Timestamp

	switch resource := req.Msg.GetResource().GetResource().(type) {
	case *internal_permission.ResourceType_Instance:
		if resource.Instance {
			member, err := s.command.RemoveInstanceMember(ctx, authz.GetInstance(ctx).InstanceID(), req.Msg.UserId)
			if err != nil {
				return nil, err
			}
			if !member.EventDate.IsZero() {
				deletionDate = timestamppb.New(member.EventDate)
			}
		}
	case *internal_permission.ResourceType_OrganizationId:
		member, err := s.command.RemoveOrgMember(ctx, resource.OrganizationId, req.Msg.UserId)
		if err != nil {
			return nil, err
		}
		if !member.EventDate.IsZero() {
			deletionDate = timestamppb.New(member.EventDate)
		}
	case *internal_permission.ResourceType_ProjectId:
		member, err := s.command.RemoveProjectMember(ctx, resource.ProjectId, req.Msg.UserId, "")
		if err != nil {
			return nil, err
		}
		if !member.EventDate.IsZero() {
			deletionDate = timestamppb.New(member.EventDate)
		}
	case *internal_permission.ResourceType_ProjectGrant_:
		member, err := s.command.RemoveProjectGrantMember(ctx, resource.ProjectGrant.ProjectId, req.Msg.UserId, resource.ProjectGrant.ProjectGrantId)
		if err != nil {
			return nil, err
		}
		if !member.EventDate.IsZero() {
			deletionDate = timestamppb.New(member.EventDate)
		}
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "ADMIN-3UOjLtuohh", "Errors.Invalid.Argument")
	}

	return connect.NewResponse(&internal_permission.DeleteAdministratorResponse{
		DeletionDate: deletionDate,
	}), nil
}
