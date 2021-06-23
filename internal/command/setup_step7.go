package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
)

type Step7 struct {
	OTP bool
}

func (s *Step7) Step() domain.Step {
	return domain.Step7
}

func (s *Step7) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep7(ctx, s)
}

func (c *Commands) SetupStep7(ctx context.Context, step *Step7) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.EventPusher, error) {
		secondFactorModel := NewIAMSecondFactorWriteModel(domain.SecondFactorTypeOTP)
		iamAgg := IAMAggregateFromWriteModel(&secondFactorModel.SecondFactorWriteModel.WriteModel)
		if !step.OTP {
			return []eventstore.EventPusher{}, nil
		}
		event, err := c.addSecondFactorToDefaultLoginPolicy(ctx, iamAgg, secondFactorModel, domain.SecondFactorTypeOTP)
		if err != nil {
			return nil, err
		}
		logging.Log("SETUP-Dggsg").Info("added Provider0 to 2FA login policy")
		return []eventstore.EventPusher{event}, nil
	}
	return c.setup(ctx, step, fn)
}
