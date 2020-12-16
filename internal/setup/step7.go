package setup

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

type Step7 struct {
	DefaultSecondFactor iam_model.SecondFactorType

	setup *Setup
}

func (step *Step7) isNil() bool {
	return step == nil
}

func (step *Step7) step() iam_model.Step {
	return iam_model.Step7
}

func (step *Step7) init(setup *Setup) {
	step.setup = setup
}

func (step *Step7) execute(ctx context.Context) (*iam_model.IAM, error) {
	iam, agg, err := step.add2FAToPolicy(ctx, step.DefaultSecondFactor)
	if err != nil {
		logging.Log("SETUP-GBD32").WithField("step", step.step()).WithError(err).Error("unable to finish setup (add default mfa to login policy)")
		return nil, err
	}
	iam, agg, push, err := step.setup.IamEvents.PrepareSetupDone(ctx, iam, agg, step.step())
	if err != nil {
		logging.Log("SETUP-BHrth").WithField("step", step.step()).WithError(err).Error("unable to finish setup (prepare setup done)")
		return nil, err
	}
	err = es_sdk.PushAggregates(ctx, push, iam.AppendEvents, agg)
	if err != nil {
		logging.Log("SETUP-k2fla").WithField("step", step.step()).WithError(err).Error("unable to finish setup")
		return nil, err
	}
	return iam_es_model.IAMToModel(iam), nil
}

func (step *Step7) add2FAToPolicy(ctx context.Context, secondFactor iam_model.SecondFactorType) (*iam_es_model.IAM, *models.Aggregate, error) {
	logging.Log("SETUP-Bew1a").Info("adding 2FA to loginPolicy")
	return step.setup.IamEvents.PrepareAddSecondFactorToLoginPolicy(ctx, step.setup.iamID, secondFactor)
}
