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
	return commandSide.SetupStep7(ctx, commandSide.iamID, s)
}

func (r *CommandSide) SetupStep7(ctx context.Context, iamID string, step *Step7) error {
	fn := func(iam *IAMWriteModel) (*iam_repo.Aggregate, error) {
		secondFactorModel := NewIAMSecondFactorWriteModel(iam.AggregateID)
		if step.OTP {
			return r.addSecondFactorToDefaultLoginPolicy(ctx, secondFactorModel, iam_model.SecondFactorTypeOTP)
		}
		return IAMAggregateFromWriteModel(&secondFactorModel.SecondFactorWriteModel.WriteModel), nil
	}
	return r.setup(ctx, iamID, step, fn)
}
