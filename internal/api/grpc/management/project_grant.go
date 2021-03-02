package management

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetProjectGrantByID(ctx context.Context, req *mgmt_pb.GetProjectGrantByIDRequest) (*mgmt_pb.GetProjectGrantByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProjectGrantByID not implemented")
}

func (s *Server) ListProjectGrants(ctx context.Context, req *mgmt_pb.ListProjectGrantsRequest) (*mgmt_pb.ListProjectGrantsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProjectGrants not implemented")
}

func (s *Server) AddProjectGrant(ctx context.Context, req *mgmt_pb.AddProjectGrantRequest) (*mgmt_pb.AddProjectGrantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddProjectGrant not implemented")
}

func (s *Server) UpdateProjectGrant(ctx context.Context, req *mgmt_pb.UpdateProjectGrantRequest) (*mgmt_pb.UpdateProjectGrantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProjectGrant not implemented")
}

func (s *Server) DeactivateProjectGrant(ctx context.Context, req *mgmt_pb.DeactivateProjectGrantRequest) (*mgmt_pb.DeactivateProjectGrantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeactivateProjectGrant not implemented")
}

func (s *Server) ReactivateProjectGrant(ctx context.Context, req *mgmt_pb.ReactivateProjectGrantRequest) (*mgmt_pb.ReactivateProjectGrantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReactivateProjectGrant not implemented")
}

func (s *Server) RemoveProjectGrant(ctx context.Context, req *mgmt_pb.RemoveProjectGrantRequest) (*mgmt_pb.RemoveProjectGrantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveProjectGrant not implemented")
}

func (s *Server) ListProjectGrantMemberRoles(ctx context.Context, req *mgmt_pb.ListProjectGrantMemberRolesRequest) (*mgmt_pb.ListProjectGrantMemberRolesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProjectGrantMemberRoles not implemented")
}

func (s *Server) ListProjectGrantMembers(ctx context.Context, req *mgmt_pb.ListProjectGrantMembersRequest) (*mgmt_pb.ListProjectGrantMembersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProjectGrantMembers not implemented")
}

func (s *Server) AddProjectGrantMember(ctx context.Context, req *mgmt_pb.AddProjectGrantMemberRequest) (*mgmt_pb.AddProjectGrantMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddProjectGrantMember not implemented")
}

func (s *Server) UpdateProjectGrantMember(ctx context.Context, req *mgmt_pb.UpdateProjectGrantMemberRequest) (*mgmt_pb.UpdateProjectGrantMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProjectGrantMember not implemented")
}

func (s *Server) RemoveProjectGrantMember(ctx context.Context, req *mgmt_pb.RemoveProjectGrantMemberRequest) (*mgmt_pb.RemoveProjectGrantMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveProjectGrantMember not implemented")
}
