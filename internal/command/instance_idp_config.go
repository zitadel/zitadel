package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddDefaultIDPConfig(ctx context.Context, config *domain.IDPConfig) (*domain.IDPConfig, error) {
	if config.OIDCConfig == nil && config.JWTConfig == nil {
		return nil, zerrors.ThrowInvalidArgument(nil, "IDP-s8nn3", "Errors.IDPConfig.Invalid")
	}
	idpConfigID, err := id_generator.Next()
	if err != nil {
		return nil, err
	}
	addedConfig := NewInstanceIDPConfigWriteModel(ctx, idpConfigID)

	instanceAgg := InstanceAggregateFromWriteModel(&addedConfig.WriteModel)
	events := []eventstore.Command{
		instance.NewIDPConfigAddedEvent(
			ctx,
			instanceAgg,
			idpConfigID,
			config.Name,
			config.Type,
			config.StylingType,
			config.AutoRegister,
		),
	}
	if config.OIDCConfig != nil {
		clientSecret, err := crypto.Encrypt([]byte(config.OIDCConfig.ClientSecretString), c.idpConfigEncryption)
		if err != nil {
			return nil, err
		}

		events = append(events, instance.NewIDPOIDCConfigAddedEvent(
			ctx,
			instanceAgg,
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
			instanceAgg,
			idpConfigID,
			config.JWTConfig.JWTEndpoint,
			config.JWTConfig.Issuer,
			config.JWTConfig.KeysEndpoint,
			config.JWTConfig.HeaderName,
		))
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedConfig, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToIDPConfig(&addedConfig.IDPConfigWriteModel), nil
}

func (c *Commands) ChangeDefaultIDPConfig(ctx context.Context, config *domain.IDPConfig) (*domain.IDPConfig, error) {
	if config.IDPConfigID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "INSTANCE-4m9gs", "Errors.IDMissing")
	}
	existingIDP, err := c.instanceIDPConfigWriteModelByID(ctx, config.IDPConfigID)
	if err != nil {
		return nil, err
	}
	if existingIDP.State == domain.IDPConfigStateRemoved || existingIDP.State == domain.IDPConfigStateUnspecified {
		return nil, zerrors.ThrowNotFound(nil, "INSTANCE-m0e3r", "Errors.IDPConfig.NotExisting")
	}

	instanceAgg := InstanceAggregateFromWriteModel(&existingIDP.WriteModel)
	changedEvent, hasChanged := existingIDP.NewChangedEvent(ctx, instanceAgg, config.IDPConfigID, config.Name, config.StylingType, config.AutoRegister)
	if !hasChanged {
		return nil, zerrors.ThrowPreconditionFailed(nil, "INSTANCE-3k0fs", "Errors.IAM.IDPConfig.NotChanged")
	}
	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingIDP, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToIDPConfig(&existingIDP.IDPConfigWriteModel), nil
}

func (c *Commands) DeactivateDefaultIDPConfig(ctx context.Context, idpID string) (*domain.ObjectDetails, error) {
	existingIDP, err := c.instanceIDPConfigWriteModelByID(ctx, idpID)
	if err != nil {
		return nil, err
	}
	if existingIDP.State != domain.IDPConfigStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "INSTANCE-2n0fs", "Errors.IAM.IDPConfig.NotActive")
	}
	instanceAgg := InstanceAggregateFromWriteModel(&existingIDP.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewIDPConfigDeactivatedEvent(ctx, instanceAgg, idpID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingIDP, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingIDP.IDPConfigWriteModel.WriteModel), nil
}

func (c *Commands) ReactivateDefaultIDPConfig(ctx context.Context, idpID string) (*domain.ObjectDetails, error) {
	existingIDP, err := c.instanceIDPConfigWriteModelByID(ctx, idpID)
	if err != nil {
		return nil, err
	}
	if existingIDP.State != domain.IDPConfigStateInactive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "INSTANCE-5Mo0d", "Errors.IAM.IDPConfig.NotInactive")
	}
	instanceAgg := InstanceAggregateFromWriteModel(&existingIDP.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewIDPConfigReactivatedEvent(ctx, instanceAgg, idpID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingIDP, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingIDP.IDPConfigWriteModel.WriteModel), nil
}

func (c *Commands) RemoveDefaultIDPConfig(ctx context.Context, idpID string, idpProviders []*domain.IDPProvider, externalIDPs ...*domain.UserIDPLink) (*domain.ObjectDetails, error) {
	existingIDP, err := c.instanceIDPConfigWriteModelByID(ctx, idpID)
	if err != nil {
		return nil, err
	}
	if existingIDP.State == domain.IDPConfigStateRemoved || existingIDP.State == domain.IDPConfigStateUnspecified {
		return nil, zerrors.ThrowNotFound(nil, "INSTANCE-4M0xy", "Errors.IDPConfig.NotExisting")
	}

	instanceAgg := InstanceAggregateFromWriteModel(&existingIDP.WriteModel)
	events := []eventstore.Command{
		instance.NewIDPConfigRemovedEvent(ctx, instanceAgg, idpID, existingIDP.Name),
	}

	for _, idpProvider := range idpProviders {
		if idpProvider.AggregateID == authz.GetInstance(ctx).InstanceID() {
			userEvents := c.removeIDPProviderFromDefaultLoginPolicy(ctx, instanceAgg, idpProvider, true, externalIDPs...)
			events = append(events, userEvents...)
		}
		orgAgg := OrgAggregateFromWriteModel(&NewOrgIdentityProviderWriteModel(idpProvider.AggregateID, idpID).WriteModel)
		orgEvents := c.removeIDPFromLoginPolicy(ctx, orgAgg, idpID, true, externalIDPs...)
		events = append(events, orgEvents...)
	}

	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingIDP, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingIDP.IDPConfigWriteModel.WriteModel), nil
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
