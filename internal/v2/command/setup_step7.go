package command

import (
	"context"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

type Step7 struct {
	OTP bool
}

func (s *Step7) Step() domain.Step {
	return domain.Step7
}

func (s *Step7) execute(ctx context.Context, commandSide *CommandSide) error {
	return commandSide.SetupStep7(ctx, s)
}

func (r *CommandSide) SetupStep7(ctx context.Context, step *Step7) error {
	fn := func(iam *IAMWriteModel) (*iam_repo.Aggregate, error) {
		secondFactorModel := NewIAMSecondFactorWriteModel()
		iamAgg := IAMAggregateFromWriteModel(&secondFactorModel.SecondFactorWriteModel.WriteModel)
		if !step.OTP {
			return iamAgg, nil
		}
		err := r.addSecondFactorToDefaultLoginPolicy(ctx, iamAgg, secondFactorModel, iam_model.SecondFactorTypeOTP)
		if err != nil {
			return nil, err
		}
		return iamAgg, nil
	}
	return r.setup(ctx, step, fn)
}
