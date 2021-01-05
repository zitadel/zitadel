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
	return commandSide.SetupStep4(ctx, commandSide.iamID, s)
}

func (r *CommandSide) SetupStep4(ctx context.Context, iamID string, step *Step4) error {
	fn := func(iam *IAMWriteModel) (*iam_repo.Aggregate, error) {
		return r.addDefaultPasswordLockoutPolicy(ctx, NewIAMPasswordLockoutPolicyWriteModel(iam.AggregateID), &iam_model.PasswordLockoutPolicy{
			MaxAttempts:         step.DefaultPasswordLockoutPolicy.MaxAttempts,
			ShowLockOutFailures: step.DefaultPasswordLockoutPolicy.ShowLockOutFailures,
		})
	}
	return r.setup(ctx, iamID, step, fn)
}
