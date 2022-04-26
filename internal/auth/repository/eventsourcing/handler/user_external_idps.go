package handler

import (
	"context"

	"github.com/caos/logging"

	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/eventstore/v1/query"
	"github.com/zitadel/zitadel/internal/eventstore/v1/spooler"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
	iam_view_model "github.com/zitadel/zitadel/internal/iam/repository/view/model"
	query2 "github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	usr_view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

const (
	externalIDPTable = "auth.user_external_idps"
)

type ExternalIDP struct {
	handler
	systemDefaults systemdefaults.SystemDefaults
	subscription   *v1.Subscription
	queries        *query2.Queries
}

func newExternalIDP(
	handler handler,
	defaults systemdefaults.SystemDefaults,
	queries *query2.Queries,
) *ExternalIDP {
	h := &ExternalIDP{
		handler:        handler,
		systemDefaults: defaults,
		queries:        queries,
	}

	h.subscribe()

	return h
}

func (i *ExternalIDP) subscribe() {
	i.subscription = i.es.Subscribe(i.AggregateTypes()...)
	go func() {
		for event := range i.subscription.Events {
			query.ReduceEvent(i, event)
		}
	}()
}

func (i *ExternalIDP) ViewModel() string {
	return externalIDPTable
}

func (i *ExternalIDP) Subscription() *v1.Subscription {
	return i.subscription
}

func (_ *ExternalIDP) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{user.AggregateType, instance.AggregateType, org.AggregateType}
}

func (i *ExternalIDP) CurrentSequence(instanceID string) (uint64, error) {
	sequence, err := i.view.GetLatestExternalIDPSequence(instanceID)
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (i *ExternalIDP) EventQuery() (*es_models.SearchQuery, error) {
	sequences, err := i.view.GetLatestExternalIDPSequences()
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

func (i *ExternalIDP) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case user.AggregateType:
		err = i.processUser(event)
	case instance.AggregateType, org.AggregateType:
		err = i.processIdpConfig(event)
	}
	return err
}

func (i *ExternalIDP) processUser(event *es_models.Event) (err error) {
	externalIDP := new(usr_view_model.ExternalIDPView)
	switch eventstore.EventType(event.Type) {
	case user.UserIDPLinkAddedType:
		err = externalIDP.AppendEvent(event)
		if err != nil {
			return err
		}
		err = i.fillData(externalIDP)
	case user.UserIDPLinkRemovedType, user.UserIDPLinkCascadeRemovedType:
		err = externalIDP.SetData(event)
		if err != nil {
			return err
		}
		return i.view.DeleteExternalIDP(externalIDP.ExternalUserID, externalIDP.IDPConfigID, externalIDP.InstanceID, event)
	case user.UserRemovedType:
		return i.view.DeleteExternalIDPsByUserID(event.AggregateID, event.InstanceID, event)
	default:
		return i.view.ProcessedExternalIDPSequence(event)
	}
	if err != nil {
		return err
	}
	return i.view.PutExternalIDP(externalIDP, event)
}

func (i *ExternalIDP) processIdpConfig(event *es_models.Event) (err error) {
	switch eventstore.EventType(event.Type) {
	case instance.IDPConfigChangedEventType, org.IDPConfigChangedEventType:
		configView := new(iam_view_model.IDPConfigView)
		config := new(query2.IDP)
		if eventstore.EventType(event.Type) == instance.IDPConfigChangedEventType {
			configView.AppendEvent(iam_model.IDPProviderTypeSystem, event)
		} else {
			configView.AppendEvent(iam_model.IDPProviderTypeOrg, event)
		}
		exterinalIDPs, err := i.view.ExternalIDPsByIDPConfigID(configView.IDPConfigID, configView.InstanceID)
		if err != nil {
			return err
		}
		if event.AggregateType == instance.AggregateType {
			config, err = i.getDefaultIDPConfig(event.InstanceID, configView.IDPConfigID)
		} else {
			config, err = i.getOrgIDPConfig(event.InstanceID, event.AggregateID, configView.IDPConfigID)
		}
		if err != nil {
			return err
		}
		for _, provider := range exterinalIDPs {
			i.fillConfigData(provider, config)
		}
		return i.view.PutExternalIDPs(event, exterinalIDPs...)
	default:
		return i.view.ProcessedExternalIDPSequence(event)
	}
	return nil
}

func (i *ExternalIDP) fillData(externalIDP *usr_view_model.ExternalIDPView) error {
	config, err := i.getOrgIDPConfig(externalIDP.InstanceID, externalIDP.ResourceOwner, externalIDP.IDPConfigID)
	if caos_errs.IsNotFound(err) {
		config, err = i.getDefaultIDPConfig(externalIDP.InstanceID, externalIDP.IDPConfigID)
	}
	if err != nil {
		return err
	}
	i.fillConfigData(externalIDP, config)
	return nil
}

func (i *ExternalIDP) fillConfigData(externalIDP *usr_view_model.ExternalIDPView, config *query2.IDP) {
	externalIDP.IDPName = config.Name
}

func (i *ExternalIDP) OnError(event *es_models.Event, err error) error {
	logging.WithFields("id", event.AggregateID).WithError(err).Warn("something went wrong in idp provider handler")
	return spooler.HandleError(event, err, i.view.GetLatestExternalIDPFailedEvent, i.view.ProcessedExternalIDPFailedEvent, i.view.ProcessedExternalIDPSequence, i.errorCountUntilSkip)
}

func (i *ExternalIDP) OnSuccess() error {
	return spooler.HandleSuccess(i.view.UpdateExternalIDPSpoolerRunTimestamp)
}

func (i *ExternalIDP) getOrgIDPConfig(instanceID, aggregateID, idpConfigID string) (*query2.IDP, error) {
	return i.queries.IDPByIDAndResourceOwner(withInstanceID(context.Background(), instanceID), idpConfigID, aggregateID)
}

func (i *ExternalIDP) getDefaultIDPConfig(instanceID, idpConfigID string) (*query2.IDP, error) {
	return i.queries.IDPByIDAndResourceOwner(withInstanceID(context.Background(), instanceID), idpConfigID, instanceID)
}
