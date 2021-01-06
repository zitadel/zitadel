package command

import (
	"context"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

type Step4 struct {
	DefaultPasswordLockoutPolicy iam_model.PasswordLockoutPolicy
}

func (s *Step4) Step() domain.Step {
	return domain.Step4
}

func (s *Step4) execute(ctx context.Context, commandSide *CommandSide) error {
	return commandSide.SetupStep4(ctx, s)
}

func (r *CommandSide) SetupStep4(ctx context.Context, step *Step4) error {
	fn := func(iam *IAMWriteModel) (*iam_repo.Aggregate, error) {
		iamAgg := IAMAggregateFromWriteModel(&iam.WriteModel)
		err := r.addDefaultPasswordLockoutPolicy(ctx, iamAgg, NewIAMPasswordLockoutPolicyWriteModel(iam.AggregateID), &domain.PasswordLockoutPolicy{
			MaxAttempts:         step.DefaultPasswordLockoutPolicy.MaxAttempts,
			ShowLockOutFailures: step.DefaultPasswordLockoutPolicy.ShowLockOutFailures,
		})
		if err != nil {
			return nil, err
		}
		return iamAgg, nil
	}
	return r.setup(ctx, step, fn)
}
