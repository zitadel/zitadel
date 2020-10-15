package setup

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

type Step4 struct {
	DefaultPasswordLockoutPolicy iam_model.PasswordLockoutPolicy

	setup *Setup
}

func (s *Step4) isNil() bool {
	return s == nil
}

func (step *Step4) step() iam_model.Step {
	return iam_model.Step4
}

func (step *Step4) init(setup *Setup) {
	step.setup = setup
}

func (step *Step4) execute(ctx context.Context) (*iam_model.IAM, error) {
	iam, agg, err := step.passwordLockoutPolicy(ctx, &step.DefaultPasswordLockoutPolicy)
	if err != nil {
		logging.Log("SETUP-xCd9i").WithField("step", step.step()).WithError(err).Error("unable to finish setup (pw age policy)")
		return nil, err
	}
	iam, agg, push, err := step.setup.IamEvents.PrepareSetupDone(ctx, iam, agg, step.step())
	if err != nil {
		logging.Log("SETUP-bVsm9").WithField("step", step.step()).WithError(err).Error("unable to finish setup (prepare setup done)")
		return nil, err
	}
	err = es_sdk.PushAggregates(ctx, push, iam.AppendEvents, agg)
	if err != nil {
		logging.Log("SETUP-wCxko").WithField("step", step.step()).WithError(err).Error("unable to finish setup")
		return nil, err
	}
	return iam_es_model.IAMToModel(iam), nil
}

func (step *Step4) passwordLockoutPolicy(ctx context.Context, policy *iam_model.PasswordLockoutPolicy) (*iam_es_model.IAM, *models.Aggregate, error) {
	logging.Log("SETUP-vSfr4").Info("setting up password complexity policy")
	policy.AggregateID = step.setup.iamID
	iam, aggregate, err := step.setup.IamEvents.PrepareAddPasswordLockoutPolicy(ctx, policy)
	if err != nil {
		return nil, nil, err
	}
	return iam, aggregate, nil
}
