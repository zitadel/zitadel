package iam

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/label"

	iam_model "github.com/caos/zitadel/internal/iam/model"
)

func (r *Repository) AddLabelPolicy(ctx context.Context, policy *iam_model.LabelPolicy) (*iam_model.LabelPolicy, error) {
	if !policy.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-5Mv0s", "Errors.IAM.LabelPolicyInvalid")
	}

	addedPolicy := label.NewWriteModel(policy.AggregateID)
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy != nil {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-2B0ps", "Errors.IAM.LabelPolicy.AlreadyExists")
	}

	iamAgg := iam_repo.AggregateFromWriteModel(&addedPolicy.Policy.WriteModel).
		PushLabelPolicyAddedEvent(ctx, policy.PrimaryColor, policy.SecondaryColor)

	err = r.eventstore.PushAggregate(ctx, addedPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToLabelPolicy(addedPolicy), nil
}

func (r *Repository) ChangeLabelPolicy(ctx context.Context, policy *iam_model.LabelPolicy) (*iam_model.LabelPolicy, error) {
	if !policy.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-6M0od", "Errors.IAM.LabelPolicyInvalid")
	}

	existingPolicy, err := r.labelPolicyWriteModelByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}

	iamAgg := iam_repo.AggregateFromWriteModel(&existingPolicy.Policy.WriteModel).
		PushLabelPolicyChangedFromExisting(ctx, existingPolicy, policy.PrimaryColor, policy.SecondaryColor)

	err = r.eventstore.PushAggregate(ctx, existingPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToLabelPolicy(existingPolicy), nil
}

func (r *Repository) labelPolicyWriteModelByID(ctx context.Context, iamID string) (policy *label.WriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := label.NewWriteModel(iamID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
