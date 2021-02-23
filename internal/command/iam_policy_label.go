package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (r *CommandSide) AddDefaultLabelPolicy(ctx context.Context, policy *domain.LabelPolicy) (*domain.LabelPolicy, error) {
	addedPolicy := NewIAMLabelPolicyWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&addedPolicy.LabelPolicyWriteModel.WriteModel)
	event, err := r.addDefaultLabelPolicy(ctx, iamAgg, addedPolicy, policy)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := r.eventstore.PushEvents(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToLabelPolicy(&addedPolicy.LabelPolicyWriteModel), nil
}

func (r *CommandSide) addDefaultLabelPolicy(ctx context.Context, iamAgg *eventstore.Aggregate, addedPolicy *IAMLabelPolicyWriteModel, policy *domain.LabelPolicy) (eventstore.EventPusher, error) {
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-2B0ps", "Errors.IAM.LabelPolicy.AlreadyExists")
	}

	return iam_repo.NewLabelPolicyAddedEvent(ctx, iamAgg, policy.PrimaryColor, policy.SecondaryColor), nil

}

func (r *CommandSide) ChangeDefaultLabelPolicy(ctx context.Context, policy *domain.LabelPolicy) (*domain.LabelPolicy, error) {
	existingPolicy, err := r.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-0K9dq", "Errors.IAM.LabelPolicy.NotFound")
	}
	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, iamAgg, policy.PrimaryColor, policy.SecondaryColor)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-4M9vs", "Errors.IAM.LabelPolicy.NotChanged")
	}

	pushedEvents, err := r.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToLabelPolicy(&existingPolicy.LabelPolicyWriteModel), nil
}

func (r *CommandSide) defaultLabelPolicyWriteModelByID(ctx context.Context) (policy *IAMLabelPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMLabelPolicyWriteModel()
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
