package handler

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	org_es "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
)

type IdpProvider struct {
	handler
	systemDefaults systemdefaults.SystemDefaults
	iamEvents      *eventsourcing.IamEventstore
	orgEvents      *org_es.OrgEventstore
}

const (
	idpProviderTable = "management.idp_providers"
)

func (m *IdpProvider) ViewModel() string {
	return idpProviderTable
}

func (m *IdpProvider) EventQuery() (*models.SearchQuery, error) {
	sequence, err := m.view.GetLatestIdpProviderSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.IamAggregate, org_es_model.OrgAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *IdpProvider) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.IamAggregate, org_es_model.OrgAggregate:
		err = m.processIdpProvider(event)
	}
	return err
}

func (m *IdpProvider) processIdpProvider(event *models.Event) (err error) {
	provider := new(iam_view_model.IdpProviderView)
	switch event.Type {
	case model.LoginPolicyIdpProviderAdded, org_es_model.LoginPolicyIdpProviderAdded:
		err = provider.AppendEvent(event)
		if err != nil {
			return err
		}
		err = m.fillData(provider)
	case model.LoginPolicyIdpProviderRemoved, model.LoginPolicyIdpProviderCascadeRemoved,
		org_es_model.LoginPolicyIdpProviderRemoved, org_es_model.LoginPolicyIdpProviderCascadeRemoved:
		err = provider.SetData(event)
		if err != nil {
			return err
		}
		return m.view.DeleteIdpProvider(event.AggregateID, provider.IdpConfigID, event.Sequence)
	case model.IdpConfigChanged, org_es_model.IdpConfigChanged:
		config := new(iam_model.IdpConfig)
		config.AppendEvent(event)
		providers, err := m.view.IdpProvidersByIdpConfigID(event.AggregateID, config.IDPConfigID)
		if err != nil {
			return err
		}
		if provider.IdpProviderType == int32(iam_model.IdpProviderTypeSystem) {
			config, err = m.iamEvents.GetIdpConfiguration(context.Background(), provider.AggregateID, config.IDPConfigID)
		} else {
			config, err = m.orgEvents.GetIdpConfiguration(context.Background(), provider.AggregateID, provider.IdpConfigID)
		}
		if err != nil {
			return err
		}
		for _, provider := range providers {
			m.fillConfigData(provider, config)
		}
		return m.view.PutIdpProviders(event.Sequence, providers...)
	case org_es_model.LoginPolicyRemoved:
		return m.view.DeleteIdpProvidersByAggregateID(event.AggregateID, event.Sequence)
	default:
		return m.view.ProcessedIdpProviderSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return m.view.PutIdpProvider(provider, provider.Sequence)
}

func (m *IdpProvider) fillData(provider *iam_view_model.IdpProviderView) (err error) {
	var config *iam_model.IdpConfig
	if provider.IdpProviderType == int32(iam_model.IdpProviderTypeSystem) {
		config, err = m.iamEvents.GetIdpConfiguration(context.Background(), m.systemDefaults.IamID, provider.IdpConfigID)
	} else {
		config, err = m.orgEvents.GetIdpConfiguration(context.Background(), provider.AggregateID, provider.IdpConfigID)
	}
	if err != nil {
		return err
	}
	m.fillConfigData(provider, config)
	return nil
}

func (m *IdpProvider) fillConfigData(provider *iam_view_model.IdpProviderView, config *iam_model.IdpConfig) {
	provider.Name = config.Name
	provider.IdpConfigType = int32(config.Type)
}

func (m *IdpProvider) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-Msj8c", "id", event.AggregateID).WithError(err).Warn("something went wrong in idp provider handler")
	return spooler.HandleError(event, err, m.view.GetLatestIdpProviderFailedEvent, m.view.ProcessedIdpProviderFailedEvent, m.view.ProcessedIdpProviderSequence, m.errorCountUntilSkip)
}
