package command

import (
	"context"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

type Step2 struct {
	DefaultPasswordComplexityPolicy iam_model.PasswordComplexityPolicy
}

func (s *Step2) Step() domain.Step {
	return domain.Step2
}

func (s *Step2) execute(ctx context.Context, commandSide *CommandSide) error {
	return commandSide.SetupStep2(ctx, s)
}

func (r *CommandSide) SetupStep2(ctx context.Context, step *Step2) error {
	fn := func(iam *IAMWriteModel) (*iam_repo.Aggregate, error) {
		iamAgg := IAMAggregateFromWriteModel(&iam.WriteModel)
		err := r.addDefaultPasswordComplexityPolicy(ctx, iamAgg, NewIAMPasswordComplexityPolicyWriteModel(iam.AggregateID), &domain.PasswordComplexityPolicy{
			MinLength:    step.DefaultPasswordComplexityPolicy.MinLength,
			HasLowercase: step.DefaultPasswordComplexityPolicy.HasLowercase,
			HasUppercase: step.DefaultPasswordComplexityPolicy.HasUppercase,
			HasNumber:    step.DefaultPasswordComplexityPolicy.HasNumber,
			HasSymbol:    step.DefaultPasswordComplexityPolicy.HasSymbol,
		})
		if err != nil {
			return nil, err
		}
		return iamAgg, err
	}
	return r.setup(ctx, step, fn)
}
