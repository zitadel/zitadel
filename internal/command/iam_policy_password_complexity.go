package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (r *CommandSide) getDefaultPasswordComplexityPolicy(ctx context.Context) (*domain.PasswordComplexityPolicy, error) {
	policyWriteModel := NewIAMPasswordComplexityPolicyWriteModel()
	err := r.eventstore.FilterToQueryReducer(ctx, policyWriteModel)
	if err != nil {
		return nil, err
	}
	policy := writeModelToPasswordComplexityPolicy(&policyWriteModel.PasswordComplexityPolicyWriteModel)
	policy.Default = true
	return policy, nil
}

func (r *CommandSide) AddDefaultPasswordComplexityPolicy(ctx context.Context, policy *domain.PasswordComplexityPolicy) (*domain.PasswordComplexityPolicy, error) {
	addedPolicy := NewIAMPasswordComplexityPolicyWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&addedPolicy.WriteModel)
	events, err := r.addDefaultPasswordComplexityPolicy(ctx, iamAgg, addedPolicy, policy)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := r.eventstore.PushEvents(ctx, events)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToPasswordComplexityPolicy(&addedPolicy.PasswordComplexityPolicyWriteModel), nil
}

func (r *CommandSide) addDefaultPasswordComplexityPolicy(ctx context.Context, iamAgg *eventstore.Aggregate, addedPolicy *IAMPasswordComplexityPolicyWriteModel, policy *domain.PasswordComplexityPolicy) (eventstore.EventPusher, error) {
	if err := policy.IsValid(); err != nil {
		return nil, err
	}

	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-Lk0dS", "Errors.IAM.PasswordComplexityPolicy.AlreadyExists")
	}

	return iam_repo.NewPasswordComplexityPolicyAddedEvent(ctx, iamAgg, policy.MinLength, policy.HasLowercase, policy.HasUppercase, policy.HasNumber, policy.HasSymbol), nil
}

func (r *CommandSide) ChangeDefaultPasswordComplexityPolicy(ctx context.Context, policy *domain.PasswordComplexityPolicy) (*domain.PasswordComplexityPolicy, error) {
	if err := policy.IsValid(); err != nil {
		return nil, err
	}

	existingPolicy, err := r.defaultPasswordComplexityPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-0oPew", "Errors.IAM.PasswordAgePolicy.NotFound")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.PasswordComplexityPolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, iamAgg, policy.MinLength, policy.HasLowercase, policy.HasUppercase, policy.HasNumber, policy.HasSymbol)
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
	return writeModelToPasswordComplexityPolicy(&existingPolicy.PasswordComplexityPolicyWriteModel), nil
}

func (r *CommandSide) defaultPasswordComplexityPolicyWriteModelByID(ctx context.Context) (policy *IAMPasswordComplexityPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMPasswordComplexityPolicyWriteModel()
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
