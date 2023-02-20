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

func (c *Commands) AddInstanceLDAPProvider(ctx context.Context, provider LDAPProvider) (string, *domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddInstanceLDAPProvider(instanceAgg, id, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateInstanceLDAPProvider(ctx context.Context, id string, provider LDAPProvider) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateInstanceLDAPProvider(instanceAgg, id, provider))
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

func (c *Commands) DeleteInstanceProvider(ctx context.Context, id string) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareDeleteInstanceProvider(instanceAgg, id))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) prepareAddInstanceLDAPProvider(a *instance.Aggregate, id string, provider LDAPProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-SAfdd", "Errors.Invalid.Argument")
		}
		if provider.Host = strings.TrimSpace(provider.Host); provider.Host == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-SDVg2", "Errors.Invalid.Argument")
		}
		if provider.BaseDN = strings.TrimSpace(provider.BaseDN); provider.BaseDN == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-sv31s", "Errors.Invalid.Argument")
		}
		if provider.UserObjectClass = strings.TrimSpace(provider.UserObjectClass); provider.UserObjectClass == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-sdgf4", "Errors.Invalid.Argument")
		}
		if provider.UserUniqueAttribute = strings.TrimSpace(provider.UserUniqueAttribute); provider.UserUniqueAttribute == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-AEG2w", "Errors.Invalid.Argument")
		}
		if provider.Admin = strings.TrimSpace(provider.Admin); provider.Admin == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-SAD5n", "Errors.Invalid.Argument")
		}
		if provider.Password = strings.TrimSpace(provider.Password); provider.Password == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-sdf5h", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewLDAPInstanceIDPWriteModel(a.InstanceID, id)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			secret, err := crypto.Encrypt([]byte(provider.Password), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				instance.NewLDAPIDPAddedEvent(
					ctx,
					&a.Aggregate,
					id,
					provider.Name,
					provider.Host,
					provider.Port,
					provider.TLS,
					provider.BaseDN,
					provider.UserObjectClass,
					provider.UserUniqueAttribute,
					provider.Admin,
					secret,
					provider.LDAPAttributes,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateInstanceLDAPProvider(a *instance.Aggregate, id string, provider LDAPProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if id = strings.TrimSpace(id); id == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-Dgdbs", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-Sffgd", "Errors.Invalid.Argument")
		}
		if provider.Host = strings.TrimSpace(provider.Host); provider.Host == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-Dz62d", "Errors.Invalid.Argument")
		}
		if provider.BaseDN = strings.TrimSpace(provider.BaseDN); provider.BaseDN == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-vb3ss", "Errors.Invalid.Argument")
		}
		if provider.UserObjectClass = strings.TrimSpace(provider.UserObjectClass); provider.UserObjectClass == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-hbere", "Errors.Invalid.Argument")
		}
		if provider.UserUniqueAttribute = strings.TrimSpace(provider.UserUniqueAttribute); provider.UserUniqueAttribute == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-ASFt6", "Errors.Invalid.Argument")
		}
		if provider.Admin = strings.TrimSpace(provider.Admin); provider.Admin == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INST-DG45z", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewLDAPInstanceIDPWriteModel(a.InstanceID, id)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "INST-ASF3F", "Errors.Instance.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				id,
				writeModel.Name,
				provider.Name,
				provider.Host,
				provider.Port,
				provider.TLS,
				provider.BaseDN,
				provider.UserObjectClass,
				provider.UserUniqueAttribute,
				provider.Admin,
				provider.Password,
				c.idpConfigEncryption,
				provider.LDAPAttributes,
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

func (c *Commands) prepareDeleteInstanceProvider(a *instance.Aggregate, id string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewInstanceIDPRemoveWriteModel(a.InstanceID, id)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "INST-Se3tg", "Errors.Instance.IDPConfig.NotExisting")
			}
			return []eventstore.Command{instance.NewIDPRemovedEvent(ctx, &a.Aggregate, id, writeModel.name)}, nil
		}, nil
	}
}
