package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/api/authz"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

type Step interface {
	Step() domain.Step
	execute(context.Context, *CommandSide) error
}

const (
	SetupUser = "SETUP"
)

func (r *CommandSide) ExecuteSetupSteps(ctx context.Context, steps []Step) error {
	iam, err := r.GetIAM(ctx)
	if err != nil && !caos_errs.IsNotFound(err) {
		return err
	}
	if iam != nil && (iam.SetUpDone == domain.StepCount-1 || iam.SetUpStarted != iam.SetUpDone) {
		logging.Log("COMMA-dgd2z").Info("all steps done")
		return nil
	}

	if iam == nil {
		iam = &domain.IAM{ObjectRoot: models.ObjectRoot{}}
	}

	ctx = setSetUpContextData(ctx)

	for _, step := range steps {
		iam, err = r.StartSetup(ctx, step.Step())
		if err != nil {
			return err
		}

		err = step.execute(ctx, r)
		if err != nil {
			return err
		}
	}
	return nil
}

func setSetUpContextData(ctx context.Context) context.Context {
	return authz.SetCtxData(ctx, authz.CtxData{UserID: SetupUser})
}

func (r *CommandSide) StartSetup(ctx context.Context, step domain.Step) (*domain.IAM, error) {
	iamWriteModel, err := r.getIAMWriteModel(ctx)
	if err != nil && !caos_errs.IsNotFound(err) {
		return nil, err
	}
	if iamWriteModel.SetUpStarted >= step || iamWriteModel.SetUpStarted != iamWriteModel.SetUpDone {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9so34", "setup error")
	}
	aggregate := IAMAggregateFromWriteModel(&iamWriteModel.WriteModel)
	pushedEvents, err := r.eventstore.PushEvents(ctx, iam_repo.NewSetupStepStartedEvent(ctx, aggregate, step))
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Grgh1", "Setup start failed")
	}
	err = AppendAndReduce(iamWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	logging.LogWithFields("SETUP-fhh21", "step", step).Info("setup step started")
	return writeModelToIAM(iamWriteModel), nil
}

func (r *CommandSide) setup(ctx context.Context, step Step, iamAggregateProvider func(*IAMWriteModel) ([]eventstore.EventPusher, error)) error {
	iam, err := r.getIAMWriteModel(ctx)
	if err != nil && !caos_errs.IsNotFound(err) {
		return err
	}
	if iam.SetUpStarted != step.Step() && iam.SetUpDone+1 != step.Step() {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-Dge32", "wrong step")
	}
	events, err := iamAggregateProvider(iam)
	if err != nil {
		return err
	}
	iamAgg := IAMAggregateFromWriteModel(&iam.WriteModel)
	events = append(events, iam_repo.NewSetupStepDoneEvent(ctx, iamAgg, step.Step()))

	_, err = r.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return caos_errs.ThrowPreconditionFailedf(nil, "EVENT-dbG31", "Setup %v failed", step.Step())
	}
	logging.LogWithFields("SETUP-Sg1t1", "step", step.Step()).Info("setup step done")
	return nil
}
