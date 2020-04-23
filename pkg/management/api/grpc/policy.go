package grpc

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) CreatePasswordComplexityPolicy(ctx context.Context, policy *PasswordComplexityPolicyCreate) (*PasswordComplexityPolicy, error) {
	policyresp, err := s.policy.CreatePasswordComplexityPolicy(ctx, passwordComplexityPolicyCreateToModel(policy))
	if err != nil {
		return nil, err
	}

	return passwordComplexityPolicyFromModel(policyresp), nil
}

func (s *Server) GetPasswordComplexityPolicy(ctx context.Context, _ *empty.Empty) (*PasswordComplexityPolicy, error) {
	policy, err := s.policy.GetPasswordComplexityPolicy(ctx)
	if err != nil {
		return nil, err
	}

	return passwordComplexityPolicyFromModel(policy), nil
}

func (s *Server) UpdatePasswordComplexityPolicy(ctx context.Context, policy *PasswordComplexityPolicyUpdate) (*PasswordComplexityPolicy, error) {
	policyresp, err := s.policy.UpdatePasswordComplexityPolicy(ctx, passwordComplexityPolicyUpdateToModel(policy))
	if err != nil {
		return nil, err
	}

	return passwordComplexityPolicyFromModel(policyresp), nil
}

func (s *Server) DeletePasswordComplexityPolicy(ctx context.Context, ID *PasswordComplexityPolicyID) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-skw3f", "Not implemented")
}

func (s *Server) CreatePasswordAgePolicy(ctx context.Context, policy *PasswordAgePolicyCreate) (*PasswordAgePolicy, error) {
	policyresp, err := s.policy.CreatePasswordAgePolicy(ctx, passwordAgePolicyCreateToModel(policy))
	if err != nil {
		return nil, err
	}

	return passwordAgePolicyFromModel(policyresp), nil
}

func (s *Server) GetPasswordAgePolicy(ctx context.Context, _ *empty.Empty) (*PasswordAgePolicy, error) {
	policy, err := s.policy.GetPasswordAgePolicy(ctx)
	if err != nil {
		return nil, err
	}

	return passwordAgePolicyFromModel(policy), nil
}

func (s *Server) UpdatePasswordAgePolicy(ctx context.Context, policy *PasswordAgePolicyUpdate) (*PasswordAgePolicy, error) {
	policyresp, err := s.policy.UpdatePasswordAgePolicy(ctx, passwordAgePolicyUpdateToModel(policy))
	if err != nil {
		return nil, err
	}

	return passwordAgePolicyFromModel(policyresp), nil
}

func (s *Server) DeletePasswordAgePolicy(ctx context.Context, ID *PasswordAgePolicyID) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-plo67", "Not implemented")
}

func (s *Server) GetPasswordLockoutPolicy(ctx context.Context, _ *empty.Empty) (*PasswordLockoutPolicy, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-GHkd9", "Not implemented")
}

func (s *Server) CreatePasswordLockoutPolicy(ctx context.Context, policy *PasswordLockoutPolicyCreate) (*PasswordLockoutPolicy, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-mdk3c", "Not implemented")
}

func (s *Server) UpdatePasswordLockoutPolicy(ctx context.Context, policy *PasswordLockoutPolicyUpdate) (*PasswordLockoutPolicy, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-8dbN4", "Not implemented")
}

func (s *Server) DeletePasswordLockoutPolicy(ctx context.Context, ID *PasswordLockoutPolicyID) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-plV53", "Not implemented")
}
