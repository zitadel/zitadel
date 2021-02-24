package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type Step4 struct {
	DefaultPasswordLockoutPolicy iam_model.PasswordLockoutPolicy
}

func (s *Step4) Step() domain.Step {
	return domain.Step4
}

func (s *Step4) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep4(ctx, s)
}

func (c *Commands) SetupStep4(ctx context.Context, step *Step4) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.EventPusher, error) {
		iamAgg := IAMAggregateFromWriteModel(&iam.WriteModel)
		event, err := c.addDefaultPasswordLockoutPolicy(ctx, iamAgg, NewIAMPasswordLockoutPolicyWriteModel(), &domain.PasswordLockoutPolicy{
			MaxAttempts:         step.DefaultPasswordLockoutPolicy.MaxAttempts,
			ShowLockOutFailures: step.DefaultPasswordLockoutPolicy.ShowLockOutFailures,
		})
		if err != nil {
			return nil, err
		}
		logging.Log("SETUP-Bfnge").Info("default password lockout policy set up")
		return []eventstore.EventPusher{event}, nil
	}
	return c.setup(ctx, step, fn)
}
