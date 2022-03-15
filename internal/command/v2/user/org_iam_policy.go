package user

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/errors"
)

func orgIAMPolicyWriteModel(ctx context.Context, filter preparation.FilterToQueryReducer) (*command.PolicyOrgIAMWriteModel, error) {
	wm, err := customOrgIAMPolicy(ctx, filter)
	if err != nil || wm != nil {
		return wm, err
	}
	wm, err = defaultOrgIAMPolicy(ctx, filter)
	if err != nil || wm != nil {
		return wm, err
	}
	return nil, errors.ThrowInternal(nil, "USER-Ggk9n", "Errors.Internal")
}

func customOrgIAMPolicy(ctx context.Context, filter preparation.FilterToQueryReducer) (*command.PolicyOrgIAMWriteModel, error) {
	policy := command.NewORGOrgIAMPolicyWriteModel(authz.GetCtxData(ctx).OrgID)
	events, err := filter(ctx, policy.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, nil
	}
	policy.AppendEvents(events...)
	err = policy.Reduce()
	return &policy.PolicyOrgIAMWriteModel, err
}

func defaultOrgIAMPolicy(ctx context.Context, filter preparation.FilterToQueryReducer) (*command.PolicyOrgIAMWriteModel, error) {
	policy := command.NewIAMOrgIAMPolicyWriteModel()
	events, err := filter(ctx, policy.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, nil
	}
	policy.AppendEvents(events...)
	err = policy.Reduce()
	return &policy.PolicyOrgIAMWriteModel, err
}
