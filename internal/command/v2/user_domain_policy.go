package command

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/errors"
)

func domainPolicyWriteModel(ctx context.Context, instanceID string, filter preparation.FilterToQueryReducer) (*command.PolicyDomainWriteModel, error) {
	wm, err := orgDomainPolicy(ctx, filter)
	if err != nil || wm != nil && wm.State.Exists() {
		return wm, err
	}
	wm, err = instanceDomainPolicy(ctx, instanceID, filter)
	if err != nil || wm != nil {
		return wm, err
	}
	return nil, errors.ThrowInternal(nil, "USER-Ggk9n", "Errors.Internal")
}

func orgDomainPolicy(ctx context.Context, filter preparation.FilterToQueryReducer) (*command.PolicyDomainWriteModel, error) {
	policy := command.NewOrgDomainPolicyWriteModel(authz.GetCtxData(ctx).OrgID)
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

func instanceDomainPolicy(ctx context.Context, instanceID string, filter preparation.FilterToQueryReducer) (*command.PolicyDomainWriteModel, error) {
	policy := command.NewInstanceDomainPolicyWriteModel(instanceID)
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
