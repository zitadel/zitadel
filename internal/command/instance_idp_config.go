package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func (c *Commands) AddDefaultIDPConfig(ctx context.Context, config *domain.IDPConfig) (string, *domain.ObjectDetails, error) {
	wm := NewInstanceIDPConfigWriteModel(authz.GetInstance(ctx).InstanceID(), config.IDPConfigID)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareAddDefaultIDPConfig(wm, c.idGenerator, c.idpConfigEncryption, config))
	if err != nil {
		return "", nil, err
	}
	err = c.pushAppendAndReduce(ctx, wm, cmds...)
	if err != nil {
		return "", nil, err
	}
	return wm.ConfigID, writeModelToObjectDetails(&wm.WriteModel), nil
}

func prepareAddDefaultIDPConfig(wm *InstanceIDPConfigWriteModel, idGenerator id.Generator, encrypt crypto.EncryptionAlgorithm, config *domain.IDPConfig) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if config.OIDCConfig == nil && config.JWTConfig == nil {
			return nil, errors.ThrowInvalidArgument(nil, "IDP-s8nn3", "Errors.IDPConfig.Invalid")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			idpConfigID, err := idGenerator.Next()
			if err != nil {
				return nil, err
			}
			if err := queryAndReduce(ctx, filter, wm); err != nil {
				return nil, err
			}
			if wm.State == domain.IDPConfigStateActive {
				return nil, errors.ThrowAlreadyExists(nil, "INST-i2nl", "Errors.IDPConfig.AlreadyExists")
			}

			events := []eventstore.Command{
				instance.NewIDPConfigAddedEvent(
					ctx,
					InstanceAggregateFromWriteModel(&wm.WriteModel),
					idpConfigID,
					config.Name,
					config.Type,
					config.StylingType,
					config.AutoRegister,
				),
			}
			if config.OIDCConfig != nil {
				clientSecret, err := crypto.Encrypt([]byte(config.OIDCConfig.ClientSecretString), encrypt)
				if err != nil {
					return nil, err
				}

				events = append(events, instance.NewIDPOIDCConfigAddedEvent(
					ctx,
					InstanceAggregateFromWriteModel(&wm.WriteModel),
					config.OIDCConfig.ClientID,
					idpConfigID,
					config.OIDCConfig.Issuer,
					config.OIDCConfig.AuthorizationEndpoint,
					config.OIDCConfig.TokenEndpoint,
					clientSecret,
					config.OIDCConfig.IDPDisplayNameMapping,
					config.OIDCConfig.UsernameMapping,
					config.OIDCConfig.Scopes...,
				))
			} else if config.JWTConfig != nil {
				events = append(events, instance.NewIDPJWTConfigAddedEvent(
					ctx,
					InstanceAggregateFromWriteModel(&wm.WriteModel),
					idpConfigID,
					config.JWTConfig.JWTEndpoint,
					config.JWTConfig.Issuer,
					config.JWTConfig.KeysEndpoint,
					config.JWTConfig.HeaderName,
				))
			}
			return events, nil
		}, nil
	}
}

func (c *Commands) ChangeDefaultIDPConfig(ctx context.Context, config *domain.IDPConfig) (*domain.ObjectDetails, error) {
	wm := NewInstanceIDPConfigWriteModel(authz.GetInstance(ctx).InstanceID(), config.IDPConfigID)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareChangeDefaultIDPConfig(wm, config))
	if err != nil {
		return nil, err
	}
	err = c.pushAppendAndReduce(ctx, wm, cmds...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func prepareChangeDefaultIDPConfig(wm *InstanceIDPConfigWriteModel, config *domain.IDPConfig) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if config.IDPConfigID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INSTANCE-4m9gs", "Errors.IDMissing")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			if err := queryAndReduce(ctx, filter, wm); err != nil {
				return nil, err
			}
			if wm.State == domain.IDPConfigStateRemoved || wm.State == domain.IDPConfigStateUnspecified {
				return nil, errors.ThrowNotFound(nil, "INSTANCE-m0e3r", "Errors.IDPConfig.NotExisting")
			}
			changedEvent, hasChanged := wm.NewChangedEvent(ctx, InstanceAggregateFromWriteModel(&wm.WriteModel), config.IDPConfigID, config.Name, config.StylingType, config.AutoRegister)
			if !hasChanged {
				return nil, errors.ThrowPreconditionFailed(nil, "INSTANCE-3k0fs", "Errors.IAM.IDPConfig.NotChanged")
			}

			return []eventstore.Command{
				changedEvent,
			}, nil
		}, nil
	}
}

func (c *Commands) DeactivateDefaultIDPConfig(ctx context.Context, idpID string) (*domain.ObjectDetails, error) {
	wm := NewInstanceIDPConfigWriteModel(authz.GetInstance(ctx).InstanceID(), idpID)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareDeactivateDefaultIDPConfig(wm, idpID))
	if err != nil {
		return nil, err
	}
	err = c.pushAppendAndReduce(ctx, wm, cmds...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func prepareDeactivateDefaultIDPConfig(wm *InstanceIDPConfigWriteModel, idpID string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if idpID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INSTANCE-apso2n", "Errors.IDMissing")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			if err := queryAndReduce(ctx, filter, wm); err != nil {
				return nil, err
			}
			if wm.State != domain.IDPConfigStateActive {
				return nil, errors.ThrowPreconditionFailed(nil, "INSTANCE-2n0fs", "Errors.IAM.IDPConfig.NotActive")
			}
			return []eventstore.Command{
				instance.NewIDPConfigDeactivatedEvent(ctx, InstanceAggregateFromWriteModel(&wm.WriteModel), idpID),
			}, nil
		}, nil
	}
}

func (c *Commands) ReactivateDefaultIDPConfig(ctx context.Context, idpID string) (*domain.ObjectDetails, error) {
	wm := NewInstanceIDPConfigWriteModel(authz.GetInstance(ctx).InstanceID(), idpID)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareReactivateDefaultIDPConfig(wm, idpID))
	if err != nil {
		return nil, err
	}
	err = c.pushAppendAndReduce(ctx, wm, cmds...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func prepareReactivateDefaultIDPConfig(wm *InstanceIDPConfigWriteModel, idpID string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if idpID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INSTANCE-arso2n", "Errors.IDMissing")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			if err := queryAndReduce(ctx, filter, wm); err != nil {
				return nil, err
			}
			if wm.State != domain.IDPConfigStateInactive {
				return nil, errors.ThrowPreconditionFailed(nil, "INSTANCE-5Mo0d", "Errors.IAM.IDPConfig.NotInactive")
			}
			return []eventstore.Command{
				instance.NewIDPConfigReactivatedEvent(ctx, InstanceAggregateFromWriteModel(&wm.WriteModel), idpID),
			}, nil
		}, nil
	}
}

func (c *Commands) RemoveDefaultIDPConfig(ctx context.Context, idpID string, idpProviders []*domain.IDPProvider, externalIDPs ...*domain.UserIDPLink) (*domain.ObjectDetails, error) {
	wm := NewInstanceIDPConfigWriteModel(authz.GetInstance(ctx).InstanceID(), idpID)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareRemoveDefaultIDPConfig(wm, idpID, idpProviders, externalIDPs...))
	if err != nil {
		return nil, err
	}
	err = c.pushAppendAndReduce(ctx, wm, cmds...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func (c *Commands) prepareRemoveDefaultIDPConfig(wm *InstanceIDPConfigWriteModel, idpID string, idpProviders []*domain.IDPProvider, externalIDPs ...*domain.UserIDPLink) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if idpID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INSTANCE-arso2n", "Errors.IDMissing")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			if err := queryAndReduce(ctx, filter, wm); err != nil {
				return nil, err
			}
			if wm.State == domain.IDPConfigStateRemoved || wm.State == domain.IDPConfigStateUnspecified {
				return nil, errors.ThrowNotFound(nil, "INSTANCE-4M0xy", "Errors.IDPConfig.NotExisting")
			}

			events := []eventstore.Command{
				instance.NewIDPConfigRemovedEvent(ctx, InstanceAggregateFromWriteModel(&wm.WriteModel), idpID, wm.Name),
			}

			for _, idpProvider := range idpProviders {
				if idpProvider.AggregateID == authz.GetInstance(ctx).InstanceID() {
					userEvents := c.removeIDPProviderFromDefaultLoginPolicy(ctx, InstanceAggregateFromWriteModel(&wm.WriteModel), idpProvider, true, externalIDPs...)
					events = append(events, userEvents...)
				}
				orgAgg := OrgAggregateFromWriteModel(&NewOrgIdentityProviderWriteModel(idpProvider.AggregateID, idpID).WriteModel)
				orgEvents := c.removeIDPFromLoginPolicy(ctx, orgAgg, idpID, true)
				events = append(events, orgEvents...)
			}
			return events, nil
		}, nil
	}
}
