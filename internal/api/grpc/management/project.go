package management

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	change_grpc "github.com/zitadel/zitadel/internal/api/grpc/change"
	member_grpc "github.com/zitadel/zitadel/internal/api/grpc/member"
	object_grpc "github.com/zitadel/zitadel/internal/api/grpc/object"
	project_grpc "github.com/zitadel/zitadel/internal/api/grpc/project"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/project"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func (s *Server) GetProjectByID(ctx context.Context, req *mgmt_pb.GetProjectByIDRequest) (*mgmt_pb.GetProjectByIDResponse, error) {
	project, err := s.query.ProjectByID(ctx, true, req.Id, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetProjectByIDResponse{
		Project: project_grpc.ProjectViewToPb(project),
	}, nil
}

func (s *Server) GetGrantedProjectByID(ctx context.Context, req *mgmt_pb.GetGrantedProjectByIDRequest) (*mgmt_pb.GetGrantedProjectByIDResponse, error) {
	grant, err := s.query.ProjectGrantByID(ctx, true, req.GrantId, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetGrantedProjectByIDResponse{
		GrantedProject: project_grpc.GrantedProjectViewToPb(grant),
	}, nil
}

func (s *Server) ListProjects(ctx context.Context, req *mgmt_pb.ListProjectsRequest) (*mgmt_pb.ListProjectsResponse, error) {
	queries, err := listProjectRequestToModel(req)
	if err != nil {
		return nil, err
	}
	err = queries.AppendMyResourceOwnerQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	err = queries.AppendPermissionQueries(authz.GetRequestPermissionsFromCtx(ctx))
	if err != nil {
		return nil, err
	}
	projects, err := s.query.SearchProjects(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListProjectsResponse{
		Result:  project_grpc.ProjectViewsToPb(projects.Projects),
		Details: object_grpc.ToListDetails(projects.Count, projects.Sequence, projects.Timestamp),
	}, nil
}

func (s *Server) ListProjectGrantChanges(ctx context.Context, req *mgmt_pb.ListProjectGrantChangesRequest) (*mgmt_pb.ListProjectGrantChangesResponse, error) {
	var (
		limit    uint64
		sequence uint64
		asc      bool
	)
	if req.Query != nil {
		limit = uint64(req.Query.Limit)
		sequence = req.Query.Sequence
		asc = req.Query.Asc
	}

	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AllowTimeTravel().
		Limit(limit).
		OrderDesc().
		ResourceOwner(authz.GetCtxData(ctx).OrgID).
		AddQuery().
		SequenceGreater(sequence).
		AggregateTypes(project.AggregateType).
		AggregateIDs(req.ProjectId).
		EventData(map[string]interface{}{
			"grantId": req.GrantId,
		}).
		Builder()
	if asc {
		query.OrderAsc()
	}

	changes, err := s.query.SearchEvents(ctx, query, s.auditLogRetention)
	if err != nil {
		return nil, err
	}

	return &mgmt_pb.ListProjectGrantChangesResponse{
		Result: change_grpc.EventsToChangesPb(changes, s.assetAPIPrefix(ctx)),
	}, nil
}

func (s *Server) ListGrantedProjects(ctx context.Context, req *mgmt_pb.ListGrantedProjectsRequest) (*mgmt_pb.ListGrantedProjectsResponse, error) {
	queries, err := listGrantedProjectsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	err = queries.AppendGrantedOrgQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	err = queries.AppendPermissionQueries(authz.GetRequestPermissionsFromCtx(ctx))
	if err != nil {
		return nil, err
	}
	projects, err := s.query.SearchProjectGrants(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListGrantedProjectsResponse{
		Result:  project_grpc.GrantedProjectViewsToPb(projects.ProjectGrants),
		Details: object_grpc.ToListDetails(projects.Count, projects.Sequence, projects.Timestamp),
	}, nil
}

func (s *Server) ListGrantedProjectRoles(ctx context.Context, req *mgmt_pb.ListGrantedProjectRolesRequest) (*mgmt_pb.ListGrantedProjectRolesResponse, error) {
	queries, err := listGrantedProjectRolesRequestToModel(req)
	if err != nil {
		return nil, err
	}
	err = queries.AppendProjectIDQuery(req.ProjectId)
	if err != nil {
		return nil, err
	}
	roles, err := s.query.SearchGrantedProjectRoles(ctx, req.GrantId, authz.GetCtxData(ctx).OrgID, queries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListGrantedProjectRolesResponse{
		Result:  project_grpc.RoleViewsToPb(roles.ProjectRoles),
		Details: object_grpc.ToListDetails(roles.Count, roles.Sequence, roles.Timestamp),
	}, nil
}

func (s *Server) ListProjectChanges(ctx context.Context, req *mgmt_pb.ListProjectChangesRequest) (*mgmt_pb.ListProjectChangesResponse, error) {
	var (
		limit    uint64
		sequence uint64
		asc      bool
	)
	if req.Query != nil {
		limit = uint64(req.Query.Limit)
		sequence = req.Query.Sequence
		asc = req.Query.Asc
	}

	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AllowTimeTravel().
		Limit(limit).
		OrderDesc().
		ResourceOwner(authz.GetCtxData(ctx).OrgID).
		AddQuery().
		SequenceGreater(sequence).
		AggregateTypes(project.AggregateType).
		AggregateIDs(req.ProjectId).
		Builder()
	if asc {
		query.OrderAsc()
	}

	changes, err := s.query.SearchEvents(ctx, query, s.auditLogRetention)
	if err != nil {
		return nil, err
	}

	return &mgmt_pb.ListProjectChangesResponse{
		Result: change_grpc.EventsToChangesPb(changes, s.assetAPIPrefix(ctx)),
	}, nil
}

func (s *Server) AddProject(ctx context.Context, req *mgmt_pb.AddProjectRequest) (*mgmt_pb.AddProjectResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	project, err := s.command.AddProject(ctx, ProjectCreateToDomain(req), ctxData.OrgID, ctxData.UserID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddProjectResponse{
		Id:      project.AggregateID,
		Details: object_grpc.AddToDetailsPb(project.Sequence, project.ChangeDate, project.ResourceOwner),
	}, nil
}

func (s *Server) UpdateProject(ctx context.Context, req *mgmt_pb.UpdateProjectRequest) (*mgmt_pb.UpdateProjectResponse, error) {
	project, err := s.command.ChangeProject(ctx, ProjectUpdateToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateProjectResponse{
		Details: object_grpc.ChangeToDetailsPb(
			project.Sequence,
			project.ChangeDate,
			project.ResourceOwner,
		),
	}, nil
}

func (s *Server) DeactivateProject(ctx context.Context, req *mgmt_pb.DeactivateProjectRequest) (*mgmt_pb.DeactivateProjectResponse, error) {
	details, err := s.command.DeactivateProject(ctx, req.Id, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.DeactivateProjectResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) ReactivateProject(ctx context.Context, req *mgmt_pb.ReactivateProjectRequest) (*mgmt_pb.ReactivateProjectResponse, error) {
	details, err := s.command.ReactivateProject(ctx, req.Id, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ReactivateProjectResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) RemoveProject(ctx context.Context, req *mgmt_pb.RemoveProjectRequest) (*mgmt_pb.RemoveProjectResponse, error) {
	projectQuery, err := query.NewUserGrantProjectIDSearchQuery(req.Id)
	if err != nil {
		return nil, err
	}
	grants, err := s.query.UserGrants(ctx, &query.UserGrantsQueries{
		Queries: []query.SearchQuery{projectQuery},
	}, true, false)
	if err != nil {
		return nil, err
	}
	details, err := s.command.RemoveProject(ctx, req.Id, authz.GetCtxData(ctx).OrgID, userGrantsToIDs(grants.UserGrants)...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveProjectResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) ListProjectRoles(ctx context.Context, req *mgmt_pb.ListProjectRolesRequest) (*mgmt_pb.ListProjectRolesResponse, error) {
	queries, err := listProjectRolesRequestToModel(req)
	if err != nil {
		return nil, err
	}
	err = queries.AppendMyResourceOwnerQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	err = queries.AppendProjectIDQuery(req.ProjectId)
	if err != nil {
		return nil, err
	}
	roles, err := s.query.SearchProjectRoles(ctx, true, queries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListProjectRolesResponse{
		Result:  project_grpc.RoleViewsToPb(roles.ProjectRoles),
		Details: object_grpc.ToListDetails(roles.Count, roles.Sequence, roles.Timestamp),
	}, nil
}

func (s *Server) AddProjectRole(ctx context.Context, req *mgmt_pb.AddProjectRoleRequest) (*mgmt_pb.AddProjectRoleResponse, error) {
	role, err := s.command.AddProjectRole(ctx, AddProjectRoleRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddProjectRoleResponse{
		Details: object_grpc.AddToDetailsPb(
			role.Sequence,
			role.ChangeDate,
			role.ResourceOwner,
		),
	}, nil
}

func (s *Server) BulkAddProjectRoles(ctx context.Context, req *mgmt_pb.BulkAddProjectRolesRequest) (*mgmt_pb.BulkAddProjectRolesResponse, error) {
	details, err := s.command.BulkAddProjectRole(ctx, req.ProjectId, authz.GetCtxData(ctx).OrgID, BulkAddProjectRolesRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.BulkAddProjectRolesResponse{
		Details: object_grpc.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateProjectRole(ctx context.Context, req *mgmt_pb.UpdateProjectRoleRequest) (*mgmt_pb.UpdateProjectRoleResponse, error) {
	role, err := s.command.ChangeProjectRole(ctx, UpdateProjectRoleRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateProjectRoleResponse{
		Details: object_grpc.ChangeToDetailsPb(
			role.Sequence,
			role.ChangeDate,
			role.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveProjectRole(ctx context.Context, req *mgmt_pb.RemoveProjectRoleRequest) (*mgmt_pb.RemoveProjectRoleResponse, error) {
	projectQuery, err := query.NewUserGrantProjectIDSearchQuery(req.ProjectId)
	if err != nil {
		return nil, err
	}
	rolesQuery, err := query.NewUserGrantRoleQuery(req.RoleKey)
	if err != nil {
		return nil, err
	}
	userGrants, err := s.query.UserGrants(ctx, &query.UserGrantsQueries{
		Queries: []query.SearchQuery{projectQuery, rolesQuery},
	}, false, false)

	if err != nil {
		return nil, err
	}
	projectGrants, err := s.query.SearchProjectGrantsByProjectIDAndRoleKey(ctx, req.ProjectId, req.RoleKey, false)
	if err != nil {
		return nil, err
	}
	details, err := s.command.RemoveProjectRole(ctx, req.ProjectId, req.RoleKey, authz.GetCtxData(ctx).OrgID, ProjectGrantsToIDs(projectGrants), userGrantsToIDs(userGrants.UserGrants)...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveProjectRoleResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) ListProjectMemberRoles(ctx context.Context, _ *mgmt_pb.ListProjectMemberRolesRequest) (*mgmt_pb.ListProjectMemberRolesResponse, error) {
	roles, err := s.query.GetProjectMemberRoles(ctx)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListProjectMemberRolesResponse{Result: roles}, nil //TODO: details
}

func (s *Server) ListProjectMembers(ctx context.Context, req *mgmt_pb.ListProjectMembersRequest) (*mgmt_pb.ListProjectMembersResponse, error) {
	queries, err := ListProjectMembersRequestToModel(ctx, req)
	if err != nil {
		return nil, err
	}
	members, err := s.query.ProjectMembers(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListProjectMembersResponse{
		Result:  member_grpc.MembersToPb(s.assetAPIPrefix(ctx), members.Members),
		Details: object_grpc.ToListDetails(members.Count, members.Sequence, members.Timestamp),
	}, nil
}

func (s *Server) AddProjectMember(ctx context.Context, req *mgmt_pb.AddProjectMemberRequest) (*mgmt_pb.AddProjectMemberResponse, error) {
	member, err := s.command.AddProjectMember(ctx, AddProjectMemberRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddProjectMemberResponse{
		Details: object_grpc.AddToDetailsPb(member.Sequence, member.ChangeDate, member.ResourceOwner),
	}, nil
}

func (s *Server) UpdateProjectMember(ctx context.Context, req *mgmt_pb.UpdateProjectMemberRequest) (*mgmt_pb.UpdateProjectMemberResponse, error) {
	member, err := s.command.ChangeProjectMember(ctx, UpdateProjectMemberRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateProjectMemberResponse{
		Details: object_grpc.ChangeToDetailsPb(
			member.Sequence,
			member.ChangeDate,
			member.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveProjectMember(ctx context.Context, req *mgmt_pb.RemoveProjectMemberRequest) (*mgmt_pb.RemoveProjectMemberResponse, error) {
	details, err := s.command.RemoveProjectMember(ctx, req.ProjectId, req.UserId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveProjectMemberResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}
