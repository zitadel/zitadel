package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type Step2 struct {
	DefaultPasswordComplexityPolicy iam_model.PasswordComplexityPolicy
}

func (s *Step2) Step() domain.Step {
	return domain.Step2
}

func (s *Step2) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep2(ctx, s)
}

func (c *Commands) SetupStep2(ctx context.Context, step *Step2) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.EventPusher, error) {
		iamAgg := IAMAggregateFromWriteModel(&iam.WriteModel)
		event, err := c.addDefaultPasswordComplexityPolicy(ctx, iamAgg, NewIAMPasswordComplexityPolicyWriteModel(), &domain.PasswordComplexityPolicy{
			MinLength:    step.DefaultPasswordComplexityPolicy.MinLength,
			HasLowercase: step.DefaultPasswordComplexityPolicy.HasLowercase,
			HasUppercase: step.DefaultPasswordComplexityPolicy.HasUppercase,
			HasNumber:    step.DefaultPasswordComplexityPolicy.HasNumber,
			HasSymbol:    step.DefaultPasswordComplexityPolicy.HasSymbol,
		})
		if err != nil {
			return nil, err
		}
		logging.Log("SETUP-ADgd2").Info("default password complexity policy set up")
		return []eventstore.EventPusher{event}, nil
	}
	return c.setup(ctx, step, fn)
}
