package handler

import (
	"github.com/caos/logging"
	org_model "github.com/caos/zitadel/internal/org/repository/view/model"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

type IdpConfig struct {
	handler
}

const (
	idpConfigTable = "management.idp_configs"
)

func (m *IdpConfig) ViewModel() string {
	return idpConfigTable
}

func (m *IdpConfig) EventQuery() (*models.SearchQuery, error) {
	sequence, err := m.view.GetLatestIdpConfigSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.OrgAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *IdpConfig) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate:
		err = m.processIdpConfig(event)
	}
	return err
}

func (m *IdpConfig) processIdpConfig(event *models.Event) (err error) {
	idp := new(org_model.IdpConfigView)
	switch event.Type {
	case model.IdpConfigAdded:
		err = idp.AppendEvent(event)
	case model.IdpConfigChanged,
		model.OidcIdpConfigChanged:
		err = idp.SetData(event)
		if err != nil {
			return err
		}
		idp, err = m.view.IdpConfigByID(idp.IdpConfigID)
		if err != nil {
			return err
		}
		err = idp.AppendEvent(event)
	case model.IdpConfigRemoved:
		err = idp.SetData(event)
		if err != nil {
			return err
		}
		return m.view.DeleteIdpConfig(idp.IdpConfigID, event.Sequence)
	default:
		return m.view.ProcessedIdpConfigSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return m.view.PutIdpConfig(idp, idp.Sequence)
}

func (m *IdpConfig) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-Nxu8s", "id", event.AggregateID).WithError(err).Warn("something went wrong in idp config handler")
	return spooler.HandleError(event, err, m.view.GetLatestIdpConfigFailedEvent, m.view.ProcessedIdpConfigFailedEvent, m.view.ProcessedIdpConfigSequence, m.errorCountUntilSkip)
}
