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
	return commandSide.SetupStep5(ctx, commandSide.iamID, s)
}

func (r *CommandSide) SetupStep5(ctx context.Context, iamID string, step *Step5) error {
	fn := func(iam *IAMWriteModel) (*iam_repo.Aggregate, error) {
		return r.addDefaultOrgIAMPolicy(ctx, NewIAMOrgIAMPolicyWriteModel(iam.AggregateID), &iam_model.OrgIAMPolicy{
			UserLoginMustBeDomain: step.DefaultOrgIAMPolicy.UserLoginMustBeDomain,
		})
	}
	return r.setup(ctx, iamID, step, fn)
}
