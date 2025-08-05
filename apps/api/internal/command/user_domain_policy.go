package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// Deprecated: User commands.domainPolicyWriteModel directly, to remove use of eventstore.Filter function
func domainPolicyWriteModel(ctx context.Context, filter preparation.FilterToQueryReducer, orgID string) (*PolicyDomainWriteModel, error) {
	wm, err := orgDomainPolicy(ctx, filter, orgID)
	if err != nil {
		return nil, err
	}
	if wm != nil && wm.State.Exists() {
		return &wm.PolicyDomainWriteModel, err
	}
	instanceWriteModel, err := instanceDomainPolicy(ctx, filter)
	if err != nil {
		return nil, err
	}
	if instanceWriteModel != nil && instanceWriteModel.State.Exists() {
		return &instanceWriteModel.PolicyDomainWriteModel, err
	}
	return nil, zerrors.ThrowInternal(nil, "USER-Ggk9n", "Errors.Internal")
}

func (c *Commands) domainPolicyWriteModel(ctx context.Context, orgID string) (*PolicyDomainWriteModel, error) {
	wm, err := c.orgDomainPolicyWriteModel(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if wm != nil && wm.State.Exists() {
		return &wm.PolicyDomainWriteModel, err
	}
	instanceWriteModel, err := c.instanceDomainPolicyWriteModel(ctx)
	if err != nil {
		return nil, err
	}
	if instanceWriteModel != nil && instanceWriteModel.State.Exists() {
		return &instanceWriteModel.PolicyDomainWriteModel, err
	}
	return nil, zerrors.ThrowInternal(nil, "USER-Ggk9n", "Errors.Internal")
}

// Deprecated: Use commands.orgDomainPolicyWriteModel directly, to remove use of eventstore.Filter function
func orgDomainPolicy(ctx context.Context, filter preparation.FilterToQueryReducer, orgID string) (*OrgDomainPolicyWriteModel, error) {
	policy := NewOrgDomainPolicyWriteModel(orgID)
	events, err := filter(ctx, policy.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return policy, nil
	}
	policy.AppendEvents(events...)
	err = policy.Reduce()
	return policy, err
}

// Deprecated: Use commands.instanceDomainPolicyWriteModel directly, to remove use of eventstore.Filter function
func instanceDomainPolicy(ctx context.Context, filter preparation.FilterToQueryReducer) (_ *InstanceDomainPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	policy := NewInstanceDomainPolicyWriteModel(ctx)
	events, err := filter(ctx, policy.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return policy, nil
	}
	policy.AppendEvents(events...)
	err = policy.Reduce()
	return policy, err
}

func domainPolicyUsernames(ctx context.Context, filter preparation.FilterToQueryReducer, orgID string) (_ *DomainPolicyUsernamesWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	policy := NewDomainPolicyUsernamesWriteModel(orgID)
	events, err := filter(ctx, policy.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return policy, nil
	}
	policy.AppendEvents(events...)
	err = policy.Reduce()
	return policy, err
}

func domainPolicyOrgs(ctx context.Context, filter preparation.FilterToQueryReducer) (*DomainPolicyOrgsWriteModel, error) {
	policy := NewDomainPolicyOrgsWriteModel()
	events, err := filter(ctx, policy.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return policy, nil
	}
	policy.AppendEvents(events...)
	err = policy.Reduce()
	return policy, err
}
