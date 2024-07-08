package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	org_repo "github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) ImportIDPConfig(ctx context.Context, config *domain.IDPConfig, idpConfigID, resourceOwner string) (_ *domain.IDPConfig, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	existingIDP, err := c.orgIDPConfigWriteModelByID(ctx, idpConfigID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingIDP.State != domain.IDPConfigStateRemoved && existingIDP.State != domain.IDPConfigStateUnspecified {
		return nil, zerrors.ThrowNotFound(nil, "Org-1J8fs", "Errors.Org.IDPConfig.AlreadyExisting")
	}
	return c.addIDPConfig(ctx, config, idpConfigID, resourceOwner)
}

func (c *Commands) AddIDPConfig(ctx context.Context, config *domain.IDPConfig, resourceOwner string) (*domain.IDPConfig, error) {
	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-0j8gs", "Errors.ResourceOwnerMissing")
	}
	if config.OIDCConfig == nil && config.JWTConfig == nil {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-eUpQU", "Errors.idp.config.notset")
	}
	idpConfigID, err := id_generator.Next()
	if err != nil {
		return nil, err
	}

	return c.addIDPConfig(ctx, config, idpConfigID, resourceOwner)
}

func (c *Commands) addIDPConfig(ctx context.Context, config *domain.IDPConfig, idpConfigID, resourceOwner string) (*domain.IDPConfig, error) {

	addedConfig := NewOrgIDPConfigWriteModel(idpConfigID, resourceOwner)

	orgAgg := OrgAggregateFromWriteModel(&addedConfig.WriteModel)
	events := []eventstore.Command{
		org_repo.NewIDPConfigAddedEvent(
			ctx,
			orgAgg,
			idpConfigID,
			config.Name,
			config.Type,
			config.StylingType,
			config.AutoRegister,
		),
	}
	if config.OIDCConfig != nil {
		clientSecret, err := crypto.Crypt([]byte(config.OIDCConfig.ClientSecretString), c.idpConfigEncryption)
		if err != nil {
			return nil, err
		}
		events = append(events, org_repo.NewIDPOIDCConfigAddedEvent(
			ctx,
			orgAgg,
			config.OIDCConfig.ClientID,
			idpConfigID,
			config.OIDCConfig.Issuer,
			config.OIDCConfig.AuthorizationEndpoint,
			config.OIDCConfig.TokenEndpoint,
			clientSecret,
			config.OIDCConfig.IDPDisplayNameMapping,
			config.OIDCConfig.UsernameMapping,
			config.OIDCConfig.Scopes...))
	} else if config.JWTConfig != nil {
		events = append(events, org_repo.NewIDPJWTConfigAddedEvent(
			ctx,
			orgAgg,
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

func (c *Commands) ChangeIDPConfig(ctx context.Context, config *domain.IDPConfig, resourceOwner string) (*domain.IDPConfig, error) {
	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-Gh8ds", "Errors.ResourceOwnerMissing")
	}
	existingIDP, err := c.orgIDPConfigWriteModelByID(ctx, config.IDPConfigID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingIDP.State == domain.IDPConfigStateRemoved || existingIDP.State == domain.IDPConfigStateUnspecified {
		return nil, zerrors.ThrowNotFound(nil, "Org-1J9fs", "Errors.Org.IDPConfig.NotExisting")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingIDP.WriteModel)
	changedEvent, hasChanged := existingIDP.NewChangedEvent(
		ctx,
		orgAgg,
		config.IDPConfigID,
		config.Name,
		config.StylingType,
		config.AutoRegister)

	if !hasChanged {
		return nil, zerrors.ThrowPreconditionFailed(nil, "Org-jf9w", "Errors.Org.IDPConfig.NotChanged")
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

func (c *Commands) DeactivateIDPConfig(ctx context.Context, idpID, orgID string) (*domain.ObjectDetails, error) {
	existingIDP, err := c.orgIDPConfigWriteModelByID(ctx, idpID, orgID)
	if err != nil {
		return nil, err
	}
	if existingIDP.State != domain.IDPConfigStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "Org-BBmd0", "Errors.Org.IDPConfig.NotActive")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingIDP.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org_repo.NewIDPConfigDeactivatedEvent(ctx, orgAgg, idpID))
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
		return nil, zerrors.ThrowPreconditionFailed(nil, "Org-5Mo0d", "Errors.Org.IDPConfig.NotInactive")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingIDP.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org_repo.NewIDPConfigReactivatedEvent(ctx, orgAgg, idpID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingIDP, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingIDP.IDPConfigWriteModel.WriteModel), nil
}

func (c *Commands) RemoveIDPConfig(ctx context.Context, idpID, orgID string, cascadeRemoveProvider bool, cascadeExternalIDPs ...*domain.UserIDPLink) (*domain.ObjectDetails, error) {
	existingIDP, err := c.orgIDPConfigWriteModelByID(ctx, idpID, orgID)
	if err != nil {
		return nil, err
	}
	events, err := c.removeIDPConfig(ctx, existingIDP, cascadeRemoveProvider, cascadeExternalIDPs...)
	if err != nil {
		return nil, err
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

func (c *Commands) removeIDPConfig(ctx context.Context, existingIDP *OrgIDPConfigWriteModel, cascadeRemoveProvider bool, cascadeExternalIDPs ...*domain.UserIDPLink) ([]eventstore.Command, error) {
	if existingIDP.State == domain.IDPConfigStateRemoved || existingIDP.State == domain.IDPConfigStateUnspecified {
		return nil, zerrors.ThrowNotFound(nil, "Org-Yx9vd", "Errors.Org.IDPConfig.NotExisting")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingIDP.WriteModel)
	events := []eventstore.Command{
		org_repo.NewIDPConfigRemovedEvent(ctx, orgAgg, existingIDP.ConfigID, existingIDP.Name),
	}

	if cascadeRemoveProvider {
		removeIDPEvents := c.removeIDPFromLoginPolicy(ctx, orgAgg, existingIDP.ConfigID, true, cascadeExternalIDPs...)
		events = append(events, removeIDPEvents...)
	}
	return events, nil
}

func (c *Commands) getOrgIDPConfigByID(ctx context.Context, idpID, orgID string) (*domain.IDPConfig, error) {
	config, err := c.orgIDPConfigWriteModelByID(ctx, idpID, orgID)
	if err != nil {
		return nil, err
	}
	if !config.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "ORG-2m90f", "Errors.Org.IDPConfig.NotExisting")
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
