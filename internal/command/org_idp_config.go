package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/telemetry/tracing"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	org_repo "github.com/caos/zitadel/internal/repository/org"
)

func (c *Commands) AddIDPConfig(ctx context.Context, config *domain.IDPConfig) (*domain.IDPConfig, error) {
	if config.OIDCConfig == nil {
		return nil, errors.ThrowInvalidArgument(nil, "Org-eUpQU", "Errors.idp.config.notset")
	}

	idpConfigID, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	addedConfig := NewOrgIDPConfigWriteModel(idpConfigID, config.AggregateID)

	clientSecret, err := crypto.Crypt([]byte(config.OIDCConfig.ClientSecretString), c.idpConfigSecretCrypto)
	if err != nil {
		return nil, err
	}

	orgAgg := OrgAggregateFromWriteModel(&addedConfig.WriteModel)
	events := []eventstore.EventPusher{
		org_repo.NewIDPConfigAddedEvent(
			ctx,
			orgAgg,
			idpConfigID,
			config.Name,
			config.Type,
			config.StylingType,
		),
		org_repo.NewIDPOIDCConfigAddedEvent(
			ctx,
			orgAgg,
			config.OIDCConfig.ClientID,
			idpConfigID,
			config.OIDCConfig.Issuer,
			clientSecret,
			config.OIDCConfig.IDPDisplayNameMapping,
			config.OIDCConfig.UsernameMapping,
			config.OIDCConfig.Scopes...),
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedConfig, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToIDPConfig(&addedConfig.IDPConfigWriteModel), nil
}

func (c *Commands) ChangeIDPConfig(ctx context.Context, config *domain.IDPConfig) (*domain.IDPConfig, error) {
	existingIDP, err := c.orgIDPConfigWriteModelByID(ctx, config.IDPConfigID, config.AggregateID)
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
		config.IDPConfigID,
		config.Name,
		config.StylingType)

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
	return writeModelToIDPConfig(&existingIDP.IDPConfigWriteModel), nil
}

func (c *Commands) DeactivateIDPConfig(ctx context.Context, idpID, orgID string) error {
	existingIDP, err := c.orgIDPConfigWriteModelByID(ctx, idpID, orgID)
	if err != nil {
		return err
	}
	if existingIDP.State != domain.IDPConfigStateActive {
		return caos_errs.ThrowPreconditionFailed(nil, "Org-4M9so", "Errors.Org.IDPConfig.NotActive")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingIDP.WriteModel)
	_, err = c.eventstore.PushEvents(ctx, org_repo.NewIDPConfigDeactivatedEvent(ctx, orgAgg, idpID))
	return err
}

func (c *Commands) ReactivateIDPConfig(ctx context.Context, idpID, orgID string) error {
	existingIDP, err := c.orgIDPConfigWriteModelByID(ctx, idpID, orgID)
	if err != nil {
		return err
	}
	if existingIDP.State != domain.IDPConfigStateInactive {
		return caos_errs.ThrowPreconditionFailed(nil, "Org-5Mo0d", "Errors.Org.IDPConfig.NotInactive")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingIDP.WriteModel)
	_, err = c.eventstore.PushEvents(ctx, org_repo.NewIDPConfigReactivatedEvent(ctx, orgAgg, idpID))
	return err
}

func (c *Commands) RemoveIDPConfig(ctx context.Context, idpID, orgID string, cascadeRemoveProvider bool, cascadeExternalIDPs ...*domain.ExternalIDP) error {
	existingIDP, err := c.orgIDPConfigWriteModelByID(ctx, idpID, orgID)
	if err != nil {
		return err
	}

	if existingIDP.State == domain.IDPConfigStateRemoved || existingIDP.State == domain.IDPConfigStateUnspecified {
		return caos_errs.ThrowNotFound(nil, "Org-Yx9vd", "Errors.Org.IDPConfig.NotExisting")
	}
	if existingIDP.State != domain.IDPConfigStateInactive {
		return caos_errs.ThrowPreconditionFailed(nil, "Org-5Mo0d", "Errors.Org.IDPConfig.NotInactive")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingIDP.WriteModel)
	events := []eventstore.EventPusher{
		org_repo.NewIDPConfigRemovedEvent(ctx, orgAgg, idpID, existingIDP.Name),
	}

	if cascadeRemoveProvider {
		removeIDPEvents := c.removeIDPProviderFromLoginPolicy(ctx, orgAgg, idpID, true, cascadeExternalIDPs...)
		events = append(events, removeIDPEvents...)
	}
	_, err = c.eventstore.PushEvents(ctx, events...)
	return err
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
