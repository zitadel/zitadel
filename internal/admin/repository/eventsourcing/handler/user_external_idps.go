package handler

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/org/repository/view"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	usr_view_model "github.com/caos/zitadel/internal/user/repository/view/model"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

const (
	externalIDPTable = "adminapi.user_external_idps"
)

type ExternalIDP struct {
	handler
	systemDefaults systemdefaults.SystemDefaults
	iamEvents      *eventsourcing.IAMEventstore
	subscription   *eventstore.Subscription
}

func newExternalIDP(
	handler handler,
	systemDefaults systemdefaults.SystemDefaults,
	iamEvents *eventsourcing.IAMEventstore,
) *ExternalIDP {
	h := &ExternalIDP{
		handler:        handler,
		systemDefaults: systemDefaults,
		iamEvents:      iamEvents,
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

func (i *ExternalIDP) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{model.UserAggregate, iam_es_model.IAMAggregate, org_es_model.OrgAggregate}
}

func (i *ExternalIDP) CurrentSequence() (uint64, error) {
	sequence, err := i.view.GetLatestExternalIDPSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (i *ExternalIDP) EventQuery() (*models.SearchQuery, error) {
	sequence, err := i.view.GetLatestExternalIDPSequence()
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
			config, err = i.getOrgIDPConfig(context.Background(), event.AggregateID, configView.IDPConfigID)
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
}

func (i *ExternalIDP) fillData(externalIDP *usr_view_model.ExternalIDPView) error {
	config, err := i.getOrgIDPConfig(context.Background(), externalIDP.ResourceOwner, externalIDP.IDPConfigID)
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

func (i *ExternalIDP) getOrgIDPConfig(ctx context.Context, aggregateID, idpConfigID string) (*iam_model.IDPConfig, error) {
	existing, err := i.getOrgByID(ctx, aggregateID)
	if err != nil {
		return nil, err
	}
	if _, i := existing.GetIDP(idpConfigID); i != nil {
		return i, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-2n8Fh", "Errors.Org.IdpNotExisting")
}

func (i *ExternalIDP) getOrgByID(ctx context.Context, orgID string) (*org_model.Org, error) {
	query, err := view.OrgByIDQuery(orgID, 0)
	if err != nil {
		return nil, err
	}

	var esOrg *org_es_model.Org
	err = es_sdk.Filter(ctx, i.Eventstore().FilterEvents, esOrg.AppendEvents, query)
	if err != nil && !caos_errs.IsNotFound(err) {
		return nil, err
	}
	if esOrg.Sequence == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-MOFMs", "Errors.Org.NotFound")
	}

	return org_es_model.OrgToModel(esOrg), nil
}
