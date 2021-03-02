package management

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	member_grpc "github.com/caos/zitadel/internal/api/grpc/member"
	object_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	proj_grpc "github.com/caos/zitadel/internal/api/grpc/project"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetProjectGrantByID(ctx context.Context, req *mgmt_pb.GetProjectGrantByIDRequest) (*mgmt_pb.GetProjectGrantByIDResponse, error) {
	grant, err := s.project.ProjectGrantByID(ctx, req.GrantId)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetProjectGrantByIDResponse{
		ProjectGrant: proj_grpc.GrantedProjectToPb(grant),
	}, nil
}

func (s *Server) ListProjectGrants(ctx context.Context, req *mgmt_pb.ListProjectGrantsRequest) (*mgmt_pb.ListProjectGrantsResponse, error) {
	queries, err := ListProjectGrantsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	domains, err := s.project.SearchProjectGrants(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListProjectGrantsResponse{
		Result: proj_grpc.GrantedProjectsToPb(domains.Result),
		MetaData: object_grpc.ToListDetails(
			domains.TotalResult,
			domains.Sequence,
			domains.Timestamp,
		),
	}, nil
}

func (s *Server) AddProjectGrant(ctx context.Context, req *mgmt_pb.AddProjectGrantRequest) (*mgmt_pb.AddProjectGrantResponse, error) {
	grant, err := s.command.AddProjectGrant(ctx, AddProjectGrantRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddProjectGrantResponse{
		GrantId: grant.GrantID,
		Details: object_grpc.ToDetailsPb(
			grant.Sequence,
			grant.ChangeDate,
			grant.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateProjectGrant(ctx context.Context, req *mgmt_pb.UpdateProjectGrantRequest) (*mgmt_pb.UpdateProjectGrantResponse, error) {
	userGrants, err := s.usergrant.UserGrantsByProjectAndGrantID(ctx, req.ProjectId, req.GrantId)
	if err != nil {
		return nil, err
	}
	grant, err := s.command.ChangeProjectGrant(ctx, UpdateProjectGrantRequestToDomain(req), authz.GetCtxData(ctx).OrgID, userGrantsToIDs(userGrants)...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateProjectGrantResponse{
		Details: object_grpc.ToDetailsPb(
			grant.Sequence,
			grant.ChangeDate,
			grant.ResourceOwner,
		),
	}, nil
}

func (s *Server) DeactivateProjectGrant(ctx context.Context, req *mgmt_pb.DeactivateProjectGrantRequest) (*mgmt_pb.DeactivateProjectGrantResponse, error) {
	err := s.command.DeactivateProjectGrant(ctx, req.ProjectId, req.GrantId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.DeactivateProjectGrantResponse{
		//TODO: details
	}, nil
}

func (s *Server) ReactivateProjectGrant(ctx context.Context, req *mgmt_pb.ReactivateProjectGrantRequest) (*mgmt_pb.ReactivateProjectGrantResponse, error) {
	err := s.command.ReactivateProjectGrant(ctx, req.ProjectId, req.GrantId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ReactivateProjectGrantResponse{
		//TODO: details
	}, nil
}

func (s *Server) RemoveProjectGrant(ctx context.Context, req *mgmt_pb.RemoveProjectGrantRequest) (*mgmt_pb.RemoveProjectGrantResponse, error) {
	err := s.command.RemoveProjectGrant(ctx, req.ProjectId, req.GrantId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveProjectGrantResponse{
		//TODO: details
	}, nil
}

func (s *Server) ListProjectGrantMemberRoles(ctx context.Context, req *mgmt_pb.ListProjectGrantMemberRolesRequest) (*mgmt_pb.ListProjectGrantMemberRolesResponse, error) {
	roles := s.project.GetProjectGrantMemberRoles()
	return &mgmt_pb.ListProjectGrantMemberRolesResponse{
		Result: roles,
		//TODO: metadata
	}, nil
}

func (s *Server) ListProjectGrantMembers(ctx context.Context, req *mgmt_pb.ListProjectGrantMembersRequest) (*mgmt_pb.ListProjectGrantMembersResponse, error) {
	response, err := s.project.SearchProjectGrantMembers(ctx, ListProjectGrantMembersRequestToModel(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListProjectGrantMembersResponse{
		Result: member_grpc.ProjectGrantMembersToPb(response.Result),
		MetaData: object_grpc.ToListDetails(
			response.TotalResult,
			response.Sequence,
			response.Timestamp,
		),
	}, nil
}

func (s *Server) AddProjectGrantMember(ctx context.Context, req *mgmt_pb.AddProjectGrantMemberRequest) (*mgmt_pb.AddProjectGrantMemberResponse, error) {
	member, err := s.command.AddProjectGrantMember(ctx, AddProjectGrantMemberRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddProjectGrantMemberResponse{
		Details: object_grpc.ToDetailsPb(
			member.Sequence,
			member.ChangeDate,
			member.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateProjectGrantMember(ctx context.Context, req *mgmt_pb.UpdateProjectGrantMemberRequest) (*mgmt_pb.UpdateProjectGrantMemberResponse, error) {
	member, err := s.command.ChangeProjectGrantMember(ctx, UpdateProjectGrantMemberRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateProjectGrantMemberResponse{
		Details: object_grpc.ToDetailsPb(
			member.Sequence,
			member.ChangeDate,
			member.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveProjectGrantMember(ctx context.Context, req *mgmt_pb.RemoveProjectGrantMemberRequest) (*mgmt_pb.RemoveProjectGrantMemberResponse, error) {
	err := s.command.RemoveProjectGrantMember(ctx, req.ProjectId, req.UserId, req.GrantId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveProjectGrantMemberResponse{
		//TODO: details
	}, nil
}
