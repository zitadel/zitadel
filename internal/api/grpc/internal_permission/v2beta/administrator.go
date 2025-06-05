package internal_permission

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/command"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	internal_permission "github.com/zitadel/zitadel/pkg/grpc/internal_permission/v2beta"
)

func (s *Server) CreateAdministrator(ctx context.Context, req *internal_permission.CreateAdministratorRequest) (*internal_permission.CreateAdministratorResponse, error) {
	var creationDate *timestamppb.Timestamp
	var id string

	switch resource := req.GetResource().(type) {
	case *internal_permission.CreateAdministratorRequest_Instance:
		if resource.Instance {
			member, err := s.command.AddInstanceMember(ctx, authz.GetInstance(ctx).InstanceID(), req.UserId, req.Roles...)
			if err != nil {
				return nil, err
			}
			if !member.CreationDate.IsZero() {
				creationDate = timestamppb.New(member.CreationDate)
			}
		}
	case *internal_permission.CreateAdministratorRequest_OrganizationId:
		member, err := s.command.AddOrgMember(ctx, resource.OrganizationId, req.UserId, req.Roles...)
		if err != nil {
			return nil, err
		}
		if !member.CreationDate.IsZero() {
			creationDate = timestamppb.New(member.CreationDate)
		}
	case *internal_permission.CreateAdministratorRequest_ProjectId:
		member, err := s.command.AddProjectMember(ctx, createAdministratorProjectToCommand(resource, req.UserId, req.Roles))
		if err != nil {
			return nil, err
		}
		if !member.CreationDate.IsZero() {
			creationDate = timestamppb.New(member.CreationDate)
		}
	}

	return &internal_permission.CreateAdministratorResponse{
		CreationDate: creationDate,
		Id:           id,
	}, nil
}
func createAdministratorProjectToCommand(req *internal_permission.CreateAdministratorRequest_ProjectId, userID string, roles []string) *command.AddProjectMember {
	return &command.AddProjectMember{
		ProjectID: req.ProjectId,
		UserID:    userID,
		Roles:     roles,
	}
}

func (s *Server) UpdateIAMMember(ctx context.Context, req *admin_pb.UpdateIAMMemberRequest) (*admin_pb.UpdateIAMMemberResponse, error) {
	member, err := s.command.ChangeInstanceMember(ctx, UpdateIAMMemberToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateIAMMemberResponse{
		Details: object.ChangeToDetailsPb(
			member.Sequence,
			member.ChangeDate,
			member.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveIAMMember(ctx context.Context, req *admin_pb.RemoveIAMMemberRequest) (*admin_pb.RemoveIAMMemberResponse, error) {
	objectDetails, err := s.command.RemoveInstanceMember(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &admin_pb.RemoveIAMMemberResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}
