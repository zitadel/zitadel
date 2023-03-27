package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func (c *Commands) AddDefaultIDPConfig(ctx context.Context, config *domain.IDPConfig) (string, *domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	return c.processWithID(ctx, prepareAddDefaultIDPConfig(instanceAgg, c.idGenerator, c.idpConfigEncryption, config))
}

func prepareAddDefaultIDPConfig(a *instance.Aggregate, idGenerator id.Generator, encrypt crypto.EncryptionAlgorithm, config *domain.IDPConfig) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if config.OIDCConfig == nil && config.JWTConfig == nil {
			return nil, errors.ThrowInvalidArgument(nil, "IDP-s8nn3", "Errors.IDPConfig.Invalid")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			idpConfigID, err := idGenerator.Next()
			if err != nil {
				return nil, err
			}

			idpWriteModel, err := getInstanceIDPConfigWriteModel(ctx, filter, idpConfigID)
			if err != nil {
				return nil, err
			}
			if idpWriteModel.State == domain.IDPConfigStateActive {
				return nil, errors.ThrowAlreadyExists(nil, "INST-i2nl", "Errors.IDPConfig.AlreadyExists")
			}

			events := []eventstore.Command{
				instance.NewIDPConfigAddedEvent(
					ctx,
					&a.Aggregate,
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
					&a.Aggregate,
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
					&a.Aggregate,
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
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	return c.processWithLast(ctx, c.prepareChangeDefaultIDPConfig(instanceAgg, config))
}

func (c *Commands) prepareChangeDefaultIDPConfig(a *instance.Aggregate, config *domain.IDPConfig) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if config.IDPConfigID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INSTANCE-4m9gs", "Errors.IDMissing")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			existingIDP, err := c.instanceIDPConfigWriteModelByID(ctx, config.IDPConfigID)
			if err != nil {
				return nil, err
			}
			if existingIDP.State == domain.IDPConfigStateRemoved || existingIDP.State == domain.IDPConfigStateUnspecified {
				return nil, errors.ThrowNotFound(nil, "INSTANCE-m0e3r", "Errors.IDPConfig.NotExisting")
			}
			changedEvent, hasChanged := existingIDP.NewChangedEvent(ctx, &a.Aggregate, config.IDPConfigID, config.Name, config.StylingType, config.AutoRegister)
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
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	return c.processWithLast(ctx, prepareDeactivateDefaultIDPConfig(instanceAgg, idpID))
}

func prepareDeactivateDefaultIDPConfig(a *instance.Aggregate, idpID string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if idpID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INSTANCE-apso2n", "Errors.IDMissing")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			existingIDP, err := getInstanceIDPConfigWriteModel(ctx, filter, idpID)
			if err != nil {
				return nil, err
			}
			if existingIDP.State != domain.IDPConfigStateActive {
				return nil, errors.ThrowPreconditionFailed(nil, "INSTANCE-2n0fs", "Errors.IAM.IDPConfig.NotActive")
			}
			return []eventstore.Command{
				instance.NewIDPConfigDeactivatedEvent(ctx, &a.Aggregate, idpID),
			}, nil
		}, nil
	}
}

func (c *Commands) ReactivateDefaultIDPConfig(ctx context.Context, idpID string) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	return c.processWithLast(ctx, prepareReactivateDefaultIDPConfig(instanceAgg, idpID))
}

func prepareReactivateDefaultIDPConfig(a *instance.Aggregate, idpID string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if idpID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INSTANCE-arso2n", "Errors.IDMissing")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			existingIDP, err := getInstanceIDPConfigWriteModel(ctx, filter, idpID)
			if err != nil {
				return nil, err
			}
			if existingIDP.State != domain.IDPConfigStateInactive {
				return nil, errors.ThrowPreconditionFailed(nil, "INSTANCE-5Mo0d", "Errors.IAM.IDPConfig.NotInactive")
			}
			return []eventstore.Command{
				instance.NewIDPConfigReactivatedEvent(ctx, &a.Aggregate, idpID),
			}, nil
		}, nil
	}
}

func (c *Commands) RemoveDefaultIDPConfig(ctx context.Context, idpID string, idpProviders []*domain.IDPProvider, externalIDPs ...*domain.UserIDPLink) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	return c.processWithLast(ctx, c.prepareRemoveDefaultIDPConfig(instanceAgg, idpID, idpProviders, externalIDPs...))
}

func (c *Commands) prepareRemoveDefaultIDPConfig(a *instance.Aggregate, idpID string, idpProviders []*domain.IDPProvider, externalIDPs ...*domain.UserIDPLink) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if idpID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INSTANCE-arso2n", "Errors.IDMissing")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			existingIDP, err := getInstanceIDPConfigWriteModel(ctx, filter, idpID)
			if err != nil {
				return nil, err
			}
			if existingIDP.State == domain.IDPConfigStateRemoved || existingIDP.State == domain.IDPConfigStateUnspecified {
				return nil, errors.ThrowNotFound(nil, "INSTANCE-4M0xy", "Errors.IDPConfig.NotExisting")
			}

			events := []eventstore.Command{
				instance.NewIDPConfigRemovedEvent(ctx, &a.Aggregate, idpID, existingIDP.Name),
			}

			for _, idpProvider := range idpProviders {
				if idpProvider.AggregateID == authz.GetInstance(ctx).InstanceID() {
					userEvents := c.removeIDPProviderFromDefaultLoginPolicy(ctx, &a.Aggregate, idpProvider, true, externalIDPs...)
					events = append(events, userEvents...)
				}
				orgAgg := OrgAggregateFromWriteModel(&NewOrgIdentityProviderWriteModel(idpProvider.AggregateID, idpID).WriteModel)
				orgEvents := c.removeIDPFromLoginPolicy(ctx, orgAgg, idpID, true)
				events = append(events, orgEvents...)
			}

			pushedEvents, err := c.eventstore.Push(ctx, events...)
			if err != nil {
				return nil, err
			}
			err = AppendAndReduce(existingIDP, pushedEvents...)

			return events, nil
		}, nil
	}
}

func (c *Commands) getInstanceIDPConfigByID(ctx context.Context, idpID string) (*domain.IDPConfig, error) {
	config, err := c.instanceIDPConfigWriteModelByID(ctx, idpID)
	if err != nil {
		return nil, err
	}
	if !config.State.Exists() {
		return nil, errors.ThrowNotFound(nil, "INSTANCE-p0pFF", "Errors.IDPConfig.NotExisting")
	}
	return writeModelToIDPConfig(&config.IDPConfigWriteModel), nil
}

func (c *Commands) instanceIDPConfigWriteModelByID(ctx context.Context, idpID string) (policy *InstanceIDPConfigWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewInstanceIDPConfigWriteModel(ctx, idpID)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func getInstanceIDPConfigWriteModel(ctx context.Context, filter preparation.FilterToQueryReducer, idpID string) (*InstanceIDPConfigWriteModel, error) {
	writeModel := NewInstanceIDPConfigWriteModel(ctx, idpID)
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return writeModel, nil
	}
	writeModel.AppendEvents(events...)
	err = writeModel.Reduce()
	return writeModel, err
}
