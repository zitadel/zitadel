package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddOrgDomainPolicy(ctx context.Context, resourceOwner string, userLoginMustBeDomain, validateOrgDomains, smtpSenderAddressMatchesInstanceDomain bool) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-4Jfsf", "Errors.ResourceOwnerMissing")
	}
	orgAgg := org.NewAggregate(resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareAddOrgDomainPolicy(orgAgg, userLoginMustBeDomain, validateOrgDomains, smtpSenderAddressMatchesInstanceDomain))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) ChangeOrgDomainPolicy(ctx context.Context, resourceOwner string, userLoginMustBeDomain, validateOrgDomains, smtpSenderAddressMatchesInstanceDomain bool) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-5H8fs", "Errors.ResourceOwnerMissing")
	}
	orgAgg := org.NewAggregate(resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareChangeOrgDomainPolicy(orgAgg, userLoginMustBeDomain, validateOrgDomains, smtpSenderAddressMatchesInstanceDomain))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) RemoveOrgDomainPolicy(ctx context.Context, orgID string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-3H8fs", "Errors.ResourceOwnerMissing")
	}
	orgAgg := org.NewAggregate(orgID)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareRemoveOrgDomainPolicy(orgAgg))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

// Deprecated: Use commands.domainPolicyWriteModel directly, to remove the domain.DomainPolicy struct
func (c *Commands) getOrgDomainPolicy(ctx context.Context, orgID string) (_ *domain.DomainPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	policy, err := c.orgDomainPolicyWriteModel(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if policy.State.Exists() {
		return orgWriteModelToDomainPolicy(policy), nil
	}
	return c.getDefaultDomainPolicy(ctx)
}

func (c *Commands) orgDomainPolicyWriteModel(ctx context.Context, orgID string) (policy *OrgDomainPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewOrgDomainPolicyWriteModel(orgID)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func prepareAddOrgDomainPolicy(
	a *org.Aggregate,
	userLoginMustBeDomain,
	validateOrgDomains,
	smtpSenderAddressMatchesInstanceDomain bool,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) (_ []eventstore.Command, err error) {
			ctx, span := tracing.NewSpan(ctx)
			defer func() { span.EndWithError(err) }()

			writeModel, err := orgDomainPolicy(ctx, filter, a.ID)
			if err != nil {
				return nil, err
			}
			if writeModel.State == domain.PolicyStateActive {
				return nil, zerrors.ThrowAlreadyExists(nil, "ORG-1M8ds", "Errors.Org.DomainPolicy.AlreadyExists")
			}
			cmds := []eventstore.Command{
				org.NewDomainPolicyAddedEvent(ctx, &a.Aggregate,
					userLoginMustBeDomain,
					validateOrgDomains,
					smtpSenderAddressMatchesInstanceDomain,
				),
			}
			instancePolicy, err := instanceDomainPolicy(ctx, filter)
			if err != nil {
				return nil, err
			}
			// regardless if the UserLoginMustBeDomain setting is true or false,
			// if it will be the same value as currently on the instance,
			// then there no further changes are needed
			if instancePolicy.UserLoginMustBeDomain == userLoginMustBeDomain {
				return cmds, nil
			}
			// the UserLoginMustBeDomain setting will be different from the instance
			// therefore get all usernames and the current primary domain
			usersWriteModel, err := domainPolicyUsernames(ctx, filter, a.ID)
			if err != nil {
				return nil, err
			}
			return append(cmds, usersWriteModel.NewUsernameChangedEvents(ctx, userLoginMustBeDomain)...), nil
		}, nil
	}
}

func prepareChangeOrgDomainPolicy(
	a *org.Aggregate,
	userLoginMustBeDomain,
	validateOrgDomains,
	smtpSenderAddressMatchesInstanceDomain bool,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := orgDomainPolicy(ctx, filter, a.ID)
			if err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "ORG-2N9sd", "Errors.Org.DomainPolicy.NotFound")
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
			// get all usernames and the primary domain
			usersWriteModel, err := domainPolicyUsernames(ctx, filter, a.ID)
			if err != nil {
				return nil, err
			}
			// to compute the username changed events
			return append(cmds, usersWriteModel.NewUsernameChangedEvents(ctx, userLoginMustBeDomain)...), nil
		}, nil
	}
}

func prepareRemoveOrgDomainPolicy(
	a *org.Aggregate,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := orgDomainPolicy(ctx, filter, a.ID)
			if err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "ORG-Dvsh3", "Errors.Org.DomainPolicy.NotFound")
			}
			instancePolicy, err := instanceDomainPolicy(ctx, filter)
			if err != nil {
				return nil, err
			}
			cmds := []eventstore.Command{
				org.NewDomainPolicyRemovedEvent(ctx, &a.Aggregate),
			}
			// regardless if the UserLoginMustBeDomain setting is true or false,
			// if it will be the same value as currently on the instance,
			// then there no further changes are needed
			if instancePolicy.UserLoginMustBeDomain == writeModel.UserLoginMustBeDomain {
				return cmds, nil
			}
			// get all usernames and the primary domain
			usersWriteModel, err := domainPolicyUsernames(ctx, filter, a.ID)
			if err != nil {
				return nil, err
			}
			// to compute the username changed events
			return append(cmds, usersWriteModel.NewUsernameChangedEvents(ctx, instancePolicy.UserLoginMustBeDomain)...), nil
		}, nil
	}
}
