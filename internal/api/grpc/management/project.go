package management

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	change_grpc "github.com/caos/zitadel/internal/api/grpc/change"
	member_grpc "github.com/caos/zitadel/internal/api/grpc/member"
	object_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	project_grpc "github.com/caos/zitadel/internal/api/grpc/project"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetProjectByID(ctx context.Context, req *mgmt_pb.GetProjectByIDRequest) (*mgmt_pb.GetProjectByIDResponse, error) {
	project, err := s.project.ProjectByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetProjectByIDResponse{
		Project: project_grpc.ProjectToPb(project),
	}, nil
}

func (s *Server) GetGrantedProjectByID(ctx context.Context, req *mgmt_pb.GetGrantedProjectByIDRequest) (*mgmt_pb.GetGrantedProjectByIDResponse, error) {
	project, err := s.project.ProjectGrantViewByID(ctx, req.GrantId)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetGrantedProjectByIDResponse{
		GrantedProject: project_grpc.GrantedProjectToPb(project),
	}, nil
}

func (s *Server) ListProjects(ctx context.Context, req *mgmt_pb.ListProjectsRequest) (*mgmt_pb.ListProjectsResponse, error) {
	queries, err := ListProjectsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	queries.AppendMyResourceOwnerQuery(authz.GetCtxData(ctx).OrgID)
	projects, err := s.project.SearchProjects(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListProjectsResponse{
		Result: project_grpc.ProjectsToPb(projects.Result),
		Details: object_grpc.ToListDetails(
			projects.TotalResult,
			projects.Sequence,
			projects.Timestamp,
		),
	}, nil
}

func (s *Server) ListGrantedProjects(ctx context.Context, req *mgmt_pb.ListGrantedProjectsRequest) (*mgmt_pb.ListGrantedProjectsResponse, error) {
	queries, err := ListGrantedProjectsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	queries.AppendMyResourceOwnerQuery(authz.GetCtxData(ctx).OrgID)
	projects, err := s.project.SearchGrantedProjects(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListGrantedProjectsResponse{
		Result: project_grpc.GrantedProjectsToPb(projects.Result),
		Details: object_grpc.ToListDetails(
			projects.TotalResult,
			projects.Sequence,
			projects.Timestamp,
		),
	}, nil
}

func (s *Server) ListProjectChanges(ctx context.Context, req *mgmt_pb.ListProjectChangesRequest) (*mgmt_pb.ListProjectChangesResponse, error) {
	offset, limit, asc := object_grpc.ListQueryToModel(req.Query)
	res, err := s.project.ProjectChanges(ctx, req.ProjectId, offset, limit, asc)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListProjectChangesResponse{
		Result: change_grpc.ProjectChangesToPb(res.Changes),
	}, nil
}

func (s *Server) AddProject(ctx context.Context, req *mgmt_pb.AddProjectRequest) (*mgmt_pb.AddProjectResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	project, err := s.command.AddProject(ctx, ProjectCreateToDomain(req), ctxData.ResourceOwner, ctxData.UserID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddProjectResponse{
		Id: project.AggregateID,
		Details: object_grpc.ToDetailsPb(
			project.Sequence,
			project.ChangeDate,
			project.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateProject(ctx context.Context, req *mgmt_pb.UpdateProjectRequest) (*mgmt_pb.UpdateProjectResponse, error) {
	project, err := s.command.ChangeProject(ctx, ProjectUpdateToDomain(req), authz.GetCtxData(ctx).ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateProjectResponse{
		Details: object_grpc.ToDetailsPb(
			project.Sequence,
			project.ChangeDate,
			project.ResourceOwner,
		),
	}, nil
}

func (s *Server) DeactivateProject(ctx context.Context, req *mgmt_pb.DeactivateProjectRequest) (*mgmt_pb.DeactivateProjectResponse, error) {
	details, err := s.command.DeactivateProject(ctx, req.Id, authz.GetCtxData(ctx).ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.DeactivateProjectResponse{
		Details: object_grpc.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) ReactivateProject(ctx context.Context, req *mgmt_pb.ReactivateProjectRequest) (*mgmt_pb.ReactivateProjectResponse, error) {
	details, err := s.command.ReactivateProject(ctx, req.Id, authz.GetCtxData(ctx).ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ReactivateProjectResponse{
		Details: object_grpc.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) RemoveProject(ctx context.Context, req *mgmt_pb.RemoveProjectRequest) (*mgmt_pb.RemoveProjectResponse, error) {
	grants, err := s.usergrant.UserGrantsByProjectID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	details, err := s.command.RemoveProject(ctx, req.Id, authz.GetCtxData(ctx).OrgID, userGrantsToIDs(grants)...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveProjectResponse{
		Details: object_grpc.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) ListProjectRoles(ctx context.Context, req *mgmt_pb.ListProjectRolesRequest) (*mgmt_pb.ListProjectRolesResponse, error) {
	queries, err := ListProjectRolesRequestToModel(req)
	if err != nil {
		return nil, err
	}
	queries.AppendMyOrgQuery(authz.GetCtxData(ctx).OrgID)
	roles, err := s.project.SearchProjectRoles(ctx, req.ProjectId, queries)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListProjectRolesResponse{
		Result: project_grpc.RolesToPb(roles.Result),
		Details: object_grpc.ToListDetails(
			roles.TotalResult,
			roles.Sequence,
			roles.Timestamp,
		),
	}, nil
}

func (s *Server) AddProjectRole(ctx context.Context, req *mgmt_pb.AddProjectRoleRequest) (*mgmt_pb.AddProjectRoleResponse, error) {
	role, err := s.command.AddProjectRole(ctx, AddProjectRoleRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddProjectRoleResponse{
		Details: object_grpc.ToDetailsPb(
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
		Details: object_grpc.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) UpdateProjectRole(ctx context.Context, req *mgmt_pb.UpdateProjectRoleRequest) (*mgmt_pb.UpdateProjectRoleResponse, error) {
	role, err := s.command.ChangeProjectRole(ctx, UpdateProjectRoleRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateProjectRoleResponse{
		Details: object_grpc.ToDetailsPb(
			role.Sequence,
			role.ChangeDate,
			role.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveProjectRole(ctx context.Context, req *mgmt_pb.RemoveProjectRoleRequest) (*mgmt_pb.RemoveProjectRoleResponse, error) {
	userGrants, err := s.usergrant.UserGrantsByProjectIDAndRoleKey(ctx, req.ProjectId, req.RoleKey)
	if err != nil {
		return nil, err
	}
	projectGrants, err := s.project.ProjectGrantsByProjectIDAndRoleKey(ctx, req.ProjectId, req.RoleKey)
	if err != nil {
		return nil, err
	}
	details, err := s.command.RemoveProjectRole(ctx, req.ProjectId, req.RoleKey, authz.GetCtxData(ctx).OrgID, ProjectGrantsToIDs(projectGrants), userGrantsToIDs(userGrants)...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveProjectRoleResponse{
		Details: object_grpc.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) ListProjectMemberRoles(ctx context.Context, req *mgmt_pb.ListProjectMemberRolesRequest) (*mgmt_pb.ListProjectMemberRolesResponse, error) {
	roles, err := s.project.GetProjectMemberRoles(ctx)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListProjectMemberRolesResponse{Result: roles}, nil //TODO: details
}

func (s *Server) ListProjectMembers(ctx context.Context, req *mgmt_pb.ListProjectMembersRequest) (*mgmt_pb.ListProjectMembersResponse, error) {
	queries, err := ListProjectMembersRequestToModel(req)
	if err != nil {
		return nil, err
	}
	queries.AppendProjectQuery(req.ProjectId)
	members, err := s.project.SearchProjectMembers(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListProjectMembersResponse{
		Result: member_grpc.ProjectMembersToPb(members.Result),
		Details: object_grpc.ToListDetails(
			members.TotalResult,
			members.Sequence,
			members.Timestamp,
		),
	}, nil
}

func (s *Server) AddProjectMember(ctx context.Context, req *mgmt_pb.AddProjectMemberRequest) (*mgmt_pb.AddProjectMemberResponse, error) {
	member, err := s.command.AddProjectMember(ctx, AddProjectMemberRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddProjectMemberResponse{
		Details: object_grpc.ToDetailsPb(
			member.Sequence,
			member.ChangeDate,
			member.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateProjectMember(ctx context.Context, req *mgmt_pb.UpdateProjectMemberRequest) (*mgmt_pb.UpdateProjectMemberResponse, error) {
	member, err := s.command.ChangeProjectMember(ctx, UpdateProjectMemberRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateProjectMemberResponse{
		Details: object_grpc.ToDetailsPb(
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
		Details: object_grpc.DomainToDetailsPb(details),
	}, nil
}
