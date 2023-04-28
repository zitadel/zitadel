package management

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	member_grpc "github.com/zitadel/zitadel/internal/api/grpc/member"
	object_grpc "github.com/zitadel/zitadel/internal/api/grpc/object"
	proj_grpc "github.com/zitadel/zitadel/internal/api/grpc/project"
	"github.com/zitadel/zitadel/internal/query"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func (s *Server) GetProjectGrantByID(ctx context.Context, req *mgmt_pb.GetProjectGrantByIDRequest) (*mgmt_pb.GetProjectGrantByIDResponse, error) {
	grant, err := s.query.ProjectGrantByID(ctx, true, req.GrantId, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetProjectGrantByIDResponse{
		ProjectGrant: proj_grpc.GrantedProjectViewToPb(grant),
	}, nil
}

func (s *Server) ListProjectGrants(ctx context.Context, req *mgmt_pb.ListProjectGrantsRequest) (*mgmt_pb.ListProjectGrantsResponse, error) {
	queries, err := listProjectGrantsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	err = queries.AppendMyResourceOwnerQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	grants, err := s.query.SearchProjectGrants(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListProjectGrantsResponse{
		Result:  proj_grpc.GrantedProjectViewsToPb(grants.ProjectGrants),
		Details: object_grpc.ToListDetails(grants.Count, grants.Sequence, grants.Timestamp),
	}, nil
}

func (s *Server) ListAllProjectGrants(ctx context.Context, req *mgmt_pb.ListAllProjectGrantsRequest) (*mgmt_pb.ListAllProjectGrantsResponse, error) {
	queries, err := listAllProjectGrantsRequestToModel(req)
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
	grants, err := s.query.SearchProjectGrants(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListAllProjectGrantsResponse{
		Result:  proj_grpc.GrantedProjectViewsToPb(grants.ProjectGrants),
		Details: object_grpc.ToListDetails(grants.Count, grants.Sequence, grants.Timestamp),
	}, nil
}

func (s *Server) AddProjectGrant(ctx context.Context, req *mgmt_pb.AddProjectGrantRequest) (*mgmt_pb.AddProjectGrantResponse, error) {
	grant, err := s.command.AddProjectGrant(ctx, AddProjectGrantRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddProjectGrantResponse{
		GrantId: grant.GrantID,
		Details: object_grpc.AddToDetailsPb(
			grant.Sequence,
			grant.ChangeDate,
			grant.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateProjectGrant(ctx context.Context, req *mgmt_pb.UpdateProjectGrantRequest) (*mgmt_pb.UpdateProjectGrantResponse, error) {
	projectQuery, err := query.NewUserGrantProjectIDSearchQuery(req.ProjectId)
	if err != nil {
		return nil, err
	}
	grantQuery, err := query.NewUserGrantGrantIDSearchQuery(req.GrantId)
	if err != nil {
		return nil, err
	}
	grants, err := s.query.UserGrants(ctx, &query.UserGrantsQueries{
		Queries: []query.SearchQuery{projectQuery, grantQuery},
	}, true, false)
	if err != nil {
		return nil, err
	}
	grant, err := s.command.ChangeProjectGrant(ctx, UpdateProjectGrantRequestToDomain(req), authz.GetCtxData(ctx).OrgID, userGrantsToIDs(grants.UserGrants)...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateProjectGrantResponse{
		Details: object_grpc.ChangeToDetailsPb(
			grant.Sequence,
			grant.ChangeDate,
			grant.ResourceOwner,
		),
	}, nil
}

func (s *Server) DeactivateProjectGrant(ctx context.Context, req *mgmt_pb.DeactivateProjectGrantRequest) (*mgmt_pb.DeactivateProjectGrantResponse, error) {
	details, err := s.command.DeactivateProjectGrant(ctx, req.ProjectId, req.GrantId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.DeactivateProjectGrantResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) ReactivateProjectGrant(ctx context.Context, req *mgmt_pb.ReactivateProjectGrantRequest) (*mgmt_pb.ReactivateProjectGrantResponse, error) {
	details, err := s.command.ReactivateProjectGrant(ctx, req.ProjectId, req.GrantId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ReactivateProjectGrantResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) RemoveProjectGrant(ctx context.Context, req *mgmt_pb.RemoveProjectGrantRequest) (*mgmt_pb.RemoveProjectGrantResponse, error) {
	projectQuery, err := query.NewUserGrantProjectIDSearchQuery(req.ProjectId)
	if err != nil {
		return nil, err
	}
	grantQuery, err := query.NewUserGrantGrantIDSearchQuery(req.GrantId)
	if err != nil {
		return nil, err
	}
	userGrants, err := s.query.UserGrants(ctx, &query.UserGrantsQueries{
		Queries: []query.SearchQuery{projectQuery, grantQuery},
	}, false, true)
	if err != nil {
		return nil, err
	}
	details, err := s.command.RemoveProjectGrant(ctx, req.ProjectId, req.GrantId, authz.GetCtxData(ctx).OrgID, userGrantsToIDs(userGrants.UserGrants)...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveProjectGrantResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) ListProjectGrantMemberRoles(ctx context.Context, req *mgmt_pb.ListProjectGrantMemberRolesRequest) (*mgmt_pb.ListProjectGrantMemberRolesResponse, error) {
	roles := s.query.GetProjectGrantMemberRoles()
	return &mgmt_pb.ListProjectGrantMemberRolesResponse{
		Result:  roles,
		Details: object_grpc.ToListDetails(uint64(len(roles)), 0, time.Now()),
	}, nil
}

func (s *Server) ListProjectGrantMembers(ctx context.Context, req *mgmt_pb.ListProjectGrantMembersRequest) (*mgmt_pb.ListProjectGrantMembersResponse, error) {
	queries, err := ListProjectGrantMembersRequestToModel(ctx, req)
	if err != nil {
		return nil, err
	}
	response, err := s.query.ProjectGrantMembers(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListProjectGrantMembersResponse{
		Result:  member_grpc.MembersToPb(s.assetAPIPrefix(ctx), response.Members),
		Details: object_grpc.ToListDetails(response.Count, response.Sequence, response.Timestamp),
	}, nil
}

func (s *Server) AddProjectGrantMember(ctx context.Context, req *mgmt_pb.AddProjectGrantMemberRequest) (*mgmt_pb.AddProjectGrantMemberResponse, error) {
	member, err := s.command.AddProjectGrantMember(ctx, AddProjectGrantMemberRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddProjectGrantMemberResponse{
		Details: object_grpc.AddToDetailsPb(
			member.Sequence,
			member.ChangeDate,
			member.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateProjectGrantMember(ctx context.Context, req *mgmt_pb.UpdateProjectGrantMemberRequest) (*mgmt_pb.UpdateProjectGrantMemberResponse, error) {
	member, err := s.command.ChangeProjectGrantMember(ctx, UpdateProjectGrantMemberRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateProjectGrantMemberResponse{
		Details: object_grpc.ChangeToDetailsPb(
			member.Sequence,
			member.ChangeDate,
			member.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveProjectGrantMember(ctx context.Context, req *mgmt_pb.RemoveProjectGrantMemberRequest) (*mgmt_pb.RemoveProjectGrantMemberResponse, error) {
	details, err := s.command.RemoveProjectGrantMember(ctx, req.ProjectId, req.UserId, req.GrantId)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveProjectGrantMemberResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}
