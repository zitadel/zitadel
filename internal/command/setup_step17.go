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
	fn := func(iam *InstanceWriteModel) ([]eventstore.Command, error) {
		iamAgg := InstanceAggregateFromWriteModel(&iam.WriteModel)
		addedPolicy := NewInstancePrivacyPolicyWriteModel()
		events, err := c.addDefaultPrivacyPolicy(ctx, iamAgg, addedPolicy, &step.PrivacyPolicy)
		if err != nil {
			return nil, err
		}

		logging.Log("SETUP-N9sq2").Info("default privacy policy set up")
		return []eventstore.Command{events}, nil
	}
	return c.setup(ctx, step, fn)
}
