package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
)

type Step18 struct {
	LockoutPolicy domain.LockoutPolicy
}

func (s *Step18) Step() domain.Step {
	return domain.Step18
}

func (s *Step18) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep18(ctx, s)
}

func (c *Commands) SetupStep18(ctx context.Context, step *Step18) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.EventPusher, error) {
		iamAgg := IAMAggregateFromWriteModel(&iam.WriteModel)
		addedPolicy := NewIAMLockoutPolicyWriteModel()
		events, err := c.addDefaultLockoutPolicy(ctx, iamAgg, addedPolicy, &step.LockoutPolicy)
		if err != nil {
			return nil, err
		}

		logging.Log("SETUP-3m99ds").Info("default lockout policy set up")
		return []eventstore.EventPusher{events}, nil
	}
	return c.setup(ctx, step, fn)
}
