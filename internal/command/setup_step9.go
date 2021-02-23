package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
)

type Step9 struct {
	Passwordless bool
}

func (s *Step9) Step() domain.Step {
	return domain.Step9
}

func (s *Step9) execute(ctx context.Context, commandSide *CommandSide) error {
	return commandSide.SetupStep9(ctx, s)
}

func (r *CommandSide) SetupStep9(ctx context.Context, step *Step9) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.EventPusher, error) {
		multiFactorModel := NewIAMMultiFactorWriteModel()
		iamAgg := IAMAggregateFromWriteModel(&multiFactorModel.MultiFactoryWriteModel.WriteModel)
		if !step.Passwordless {
			return []eventstore.EventPusher{}, nil
		}
		passwordlessEvent, err := setPasswordlessAllowedInPolicy(ctx, r, iamAgg)
		if err != nil {
			return nil, err
		}
		logging.Log("SETUP-AEG2t").Info("allowed passwordless in login policy")
		multifactorEvent, err := r.addMultiFactorToDefaultLoginPolicy(ctx, iamAgg, multiFactorModel, domain.MultiFactorTypeU2FWithPIN)
		if err != nil {
			return nil, err
		}
		logging.Log("SETUP-ADfng").Info("added passwordless to MFA login policy")
		return []eventstore.EventPusher{passwordlessEvent, multifactorEvent}, nil
	}
	return r.setup(ctx, step, fn)
}

func setPasswordlessAllowedInPolicy(ctx context.Context, c *CommandSide, iamAgg *eventstore.Aggregate) (eventstore.EventPusher, error) {
	policy, err := c.getDefaultLoginPolicy(ctx)
	if err != nil {
		return nil, err
	}
	policy.PasswordlessType = domain.PasswordlessTypeAllowed
	return c.changeDefaultLoginPolicy(ctx, iamAgg, NewIAMLoginPolicyWriteModel(), policy)
}
