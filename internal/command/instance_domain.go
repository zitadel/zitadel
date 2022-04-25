package command

import (
	"context"
	"strings"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/command/preparation"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/project"
)

func (c *Commands) AddInstanceDomain(ctx context.Context, instanceDomain string) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	validation := c.addInstanceDomain(instanceAgg, instanceDomain, false)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: events[len(events)-1].Aggregate().InstanceID,
	}, nil
}

func (c *Commands) SetPrimaryInstanceDomain(ctx context.Context, instanceDomain string) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	validation := setPrimaryInstanceDomain(instanceAgg, instanceDomain)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: events[len(events)-1].Aggregate().InstanceID,
	}, nil
}

func (c *Commands) RemoveInstanceDomain(ctx context.Context, instanceDomain string) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	validation := removeInstanceDomain(instanceAgg, instanceDomain)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: events[len(events)-1].Aggregate().InstanceID,
	}, nil
}

func (c *Commands) addGeneratedInstanceDomain(ctx context.Context, a *instance.Aggregate, instanceName string) preparation.Validation {
	domain := domain.NewGeneratedInstanceDomain(instanceName, authz.GetInstance(ctx).RequestedDomain())
	return c.addInstanceDomain(a, domain, true)
}

func (c *Commands) addInstanceDomain(a *instance.Aggregate, instanceDomain string, generated bool) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if instanceDomain = strings.TrimSpace(instanceDomain); instanceDomain == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INST-28nlD", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			domainWriteModel, err := getInstanceDomainWriteModel(ctx, filter, instanceDomain)
			if err != nil {
				return nil, err
			}
			if domainWriteModel.State == domain.InstanceDomainStateActive {
				return nil, errors.ThrowAlreadyExists(nil, "INST-i2nl", "Errors.Instance.Domain.AlreadyExists")
			}
			events := []eventstore.Command{
				instance.NewDomainAddedEvent(ctx, &a.Aggregate, instanceDomain, generated),
			}
			appWriteModel, err := getOIDCAppWriteModel(ctx, filter, authz.GetInstance(ctx).ProjectID(), authz.GetInstance(ctx).ConsoleApplicationID(), "")
			if err != nil {
				return nil, err
			}
			if appWriteModel.State.Exists() {
				redirectUrls := append(appWriteModel.RedirectUris, http.BuildHTTP(instanceDomain, c.externalPort, c.externalSecure)+consoleRedirectPath)
				logoutUrls := append(appWriteModel.PostLogoutRedirectUris, http.BuildOrigin(instanceDomain, c.externalSecure)+consolePostLogoutPath)
				consoleChangeEvent, err := project.NewOIDCConfigChangedEvent(
					ctx,
					ProjectAggregateFromWriteModel(&appWriteModel.WriteModel),
					appWriteModel.AppID,
					[]project.OIDCConfigChanges{
						project.ChangeRedirectURIs(redirectUrls),
						project.ChangePostLogoutRedirectURIs(logoutUrls),
					},
				)
				if err != nil {
					return nil, err
				}
				events = append(events, consoleChangeEvent)
			}

			return events, nil
		}, nil
	}
}

func setPrimaryInstanceDomain(a *instance.Aggregate, instanceDomain string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if instanceDomain = strings.TrimSpace(instanceDomain); instanceDomain == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INST-9mWjf", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			domainWriteModel, err := getInstanceDomainWriteModel(ctx, filter, instanceDomain)
			if err != nil {
				return nil, err
			}
			if !domainWriteModel.State.Exists() {
				return nil, errors.ThrowNotFound(nil, "INSTANCE-9nkWf", "Errors.Instance.Domain.NotFound")
			}
			return []eventstore.Command{instance.NewDomainPrimarySetEvent(ctx, &a.Aggregate, instanceDomain)}, nil
		}, nil
	}
}

func removeInstanceDomain(a *instance.Aggregate, instanceDomain string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if instanceDomain = strings.TrimSpace(instanceDomain); instanceDomain == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INST-39nls", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			domainWriteModel, err := getInstanceDomainWriteModel(ctx, filter, instanceDomain)
			if err != nil {
				return nil, err
			}
			if domainWriteModel.State != domain.InstanceDomainStateActive {
				return nil, errors.ThrowNotFound(nil, "INSTANCE-8ls9f", "Errors.Instance.Domain.NotFound")
			}
			if domainWriteModel.Generated {
				return nil, errors.ThrowPreconditionFailed(nil, "INSTANCE-9hn3n", "Errors.Instance.Domain.GeneratedNotRemovable")
			}
			return []eventstore.Command{instance.NewDomainRemovedEvent(ctx, &a.Aggregate, instanceDomain)}, nil
		}, nil
	}
}

func getInstanceDomainWriteModel(ctx context.Context, filter preparation.FilterToQueryReducer, domain string) (*InstanceDomainWriteModel, error) {
	domainWriteModel := NewInstanceDomainWriteModel(ctx, domain)
	events, err := filter(ctx, domainWriteModel.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return domainWriteModel, nil
	}
	domainWriteModel.AppendEvents(events...)
	err = domainWriteModel.Reduce()
	return domainWriteModel, err
}
