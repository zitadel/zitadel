package command

import (
	"context"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

type Step5 struct {
	DefaultOrgIAMPolicy iam_model.OrgIAMPolicy
}

func (s *Step5) Step() domain.Step {
	return domain.Step5
}

func (s *Step5) execute(ctx context.Context, commandSide *CommandSide) error {
	return commandSide.SetupStep5(ctx, s)
}

func (r *CommandSide) SetupStep5(ctx context.Context, step *Step5) error {
	fn := func(iam *IAMWriteModel) (*iam_repo.Aggregate, error) {
		iamAgg := IAMAggregateFromWriteModel(&iam.WriteModel)
		err := r.addDefaultOrgIAMPolicy(ctx, iamAgg, NewIAMOrgIAMPolicyWriteModel(iam.AggregateID), &domain.OrgIAMPolicy{
			UserLoginMustBeDomain: step.DefaultOrgIAMPolicy.UserLoginMustBeDomain,
		})
		if err != nil {
			return nil, err
		}
		return iamAgg, nil
	}
	return r.setup(ctx, step, fn)
}
