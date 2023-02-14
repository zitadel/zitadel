package command

import (
	"context"
	"strings"

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
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddInstanceOAuthProvider(instanceAgg, id, provider))
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

func (c *Commands) AddInstanceGenericOIDCProvider(ctx context.Context, provider GenericOIDCProvider) (string, *domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddInstanceOIDCProvider(instanceAgg, id, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateInstanceGenericOIDCProvider(ctx context.Context, id string, provider GenericOIDCProvider) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateInstanceOIDCProvider(instanceAgg, id, provider))
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

func (c *Commands) AddInstanceJWTProvider(ctx context.Context, provider JWTProvider) (string, *domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddInstanceJWTProvider(instanceAgg, id, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateInstanceJWTProvider(ctx context.Context, id string, provider JWTProvider) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateInstanceJWTProvider(instanceAgg, id, provider))
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
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-D32ef", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-Dbgzf", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-DF4ga", "Errors.Invalid.Argument")
		}
		if provider.AuthorizationEndpoint = strings.TrimSpace(provider.AuthorizationEndpoint); provider.AuthorizationEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-B23bs", "Errors.Invalid.Argument")
		}
		if provider.TokenEndpoint = strings.TrimSpace(provider.TokenEndpoint); provider.TokenEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-D2gj8", "Errors.Invalid.Argument")
		}
		if provider.UserEndpoint = strings.TrimSpace(provider.UserEndpoint); provider.UserEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-Fb8jk", "Errors.Invalid.Argument")
		}
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

func (c *Commands) prepareUpdateInstanceOAuthProvider(a *instance.Aggregate, id string, provider GenericOAuthProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-D32ef", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-Dbgzf", "Errors.Invalid.Argument")
		}
		if provider.AuthorizationEndpoint = strings.TrimSpace(provider.AuthorizationEndpoint); provider.AuthorizationEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-B23bs", "Errors.Invalid.Argument")
		}
		if provider.TokenEndpoint = strings.TrimSpace(provider.TokenEndpoint); provider.TokenEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-D2gj8", "Errors.Invalid.Argument")
		}
		if provider.UserEndpoint = strings.TrimSpace(provider.UserEndpoint); provider.UserEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-Fb8jk", "Errors.Invalid.Argument")
		}
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

func (c *Commands) prepareAddInstanceOIDCProvider(a *instance.Aggregate, id string, provider GenericOIDCProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-Sgtj5", "Errors.Invalid.Argument")
		}
		if provider.Issuer = strings.TrimSpace(provider.Issuer); provider.Issuer == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-Hz6zj", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-fb5jm", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-Sfdf4", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewOIDCInstanceIDPWriteModel(a.InstanceID, id)
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
				instance.NewOIDCIDPAddedEvent(
					ctx,
					&a.Aggregate,
					id,
					provider.Name,
					provider.Issuer,
					provider.ClientID,
					secret,
					provider.Scopes,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateInstanceOIDCProvider(a *instance.Aggregate, id string, provider GenericOIDCProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if id = strings.TrimSpace(id); id == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-SAfd3", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-Dvf4f", "Errors.Invalid.Argument")
		}
		if provider.Issuer = strings.TrimSpace(provider.Issuer); provider.Issuer == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-BDfr3", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-Db3bs", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewOIDCInstanceIDPWriteModel(a.InstanceID, id)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "INST-Dg331", "Errors.Instance.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				id,
				writeModel.Name,
				provider.Name,
				provider.Issuer,
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

func (c *Commands) prepareAddInstanceJWTProvider(a *instance.Aggregate, id string, provider JWTProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-JLKef", "Errors.Invalid.Argument")
		}
		if provider.Issuer = strings.TrimSpace(provider.Issuer); provider.Issuer == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-WNJK3", "Errors.Invalid.Argument")
		}
		if provider.JWTEndpoint = strings.TrimSpace(provider.JWTEndpoint); provider.JWTEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-NJKSD", "Errors.Invalid.Argument")
		}
		if provider.KeyEndpoint = strings.TrimSpace(provider.KeyEndpoint); provider.KeyEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-NJKE3", "Errors.Invalid.Argument")
		}
		if provider.HeaderName = strings.TrimSpace(provider.HeaderName); provider.HeaderName == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-2rlks", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewJWTInstanceIDPWriteModel(a.InstanceID, id)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			return []eventstore.Command{
				instance.NewJWTIDPAddedEvent(
					ctx,
					&a.Aggregate,
					id,
					provider.Name,
					provider.Issuer,
					provider.JWTEndpoint,
					provider.KeyEndpoint,
					provider.HeaderName,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateInstanceJWTProvider(a *instance.Aggregate, id string, provider JWTProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if id = strings.TrimSpace(id); id == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-HUe3q", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-JKLS2", "Errors.Invalid.Argument")
		}
		if provider.Issuer = strings.TrimSpace(provider.Issuer); provider.Issuer == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-JKs3f", "Errors.Invalid.Argument")
		}
		if provider.JWTEndpoint = strings.TrimSpace(provider.JWTEndpoint); provider.JWTEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-NJKS2", "Errors.Invalid.Argument")
		}
		if provider.KeyEndpoint = strings.TrimSpace(provider.KeyEndpoint); provider.KeyEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-SJk2d", "Errors.Invalid.Argument")
		}
		if provider.HeaderName = strings.TrimSpace(provider.HeaderName); provider.HeaderName == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-SJK2f", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewJWTInstanceIDPWriteModel(a.InstanceID, id)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "INST-Bhju5", "Errors.Instance.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				id,
				writeModel.Name,
				provider.Name,
				provider.Issuer,
				provider.JWTEndpoint,
				provider.KeyEndpoint,
				provider.HeaderName,
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
