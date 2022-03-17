package command

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
)

type Step12 struct {
	TierName                 string
	TierDescription          string
	AuditLogRetention        time.Duration
	LoginPolicyFactors       bool
	LoginPolicyIDP           bool
	LoginPolicyPasswordless  bool
	LoginPolicyRegistration  bool
	LoginPolicyUsernameLogin bool
	PasswordComplexityPolicy bool
	LabelPolicy              bool
	CustomDomain             bool
}

func (s *Step12) Step() domain.Step {
	return domain.Step12
}

func (s *Step12) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep12(ctx, s)
}

func (c *Commands) SetupStep12(ctx context.Context, step *Step12) error {
	fn := func(iam *InstanceWriteModel) ([]eventstore.Command, error) {
		featuresWriteModel := NewInstanceFeaturesWriteModel()
		featuresEvent, err := c.setDefaultFeatures(ctx, featuresWriteModel, &domain.Features{
			TierName:                 step.TierName,
			TierDescription:          step.TierDescription,
			State:                    domain.FeaturesStateActive,
			AuditLogRetention:        step.AuditLogRetention,
			LoginPolicyFactors:       step.LoginPolicyFactors,
			LoginPolicyIDP:           step.LoginPolicyIDP,
			LoginPolicyPasswordless:  step.LoginPolicyPasswordless,
			LoginPolicyRegistration:  step.LoginPolicyRegistration,
			LoginPolicyUsernameLogin: step.LoginPolicyUsernameLogin,
			PasswordComplexityPolicy: step.PasswordComplexityPolicy,
			LabelPolicyPrivateLabel:  step.LabelPolicy,
			CustomDomain:             step.CustomDomain,
		})
		if err != nil {
			return nil, err
		}
		return []eventstore.Command{featuresEvent}, nil
	}
	return c.setup(ctx, step, fn)
}
