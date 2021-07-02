package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
)

type Step17 struct {
	PrivacyPolicy domain.PrivacyPolicy
}

func (s *Step17) Step() domain.Step {
	return domain.Step17
}

func (s *Step17) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep17(ctx, s)
}

func (c *Commands) SetupStep17(ctx context.Context, step *Step17) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.EventPusher, error) {
		iamAgg := IAMAggregateFromWriteModel(&iam.WriteModel)
		addedPolicy := NewIAMPrivacyPolicyWriteModel()
		events, err := c.addDefaultPrivacyPolicy(ctx, iamAgg, addedPolicy, &step.PrivacyPolicy)
		if err != nil {
			return nil, err
		}

		logging.Log("SETUP-4k0LL").Info("default message text set up")
		return []eventstore.EventPusher{events}, nil
	}
	return c.setup(ctx, step, fn)
}
