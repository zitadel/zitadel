package handler

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	org_events "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
)

type IDPProvider struct {
	handler
	systemDefaults systemdefaults.SystemDefaults
	iamEvents      *eventsourcing.IamEventstore
	orgEvents      *org_events.OrgEventstore
}

const (
	idpProviderTable = "adminapi.idp_providers"
)

func (m *IDPProvider) ViewModel() string {
	return idpProviderTable
}

func (m *IDPProvider) EventQuery() (*models.SearchQuery, error) {
	sequence, err := m.view.GetLatestIDPProviderSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.IamAggregate, org_es_model.OrgAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *IDPProvider) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.IamAggregate, org_es_model.OrgAggregate:
		err = m.processIdpProvider(event)
	}
	return err
}

func (m *IDPProvider) processIdpProvider(event *models.Event) (err error) {
	provider := new(iam_view_model.IDPProviderView)
	switch event.Type {
	case model.LoginPolicyIDPProviderAdded, org_es_model.LoginPolicyIDPProviderAdded:
		err = provider.AppendEvent(event)
		if err != nil {
			return err
		}
		err = m.fillData(provider)
	case model.LoginPolicyIDPProviderRemoved, model.LoginPolicyIDPProviderCascadeRemoved,
		org_es_model.LoginPolicyIDPProviderRemoved, org_es_model.LoginPolicyIDPProviderCascadeRemoved:
		err = provider.SetData(event)
		if err != nil {
			return err
		}
		return m.view.DeleteIDPProvider(event.AggregateID, provider.IDPConfigID, event.Sequence)
	case model.IDPConfigChanged, org_es_model.IDPConfigChanged:
		config := new(iam_model.IDPConfig)
		config.AppendEvent(event)
		providers, err := m.view.IDPProvidersByIdpConfigID(config.IDPConfigID)
		if err != nil {
			return err
		}
		config, err = m.iamEvents.GetIDPConfig(context.Background(), provider.AggregateID, config.IDPConfigID)
		if err != nil {
			return err
		}
		for _, provider := range providers {
			m.fillConfigData(provider, config)
		}
		return m.view.PutIDPProviders(event.Sequence, providers...)
	default:
		return m.view.ProcessedIDPProviderSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return m.view.PutIDPProvider(provider, provider.Sequence)
}

func (m *IDPProvider) fillData(provider *iam_view_model.IDPProviderView) (err error) {
	var config *iam_model.IDPConfig
	if provider.IDPProviderType == int32(iam_model.IDPProviderTypeSystem) {
		config, err = m.iamEvents.GetIDPConfig(context.Background(), m.systemDefaults.IamID, provider.IDPConfigID)
	} else {
		config, err = m.orgEvents.GetIDPConfig(context.Background(), provider.AggregateID, provider.IDPConfigID)
	}
	if err != nil {
		return err
	}
	m.fillConfigData(provider, config)
	return nil
}

func (m *IDPProvider) fillConfigData(provider *iam_view_model.IDPProviderView, config *iam_model.IDPConfig) {
	provider.Name = config.Name
	provider.IDPConfigType = int32(config.Type)
}

func (m *IDPProvider) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-Msj8c", "id", event.AggregateID).WithError(err).Warn("something went wrong in idp provider handler")
	return spooler.HandleError(event, err, m.view.GetLatestIDPProviderFailedEvent, m.view.ProcessedIDPProviderFailedEvent, m.view.ProcessedIDPProviderSequence, m.errorCountUntilSkip)
}
