package command

import (
	"context"

	"github.com/caos/logging"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

type Step3 struct {
	DefaultPasswordAgePolicy iam_model.PasswordAgePolicy
}

func (s *Step3) Step() domain.Step {
	return domain.Step3
}

func (s *Step3) execute(ctx context.Context, commandSide *CommandSide) error {
	return commandSide.SetupStep3(ctx, s)
}

func (r *CommandSide) SetupStep3(ctx context.Context, step *Step3) error {
	fn := func(iam *IAMWriteModel) (*iam_repo.Aggregate, error) {
		iamAgg := IAMAggregateFromWriteModel(&iam.WriteModel)
		err := r.addDefaultPasswordAgePolicy(ctx, iamAgg, NewIAMPasswordAgePolicyWriteModel(), &domain.PasswordAgePolicy{
			MaxAgeDays:     step.DefaultPasswordAgePolicy.MaxAgeDays,
			ExpireWarnDays: step.DefaultPasswordAgePolicy.ExpireWarnDays,
		})
		if err != nil {
			return nil, err
		}
		logging.Log("SETUP-DBqgq").Info("default password age policy set up")
		return iamAgg, nil
	}
	return r.setup(ctx, step, fn)
}
