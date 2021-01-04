package handler

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	org_es "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	usr_view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

const (
	externalIDPTable = "auth.user_external_idps"
)

type ExternalIDP struct {
	handler
	systemDefaults systemdefaults.SystemDefaults
	iamEvents      *eventsourcing.IAMEventstore
	orgEvents      *org_es.OrgEventstore
	subscription   *eventstore.Subscription
}

func newExternalIDP(
	handler handler,
	defaults systemdefaults.SystemDefaults,
	iamEvents *eventsourcing.IAMEventstore,
	orgEvents *org_es.OrgEventstore,
) *ExternalIDP {
	h := &ExternalIDP{
		handler:        handler,
		systemDefaults: defaults,
		iamEvents:      iamEvents,
		orgEvents:      orgEvents,
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

func (_ *ExternalIDP) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.UserAggregate, iam_es_model.IAMAggregate, org_es_model.OrgAggregate}
}

func (i *ExternalIDP) CurrentSequence(event *models.Event) (uint64, error) {
	sequence, err := i.view.GetLatestExternalIDPSequence(string(event.AggregateType))
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (i *ExternalIDP) EventQuery() (*models.SearchQuery, error) {
	sequence, err := i.view.GetLatestExternalIDPSequence("")
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(i.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (i *ExternalIDP) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.UserAggregate:
		err = i.processUser(event)
	case iam_es_model.IAMAggregate, org_es_model.OrgAggregate:
		err = i.processIdpConfig(event)
	}
	return err
}

func (i *ExternalIDP) processUser(event *models.Event) (err error) {
	externalIDP := new(usr_view_model.ExternalIDPView)
	switch event.Type {
	case model.HumanExternalIDPAdded:
		err = externalIDP.AppendEvent(event)
		if err != nil {
			return err
		}
		err = i.fillData(externalIDP)
	case model.HumanExternalIDPRemoved, model.HumanExternalIDPCascadeRemoved:
		err = externalIDP.SetData(event)
		if err != nil {
			return err
		}
		return i.view.DeleteExternalIDP(externalIDP.ExternalUserID, externalIDP.IDPConfigID, event)
	case model.UserRemoved:
		return i.view.DeleteExternalIDPsByUserID(event.AggregateID, event)
	default:
		return i.view.ProcessedExternalIDPSequence(event)
	}
	if err != nil {
		return err
	}
	return i.view.PutExternalIDP(externalIDP, event)
}

func (i *ExternalIDP) processIdpConfig(event *models.Event) (err error) {
	switch event.Type {
	case iam_es_model.IDPConfigChanged, org_es_model.IDPConfigChanged:
		configView := new(iam_view_model.IDPConfigView)
		config := new(iam_model.IDPConfig)
		if event.Type == iam_es_model.IDPConfigChanged {
			configView.AppendEvent(iam_model.IDPProviderTypeSystem, event)
		} else {
			configView.AppendEvent(iam_model.IDPProviderTypeOrg, event)
		}
		exterinalIDPs, err := i.view.ExternalIDPsByIDPConfigID(configView.IDPConfigID)
		if err != nil {
			return err
		}
		if event.AggregateType == iam_es_model.IAMAggregate {
			config, err = i.iamEvents.GetIDPConfig(context.Background(), event.AggregateID, configView.IDPConfigID)
		} else {
			config, err = i.orgEvents.GetIDPConfig(context.Background(), event.AggregateID, configView.IDPConfigID)
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
	config, err := i.orgEvents.GetIDPConfig(context.Background(), externalIDP.ResourceOwner, externalIDP.IDPConfigID)
	if caos_errs.IsNotFound(err) {
		config, err = i.iamEvents.GetIDPConfig(context.Background(), i.systemDefaults.IamID, externalIDP.IDPConfigID)
	}
	if err != nil {
		return err
	}
	i.fillConfigData(externalIDP, config)
	return nil
}

func (i *ExternalIDP) fillConfigData(externalIDP *usr_view_model.ExternalIDPView, config *iam_model.IDPConfig) {
	externalIDP.IDPName = config.Name
}

func (i *ExternalIDP) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-4Rsu8", "id", event.AggregateID).WithError(err).Warn("something went wrong in idp provider handler")
	return spooler.HandleError(event, err, i.view.GetLatestExternalIDPFailedEvent, i.view.ProcessedExternalIDPFailedEvent, i.view.ProcessedExternalIDPSequence, i.errorCountUntilSkip)
}

func (i *ExternalIDP) OnSuccess() error {
	return spooler.HandleSuccess(i.view.UpdateExternalIDPSpoolerRunTimestamp)
}
