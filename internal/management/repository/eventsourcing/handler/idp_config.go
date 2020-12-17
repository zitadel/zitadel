package handler

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

type IDPConfig struct {
	handler
}

const (
	idpConfigTable = "management.idp_configs"
)

func (m *IDPConfig) ViewModel() string {
	return idpConfigTable
}

func (_ *IDPConfig) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (m *IDPConfig) CurrentSequence(event *models.Event) (uint64, error) {
	sequence, err := m.view.GetLatestIDPConfigSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (m *IDPConfig) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := m.view.GetLatestIDPConfigSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(m.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *IDPConfig) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate:
		err = m.processIdpConfig(iam_model.IDPProviderTypeOrg, event)
	case iam_es_model.IAMAggregate:
		err = m.processIdpConfig(iam_model.IDPProviderTypeSystem, event)
	}
	return err
}

func (m *IDPConfig) processIdpConfig(providerType iam_model.IDPProviderType, event *es_models.Event) (err error) {
	idp := new(iam_view_model.IDPConfigView)
	switch event.Type {
	case model.IDPConfigAdded,
		iam_es_model.IDPConfigAdded:
		err = idp.AppendEvent(providerType, event)
	case model.IDPConfigChanged, iam_es_model.IDPConfigChanged,
		model.OIDCIDPConfigAdded, iam_es_model.OIDCIDPConfigAdded,
		model.OIDCIDPConfigChanged, iam_es_model.OIDCIDPConfigChanged:
		err = idp.SetData(event)
		if err != nil {
			return err
		}
		idp, err = m.view.IDPConfigByID(idp.IDPConfigID)
		if err != nil {
			return err
		}
		err = idp.AppendEvent(providerType, event)
	case model.IDPConfigRemoved, iam_es_model.IDPConfigRemoved:
		err = idp.SetData(event)
		if err != nil {
			return err
		}
		return m.view.DeleteIDPConfig(idp.IDPConfigID, event.Sequence, event.CreationDate)
	default:
		return m.view.ProcessedIDPConfigSequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return m.view.PutIDPConfig(idp, idp.Sequence, event.CreationDate)
}

func (i *IDPConfig) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-Nxu8s", "id", event.AggregateID).WithError(err).Warn("something went wrong in idp config handler")
	return spooler.HandleError(event, err, i.view.GetLatestIDPConfigFailedEvent, i.view.ProcessedIDPConfigFailedEvent, i.view.ProcessedIDPConfigSequence, i.errorCountUntilSkip)
}

func (i *IDPConfig) OnSuccess() error {
	return spooler.HandleSuccess(i.view.UpdateIDPConfigSpoolerRunTimestamp)
}
