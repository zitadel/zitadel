package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) getDefaultPrivacyPolicy(ctx context.Context) (*domain.PrivacyPolicy, error) {
	policyWriteModel := NewIAMPrivacyPolicyWriteModel()
	err := c.eventstore.FilterToQueryReducer(ctx, policyWriteModel)
	if err != nil {
		return nil, err
	}
	if !policyWriteModel.State.Exists() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-559os", "Errors.IAM.OrgIAMPolicy.NotFound")
	}
	policy := writeModelToPrivacyPolicy(&policyWriteModel.PrivacyPolicyWriteModel)
	policy.Default = true
	return policy, nil
}

func (c *Commands) AddDefaultPrivacyPolicy(ctx context.Context, policy *domain.PrivacyPolicy) (*domain.PrivacyPolicy, error) {
	addedPolicy := NewIAMPrivacyPolicyWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&addedPolicy.WriteModel)
	events, err := c.addDefaultPrivacyPolicy(ctx, iamAgg, addedPolicy, policy)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, events)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToPrivacyPolicy(&addedPolicy.PrivacyPolicyWriteModel), nil
}

func (c *Commands) addDefaultPrivacyPolicy(ctx context.Context, iamAgg *eventstore.Aggregate, addedPolicy *IAMPrivacyPolicyWriteModel, policy *domain.PrivacyPolicy) (eventstore.EventPusher, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-M00rJ", "Errors.IAM.PrivacyPolicy.AlreadyExists")
	}

	return iam_repo.NewPrivacyPolicyAddedEvent(ctx, iamAgg, policy.TOSLink, policy.PrivacyLink), nil
}

func (c *Commands) ChangeDefaultPrivacyPolicy(ctx context.Context, policy *domain.PrivacyPolicy) (*domain.PrivacyPolicy, error) {
	existingPolicy, err := c.defaultPrivacyPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-0oPew", "Errors.IAM.PasswordAgePolicy.NotFound")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.PrivacyPolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, iamAgg, policy.TOSLink, policy.PrivacyLink)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-4M9vs", "Errors.IAM.LabelPolicy.NotChanged")
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToPrivacyPolicy(&existingPolicy.PrivacyPolicyWriteModel), nil
}

func (c *Commands) defaultPrivacyPolicyWriteModelByID(ctx context.Context) (policy *IAMPrivacyPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMPrivacyPolicyWriteModel()
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
