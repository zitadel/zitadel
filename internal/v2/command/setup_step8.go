package command

import (
	"context"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

type Step8 struct {
	U2F bool
}

func (s *Step8) Step() domain.Step {
	return domain.Step8
}

func (s *Step8) execute(ctx context.Context, commandSide *CommandSide) error {
	return commandSide.SetupStep8(ctx, commandSide.iamID, s)
}

func (r *CommandSide) SetupStep8(ctx context.Context, iamID string, step *Step8) error {
	fn := func(iam *IAMWriteModel) (*iam_repo.Aggregate, error) {
		secondFactorModel := NewIAMSecondFactorWriteModel(iam.AggregateID)
		if step.U2F {
			return r.addSecondFactorToDefaultLoginPolicy(ctx, secondFactorModel, iam_model.SecondFactorTypeU2F)
		}
		return IAMAggregateFromWriteModel(&secondFactorModel.SecondFactorWriteModel.WriteModel), nil
	}
	return r.setup(ctx, iamID, step, fn)
}
