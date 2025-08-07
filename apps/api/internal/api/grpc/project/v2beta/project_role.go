package project

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	project_pb "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

func (s *Server) AddProjectRole(ctx context.Context, req *connect.Request[project_pb.AddProjectRoleRequest]) (*connect.Response[project_pb.AddProjectRoleResponse], error) {
	role, err := s.command.AddProjectRole(ctx, addProjectRoleRequestToCommand(req.Msg))
	if err != nil {
		return nil, err
	}
	var creationDate *timestamppb.Timestamp
	if !role.EventDate.IsZero() {
		creationDate = timestamppb.New(role.EventDate)
	}
	return connect.NewResponse(&project_pb.AddProjectRoleResponse{
		CreationDate: creationDate,
	}), nil
}

func addProjectRoleRequestToCommand(req *project_pb.AddProjectRoleRequest) *command.AddProjectRole {
	group := ""
	if req.Group != nil {
		group = *req.Group
	}

	return &command.AddProjectRole{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.ProjectId,
		},
		Key:         req.RoleKey,
		DisplayName: req.DisplayName,
		Group:       group,
	}
}

func (s *Server) UpdateProjectRole(ctx context.Context, req *connect.Request[project_pb.UpdateProjectRoleRequest]) (*connect.Response[project_pb.UpdateProjectRoleResponse], error) {
	role, err := s.command.ChangeProjectRole(ctx, updateProjectRoleRequestToCommand(req.Msg))
	if err != nil {
		return nil, err
	}
	var changeDate *timestamppb.Timestamp
	if !role.EventDate.IsZero() {
		changeDate = timestamppb.New(role.EventDate)
	}
	return connect.NewResponse(&project_pb.UpdateProjectRoleResponse{
		ChangeDate: changeDate,
	}), nil
}

func updateProjectRoleRequestToCommand(req *project_pb.UpdateProjectRoleRequest) *command.ChangeProjectRole {
	displayName := ""
	if req.DisplayName != nil {
		displayName = *req.DisplayName
	}
	group := ""
	if req.Group != nil {
		group = *req.Group
	}

	return &command.ChangeProjectRole{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.ProjectId,
		},
		Key:         req.RoleKey,
		DisplayName: displayName,
		Group:       group,
	}
}

func (s *Server) RemoveProjectRole(ctx context.Context, req *connect.Request[project_pb.RemoveProjectRoleRequest]) (*connect.Response[project_pb.RemoveProjectRoleResponse], error) {
	userGrantIDs, err := s.userGrantsFromProjectAndRole(ctx, req.Msg.GetProjectId(), req.Msg.GetRoleKey())
	if err != nil {
		return nil, err
	}
	projectGrantIDs, err := s.projectGrantsFromProjectAndRole(ctx, req.Msg.GetProjectId(), req.Msg.GetRoleKey())
	if err != nil {
		return nil, err
	}
	details, err := s.command.RemoveProjectRole(ctx, req.Msg.GetProjectId(), req.Msg.GetRoleKey(), "", projectGrantIDs, userGrantIDs...)
	if err != nil {
		return nil, err
	}
	var deletionDate *timestamppb.Timestamp
	if !details.EventDate.IsZero() {
		deletionDate = timestamppb.New(details.EventDate)
	}
	return connect.NewResponse(&project_pb.RemoveProjectRoleResponse{
		RemovalDate: deletionDate,
	}), nil
}

func (s *Server) userGrantsFromProjectAndRole(ctx context.Context, projectID, roleKey string) ([]string, error) {
	projectQuery, err := query.NewUserGrantProjectIDSearchQuery(projectID)
	if err != nil {
		return nil, err
	}
	rolesQuery, err := query.NewUserGrantRoleQuery(roleKey)
	if err != nil {
		return nil, err
	}
	userGrants, err := s.query.UserGrants(ctx, &query.UserGrantsQueries{
		Queries: []query.SearchQuery{projectQuery, rolesQuery},
	}, false, nil)
	if err != nil {
		return nil, err
	}
	return userGrantsToIDs(userGrants.UserGrants), nil
}

func (s *Server) projectGrantsFromProjectAndRole(ctx context.Context, projectID, roleKey string) ([]string, error) {
	projectGrants, err := s.query.SearchProjectGrantsByProjectIDAndRoleKey(ctx, projectID, roleKey)
	if err != nil {
		return nil, err
	}
	return projectGrantsToIDs(projectGrants), nil
}

func projectGrantsToIDs(projectGrants *query.ProjectGrants) []string {
	converted := make([]string, len(projectGrants.ProjectGrants))
	for i, grant := range projectGrants.ProjectGrants {
		converted[i] = grant.GrantID
	}
	return converted
}
