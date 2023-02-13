package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func (c *Commands) AddInstanceGenericOAuthProvider(ctx context.Context, provider GenericOAuthProvider) (string, *domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddInstanceOAuthProvider(
		instanceAgg,
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

func (c *Commands) UpdateInstanceGenericOAuthProvider(ctx context.Context, id string, provider GenericOAuthProvider) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateInstanceOAuthProvider(instanceAgg, id, provider))
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

func (c *Commands) AddInstanceGoogleProvider(ctx context.Context, provider GoogleProvider) (string, *domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddInstanceGoogleProvider(instanceAgg, id, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateInstanceGoogleProvider(ctx context.Context, id string, provider GoogleProvider) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateInstanceGoogleProvider(instanceAgg, id, provider))
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

func (c *Commands) prepareAddInstanceOAuthProvider(a *instance.Aggregate, id string, provider GenericOAuthProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewOAuthInstanceIDPWriteModel(a.InstanceID, id)
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
				instance.NewOAuthIDPAddedEvent(
					ctx,
					&a.Aggregate,
					id,
					provider.Name,
					provider.ClientID,
					secret,
					provider.AuthorizationEndpoint,
					provider.TokenEndpoint,
					provider.UserEndpoint,
					provider.Scopes,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateInstanceOAuthProvider(
	a *instance.Aggregate,
	id string,
	provider GenericOAuthProvider,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewOAuthInstanceIDPWriteModel(a.InstanceID, id)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "INST-D3r1s", "Errors.Instance.IDPConfig.NotExisting")
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

func (c *Commands) prepareAddInstanceGoogleProvider(a *instance.Aggregate, id string, provider GoogleProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewGoogleInstanceIDPWriteModel(a.InstanceID, id)
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
				instance.NewGoogleIDPAddedEvent(ctx, &a.Aggregate, id, provider.ClientID, secret, provider.Scopes, provider.IDPOptions),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateInstanceGoogleProvider(a *instance.Aggregate, id string, provider GoogleProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewGoogleInstanceIDPWriteModel(a.InstanceID, id)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "INST-D3r1s", "Errors.Instance.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				id,
				provider.ClientID,
				provider.ClientSecret,
				c.idpConfigEncryption,
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
