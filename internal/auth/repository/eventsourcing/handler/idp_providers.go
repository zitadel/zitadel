package handler

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/eventstore/v1/query"
	"github.com/zitadel/zitadel/internal/eventstore/v1/spooler"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
	"github.com/zitadel/zitadel/internal/iam/repository/eventsourcing/model"
	iam_view_model "github.com/zitadel/zitadel/internal/iam/repository/view/model"
	org_es_model "github.com/zitadel/zitadel/internal/org/repository/eventsourcing/model"
	query2 "github.com/zitadel/zitadel/internal/query"
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

func (_ *IDPProvider) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.IAMAggregate, org_es_model.OrgAggregate}
}

func (i *IDPProvider) CurrentSequence() (uint64, error) {
	sequence, err := i.view.GetLatestIDPProviderSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (i *IDPProvider) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := i.view.GetLatestIDPProviderSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(i.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (i *IDPProvider) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case model.IAMAggregate, org_es_model.OrgAggregate:
		err = i.processIdpProvider(event)
	}
	return err
}

func (i *IDPProvider) processIdpProvider(event *es_models.Event) (err error) {
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
		return i.view.DeleteIDPProvider(event.AggregateID, provider.IDPConfigID, event)
	case model.IDPConfigChanged, org_es_model.IDPConfigChanged:
		esConfig := new(iam_view_model.IDPConfigView)
		providerType := iam_model.IDPProviderTypeSystem
		if event.AggregateID != i.systemDefaults.IamID {
			providerType = iam_model.IDPProviderTypeOrg
		}
		esConfig.AppendEvent(providerType, event)
		providers, err := i.view.IDPProvidersByIDPConfigID(esConfig.IDPConfigID)
		if err != nil {
			return err
		}
		config := new(query2.IDP)
		if event.AggregateID == i.systemDefaults.IamID {
			config, err = i.getDefaultIDPConfig(context.TODO(), esConfig.IDPConfigID)
		} else {
			config, err = i.getOrgIDPConfig(context.TODO(), event.AggregateID, esConfig.IDPConfigID)
		}
		if err != nil {
			return err
		}
		for _, provider := range providers {
			i.fillConfigData(provider, config)
		}
		return i.view.PutIDPProviders(event, providers...)
	case org_es_model.LoginPolicyRemoved:
		return i.view.DeleteIDPProvidersByAggregateID(event.AggregateID, event)
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
		config, err = i.getDefaultIDPConfig(context.Background(), provider.IDPConfigID)
	} else {
		config, err = i.getOrgIDPConfig(context.Background(), provider.AggregateID, provider.IDPConfigID)
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
	logging.LogWithFields("SPOOL-Fjd89", "id", event.AggregateID).WithError(err).Warn("something went wrong in idp provider handler")
	return spooler.HandleError(event, err, i.view.GetLatestIDPProviderFailedEvent, i.view.ProcessedIDPProviderFailedEvent, i.view.ProcessedIDPProviderSequence, i.errorCountUntilSkip)
}

func (i *IDPProvider) OnSuccess() error {
	return spooler.HandleSuccess(i.view.UpdateIDPProviderSpoolerRunTimestamp)
}

func (i *IDPProvider) getOrgIDPConfig(ctx context.Context, aggregateID, idpConfigID string) (*query2.IDP, error) {
	return i.queries.IDPByIDAndResourceOwner(ctx, false, idpConfigID, aggregateID)
}

func (u *IDPProvider) getDefaultIDPConfig(ctx context.Context, idpConfigID string) (*query2.IDP, error) {
	return u.queries.IDPByIDAndResourceOwner(ctx, false, idpConfigID, domain.IAMID)
}
