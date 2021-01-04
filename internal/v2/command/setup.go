package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *CommandSide) StartSetup(ctx context.Context, iamID string, step domain.Step) (*iam_model.IAM, error) {
	iamWriteModel, err := r.iamByID(ctx, iamID)
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

//func (r *CommandSide) setupDone(ctx context.Context, iamAgg *iam_repo.Aggregate, event eventstore.EventPusher, aggregates ...eventstore.Aggregater) error {
//	aggregate := iamAgg.PushEvents(event)
//
//	aggregates = append(aggregates, aggregate)
//	_, err := r.eventstore.PushAggregates(ctx, aggregates...)
//	if err != nil {
//		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-Dgd2", "Setup done failed")
//	}
//	return nil
//}

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
