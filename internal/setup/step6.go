package setup

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

type Step6 struct {
	DefaultLabelPolicy iam_model.LabelPolicy

	setup *Setup
}

func (s *Step6) isNil() bool {
	return s == nil
}

func (step *Step6) step() iam_model.Step {
	return iam_model.Step6
}

func (step *Step6) init(setup *Setup) {
	step.setup = setup
}

func (step *Step6) execute(ctx context.Context) (*iam_model.IAM, error) {
	iam, agg, err := step.labelPolicy(ctx, &step.DefaultLabelPolicy)
	if err != nil {
		logging.Log("SETUP-ZTuS1").WithField("step", step.step()).WithError(err).Error("unable to finish setup (Label policy)")
		return nil, err
	}
	iam, agg, push, err := step.setup.IamEvents.PrepareSetupDone(ctx, iam, agg, step.step())
	if err != nil {
		logging.Log("SETUP-OkF8o").WithField("step", step.step()).WithError(err).Error("unable to finish setup (prepare setup done)")
		return nil, err
	}
	err = es_sdk.PushAggregates(ctx, push, iam.AppendEvents, agg)
	if err != nil {
		logging.Log("SETUP-YbQ6T").WithField("step", step.step()).WithError(err).Error("unable to finish setup")
		return nil, err
	}
	return iam_es_model.IAMToModel(iam), nil
}

func (step *Step6) labelPolicy(ctx context.Context, policy *iam_model.LabelPolicy) (*iam_es_model.IAM, *models.Aggregate, error) {
	logging.Log("SETUP-geMuZ").Info("setting up labelpolicy")
	policy.AggregateID = step.setup.iamID
	iam, aggregate, err := step.setup.IamEvents.PrepareAddLabelPolicy(ctx, policy)
	if err != nil {
		return nil, nil, err
	}
	return iam, aggregate, nil
}
