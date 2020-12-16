package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/business/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *CommandSide) StartSetup(ctx context.Context, iamID string, step domain.Step) (*iam_model.IAM, error) {
	iamWriteModel, err := r.iamByID(ctx, iamID)
	if err != nil && !caos_errs.IsNotFound(err) {
		return nil, err
	}
	iam, err := r.setup(ctx, iamWriteModel, step, iam_repo.NewSetupStepStartedEvent(ctx, step))
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-zx03n", "Setup start failed")
	}
	return iam, nil
}

func (r *CommandSide) setup(ctx context.Context, iam *IAMWriteModel, step domain.Step, event eventstore.EventPusher) (*iam_model.IAM, error) {
	if iam != nil && (iam.SetUpStarted >= step || iam.SetUpStarted != iam.SetUpDone) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9so34", "setup error")
	}

	aggregate := AggregateFromWriteModel(&iam.WriteModel).PushEvents(event)

	err := r.eventstore.PushAggregate(ctx, iam, aggregate)
	if err != nil {
		return nil, err
	}
	return writeModelToIAM(iam), nil
}

//
////TODO: should not use readmodel
//func (r *CommandSide) setup(ctx context.Context, iamID string, step iam_repo.Step, event eventstore.EventPusher) (*iam_model.IAM, error) {
//	iam, err := r.iamByID(ctx, iamID)
//	if err != nil && !caos_errs.IsNotFound(err) {
//		return nil, err
//	}
//
//	if iam != nil && (iam.SetUpStarted >= iam_repo.Step(step) || iam.SetUpStarted != iam.SetUpDone) {
//		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9so34", "setup error")
//	}
//
//	aggregate := query.AggregateFromReadModel(iam).
//		PushEvents(event)
//
//	events, err := r.eventstore.PushAggregates(ctx, aggregate)
//	if err != nil {
//		return nil, err
//	}
//
//	if err = iam.AppendAndReduce(events...); err != nil {
//		return nil, err
//	}
//	return nil, nil
//	//TODO: return write model
//	//return readModelToIAM(iam), nil
//}
