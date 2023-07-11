package command

import (
	"context"
	"regexp"
	"strings"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/project"
)

var (
	allowDomainRunes = regexp.MustCompile("^[a-zA-Z0-9\\.\\-]+$")
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

func (c *Commands) addGeneratedInstanceDomain(ctx context.Context, a *instance.Aggregate, instanceName string) ([]preparation.Validation, error) {
	domain, err := domain.NewGeneratedInstanceDomain(instanceName, authz.GetInstance(ctx).RequestedDomain())
	if err != nil {
		return nil, err
	}
	return []preparation.Validation{
		c.addInstanceDomain(a, domain, true),
		setPrimaryInstanceDomain(a, domain),
	}, nil
}

func (c *Commands) addInstanceDomain(a *instance.Aggregate, instanceDomain string, generated bool) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if instanceDomain = strings.TrimSpace(instanceDomain); instanceDomain == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INST-28nlD", "Errors.Invalid.Argument")
		}
		if !allowDomainRunes.MatchString(instanceDomain) {
			return nil, errors.ThrowInvalidArgument(nil, "INST-S3v3w", "Errors.Instance.Domain.InvalidCharacter")
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
			consoleChangeEvent, err := c.updateConsoleRedirectURIs(ctx, filter, instanceDomain)
			if err != nil {
				return nil, err
			}
			if consoleChangeEvent == nil {
				return events, nil
			}
			return append(events, consoleChangeEvent), nil
		}, nil
	}
}

func (c *Commands) prepareUpdateConsoleRedirectURIs(instanceDomain string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if instanceDomain = strings.TrimSpace(instanceDomain); instanceDomain == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INST-E3j3s", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			consoleChangeEvent, err := c.updateConsoleRedirectURIs(ctx, filter, instanceDomain)
			if err != nil {
				return nil, err
			}
			if consoleChangeEvent == nil {
				return nil, nil
			}
			return []eventstore.Command{
				consoleChangeEvent,
			}, nil
		}, nil
	}
}

func (c *Commands) updateConsoleRedirectURIs(ctx context.Context, filter preparation.FilterToQueryReducer, instanceDomain string) (*project.OIDCConfigChangedEvent, error) {
	appWriteModel, err := getOIDCAppWriteModel(ctx, filter, authz.GetInstance(ctx).ProjectID(), authz.GetInstance(ctx).ConsoleApplicationID(), "")
	if err != nil {
		return nil, err
	}
	if !appWriteModel.State.Exists() {
		return nil, nil
	}
	redirectURI := http.BuildHTTP(instanceDomain, c.externalPort, c.externalSecure) + consoleRedirectPath
	changes := make([]project.OIDCConfigChanges, 0, 2)
	if !containsURI(appWriteModel.RedirectUris, redirectURI) {
		changes = append(changes, project.ChangeRedirectURIs(append(appWriteModel.RedirectUris, redirectURI)))
	}
	postLogoutRedirectURI := http.BuildHTTP(instanceDomain, c.externalPort, c.externalSecure) + consolePostLogoutPath
	if !containsURI(appWriteModel.PostLogoutRedirectUris, postLogoutRedirectURI) {
		changes = append(changes, project.ChangePostLogoutRedirectURIs(append(appWriteModel.PostLogoutRedirectUris, postLogoutRedirectURI)))
	}
	if len(changes) == 0 {
		return nil, nil
	}
	return project.NewOIDCConfigChangedEvent(
		ctx,
		ProjectAggregateFromWriteModel(&appWriteModel.WriteModel),
		appWriteModel.AppID,
		changes,
	)
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

func containsURI(uris []string, uri string) bool {
	for _, u := range uris {
		if u == uri {
			return true
		}
	}
	return false
}

func (c *Commands) getInstanceDomainsWriteModel(ctx context.Context, instanceID string) (*InstanceDomainsWriteModel, error) {
	domainsWriteModel := NewInstanceDomainsWriteModel(instanceID)
	err := c.eventstore.FilterToQueryReducer(ctx, domainsWriteModel)
	if err != nil {
		return nil, err
	}
	return domainsWriteModel, nil
}
