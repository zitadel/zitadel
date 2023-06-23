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

func (c *Commands) SetPrimaryInstanceDomain(ctx context.Context, instanceDomain string) (*domain.ObjectDetails, error) {
	wm := NewInstanceDomainWriteModel(authz.GetInstance(ctx).InstanceID(), instanceDomain)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareSetPrimaryInstanceDomain(wm))
	if err != nil {
		return nil, err
	}
	err = c.pushAppendAndReduce(ctx, wm, cmds...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func (c *Commands) RemoveInstanceDomain(ctx context.Context, instanceDomain string) (*domain.ObjectDetails, error) {
	wm := NewInstanceDomainWriteModel(authz.GetInstance(ctx).InstanceID(), instanceDomain)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareRemoveInstanceDomain(wm))
	if err != nil {
		return nil, err
	}
	err = c.pushAppendAndReduce(ctx, wm, cmds...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func (c *Commands) addGeneratedInstanceDomain(ctx context.Context, instanceID, instanceName string) ([]preparation.Validation, error) {
	domain, err := domain.NewGeneratedInstanceDomain(instanceName, authz.GetInstance(ctx).RequestedDomain())
	if err != nil {
		return nil, err
	}
	return []preparation.Validation{
		c.prepareAddInstanceDomain(NewInstanceDomainWriteModel(instanceID, domain), true),
		prepareSetPrimaryInstanceDomain(NewInstanceDomainWriteModel(instanceID, domain)),
	}, nil
}

func (c *Commands) AddInstanceDomain(ctx context.Context, instanceDomain string) (*domain.ObjectDetails, error) {
	wm := NewInstanceDomainWriteModel(authz.GetInstance(ctx).InstanceID(), instanceDomain)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddInstanceDomain(wm, false))
	if err != nil {
		return nil, err
	}
	err = c.pushAppendAndReduce(ctx, wm, cmds...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func (c *Commands) prepareAddInstanceDomain(wm *InstanceDomainWriteModel, generated bool) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if wm.Domain = strings.TrimSpace(wm.Domain); wm.Domain == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INST-28nlD", "Errors.Invalid.Argument")
		}
		if !allowDomainRunes.MatchString(wm.Domain) {
			return nil, errors.ThrowInvalidArgument(nil, "INST-S3v3w", "Errors.Instance.Domain.InvalidCharacter")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			if err := queryAndReduce(ctx, filter, wm); err != nil {
				return nil, err
			}
			if wm.State == domain.InstanceDomainStateActive {
				return nil, errors.ThrowAlreadyExists(nil, "INST-i2nl", "Errors.Instance.Domain.AlreadyExists")
			}
			cmds := []eventstore.Command{
				instance.NewDomainAddedEvent(ctx, InstanceAggregateFromWriteModel(&wm.WriteModel), wm.Domain, generated),
			}
			consoleChangeEvent, err := c.updateConsoleRedirectURIs(ctx, filter, wm.Domain)
			if err != nil {
				return nil, err
			}
			return append(cmds, consoleChangeEvent), nil
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
	return project.NewOIDCConfigChangedEvent(
		ctx,
		ProjectAggregateFromWriteModel(&appWriteModel.WriteModel),
		appWriteModel.AppID,
		changes,
	)
}

// checkUpdateConsoleRedirectURIs validates if the required console uri is present in the redirect_uris and post_logout_redirect_uris
// it will return true only if present in both list, otherwise false
func (c *Commands) checkUpdateConsoleRedirectURIs(instanceDomain string, redirectURIs, postLogoutRedirectURIs []string) bool {
	redirectURI := http.BuildHTTP(instanceDomain, c.externalPort, c.externalSecure) + consoleRedirectPath
	if !containsURI(redirectURIs, redirectURI) {
		return false
	}
	postLogoutRedirectURI := http.BuildHTTP(instanceDomain, c.externalPort, c.externalSecure) + consolePostLogoutPath
	return containsURI(postLogoutRedirectURIs, postLogoutRedirectURI)
}

func prepareSetPrimaryInstanceDomain(wm *InstanceDomainWriteModel) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if wm.Domain = strings.TrimSpace(wm.Domain); wm.Domain == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INST-9mWjf", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			if err := queryAndReduce(ctx, filter, wm); err != nil {
				return nil, err
			}
			if !wm.State.Exists() {
				return nil, errors.ThrowNotFound(nil, "INSTANCE-9nkWf", "Errors.Instance.Domain.NotFound")
			}
			return []eventstore.Command{instance.NewDomainPrimarySetEvent(ctx, InstanceAggregateFromWriteModel(&wm.WriteModel), wm.Domain)}, nil
		}, nil
	}
}

func prepareRemoveInstanceDomain(wm *InstanceDomainWriteModel) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if wm.Domain = strings.TrimSpace(wm.Domain); wm.Domain == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INST-39nls", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			if err := queryAndReduce(ctx, filter, wm); err != nil {
				return nil, err
			}
			if wm.State != domain.InstanceDomainStateActive {
				return nil, errors.ThrowNotFound(nil, "INSTANCE-8ls9f", "Errors.Instance.Domain.NotFound")
			}
			if wm.Generated {
				return nil, errors.ThrowPreconditionFailed(nil, "INSTANCE-9hn3n", "Errors.Instance.Domain.GeneratedNotRemovable")
			}
			return []eventstore.Command{instance.NewDomainRemovedEvent(ctx, InstanceAggregateFromWriteModel(&wm.WriteModel), wm.Domain)}, nil
		}, nil
	}
}

func containsURI(uris []string, uri string) bool {
	for _, u := range uris {
		if u == uri {
			return true
		}
	}
	return false
}
