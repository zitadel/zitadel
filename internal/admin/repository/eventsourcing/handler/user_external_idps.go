package handler

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	org_es "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	usr_view_model "github.com/caos/zitadel/internal/user/repository/view/model"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

type ExternalIDP struct {
	handler
	systemDefaults systemdefaults.SystemDefaults
	iamEvents      *eventsourcing.IAMEventstore
	orgEvents      *org_es.OrgEventstore
}

const (
	externalIDPTable = "adminapi.user_external_idps"
)

func (m *ExternalIDP) ViewModel() string {
	return externalIDPTable
}

func (m *ExternalIDP) EventQuery() (*models.SearchQuery, error) {
	sequence, err := m.view.GetLatestExternalIDPSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.UserAggregate, iam_es_model.IAMAggregate, org_es_model.OrgAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *ExternalIDP) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.UserAggregate:
		err = m.processUser(event)
	case iam_es_model.IAMAggregate, org_es_model.OrgAggregate:
		err = m.processIdpConfig(event)
	}
	return err
}

func (m *ExternalIDP) processUser(event *models.Event) (err error) {
	externalIDP := new(usr_view_model.ExternalIDPView)
	switch event.Type {
	case model.HumanExternalIDPAdded:
		err = externalIDP.AppendEvent(event)
		if err != nil {
			return err
		}
		err = m.fillData(externalIDP)
	case model.HumanExternalIDPRemoved, model.HumanExternalIDPCascadeRemoved:
		err = externalIDP.SetData(event)
		if err != nil {
			return err
		}
		return m.view.DeleteExternalIDP(externalIDP.ExternalUserID, externalIDP.IDPConfigID, event.Sequence)
	case model.UserRemoved:
		return m.view.DeleteExternalIDPsByUserID(event.AggregateID, event.Sequence)
	default:
		return m.view.ProcessedExternalIDPSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return m.view.PutExternalIDP(externalIDP, externalIDP.Sequence)
}

func (m *ExternalIDP) processIdpConfig(event *models.Event) (err error) {
	switch event.Type {
	case iam_es_model.IDPConfigChanged, org_es_model.IDPConfigChanged:
		configView := new(iam_view_model.IDPConfigView)
		config := new(iam_model.IDPConfig)
		if event.Type == iam_es_model.IDPConfigChanged {
			configView.AppendEvent(iam_model.IDPProviderTypeSystem, event)
		} else {
			configView.AppendEvent(iam_model.IDPProviderTypeOrg, event)
		}
		exterinalIDPs, err := m.view.ExternalIDPsByIDPConfigID(configView.IDPConfigID)
		if err != nil {
			return err
		}
		if event.AggregateType == iam_es_model.IAMAggregate {
			config, err = m.iamEvents.GetIDPConfig(context.Background(), event.AggregateID, configView.IDPConfigID)
		} else {
			config, err = m.orgEvents.GetIDPConfig(context.Background(), event.AggregateID, configView.IDPConfigID)
		}
		if err != nil {
			return err
		}
		for _, provider := range exterinalIDPs {
			m.fillConfigData(provider, config)
		}
		return m.view.PutExternalIDPs(event.Sequence, exterinalIDPs...)
	default:
		return m.view.ProcessedExternalIDPSequence(event.Sequence)
	}
	return nil
}

func (m *ExternalIDP) fillData(externalIDP *usr_view_model.ExternalIDPView) error {
	config, err := m.orgEvents.GetIDPConfig(context.Background(), externalIDP.ResourceOwner, externalIDP.IDPConfigID)
	if caos_errs.IsNotFound(err) {
		config, err = m.iamEvents.GetIDPConfig(context.Background(), m.systemDefaults.IamID, externalIDP.IDPConfigID)
	}
	if err != nil {
		return err
	}
	m.fillConfigData(externalIDP, config)
	return nil
}

func (m *ExternalIDP) fillConfigData(externalIDP *usr_view_model.ExternalIDPView, config *iam_model.IDPConfig) {
	externalIDP.IDPName = config.Name
}

func (m *ExternalIDP) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-4Rsu8", "id", event.AggregateID).WithError(err).Warn("something went wrong in idp provider handler")
	return spooler.HandleError(event, err, m.view.GetLatestExternalIDPFailedEvent, m.view.ProcessedExternalIDPFailedEvent, m.view.ProcessedExternalIDPSequence, m.errorCountUntilSkip)
}
