package handler

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	query2 "github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/org"
)

const (
	idpProviderTable = "auth.idp_providers"
)

type IDPProvider struct {
	handler
	systemDefaults systemdefaults.SystemDefaults
	subscription   *v1.Subscription
	queries        *query2.Queries
}

func newIDPProvider(
	h handler,
	defaults systemdefaults.SystemDefaults,
	queries *query2.Queries,
) *IDPProvider {
	idpProvider := &IDPProvider{
		handler:        h,
		systemDefaults: defaults,
		queries:        queries,
	}

	idpProvider.subscribe()

	return idpProvider
}

func (i *IDPProvider) subscribe() {
	i.subscription = i.es.Subscribe(i.AggregateTypes()...)
	go func() {
		for event := range i.subscription.Events {
			query.ReduceEvent(i, event)
		}
	}()
}

func (i *IDPProvider) ViewModel() string {
	return idpProviderTable
}

func (i *IDPProvider) Subscription() *v1.Subscription {
	return i.subscription
}

func (_ *IDPProvider) AggregateTypes() []models.AggregateType {
	return []es_models.AggregateType{instance.AggregateType, org.AggregateType}
}

func (i *IDPProvider) CurrentSequence(instanceID string) (uint64, error) {
	sequence, err := i.view.GetLatestIDPProviderSequence(instanceID)
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (i *IDPProvider) EventQuery() (*models.SearchQuery, error) {
	sequences, err := i.view.GetLatestIDPProviderSequences()
	if err != nil {
		return nil, err
	}
	query := es_models.NewSearchQuery()
	instances := make([]string, 0)
	for _, sequence := range sequences {
		for _, instance := range instances {
			if sequence.InstanceID == instance {
				break
			}
		}
		instances = append(instances, sequence.InstanceID)
		query.AddQuery().
			AggregateTypeFilter(i.AggregateTypes()...).
			LatestSequenceFilter(sequence.CurrentSequence).
			InstanceIDFilter(sequence.InstanceID)
	}
	return query.AddQuery().
		AggregateTypeFilter(i.AggregateTypes()...).
		LatestSequenceFilter(0).
		ExcludedInstanceIDsFilter(instances...).
		SearchQuery(), nil
}

func (i *IDPProvider) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case instance.AggregateType, org.AggregateType:
		err = i.processIdpProvider(event)
	}
	return err
}

func (i *IDPProvider) processIdpProvider(event *models.Event) (err error) {
	provider := new(iam_view_model.IDPProviderView)
	switch eventstore.EventType(event.Type) {
	case instance.LoginPolicyIDPProviderAddedEventType, org.LoginPolicyIDPProviderAddedEventType:
		err = provider.AppendEvent(event)
		if err != nil {
			return err
		}
		err = i.fillData(provider)
	case instance.LoginPolicyIDPProviderRemovedEventType, instance.LoginPolicyIDPProviderCascadeRemovedEventType,
		org.LoginPolicyIDPProviderRemovedEventType, org.LoginPolicyIDPProviderCascadeRemovedEventType:
		err = provider.SetData(event)
		if err != nil {
			return err
		}
		return i.view.DeleteIDPProvider(event.AggregateID, provider.IDPConfigID, event.InstanceID, event)
	case instance.IDPConfigChangedEventType, org.IDPConfigChangedEventType:
		esConfig := new(iam_view_model.IDPConfigView)
		providerType := iam_model.IDPProviderTypeSystem
		if event.AggregateID != event.InstanceID {
			providerType = iam_model.IDPProviderTypeOrg
		}
		esConfig.AppendEvent(providerType, event)
		providers, err := i.view.IDPProvidersByIDPConfigID(esConfig.IDPConfigID, esConfig.InstanceID)
		if err != nil {
			return err
		}
		config := new(query2.IDP)
		if event.AggregateID == event.InstanceID {
			config, err = i.getDefaultIDPConfig(event.InstanceID, esConfig.IDPConfigID)
		} else {
			config, err = i.getOrgIDPConfig(event.InstanceID, event.AggregateID, esConfig.IDPConfigID)
		}
		if err != nil {
			return err
		}
		for _, provider := range providers {
			i.fillConfigData(provider, config)
		}
		return i.view.PutIDPProviders(event, providers...)
	case org.LoginPolicyRemovedEventType:
		return i.view.DeleteIDPProvidersByAggregateID(event.AggregateID, event.InstanceID, event)
	default:
		return i.view.ProcessedIDPProviderSequence(event)
	}
	if err != nil {
		return err
	}
	return i.view.PutIDPProvider(provider, event)
}

func (i *IDPProvider) fillData(provider *iam_view_model.IDPProviderView) (err error) {
	var config *query2.IDP
	if provider.IDPProviderType == int32(iam_model.IDPProviderTypeSystem) {
		config, err = i.getDefaultIDPConfig(provider.InstanceID, provider.IDPConfigID)
	} else {
		config, err = i.getOrgIDPConfig(provider.InstanceID, provider.AggregateID, provider.IDPConfigID)
	}
	if err != nil {
		return err
	}
	i.fillConfigData(provider, config)
	return nil
}

func (i *IDPProvider) fillConfigData(provider *iam_view_model.IDPProviderView, config *query2.IDP) {
	provider.Name = config.Name
	provider.StylingType = int32(config.StylingType)
	if config.OIDCIDP != nil {
		provider.IDPConfigType = int32(domain.IDPConfigTypeOIDC)
	} else if config.JWTIDP != nil {
		provider.IDPConfigType = int32(domain.IDPConfigTypeJWT)
	}
	switch config.State {
	case domain.IDPConfigStateActive:
		provider.IDPState = int32(iam_model.IDPConfigStateActive)
	case domain.IDPConfigStateInactive:
		provider.IDPState = int32(iam_model.IDPConfigStateActive)
	case domain.IDPConfigStateRemoved:
		provider.IDPState = int32(iam_model.IDPConfigStateRemoved)
	default:
		provider.IDPState = int32(iam_model.IDPConfigStateActive)
	}
}

func (i *IDPProvider) OnError(event *es_models.Event, err error) error {
	logging.WithFields("id", event.AggregateID).WithError(err).Warn("something went wrong in idp provider handler")
	return spooler.HandleError(event, err, i.view.GetLatestIDPProviderFailedEvent, i.view.ProcessedIDPProviderFailedEvent, i.view.ProcessedIDPProviderSequence, i.errorCountUntilSkip)
}

func (i *IDPProvider) OnSuccess() error {
	return spooler.HandleSuccess(i.view.UpdateIDPProviderSpoolerRunTimestamp)
}

func (i *IDPProvider) getOrgIDPConfig(instanceID, aggregateID, idpConfigID string) (*query2.IDP, error) {
	return i.queries.IDPByIDAndResourceOwner(withInstanceID(context.Background(), instanceID), idpConfigID, aggregateID)
}

func (u *IDPProvider) getDefaultIDPConfig(instanceID, idpConfigID string) (*query2.IDP, error) {
	return u.queries.IDPByIDAndResourceOwner(withInstanceID(context.Background(), instanceID), idpConfigID, instanceID)
}
