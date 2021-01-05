package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

type Step2 struct {
	DefaultPasswordComplexityPolicy iam_model.PasswordComplexityPolicy
}

func (r *CommandSide) SetupStep2(ctx context.Context, iamID string, step Step2) error {
	iam, err := r.iamByID(ctx, iamID)
	if err != nil && !caos_errs.IsNotFound(err) {
		return err
	}
	iamAgg, err := r.addDefaultPasswordComplexityPolicy(ctx, NewIAMPasswordComplexityPolicyWriteModel(iam.AggregateID), &domain.PasswordComplexityPolicy{
		MinLength:    step.DefaultPasswordComplexityPolicy.MinLength,
		HasLowercase: step.DefaultPasswordComplexityPolicy.HasLowercase,
		HasUppercase: step.DefaultPasswordComplexityPolicy.HasUppercase,
		HasNumber:    step.DefaultPasswordComplexityPolicy.HasNumber,
		HasSymbol:    step.DefaultPasswordComplexityPolicy.HasSymbol,
	})
	if err != nil {
		return err
	}
	iamAgg.PushEvents(iam_repo.NewSetupStepDoneEvent(ctx, domain.Step1))

	_, err = r.eventstore.PushAggregates(ctx, iamAgg)
	if err != nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-HR2na", "Setup Step2 failed")
	}
	return nil
}
