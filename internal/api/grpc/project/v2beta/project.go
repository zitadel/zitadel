package project

import (
	"context"

	"connectrpc.com/connect"
	"github.com/muhlemmer/gu"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	project_pb "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

func (s *Server) CreateProject(ctx context.Context, req *connect.Request[project_pb.CreateProjectRequest]) (*connect.Response[project_pb.CreateProjectResponse], error) {
	add := projectCreateToCommand(req.Msg)
	project, err := s.command.AddProject(ctx, add)
	if err != nil {
		return nil, err
	}
	var creationDate *timestamppb.Timestamp
	if !project.EventDate.IsZero() {
		creationDate = timestamppb.New(project.EventDate)
	}
	return connect.NewResponse(&project_pb.CreateProjectResponse{
		Id:           add.AggregateID,
		CreationDate: creationDate,
	}), nil
}

func projectCreateToCommand(req *project_pb.CreateProjectRequest) *command.AddProject {
	var aggregateID string
	if req.Id != nil {
		aggregateID = *req.Id
	}
	return &command.AddProject{
		ObjectRoot: models.ObjectRoot{
			ResourceOwner: req.OrganizationId,
			AggregateID:   aggregateID,
		},
		Name:                   req.Name,
		ProjectRoleAssertion:   req.ProjectRoleAssertion,
		ProjectRoleCheck:       req.AuthorizationRequired,
		HasProjectCheck:        req.ProjectAccessRequired,
		PrivateLabelingSetting: privateLabelingSettingToDomain(req.PrivateLabelingSetting),
	}
}

func privateLabelingSettingToDomain(setting project_pb.PrivateLabelingSetting) domain.PrivateLabelingSetting {
	switch setting {
	case project_pb.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_ALLOW_LOGIN_USER_RESOURCE_OWNER_POLICY:
		return domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy
	case project_pb.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY:
		return domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy
	case project_pb.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_UNSPECIFIED:
		return domain.PrivateLabelingSettingUnspecified
	default:
		return domain.PrivateLabelingSettingUnspecified
	}
}

func (s *Server) UpdateProject(ctx context.Context, req *connect.Request[project_pb.UpdateProjectRequest]) (*connect.Response[project_pb.UpdateProjectResponse], error) {
	project, err := s.command.ChangeProject(ctx, projectUpdateToCommand(req.Msg))
	if err != nil {
		return nil, err
	}
	var changeDate *timestamppb.Timestamp
	if !project.EventDate.IsZero() {
		changeDate = timestamppb.New(project.EventDate)
	}
	return connect.NewResponse(&project_pb.UpdateProjectResponse{
		ChangeDate: changeDate,
	}), nil
}

func projectUpdateToCommand(req *project_pb.UpdateProjectRequest) *command.ChangeProject {
	var labeling *domain.PrivateLabelingSetting
	if req.PrivateLabelingSetting != nil {
		labeling = gu.Ptr(privateLabelingSettingToDomain(*req.PrivateLabelingSetting))
	}
	return &command.ChangeProject{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.Id,
		},
		Name:                   req.Name,
		ProjectRoleAssertion:   req.ProjectRoleAssertion,
		ProjectRoleCheck:       req.ProjectRoleCheck,
		HasProjectCheck:        req.HasProjectCheck,
		PrivateLabelingSetting: labeling,
	}
}

func (s *Server) DeleteProject(ctx context.Context, req *connect.Request[project_pb.DeleteProjectRequest]) (*connect.Response[project_pb.DeleteProjectResponse], error) {
	userGrantIDs, err := s.userGrantsFromProject(ctx, req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	deletedAt, err := s.command.DeleteProject(ctx, req.Msg.GetId(), "", userGrantIDs...)
	if err != nil {
		return nil, err
	}
	var deletionDate *timestamppb.Timestamp
	if !deletedAt.IsZero() {
		deletionDate = timestamppb.New(deletedAt)
	}
	return connect.NewResponse(&project_pb.DeleteProjectResponse{
		DeletionDate: deletionDate,
	}), nil
}

func (s *Server) userGrantsFromProject(ctx context.Context, projectID string) ([]string, error) {
	projectQuery, err := query.NewUserGrantProjectIDSearchQuery(projectID)
	if err != nil {
		return nil, err
	}
	userGrants, err := s.query.UserGrants(ctx, &query.UserGrantsQueries{
		Queries: []query.SearchQuery{projectQuery},
	}, false, nil)
	if err != nil {
		return nil, err
	}
	return userGrantsToIDs(userGrants.UserGrants), nil
}

func (s *Server) DeactivateProject(ctx context.Context, req *connect.Request[project_pb.DeactivateProjectRequest]) (*connect.Response[project_pb.DeactivateProjectResponse], error) {
	details, err := s.command.DeactivateProject(ctx, req.Msg.GetId(), "")
	if err != nil {
		return nil, err
	}
	var changeDate *timestamppb.Timestamp
	if !details.EventDate.IsZero() {
		changeDate = timestamppb.New(details.EventDate)
	}
	return connect.NewResponse(&project_pb.DeactivateProjectResponse{
		ChangeDate: changeDate,
	}), nil
}

func (s *Server) ActivateProject(ctx context.Context, req *connect.Request[project_pb.ActivateProjectRequest]) (*connect.Response[project_pb.ActivateProjectResponse], error) {
	details, err := s.command.ReactivateProject(ctx, req.Msg.GetId(), "")
	if err != nil {
		return nil, err
	}
	var changeDate *timestamppb.Timestamp
	if !details.EventDate.IsZero() {
		changeDate = timestamppb.New(details.EventDate)
	}
	return connect.NewResponse(&project_pb.ActivateProjectResponse{
		ChangeDate: changeDate,
	}), nil
}

func userGrantsToIDs(userGrants []*query.UserGrant) []string {
	converted := make([]string, len(userGrants))
	for i, grant := range userGrants {
		converted[i] = grant.ID
	}
	return converted
}
