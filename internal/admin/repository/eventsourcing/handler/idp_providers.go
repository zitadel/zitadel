package handler

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1"
	es_sdk "github.com/caos/zitadel/internal/eventstore/v1/sdk"
	iam_view "github.com/caos/zitadel/internal/iam/repository/view"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/org/repository/view"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

const (
	idpProviderTable = "adminapi.idp_providers"
)

type IDPProvider struct {
	handler
	systemDefaults systemdefaults.SystemDefaults
	subscription   *v1.Subscription
}

func newIDPProvider(
	handler handler,
	systemDefaults systemdefaults.SystemDefaults,
) *IDPProvider {
	h := &IDPProvider{
		handler:        handler,
		systemDefaults: systemDefaults,
	}

	h.subscribe()

	return h
}

func (i *IDPProvider) subscribe() {
	i.subscription = i.es.Subscribe(i.AggregateTypes()...)
	go func() {
		for event := range i.subscription.Events {
			query.ReduceEvent(i, event)
		}
	}()
}

func (i *IDPProvider) Subscription() *v1.Subscription {
	return i.subscription
}

func (i *IDPProvider) ViewModel() string {
	return idpProviderTable
}

func (i *IDPProvider) AggregateTypes() []es_models.AggregateType {
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
		if err = esConfig.AppendEvent(providerType, event); err != nil {
			return err
		}
		providers, err := i.view.IDPProvidersByIDPConfigID(esConfig.IDPConfigID)
		if err != nil {
			return err
		}
		config, err := i.getDefaultIDPConfig(context.Background(), esConfig.IDPConfigID)
		if err != nil {
			return err
		}
		for _, provider := range providers {
			i.fillConfigData(provider, config)
		}
		return i.view.PutIDPProviders(event, providers...)
	case model.IDPConfigDeactivated,
		model.IDPConfigReactivated:
		esConfig := new(iam_view_model.IDPConfigView)
		providerType := iam_model.IDPProviderTypeSystem
		if event.AggregateID != i.systemDefaults.IamID {
			providerType = iam_model.IDPProviderTypeOrg
		}
		if err := esConfig.AppendEvent(providerType, event); err != nil {
			return err
		}
		providers, err := i.view.IDPProvidersByIDPConfigID(esConfig.IDPConfigID)
		if err != nil {
			return err
		}
		for _, provider := range providers {
			provider.IDPState = esConfig.IDPState
		}
		return i.view.PutIDPProviders(event, providers...)
	default:
		return i.view.ProcessedIDPProviderSequence(event)
	}
	if err != nil {
		return err
	}
	return i.view.PutIDPProvider(provider, event)
}

func (i *IDPProvider) fillData(provider *iam_view_model.IDPProviderView) (err error) {
	var config *iam_model.IDPConfig
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

func (i *IDPProvider) fillConfigData(provider *iam_view_model.IDPProviderView, config *iam_model.IDPConfig) {
	provider.Name = config.Name
	provider.StylingType = int32(config.StylingType)
	provider.IDPConfigType = int32(config.Type)
	provider.IDPState = int32(config.State)
}

func (i *IDPProvider) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-Msj8c", "id", event.AggregateID).WithError(err).Warn("something went wrong in idp provider handler")
	return spooler.HandleError(event, err, i.view.GetLatestIDPProviderFailedEvent, i.view.ProcessedIDPProviderFailedEvent, i.view.ProcessedIDPProviderSequence, i.errorCountUntilSkip)
}

func (i *IDPProvider) OnSuccess() error {
	return spooler.HandleSuccess(i.view.UpdateIDPProviderSpoolerRunTimestamp)
}

func (i *IDPProvider) getOrgIDPConfig(ctx context.Context, aggregateID, idpConfigID string) (*iam_model.IDPConfig, error) {
	existing, err := i.getOrgByID(ctx, aggregateID)
	if err != nil {
		return nil, err
	}
	if _, i := existing.GetIDP(idpConfigID); i != nil {
		return i, nil
	}
	return nil, errors.ThrowNotFound(nil, "EVENT-4m0fs", "Errors.IDP.NotExisting")
}

func (i *IDPProvider) getOrgByID(ctx context.Context, orgID string) (*org_model.Org, error) {
	query, err := view.OrgByIDQuery(orgID, 0)
	if err != nil {
		return nil, err
	}

	esOrg := &org_es_model.Org{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: orgID,
		},
	}
	err = es_sdk.Filter(ctx, i.Eventstore().FilterEvents, esOrg.AppendEvents, query)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if esOrg.Sequence == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-4m9gs", "Errors.Org.NotFound")
	}

	return org_es_model.OrgToModel(esOrg), nil
}

func (u *IDPProvider) getIAMByID(ctx context.Context) (*iam_model.IAM, error) {
	query, err := iam_view.IAMByIDQuery(domain.IAMID, 0)
	if err != nil {
		return nil, err
	}
	iam := &model.IAM{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: domain.IAMID,
		},
	}
	err = es_sdk.Filter(ctx, u.Eventstore().FilterEvents, iam.AppendEvents, query)
	if err != nil && errors.IsNotFound(err) && iam.Sequence == 0 {
		return nil, err
	}
	return model.IAMToModel(iam), nil
}

func (u *IDPProvider) getDefaultIDPConfig(ctx context.Context, idpConfigID string) (*iam_model.IDPConfig, error) {
	existing, err := u.getIAMByID(ctx)
	if err != nil {
		return nil, err
	}
	if _, existingIDP := existing.GetIDP(idpConfigID); existingIDP != nil {
		return existingIDP, nil
	}
	return nil, errors.ThrowNotFound(nil, "EVENT-4M=Fs", "Errors.IDP.NotExisting")
}
