package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"

	"github.com/caos/logging"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
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
	fn := func(iam *IAMWriteModel) ([]eventstore.EventPusher, error) {
		iamAgg := IAMAggregateFromWriteModel(&iam.WriteModel)
		event, err := r.addDefaultOrgIAMPolicy(ctx, iamAgg, NewIAMOrgIAMPolicyWriteModel(), &domain.OrgIAMPolicy{
			UserLoginMustBeDomain: step.DefaultOrgIAMPolicy.UserLoginMustBeDomain,
		})
		if err != nil {
			return nil, err
		}
		logging.Log("SETUP-ADgd2").Info("default org iam policy set up")
		return []eventstore.EventPusher{event}, nil
	}
	return r.setup(ctx, step, fn)
}
