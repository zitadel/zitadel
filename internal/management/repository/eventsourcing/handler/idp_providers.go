package handler

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_view "github.com/caos/zitadel/internal/iam/repository/view"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/org/repository/view"
	"github.com/caos/zitadel/internal/v2/domain"
)

const (
	idpProviderTable = "management.idp_providers"
)

type IDPProvider struct {
	handler
	systemDefaults systemdefaults.SystemDefaults
	subscription   *eventstore.Subscription
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

func (m *IDPProvider) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

func (m *IDPProvider) ViewModel() string {
	return idpProviderTable
}

func (_ *IDPProvider) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.IAMAggregate, org_es_model.OrgAggregate}
}

func (m *IDPProvider) CurrentSequence() (uint64, error) {
	sequence, err := m.view.GetLatestIDPProviderSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (m *IDPProvider) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := m.view.GetLatestIDPProviderSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(m.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *IDPProvider) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case model.IAMAggregate, org_es_model.OrgAggregate:
		err = m.processIdpProvider(event)
	}
	return err
}

func (m *IDPProvider) processIdpProvider(event *es_models.Event) (err error) {
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
		return m.view.DeleteIDPProvider(event.AggregateID, provider.IDPConfigID, event)
	case model.IDPConfigChanged, org_es_model.IDPConfigChanged:
		esConfig := new(iam_view_model.IDPConfigView)
		providerType := iam_model.IDPProviderTypeSystem
		if event.AggregateID != m.systemDefaults.IamID {
			providerType = iam_model.IDPProviderTypeOrg
		}
		esConfig.AppendEvent(providerType, event)
		providers, err := m.view.IDPProvidersByIdpConfigID(event.AggregateID, esConfig.IDPConfigID)
		if err != nil {
			return err
		}
		config := new(iam_model.IDPConfig)
		if event.AggregateID == m.systemDefaults.IamID {
			config, err = m.getDefaultIDPConfig(context.Background(), esConfig.IDPConfigID)
		} else {
			config, err = m.getOrgIDPConfig(context.Background(), event.AggregateID, esConfig.IDPConfigID)
		}
		if err != nil {
			return err
		}
		for _, provider := range providers {
			m.fillConfigData(provider, config)
		}
		return m.view.PutIDPProviders(event, providers...)
	case org_es_model.LoginPolicyRemoved:
		return m.view.DeleteIDPProvidersByAggregateID(event.AggregateID, event)
	default:
		return m.view.ProcessedIDPProviderSequence(event)
	}
	if err != nil {
		return err
	}
	return m.view.PutIDPProvider(provider, event)
}

func (m *IDPProvider) fillData(provider *iam_view_model.IDPProviderView) (err error) {
	var config *iam_model.IDPConfig
	if provider.IDPProviderType == int32(iam_model.IDPProviderTypeSystem) {
		config, err = m.getDefaultIDPConfig(context.Background(), provider.IDPConfigID)
	} else {
		config, err = m.getOrgIDPConfig(context.Background(), provider.AggregateID, provider.IDPConfigID)
	}
	if err != nil {
		return err
	}
	m.fillConfigData(provider, config)
	return nil
}

func (m *IDPProvider) fillConfigData(provider *iam_view_model.IDPProviderView, config *iam_model.IDPConfig) {
	provider.Name = config.Name
	provider.StylingType = int32(config.StylingType)
	provider.IDPConfigType = int32(config.Type)
	provider.IDPState = int32(config.State)
}

func (m *IDPProvider) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-Msj8c", "id", event.AggregateID).WithError(err).Warn("something went wrong in idp provider handler")
	return spooler.HandleError(event, err, m.view.GetLatestIDPProviderFailedEvent, m.view.ProcessedIDPProviderFailedEvent, m.view.ProcessedIDPProviderSequence, m.errorCountUntilSkip)
}

func (m *IDPProvider) OnSuccess() error {
	return spooler.HandleSuccess(m.view.UpdateIDPProviderSpoolerRunTimestamp)
}

func (i *IDPProvider) getOrgIDPConfig(ctx context.Context, aggregateID, idpConfigID string) (*iam_model.IDPConfig, error) {
	existing, err := i.getOrgByID(ctx, aggregateID)
	if err != nil {
		return nil, err
	}
	if _, i := existing.GetIDP(idpConfigID); i != nil {
		return i, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-2m9dS", "Errors.Org.IdpNotExisting")
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
	if err != nil && !caos_errs.IsNotFound(err) {
		return nil, err
	}
	if esOrg.Sequence == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-4n9fs", "Errors.Org.NotFound")
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
	if err != nil && caos_errs.IsNotFound(err) && iam.Sequence == 0 {
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
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-4M0fs", "Errors.IAM.IdpNotExisting")
}
