package command

import (
	"context"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

type Step6 struct {
	DefaultLabelPolicy iam_model.LabelPolicy
}

func (s *Step6) Step() domain.Step {
	return domain.Step6
}

func (s *Step6) execute(ctx context.Context, commandSide *CommandSide) error {
	return commandSide.SetupStep6(ctx, s)
}

func (r *CommandSide) SetupStep6(ctx context.Context, step *Step6) error {
	fn := func(iam *IAMWriteModel) (*iam_repo.Aggregate, error) {
		iamAgg := IAMAggregateFromWriteModel(&iam.WriteModel)
		err := r.addDefaultLabelPolicy(ctx, iamAgg, NewIAMLabelPolicyWriteModel(iam.AggregateID), &domain.LabelPolicy{
			PrimaryColor:   step.DefaultLabelPolicy.PrimaryColor,
			SecondaryColor: step.DefaultLabelPolicy.SecondaryColor,
		})
		if err != nil {
			return nil, err
		}
		return iamAgg, nil
	}
	return r.setup(ctx, step, fn)
}
