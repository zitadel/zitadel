package command

import (
	"context"

	"github.com/caos/logging"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
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
	fn := func(iam *IAMWriteModel) (*iam_repo.Aggregate, error) {
		multiFactorModel := NewIAMMultiFactorWriteModel()
		iamAgg := IAMAggregateFromWriteModel(&multiFactorModel.MultiFactoryWriteModel.WriteModel)
		if !step.Passwordless {
			return iamAgg, nil
		}
		err := setPasswordlessAllowedInPolicy(ctx, r, iamAgg)
		if err != nil {
			return nil, err
		}
		logging.Log("SETUP-AEG2t").Info("allowed passwordless in login policy")
		err = r.addMultiFactorToDefaultLoginPolicy(ctx, iamAgg, multiFactorModel, iam_model.MultiFactorTypeU2FWithPIN)
		if err != nil {
			return nil, err
		}
		logging.Log("SETUP-ADfng").Info("added passwordless to MFA login policy")
		return iamAgg, err
	}
	return r.setup(ctx, step, fn)
}

func setPasswordlessAllowedInPolicy(ctx context.Context, c *CommandSide, iamAgg *iam_repo.Aggregate) error {
	policy, err := c.getDefaultLoginPolicy(ctx)
	if err != nil {
		return err
	}
	policy.PasswordlessType = domain.PasswordlessTypeAllowed
	return c.changeDefaultLoginPolicy(ctx, iamAgg, NewIAMLoginPolicyWriteModel(), policy)
}
