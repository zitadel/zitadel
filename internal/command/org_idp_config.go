package command

import (
	"context"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	org_repo "github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddIDPConfig(ctx context.Context, config domain.IDPConfig, resourceOwner string) (string, *domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return "", nil, caos_errs.ThrowInvalidArgument(nil, "Org-0j8gs", "Errors.ResourceOwnerMissing")
	}

	idpConfigID, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	addedConfig := NewOrgIDPConfigWriteModel(idpConfigID, resourceOwner)

	orgAgg := OrgAggregateFromWriteModel(&addedConfig.WriteModel)
	events := []eventstore.EventPusher{
		org_repo.NewIDPConfigAddedEvent(
			ctx,
			orgAgg,
			idpConfigID,
			config.IDPConfigName(),
			config.IDPConfigType(),
			config.IDPConfigStylingType(),
		),
	}

	var configEventCreator configEventCreator
	switch conf := config.(type) {
	case *domain.OIDCIDPConfig:
		configEventCreator = c.addOIDCIDPConfig(conf)
	case *domain.AuthConnectorIDPConfig:
		configEventCreator = c.addAuthConnectorIDPConfig(conf)
	default:
		return "", nil, errors.ThrowInvalidArgument(nil, "\"Org-eUpQU", "Errors.idp.config.notset")
	}
	configEvent, err := configEventCreator(ctx, orgAgg, idpConfigID)
	if err != nil {
		return "", nil, err
	}
	events = append(events, configEvent)
	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return "", nil, err
	}
	err = AppendAndReduce(addedConfig, pushedEvents...)
	if err != nil {
		return "", nil, err
	}
	return idpConfigID, writeModelToObjectDetails(&addedConfig.IDPConfigWriteModel.WriteModel), nil
}

func (c *Commands) addOIDCIDPConfig(config *domain.OIDCIDPConfig) configEventCreator {
	return func(ctx context.Context, agg *eventstore.Aggregate, idpConfigID string) (eventstore.EventPusher, error) {
		clientSecret, err := crypto.Encrypt([]byte(config.ClientSecretString), c.idpConfigSecretCrypto)
		if err != nil {
			return nil, err
		}

		return org_repo.NewIDPOIDCConfigAddedEvent(
			ctx,
			agg,
			config.ClientID,
			idpConfigID,
			config.Issuer,
			clientSecret,
			config.IDPDisplayNameMapping,
			config.UsernameMapping,
			config.Scopes...,
		), nil
	}
}

func (c *Commands) addAuthConnectorIDPConfig(config *domain.AuthConnectorIDPConfig) configEventCreator {
	return func(ctx context.Context, agg *eventstore.Aggregate, idpConfigID string) (eventstore.EventPusher, error) {
		return org_repo.NewIDPAuthConnectorConfigAddedEvent(
			ctx,
			agg,
			idpConfigID,
			config.BaseURL,
			config.BackendConnectorID,
		), nil
	}
}

func (c *Commands) ChangeIDPConfig(ctx context.Context, config domain.IDPConfig, resourceOwner string) (*domain.ObjectDetails, error) {
	if config.ID() == "" {
		return nil, errors.ThrowInvalidArgument(nil, "Org-Gf9gs", "Errors.IDMissing")
	}
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-Gh8ds", "Errors.ResourceOwnerMissing")
	}
	existingIDP, err := c.orgIDPConfigWriteModelByID(ctx, config.ID(), resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingIDP.State == domain.IDPConfigStateRemoved || existingIDP.State == domain.IDPConfigStateUnspecified {
		return nil, caos_errs.ThrowNotFound(nil, "Org-4M9so", "Errors.Org.IDPConfig.NotExisting")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingIDP.WriteModel)
	changedEvent, hasChanged := existingIDP.NewChangedEvent(
		ctx,
		orgAgg,
		config.ID(),
		config.IDPConfigName(),
		config.IDPConfigStylingType())

	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-4M9vs", "Errors.Org.LabelPolicy.NotChanged")
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingIDP, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingIDP.IDPConfigWriteModel.WriteModel), nil
}

func (c *Commands) DeactivateIDPConfig(ctx context.Context, idpID, orgID string) (*domain.ObjectDetails, error) {
	existingIDP, err := c.orgIDPConfigWriteModelByID(ctx, idpID, orgID)
	if err != nil {
		return nil, err
	}
	if existingIDP.State != domain.IDPConfigStateActive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-4M9so", "Errors.Org.IDPConfig.NotActive")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingIDP.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org_repo.NewIDPConfigDeactivatedEvent(ctx, orgAgg, idpID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingIDP, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingIDP.IDPConfigWriteModel.WriteModel), nil
}

func (c *Commands) ReactivateIDPConfig(ctx context.Context, idpID, orgID string) (*domain.ObjectDetails, error) {
	existingIDP, err := c.orgIDPConfigWriteModelByID(ctx, idpID, orgID)
	if err != nil {
		return nil, err
	}
	if existingIDP.State != domain.IDPConfigStateInactive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-5Mo0d", "Errors.Org.IDPConfig.NotInactive")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingIDP.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org_repo.NewIDPConfigReactivatedEvent(ctx, orgAgg, idpID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingIDP, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingIDP.IDPConfigWriteModel.WriteModel), nil
}

func (c *Commands) RemoveIDPConfig(ctx context.Context, idpID, orgID string, cascadeRemoveProvider bool, cascadeExternalIDPs ...*domain.ExternalIDP) (*domain.ObjectDetails, error) {
	existingIDP, err := c.orgIDPConfigWriteModelByID(ctx, idpID, orgID)
	if err != nil {
		return nil, err
	}
	events, err := c.removeIDPConfig(ctx, existingIDP, cascadeRemoveProvider, cascadeExternalIDPs...)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingIDP, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingIDP.IDPConfigWriteModel.WriteModel), nil
}

func (c *Commands) removeIDPConfig(ctx context.Context, existingIDP *OrgIDPConfigWriteModel, cascadeRemoveProvider bool, cascadeExternalIDPs ...*domain.ExternalIDP) ([]eventstore.EventPusher, error) {
	if existingIDP.State == domain.IDPConfigStateRemoved || existingIDP.State == domain.IDPConfigStateUnspecified {
		return nil, caos_errs.ThrowNotFound(nil, "Org-Yx9vd", "Errors.Org.IDPConfig.NotExisting")
	}
	if existingIDP.State != domain.IDPConfigStateInactive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-5Mo0d", "Errors.Org.IDPConfig.NotInactive")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingIDP.WriteModel)
	events := []eventstore.EventPusher{
		org_repo.NewIDPConfigRemovedEvent(ctx, orgAgg, existingIDP.AggregateID, existingIDP.Name),
	}

	if cascadeRemoveProvider {
		removeIDPEvents := c.removeIDPProviderFromLoginPolicy(ctx, orgAgg, existingIDP.AggregateID, true, cascadeExternalIDPs...)
		events = append(events, removeIDPEvents...)
	}
	return events, nil
}

func (c *Commands) getOrgIDPConfigByID(ctx context.Context, idpID, orgID string) (domain.IDPConfig, error) {
	config, err := c.orgIDPConfigWriteModelByID(ctx, idpID, orgID)
	if err != nil {
		return nil, err
	}
	if !config.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-4M9so", "Errors.Org.IDPConfig.NotExisting")
	}
	return writeModelToIDPConfig(&config.IDPConfigWriteModel), nil
}

func (c *Commands) orgIDPConfigWriteModelByID(ctx context.Context, idpID, orgID string) (policy *OrgIDPConfigWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewOrgIDPConfigWriteModel(idpID, orgID)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
