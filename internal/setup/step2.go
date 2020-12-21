package setup

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/v2/business/command"
)

type Step2 struct {
	DefaultPasswordComplexityPolicy iam_model.PasswordComplexityPolicy

	setup *Setup
}

func (s *Step2) isNil() bool {
	return s == nil
}

func (step *Step2) step() iam_model.Step {
	return iam_model.Step2
}

func (step *Step2) init(setup *Setup) {
	step.setup = setup
}

func (step *Step2) execute(ctx context.Context, commands command.CommandSide) error {
	//commands.SetupStep2(ctx, )
	//iam, agg, err := step.passwordComplexityPolicy(ctx, &step.DefaultPasswordComplexityPolicy)
	//if err != nil {
	//	logging.Log("SETUP-Ms9fl").WithField("step", step.step()).WithError(err).Error("unable to finish setup (pw complexity policy)")
	//	return nil, err
	//}
	//iam, agg, push, err := step.setup.IamEvents.PrepareSetupDone(ctx, iam, agg, step.step())
	//if err != nil {
	//	logging.Log("SETUP-V8sui").WithField("step", step.step()).WithError(err).Error("unable to finish setup (prepare setup done)")
	//	return nil, err
	//}
	//err = es_sdk.PushAggregates(ctx, push, iam.AppendEvents, agg)
	//if err != nil {
	//	logging.Log("SETUP-V8sui").WithField("step", step.step()).WithError(err).Error("unable to finish setup")
	//	return nil, err
	//}
	//return iam_es_model.IAMToModel(iam), nil
	return nil
}

func (step *Step2) passwordComplexityPolicy(ctx context.Context, policy *iam_model.PasswordComplexityPolicy) (*iam_es_model.IAM, *models.Aggregate, error) {
	logging.Log("SETUP-Bs8id").Info("setting up password complexity policy")
	policy.AggregateID = step.setup.iamID
	iam, aggregate, err := step.setup.IamEvents.PrepareAddPasswordComplexityPolicy(ctx, policy)
	if err != nil {
		return nil, nil, err
	}
	return iam, aggregate, nil
}
