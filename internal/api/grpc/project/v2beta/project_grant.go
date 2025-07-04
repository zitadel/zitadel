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

func (s *Server) CreateProjectGrant(ctx context.Context, req *connect.Request[project_pb.CreateProjectGrantRequest]) (*connect.Response[project_pb.CreateProjectGrantResponse], error) {
	add := projectGrantCreateToCommand(req.Msg)
	project, err := s.command.AddProjectGrant(ctx, add)
	if err != nil {
		return nil, err
	}
	var creationDate *timestamppb.Timestamp
	if !project.EventDate.IsZero() {
		creationDate = timestamppb.New(project.EventDate)
	}
	return connect.NewResponse(&project_pb.CreateProjectGrantResponse{
		CreationDate: creationDate,
	}), nil
}

func projectGrantCreateToCommand(req *project_pb.CreateProjectGrantRequest) *command.AddProjectGrant {
	return &command.AddProjectGrant{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.ProjectId,
		},
		GrantID:      req.GrantedOrganizationId,
		GrantedOrgID: req.GrantedOrganizationId,
		RoleKeys:     req.RoleKeys,
	}
}

func (s *Server) UpdateProjectGrant(ctx context.Context, req *connect.Request[project_pb.UpdateProjectGrantRequest]) (*connect.Response[project_pb.UpdateProjectGrantResponse], error) {
	project, err := s.command.ChangeProjectGrant(ctx, projectGrantUpdateToCommand(req.Msg))
	if err != nil {
		return nil, err
	}
	var changeDate *timestamppb.Timestamp
	if !project.EventDate.IsZero() {
		changeDate = timestamppb.New(project.EventDate)
	}
	return connect.NewResponse(&project_pb.UpdateProjectGrantResponse{
		ChangeDate: changeDate,
	}), nil
}

func projectGrantUpdateToCommand(req *project_pb.UpdateProjectGrantRequest) *command.ChangeProjectGrant {
	return &command.ChangeProjectGrant{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.ProjectId,
		},
		GrantedOrgID: req.GrantedOrganizationId,
		RoleKeys:     req.RoleKeys,
	}
}

func (s *Server) DeactivateProjectGrant(ctx context.Context, req *connect.Request[project_pb.DeactivateProjectGrantRequest]) (*connect.Response[project_pb.DeactivateProjectGrantResponse], error) {
	details, err := s.command.DeactivateProjectGrant(ctx, req.Msg.GetProjectId(), "", req.Msg.GetGrantedOrganizationId(), "")
	if err != nil {
		return nil, err
	}
	var changeDate *timestamppb.Timestamp
	if !details.EventDate.IsZero() {
		changeDate = timestamppb.New(details.EventDate)
	}
	return connect.NewResponse(&project_pb.DeactivateProjectGrantResponse{
		ChangeDate: changeDate,
	}), nil
}

func (s *Server) ActivateProjectGrant(ctx context.Context, req *connect.Request[project_pb.ActivateProjectGrantRequest]) (*connect.Response[project_pb.ActivateProjectGrantResponse], error) {
	details, err := s.command.ReactivateProjectGrant(ctx, req.Msg.GetProjectId(), "", req.Msg.GetGrantedOrganizationId(), "")
	if err != nil {
		return nil, err
	}
	var changeDate *timestamppb.Timestamp
	if !details.EventDate.IsZero() {
		changeDate = timestamppb.New(details.EventDate)
	}
	return connect.NewResponse(&project_pb.ActivateProjectGrantResponse{
		ChangeDate: changeDate,
	}), nil
}

func (s *Server) DeleteProjectGrant(ctx context.Context, req *connect.Request[project_pb.DeleteProjectGrantRequest]) (*connect.Response[project_pb.DeleteProjectGrantResponse], error) {
	userGrantIDs, err := s.userGrantsFromProjectGrant(ctx, req.Msg.GetProjectId(), req.Msg.GetGrantedOrganizationId())
	if err != nil {
		return nil, err
	}
	details, err := s.command.DeleteProjectGrant(ctx, req.Msg.GetProjectId(), "", req.Msg.GetGrantedOrganizationId(), "", userGrantIDs...)
	if err != nil {
		return nil, err
	}
	var deletionDate *timestamppb.Timestamp
	if !details.EventDate.IsZero() {
		deletionDate = timestamppb.New(details.EventDate)
	}
	return connect.NewResponse(&project_pb.DeleteProjectGrantResponse{
		DeletionDate: deletionDate,
	}), nil
}

func (s *Server) userGrantsFromProjectGrant(ctx context.Context, projectID, grantedOrganizationID string) ([]string, error) {
	projectQuery, err := query.NewUserGrantProjectIDSearchQuery(projectID)
	if err != nil {
		return nil, err
	}
	grantQuery, err := query.NewUserGrantGrantIDSearchQuery(grantedOrganizationID)
	if err != nil {
		return nil, err
	}
	userGrants, err := s.query.UserGrants(ctx, &query.UserGrantsQueries{
		Queries: []query.SearchQuery{projectQuery, grantQuery},
	}, false, nil)
	if err != nil {
		return nil, err
	}
	return userGrantsToIDs(userGrants.UserGrants), nil
}
