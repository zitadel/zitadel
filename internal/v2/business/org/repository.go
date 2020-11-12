package org

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type Repository struct {
	eventstore *eventstore.Eventstore
}

type Config struct {
	Eventstore *eventstore.Eventstore
}

func StartRepository(config *Config) *Repository {
	return &Repository{
		eventstore: config.Eventstore,
	}
}

func (r *Repository) GetMyOrgIamPolicy(ctx context.Context) (*iam_model.OrgIAMPolicyView, error) {
	return nil, nil
}

func (r *Repository) GetLoginPolicy(ctx context.Context) (*iam_model.LoginPolicyView, error) {
	return nil, nil
}
func (r *Repository) GetDefaultLoginPolicy(ctx context.Context) (*iam_model.LoginPolicyView, error) {
	return nil, nil
}
func (r *Repository) AddLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	return nil, nil
}
func (r *Repository) ChangeLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	return nil, nil
}
func (r *Repository) RemoveLoginPolicy(ctx context.Context) error {
	return nil
}
func (r *Repository) SearchIDPProviders(ctx context.Context, request *iam_model.IDPProviderSearchRequest) (*iam_model.IDPProviderSearchResponse, error) {
	return nil, nil
}
func (r *Repository) AddIDPProviderToLoginPolicy(ctx context.Context, provider *iam_model.IDPProvider) (*iam_model.IDPProvider, error) {
	return nil, nil
}
func (r *Repository) RemoveIDPProviderFromLoginPolicy(ctx context.Context, provider *iam_model.IDPProvider) error {
	return nil
}

func (r *Repository) GetPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error) {
	return nil, nil
}
func (r *Repository) GetDefaultPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error) {
	return nil, nil
}
func (r *Repository) AddPasswordComplexityPolicy(ctx context.Context, policy *iam_model.PasswordComplexityPolicy) (*iam_model.PasswordComplexityPolicy, error) {
	return nil, nil
}
func (r *Repository) ChangePasswordComplexityPolicy(ctx context.Context, policy *iam_model.PasswordComplexityPolicy) (*iam_model.PasswordComplexityPolicy, error) {
	return nil, nil
}
func (r *Repository) RemovePasswordComplexityPolicy(ctx context.Context) error {
	return nil
}

func (r *Repository) GetPasswordAgePolicy(ctx context.Context) (*iam_model.PasswordAgePolicyView, error) {
	return nil, nil
}
func (r *Repository) GetDefaultPasswordAgePolicy(ctx context.Context) (*iam_model.PasswordAgePolicyView, error) {
	return nil, nil
}
func (r *Repository) AddPasswordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) (*iam_model.PasswordAgePolicy, error) {
	return nil, nil
}
func (r *Repository) ChangePasswordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) (*iam_model.PasswordAgePolicy, error) {
	return nil, nil
}
func (r *Repository) RemovePasswordAgePolicy(ctx context.Context) error {
	return nil
}

func (r *Repository) GetPasswordLockoutPolicy(ctx context.Context) (*iam_model.PasswordLockoutPolicyView, error) {
	return nil, nil
}
func (r *Repository) GetDefaultPasswordLockoutPolicy(ctx context.Context) (*iam_model.PasswordLockoutPolicyView, error) {
	return nil, nil
}
func (r *Repository) AddPasswordLockoutPolicy(ctx context.Context, policy *iam_model.PasswordLockoutPolicy) (*iam_model.PasswordLockoutPolicy, error) {
	return nil, nil
}
func (r *Repository) ChangePasswordLockoutPolicy(ctx context.Context, policy *iam_model.PasswordLockoutPolicy) (*iam_model.PasswordLockoutPolicy, error) {
	return nil, nil
}
func (r *Repository) RemovePasswordLockoutPolicy(ctx context.Context) error {
	return nil
}
