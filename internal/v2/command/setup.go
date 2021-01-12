package command

import (
	"context"

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
	aggregate := IAMAggregateFromWriteModel(&iamWriteModel.WriteModel).PushEvents(iam_repo.NewSetupStepStartedEvent(ctx, step))
	err = r.eventstore.PushAggregate(ctx, iamWriteModel, aggregate)
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Grgh1", "Setup start failed")
	}
	return writeModelToIAM(iamWriteModel), nil
}

func (r *CommandSide) setup(ctx context.Context, step Step, iamAggregateProvider func(*IAMWriteModel) (*iam_repo.Aggregate, error)) error {
	iam, err := r.getIAMWriteModel(ctx)
	if err != nil && !caos_errs.IsNotFound(err) {
		return err
	}
	if iam.SetUpStarted != step.Step() && iam.SetUpDone+1 != step.Step() {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-Dge32", "wrong step")
	}
	iamAgg, err := iamAggregateProvider(iam)
	if err != nil {
		return err
	}
	iamAgg.PushEvents(iam_repo.NewSetupStepDoneEvent(ctx, step.Step()))

	_, err = r.eventstore.PushAggregates(ctx, iamAgg)
	if err != nil {
		return caos_errs.ThrowPreconditionFailedf(nil, "EVENT-dbG31", "Setup %s failed", step.Step())
	}
	return nil
}
