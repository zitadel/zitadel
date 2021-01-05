package setup

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

type Step9 struct {
	Passwordless bool

	setup *Setup
}

func (step *Step9) isNil() bool {
	return step == nil
}

func (step *Step9) step() iam_model.Step {
	return iam_model.Step9
}

func (step *Step9) init(setup *Setup) {
	step.setup = setup
}

func (step *Step9) execute(ctx context.Context) (*iam_model.IAM, error) {
	if !step.Passwordless {
		return step.setup.IamEvents.IAMByID(ctx, step.setup.iamID)
	}
	iam, agg, err := step.setPasswordlessAllowedInPolicy(ctx)
	if err != nil {
		logging.Log("SETUP-Gdbjq").WithField("step", step.step()).WithError(err).Error("unable to finish setup (add default mfa to login policy)")
		return nil, err
	}
	iam, agg2, err := step.addMFAToPolicy(ctx)
	if err != nil {
		logging.Log("SETUP-Gdbjq").WithField("step", step.step()).WithError(err).Error("unable to finish setup (add default mfa to login policy)")
		return nil, err
	}
	agg.Events = append(agg.Events, agg2.Events...)
	iam, agg, push, err := step.setup.IamEvents.PrepareSetupDone(ctx, iam, agg, step.step())
	if err != nil {
		logging.Log("SETUP-Cnf21").WithField("step", step.step()).WithError(err).Error("unable to finish setup (prepare setup done)")
		return nil, err
	}
	err = es_sdk.PushAggregates(ctx, push, iam.AppendEvents, agg)
	if err != nil {
		logging.Log("SETUP-NFq21").WithField("step", step.step()).WithError(err).Error("unable to finish setup")
		return nil, err
	}
	return iam_es_model.IAMToModel(iam), nil
}

func (step *Step9) setPasswordlessAllowedInPolicy(ctx context.Context) (*iam_es_model.IAM, *models.Aggregate, error) {
	logging.Log("SETUP-DAd1h").Info("enabling passwordless in loginPolicy")
	iam, err := step.setup.IamEvents.IAMByID(ctx, step.setup.iamID)
	if err != nil {
		return nil, nil, err
	}
	iam.DefaultLoginPolicy.AggregateID = step.setup.iamID
	iam.DefaultLoginPolicy.PasswordlessType = iam_model.PasswordlessTypeAllowed
	return step.setup.IamEvents.PrepareChangeLoginPolicy(ctx, iam.DefaultLoginPolicy)
}

func (step *Step9) addMFAToPolicy(ctx context.Context) (*iam_es_model.IAM, *models.Aggregate, error) {
	logging.Log("SETUP-DAd1h").Info("adding MFA to loginPolicy")
	return step.setup.IamEvents.PrepareAddMultiFactorToLoginPolicy(ctx, step.setup.iamID, iam_model.MultiFactorTypeU2FWithPIN)
}
