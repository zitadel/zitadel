package management

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetPasswordComplexityPolicy(ctx context.Context, req *mgmt_pb.GetPasswordComplexityPolicyRequest) (*mgmt_pb.GetPasswordComplexityPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPasswordComplexityPolicy not implemented")
}

func (s *Server) GetDefaultPasswordComplexityPolicy(ctx context.Context, req *mgmt_pb.GetDefaultPasswordComplexityPolicyRequest) (*mgmt_pb.GetDefaultPasswordComplexityPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDefaultPasswordComplexityPolicy not implemented")
}

func (s *Server) AddCustomPasswordComplexityPolicy(ctx context.Context, req *mgmt_pb.AddCustomPasswordComplexityPolicyRequest) (*mgmt_pb.AddCustomPasswordComplexityPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCustomPasswordComplexityPolicy not implemented")
}

func (s *Server) UpdateCustomPasswordComplexityPolicy(ctx context.Context, req *mgmt_pb.UpdateCustomPasswordComplexityPolicyRequest) (*mgmt_pb.UpdateCustomPasswordComplexityPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCustomPasswordComplexityPolicy not implemented")
}

func (s *Server) ResetPasswordComplexityPolicyToDefault(ctx context.Context, req *mgmt_pb.ResetPasswordComplexityPolicyToDefaultRequest) (*mgmt_pb.ResetPasswordComplexityPolicyToDefaultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetPasswordComplexityPolicyToDefault not implemented")
}
