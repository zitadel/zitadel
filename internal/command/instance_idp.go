package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func (c *Commands) AddInstanceGenericOAuthProvider(ctx context.Context, resourceOwner string, provider GenericOAuthProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddInstanceOAuthProvider(
		orgAgg,
		resourceOwner,
		id,
		provider,
	))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateInstanceGenericOAuthProvider(ctx context.Context, resourceOwner, id string, provider GenericOAuthProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateInstanceOAuthProvider(orgAgg, resourceOwner, id, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) AddInstanceGoogleProvider(ctx context.Context, resourceOwner, clientID, clientSecret string, options idp.Options) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddInstanceGoogleProvider(
		orgAgg,
		resourceOwner,
		id,
		clientID,
		clientSecret,
		options,
	))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateInstanceGoogleProvider(ctx context.Context, resourceOwner, id, clientID, clientSecret string, options idp.Options) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateInstanceGoogleProvider(orgAgg, resourceOwner, id, clientID, clientSecret, options))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) prepareAddInstanceOAuthProvider(a *org.Aggregate, resourceOwner, id string, provider GenericOAuthProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewOAuthOrgIDPWriteModel(resourceOwner, id)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			secret, err := crypto.Encrypt([]byte(provider.ClientSecret), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewOAuthIDPAddedEvent(ctx, &a.Aggregate, id, provider.Name, provider.ClientID, secret, provider.IDPOptions),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateInstanceOAuthProvider(
	a *org.Aggregate,
	resourceOwner,
	id string,
	provider GenericOAuthProvider,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewOAuthOrgIDPWriteModel(resourceOwner, id)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "ORG-D3r1s", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				id,
				writeModel.Name,
				provider.Name,
				provider.ClientID,
				provider.ClientSecret,
				c.idpConfigEncryption,
				provider.AuthorizationEndpoint,
				provider.TokenEndpoint,
				provider.UserEndpoint,
				provider.Scopes,
				provider.IDPOptions,
			)
			if err != nil {
				return nil, err
			}
			if event == nil {
				return nil, nil
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareAddInstanceGoogleProvider(a *org.Aggregate, resourceOwner, id, clientID, clientSecret string, options idp.Options) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewGoogleOrgIDPWriteModel(resourceOwner, id)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			secret, err := crypto.Encrypt([]byte(clientSecret), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewGoogleIDPAddedEvent(ctx, &a.Aggregate, id, clientID, secret, idp.Options{IsAutoUpdate: options.IsAutoUpdate}),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateInstanceGoogleProvider(
	a *org.Aggregate,
	resourceOwner,
	id,
	clientID,
	clientSecret string,
	options idp.Options,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewGoogleOrgIDPWriteModel(resourceOwner, id)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "ORG-D3r1s", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				id,
				clientID,
				clientSecret,
				c.idpConfigEncryption,
				options,
			)
			if err != nil {
				return nil, err
			}
			if event == nil {
				return nil, nil
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}
