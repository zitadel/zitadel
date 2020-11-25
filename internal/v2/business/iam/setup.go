package iam

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *Repository) StartSetup(ctx context.Context, iamID string, step iam_model.Step) (*iam_model.IAM, error) {
	iam, err := r.setup(ctx, iamID, iam_repo.Step(step), iam_repo.NewSetupStepStartedEvent(ctx, iam_repo.Step(step)))
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-zx03n", "Setup start failed")
	}
	return iam, nil
}

func (r *Repository) SetupDone(ctx context.Context, iamID string, step iam_model.Step) (*iam_model.IAM, error) {
	iam, err := r.setup(ctx, iamID, iam_repo.Step(step), iam_repo.NewSetupStepDoneEvent(ctx, iam_repo.Step(step)))
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-zx03n", "Setup start failed")
	}
	return iam, nil
}

func (r *Repository) setup(ctx context.Context, iamID string, step iam_repo.Step, event eventstore.EventPusher) (*iam_model.IAM, error) {
	iam, err := r.iamByID(ctx, iamID)
	if err != nil && !caos_errs.IsNotFound(err) {
		return nil, err
	}

	if iam != nil && (iam.SetUpStarted >= iam_repo.Step(step) || iam.SetUpStarted != iam.SetUpDone) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9so34", "setup error")
	}

	aggregate := iam_repo.AggregateFromReadModel(iam).
		PushEvents(event)

	events, err := r.eventstore.PushAggregates(ctx, aggregate)
	if err != nil {
		return nil, err
	}

	if err = iam.AppendAndReduce(events...); err != nil {
		return nil, err
	}
	return readModelToIAM(iam), nil
}
