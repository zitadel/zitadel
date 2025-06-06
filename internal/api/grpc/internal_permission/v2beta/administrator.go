package internal_permission

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	internal_permission "github.com/zitadel/zitadel/pkg/grpc/internal_permission/v2beta"
)

func (s *Server) CreateAdministrator(ctx context.Context, req *internal_permission.CreateAdministratorRequest) (*internal_permission.CreateAdministratorResponse, error) {
	var creationDate *timestamppb.Timestamp

	switch resource := req.GetResource().Resource.(type) {
	case *internal_permission.ResourceType_Instance:
		if resource.Instance {
			member, err := s.command.AddInstanceMember(ctx, createAdministratorInstanceToCommand(authz.GetInstance(ctx).InstanceID(), req.UserId, req.Roles))
			if err != nil {
				return nil, err
			}
			if !member.CreationDate.IsZero() {
				creationDate = timestamppb.New(member.CreationDate)
			}
		}
	case *internal_permission.ResourceType_OrganizationId:
		member, err := s.command.AddOrgMember(ctx, createAdministratorOrganizationToCommand(resource, req.UserId, req.Roles))
		if err != nil {
			return nil, err
		}
		if !member.CreationDate.IsZero() {
			creationDate = timestamppb.New(member.CreationDate)
		}
	case *internal_permission.ResourceType_ProjectId:
		member, err := s.command.AddProjectMember(ctx, createAdministratorProjectToCommand(resource, req.UserId, req.Roles))
		if err != nil {
			return nil, err
		}
		if !member.CreationDate.IsZero() {
			creationDate = timestamppb.New(member.CreationDate)
		}
	case *internal_permission.ResourceType_ProjectGrant_:
		member, err := s.command.AddProjectGrantMember(ctx, createAdministratorProjectGrantToCommand(resource, req.UserId, req.Roles))
		if err != nil {
			return nil, err
		}
		if !member.CreationDate.IsZero() {
			creationDate = timestamppb.New(member.CreationDate)
		}
	}

	return &internal_permission.CreateAdministratorResponse{
		CreationDate: creationDate,
	}, nil
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

func (s *Server) UpdateAdministrator(ctx context.Context, req *internal_permission.UpdateAdministratorRequest) (*internal_permission.UpdateAdministratorResponse, error) {
	var changeDate *timestamppb.Timestamp

	switch resource := req.GetResource().Resource.(type) {
	case *internal_permission.ResourceType_Instance:
		if resource.Instance {
			member, err := s.command.ChangeInstanceMember(ctx, updateAdministratorInstanceToCommand(authz.GetInstance(ctx).InstanceID(), req.UserId, req.Roles))
			if err != nil {
				return nil, err
			}
			if !member.EventDate.IsZero() {
				changeDate = timestamppb.New(member.CreationDate)
			}
		}
	case *internal_permission.ResourceType_OrganizationId:
		member, err := s.command.ChangeOrgMember(ctx, updateAdministratorOrganizationToCommand(resource, req.UserId, req.Roles))
		if err != nil {
			return nil, err
		}
		if !member.EventDate.IsZero() {
			changeDate = timestamppb.New(member.CreationDate)
		}
	case *internal_permission.ResourceType_ProjectId:
		member, err := s.command.ChangeProjectMember(ctx, updateAdministratorProjectToCommand(resource, req.UserId, req.Roles))
		if err != nil {
			return nil, err
		}
		if !member.EventDate.IsZero() {
			changeDate = timestamppb.New(member.CreationDate)
		}
	case *internal_permission.ResourceType_ProjectGrant_:
		member, err := s.command.ChangeProjectGrantMember(ctx, updateAdministratorProjectGrantToCommand(resource, req.UserId, req.Roles))
		if err != nil {
			return nil, err
		}
		if !member.EventDate.IsZero() {
			changeDate = timestamppb.New(member.CreationDate)
		}
	}

	return &internal_permission.UpdateAdministratorResponse{
		ChangeDate: changeDate,
	}, nil
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
