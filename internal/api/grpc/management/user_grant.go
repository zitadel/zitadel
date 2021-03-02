package management

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetUserGrantByID(ctx context.Context, req *mgmt_pb.GetUserGrantByIDRequest) (*mgmt_pb.GetUserGrantByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserGrantByID not implemented")
}

func (s *Server) ListUserGrants(ctx context.Context, req *mgmt_pb.ListUserGrantRequest) (*mgmt_pb.ListUserGrantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUserGrants not implemented")
}

func (s *Server) AddUserGrant(ctx context.Context, req *mgmt_pb.AddUserGrantRequest) (*mgmt_pb.AddUserGrantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUserGrant not implemented")
}

func (s *Server) UpdateUserGrant(ctx context.Context, req *mgmt_pb.UpdateUserGrantRequest) (*mgmt_pb.UpdateUserGrantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserGrant not implemented")
}

func (s *Server) DeactivateUserGrant(ctx context.Context, req *mgmt_pb.DeactivateUserGrantRequest) (*mgmt_pb.DeactivateUserGrantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeactivateUserGrant not implemented")
}

func (s *Server) ReactivateUserGrant(ctx context.Context, req *mgmt_pb.ReactivateUserGrantRequest) (*mgmt_pb.ReactivateUserGrantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReactivateUserGrant not implemented")
}

func (s *Server) RemoveUserGrant(ctx context.Context, req *mgmt_pb.RemoveUserGrantRequest) (*mgmt_pb.RemoveUserGrantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveUserGrant not implemented")
}

func (s *Server) BulkRemoveUserGrant(ctx context.Context, req *mgmt_pb.BulkRemoveUserGrantRequest) (*mgmt_pb.BulkRemoveUserGrantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BulkRemoveUserGrant not implemented")
}
