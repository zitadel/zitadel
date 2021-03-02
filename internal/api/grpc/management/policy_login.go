package management

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetLoginPolicy(ctx context.Context, req *mgmt_pb.GetLoginPolicyRequest) (*mgmt_pb.GetLoginPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLoginPolicy not implemented")
}

func (s *Server) GetDefaultLoginPolicy(ctx context.Context, req *mgmt_pb.GetDefaultLoginPolicyRequest) (*mgmt_pb.GetDefaultLoginPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDefaultLoginPolicy not implemented")
}

func (s *Server) AddCustomLoginPolicy(ctx context.Context, req *mgmt_pb.AddCustomLoginPolicyRequest) (*mgmt_pb.AddCustomLoginPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCustomLoginPolicy not implemented")
}

func (s *Server) UpdateCustomLoginPolicy(ctx context.Context, req *mgmt_pb.UpdateCustomLoginPolicyRequest) (*mgmt_pb.UpdateCustomLoginPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCustomLoginPolicy not implemented")
}

func (s *Server) ResetLoginPolicyToDefault(ctx context.Context, req *mgmt_pb.ResetLoginPolicyToDefaultRequest) (*mgmt_pb.ResetLoginPolicyToDefaultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetLoginPolicyToDefault not implemented")
}

func (s *Server) ListLoginPolicyIDPs(ctx context.Context, req *mgmt_pb.ListLoginPolicyIDPsRequest) (*mgmt_pb.ListLoginPolicyIDPsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLoginPolicyIDPs not implemented")
}

func (s *Server) AddIDPToLoginPolicy(ctx context.Context, req *mgmt_pb.AddIDPToLoginPolicyRequest) (*mgmt_pb.AddIDPToLoginPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddIDPToLoginPolicy not implemented")
}

func (s *Server) RemoveIDPFromLoginPolicy(ctx context.Context, req *mgmt_pb.RemoveIDPFromLoginPolicyRequest) (*mgmt_pb.RemoveIDPFromLoginPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveIDPFromLoginPolicy not implemented")
}

func (s *Server) ListLoginPolicySecondFactors(ctx context.Context, req *mgmt_pb.ListLoginPolicySecondFactorsRequest) (*mgmt_pb.ListLoginPolicySecondFactorsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLoginPolicySecondFactors not implemented")
}

func (s *Server) AddSecondFactorToLoginPolicy(ctx context.Context, req *mgmt_pb.AddSecondFactorToLoginPolicyRequest) (*mgmt_pb.AddSecondFactorToLoginPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddSecondFactorToLoginPolicy not implemented")
}

func (s *Server) RemoveSecondFactorFromLoginPolicy(ctx context.Context, req *mgmt_pb.RemoveSecondFactorFromLoginPolicyRequest) (*mgmt_pb.RemoveSecondFactorFromLoginPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveSecondFactorFromLoginPolicy not implemented")
}

func (s *Server) ListLoginPolicyMultiFactors(ctx context.Context, req *mgmt_pb.ListLoginPolicyMultiFactorsRequest) (*mgmt_pb.ListLoginPolicyMultiFactorsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLoginPolicyMultiFactors not implemented")
}

func (s *Server) AddMultiFactorToLoginPolicy(ctx context.Context, req *mgmt_pb.AddMultiFactorToLoginPolicyRequest) (*mgmt_pb.AddMultiFactorToLoginPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddMultiFactorToLoginPolicy not implemented")
}

func (s *Server) RemoveMultiFactorFromLoginPolicy(ctx context.Context, req *mgmt_pb.RemoveMultiFactorFromLoginPolicyRequest) (*mgmt_pb.RemoveMultiFactorFromLoginPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveMultiFactorFromLoginPolicy not implemented")
}
