package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/errors"
)

func domainPolicyWriteModel(ctx context.Context, filter preparation.FilterToQueryReducer) (*PolicyDomainWriteModel, error) {
	wm, err := orgDomainPolicy(ctx, filter)
	if err != nil || wm != nil && wm.State.Exists() {
		return wm, err
	}
	wm, err = instanceDomainPolicy(ctx, filter)
	if err != nil || wm != nil {
		return wm, err
	}
	return nil, errors.ThrowInternal(nil, "USER-Ggk9n", "Errors.Internal")
}

func orgDomainPolicy(ctx context.Context, filter preparation.FilterToQueryReducer) (*PolicyDomainWriteModel, error) {
	policy := NewOrgDomainPolicyWriteModel(authz.GetCtxData(ctx).OrgID)
	events, err := filter(ctx, policy.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, nil
	}
	policy.AppendEvents(events...)
	err = policy.Reduce()
	return &policy.PolicyDomainWriteModel, err
}

func instanceDomainPolicy(ctx context.Context, filter preparation.FilterToQueryReducer) (*PolicyDomainWriteModel, error) {
	policy := NewInstanceDomainPolicyWriteModel(ctx)
	events, err := filter(ctx, policy.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, nil
	}
	policy.AppendEvents(events...)
	err = policy.Reduce()
	return &policy.PolicyDomainWriteModel, err
}
