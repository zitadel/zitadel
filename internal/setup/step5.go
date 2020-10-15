package setup

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

type Step5 struct {
	DefaultOrgIAMPolicy iam_model.OrgIAMPolicy

	setup *Setup
}

func (s *Step5) isNil() bool {
	return s == nil
}

func (step *Step5) step() iam_model.Step {
	return iam_model.Step5
}

func (step *Step5) init(setup *Setup) {
	step.setup = setup
}

func (step *Step5) execute(ctx context.Context) (*iam_model.IAM, error) {
	iam, agg, err := step.orgIAMPolicy(ctx, &step.DefaultOrgIAMPolicy)
	if err != nil {
		logging.Log("SETUP-3nKd9").WithField("step", step.step()).WithError(err).Error("unable to finish setup (org iam policy)")
		return nil, err
	}
	iam, agg, push, err := step.setup.IamEvents.PrepareSetupDone(ctx, iam, agg, step.step())
	if err != nil {
		logging.Log("SETUP-5h8Ds").WithField("step", step.step()).WithError(err).Error("unable to finish setup (prepare setup done)")
		return nil, err
	}
	err = es_sdk.PushAggregates(ctx, push, iam.AppendEvents, agg)
	if err != nil {
		logging.Log("SETUP-3fGk0").WithField("step", step.step()).WithError(err).Error("unable to finish setup")
		return nil, err
	}
	return iam_es_model.IAMToModel(iam), nil
}

func (step *Step5) orgIAMPolicy(ctx context.Context, policy *iam_model.OrgIAMPolicy) (*iam_es_model.IAM, *models.Aggregate, error) {
	logging.Log("SETUP-5Gn8s").Info("setting up org iam policy")
	policy.AggregateID = step.setup.iamID
	iam, aggregate, err := step.setup.IamEvents.PrepareAddOrgIAMPolicy(ctx, policy)
	if err != nil {
		return nil, nil, err
	}
	return iam, aggregate, nil
}
