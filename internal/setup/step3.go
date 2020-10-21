package setup

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

type Step3 struct {
	DefaultPasswordAgePolicy iam_model.PasswordAgePolicy

	setup *Setup
}

func (s *Step3) isNil() bool {
	return s == nil
}

func (step *Step3) step() iam_model.Step {
	return iam_model.Step3
}

func (step *Step3) init(setup *Setup) {
	step.setup = setup
}

func (step *Step3) execute(ctx context.Context) (*iam_model.IAM, error) {
	iam, agg, err := step.passwordAgePolicy(ctx, &step.DefaultPasswordAgePolicy)
	if err != nil {
		logging.Log("SETUP-Mski9").WithField("step", step.step()).WithError(err).Error("unable to finish setup (pw age policy)")
		return nil, err
	}
	iam, agg, push, err := step.setup.IamEvents.PrepareSetupDone(ctx, iam, agg, step.step())
	if err != nil {
		logging.Log("SETUP-4Gsny").WithField("step", step.step()).WithError(err).Error("unable to finish setup (prepare setup done)")
		return nil, err
	}
	err = es_sdk.PushAggregates(ctx, push, iam.AppendEvents, agg)
	if err != nil {
		logging.Log("SETUP-Yc8ui").WithField("step", step.step()).WithError(err).Error("unable to finish setup")
		return nil, err
	}
	return iam_es_model.IAMToModel(iam), nil
}

func (step *Step3) passwordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) (*iam_es_model.IAM, *models.Aggregate, error) {
	logging.Log("SETUP-bVs8i").Info("setting up password complexity policy")
	policy.AggregateID = step.setup.iamID
	iam, aggregate, err := step.setup.IamEvents.PrepareAddPasswordAgePolicy(ctx, policy)
	if err != nil {
		return nil, nil, err
	}
	return iam, aggregate, nil
}
