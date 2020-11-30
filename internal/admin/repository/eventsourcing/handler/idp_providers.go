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
	iamEvents      *eventsourcing.IAMEventstore
	orgEvents      *org_events.OrgEventstore
}

const (
	idpProviderTable = "adminapi.idp_providers"
)

func (i *IDPProvider) ViewModel() string {
	return idpProviderTable
}

func (i *IDPProvider) EventQuery() (*models.SearchQuery, error) {
	sequence, err := i.view.GetLatestIDPProviderSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.IAMAggregate, org_es_model.OrgAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (i *IDPProvider) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.IAMAggregate, org_es_model.OrgAggregate:
		err = i.processIdpProvider(event)
	}
	return err
}

func (i *IDPProvider) processIdpProvider(event *models.Event) (err error) {
	provider := new(iam_view_model.IDPProviderView)
	switch event.Type {
	case model.LoginPolicyIDPProviderAdded, org_es_model.LoginPolicyIDPProviderAdded:
		err = provider.AppendEvent(event)
		if err != nil {
			return err
		}
		err = i.fillData(provider)
	case model.LoginPolicyIDPProviderRemoved, model.LoginPolicyIDPProviderCascadeRemoved,
		org_es_model.LoginPolicyIDPProviderRemoved, org_es_model.LoginPolicyIDPProviderCascadeRemoved:
		err = provider.SetData(event)
		if err != nil {
			return err
		}
		return i.view.DeleteIDPProvider(event.AggregateID, provider.IDPConfigID, event.Sequence, event.CreationDate)
	case model.IDPConfigChanged, org_es_model.IDPConfigChanged:
		esConfig := new(iam_view_model.IDPConfigView)
		providerType := iam_model.IDPProviderTypeSystem
		if event.AggregateID != i.systemDefaults.IamID {
			providerType = iam_model.IDPProviderTypeOrg
		}
		esConfig.AppendEvent(providerType, event)
		providers, err := i.view.IDPProvidersByIdpConfigID(esConfig.IDPConfigID)
		if err != nil {
			return err
		}
		config, err := i.iamEvents.GetIDPConfig(context.Background(), event.AggregateID, esConfig.IDPConfigID)
		if err != nil {
			return err
		}
		for _, provider := range providers {
			i.fillConfigData(provider, config)
		}
		return i.view.PutIDPProviders(event.Sequence, event.CreationDate, providers...)
	default:
		return i.view.ProcessedIDPProviderSequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return i.view.PutIDPProvider(provider, provider.Sequence, event.CreationDate)
}

func (i *IDPProvider) fillData(provider *iam_view_model.IDPProviderView) (err error) {
	var config *iam_model.IDPConfig
	if provider.IDPProviderType == int32(iam_model.IDPProviderTypeSystem) {
		config, err = i.iamEvents.GetIDPConfig(context.Background(), i.systemDefaults.IamID, provider.IDPConfigID)
	} else {
		config, err = i.orgEvents.GetIDPConfig(context.Background(), provider.AggregateID, provider.IDPConfigID)
	}
	if err != nil {
		return err
	}
	i.fillConfigData(provider, config)
	return nil
}

func (i *IDPProvider) fillConfigData(provider *iam_view_model.IDPProviderView, config *iam_model.IDPConfig) {
	provider.Name = config.Name
	provider.StylingType = int32(config.StylingType)
	provider.IDPConfigType = int32(config.Type)
	provider.IDPState = int32(config.State)
}

func (i *IDPProvider) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-Msj8c", "id", event.AggregateID).WithError(err).Warn("something went wrong in idp provider handler")
	return spooler.HandleError(event, err, i.view.GetLatestIDPProviderFailedEvent, i.view.ProcessedIDPProviderFailedEvent, i.view.ProcessedIDPProviderSequence, i.errorCountUntilSkip)
}

func (i *IDPProvider) OnSuccess() error {
	return spooler.HandleSuccess(i.view.UpdateIDPConfigSpoolerRunTimestamp)
}
