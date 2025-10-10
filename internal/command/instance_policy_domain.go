package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddDefaultDomainPolicy(ctx context.Context, userLoginMustBeDomain, validateOrgDomains, smtpSenderAddressMatchesInstanceDomain bool) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareAddDefaultDomainPolicy(instanceAgg, userLoginMustBeDomain, validateOrgDomains, smtpSenderAddressMatchesInstanceDomain))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) ChangeDefaultDomainPolicy(ctx context.Context, userLoginMustBeDomain, validateOrgDomains, smtpSenderAddressMatchesInstanceDomain bool) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareChangeDefaultDomainPolicy(instanceAgg, userLoginMustBeDomain, validateOrgDomains, smtpSenderAddressMatchesInstanceDomain))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	// returning the values of the first event as this is the one from the instance
	return &domain.ObjectDetails{
		Sequence:      pushedEvents[0].Sequence(),
		EventDate:     pushedEvents[0].CreatedAt(),
		ResourceOwner: pushedEvents[0].Aggregate().ResourceOwner,
	}, nil
}

func (c *Commands) getDefaultDomainPolicy(ctx context.Context) (*domain.DomainPolicy, error) {
	policyWriteModel, err := c.instanceDomainPolicyWriteModel(ctx)
	if err != nil {
		return nil, err
	}
	if !policyWriteModel.State.Exists() {
		return nil, zerrors.ThrowInvalidArgument(nil, "INSTANCE-3n8fs", "Errors.IAM.PasswordComplexityPolicy.NotFound")
	}
	policy := writeModelToDomainPolicy(policyWriteModel)
	policy.Default = true
	return policy, nil
}

func (c *Commands) instanceDomainPolicyWriteModel(ctx context.Context) (policy *InstanceDomainPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewInstanceDomainPolicyWriteModel(ctx)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func prepareAddDefaultDomainPolicy(
	a *instance.Aggregate,
	userLoginMustBeDomain,
	validateOrgDomains,
	smtpSenderAddressMatchesInstanceDomain bool,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := instanceDomainPolicy(ctx, filter)
			if err != nil {
				return nil, err
			}
			if writeModel.State == domain.PolicyStateActive {
				return nil, zerrors.ThrowAlreadyExists(nil, "INSTANCE-Lk0dS", "Errors.Instance.DomainPolicy.AlreadyExists")
			}
			return []eventstore.Command{
				instance.NewDomainPolicyAddedEvent(ctx, &a.Aggregate,
					userLoginMustBeDomain,
					validateOrgDomains,
					smtpSenderAddressMatchesInstanceDomain,
				),
			}, nil
		}, nil
	}
}

func prepareChangeDefaultDomainPolicy(
	a *instance.Aggregate,
	userLoginMustBeDomain,
	validateOrgDomains,
	smtpSenderAddressMatchesInstanceDomain bool,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := instanceDomainPolicy(ctx, filter)
			if err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "INSTANCE-0Pl0d", "Errors.Instance.DomainPolicy.NotFound")
			}
			changedEvent, usernameChange, err := writeModel.NewChangedEvent(ctx, &a.Aggregate,
				userLoginMustBeDomain,
				validateOrgDomains,
				smtpSenderAddressMatchesInstanceDomain,
			)
			if err != nil {
				return nil, err
			}
			cmds := []eventstore.Command{changedEvent}
			// if the UserLoginMustBeDomain has not changed, no further changes are needed
			if !usernameChange {
				return cmds, err
			}
			// get all organisations without a custom domain policy
			orgsWriteModel, err := domainPolicyOrgs(ctx, filter)
			if err != nil {
				return nil, err
			}
			// loop over all found organisations to get their usernames
			// and to compute the username changed events
			for _, orgID := range orgsWriteModel.OrgIDs {
				organizationScopedUsernames, err := checkOrganizationScopedUsernames(ctx, filter, a.ID, nil)
				if err != nil {
					return nil, err
				}

				usersWriteModel, err := domainPolicyUsernames(ctx, filter, orgID)
				if err != nil {
					return nil, err
				}
				cmds = append(cmds, usersWriteModel.NewUsernameChangedEvents(ctx,
					userLoginMustBeDomain,
					organizationScopedUsernames,
					writeModel.UserLoginMustBeDomain,
				)...)
			}
			return cmds, nil
		}, nil
	}
}
