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
	return commandSide.SetupStep6(ctx, commandSide.iamID, s)
}

func (r *CommandSide) SetupStep6(ctx context.Context, iamID string, step *Step6) error {
	fn := func(iam *IAMWriteModel) (*iam_repo.Aggregate, error) {
		return r.addDefaultLabelPolicy(ctx, NewIAMLabelPolicyWriteModel(iam.AggregateID), &iam_model.LabelPolicy{
			PrimaryColor:   step.DefaultLabelPolicy.PrimaryColor,
			SecondaryColor: step.DefaultLabelPolicy.SecondaryColor,
		})
	}
	return r.setup(ctx, iamID, step, fn)
}
