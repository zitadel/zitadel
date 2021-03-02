package management

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetPasswordLockoutPolicy(ctx context.Context, req *mgmt_pb.GetPasswordLockoutPolicyRequest) (*mgmt_pb.GetPasswordLockoutPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPasswordLockoutPolicy not implemented")
}

func (s *Server) GetDefaultPasswordLockoutPolicy(ctx context.Context, req *mgmt_pb.GetDefaultPasswordLockoutPolicyRequest) (*mgmt_pb.GetDefaultPasswordLockoutPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDefaultPasswordLockoutPolicy not implemented")
}

func (s *Server) AddCustomPasswordLockoutPolicy(ctx context.Context, req *mgmt_pb.AddCustomPasswordLockoutPolicyRequest) (*mgmt_pb.AddCustomPasswordLockoutPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCustomPasswordLockoutPolicy not implemented")
}

func (s *Server) UpdateCustomPasswordLockoutPolicy(ctx context.Context, req *mgmt_pb.UpdateCustomPasswordLockoutPolicyRequest) (*mgmt_pb.UpdateCustomPasswordLockoutPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCustomPasswordLockoutPolicy not implemented")
}

func (s *Server) ResetPasswordLockoutPolicyToDefault(ctx context.Context, req *mgmt_pb.ResetPasswordLockoutPolicyToDefaultRequest) (*mgmt_pb.ResetPasswordLockoutPolicyToDefaultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetPasswordLockoutPolicyToDefault not implemented")
}
