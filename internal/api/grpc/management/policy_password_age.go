package management

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetPasswordAgePolicy(ctx context.Context, req *mgmt_pb.GetPasswordAgePolicyRequest) (*mgmt_pb.GetPasswordAgePolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPasswordAgePolicy not implemented")
}

func (s *Server) GetDefaultPasswordAgePolicy(ctx context.Context, req *mgmt_pb.GetDefaultPasswordAgePolicyRequest) (*mgmt_pb.GetDefaultPasswordAgePolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDefaultPasswordAgePolicy not implemented")
}

func (s *Server) AddCustomPasswordAgePolicy(ctx context.Context, req *mgmt_pb.AddCustomPasswordAgePolicyRequest) (*mgmt_pb.AddCustomPasswordAgePolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCustomPasswordAgePolicy not implemented")
}

func (s *Server) UpdateCustomPasswordAgePolicy(ctx context.Context, req *mgmt_pb.UpdateCustomPasswordAgePolicyRequest) (*mgmt_pb.UpdateCustomPasswordAgePolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCustomPasswordAgePolicy not implemented")
}

func (s *Server) ResetPasswordAgePolicyToDefault(ctx context.Context, req *mgmt_pb.ResetPasswordAgePolicyToDefaultRequest) (*mgmt_pb.ResetPasswordAgePolicyToDefaultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetPasswordAgePolicyToDefault not implemented")
}
