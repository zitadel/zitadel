package command

import (
	"context"

	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
)

type Step12 struct {
	TierName                 string
	TierDescription          string
	AuditLogRetention        types.Duration
	LoginPolicyFactors       bool
	LoginPolicyIDP           bool
	LoginPolicyPasswordless  bool
	LoginPolicyRegistration  bool
	LoginPolicyUsernameLogin bool
	PasswordComplexityPolicy bool
}

func (s *Step12) Step() domain.Step {
	return domain.Step12
}

func (s *Step12) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep12(ctx, s)
}

func (c *Commands) SetupStep12(ctx context.Context, step *Step12) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.EventPusher, error) {
		featuresWriteModel := NewIAMFeaturesWriteModel()
		featuresEvent, err := c.setDefaultFeatures(ctx, featuresWriteModel, &domain.Features{
			TierName:                 step.TierName,
			TierDescription:          step.TierDescription,
			TierState:                domain.FeaturesStateActive,
			AuditLogRetention:        step.AuditLogRetention.Duration,
			LoginPolicyFactors:       step.LoginPolicyFactors,
			LoginPolicyIDP:           step.LoginPolicyIDP,
			LoginPolicyPasswordless:  step.LoginPolicyPasswordless,
			LoginPolicyRegistration:  step.LoginPolicyRegistration,
			LoginPolicyUsernameLogin: step.LoginPolicyUsernameLogin,
			PasswordComplexityPolicy: step.PasswordComplexityPolicy,
		})
		if err != nil {
			return nil, err
		}
		return []eventstore.EventPusher{featuresEvent}, nil
	}
	return c.setup(ctx, step, fn)
}
