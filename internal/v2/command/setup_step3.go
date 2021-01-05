package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
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
	return commandSide.SetupStep3(ctx, commandSide.iamID, s)
}

func (r *CommandSide) SetupStep3(ctx context.Context, iamID string, step *Step3) error {
	fn := func(iam *IAMWriteModel) (*iam_repo.Aggregate, error) {
		return r.addDefaultPasswordAgePolicy(ctx, NewIAMPasswordAgePolicyWriteModel(iam.AggregateID), &iam_model.PasswordAgePolicy{
			MaxAgeDays:     step.DefaultPasswordAgePolicy.MaxAgeDays,
			ExpireWarnDays: step.DefaultPasswordAgePolicy.ExpireWarnDays,
		})
	}
	return r.setup(ctx, iamID, step, fn)
}

func (r *CommandSide) setup(ctx context.Context, iamID string, step Step, fn func(*IAMWriteModel) (*iam_repo.Aggregate, error)) error {
	iam, err := r.iamByID(ctx, iamID)
	if err != nil && !caos_errs.IsNotFound(err) {
		return err
	}
	if iam.SetUpStarted != step.Step() && iam.SetUpDone+1 != step.Step() {

	}
	iamAgg, err := fn(iam)
	if err != nil {
		return err
	}
	iamAgg.PushEvents(iam_repo.NewSetupStepDoneEvent(ctx, step.Step()))

	_, err = r.eventstore.PushAggregates(ctx, iamAgg)
	if err != nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-dbG31", "Setup Step3 failed")
	}
	return nil
}
