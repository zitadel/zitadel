package handler

import (
	"github.com/caos/logging"
	iam_model "github.com/caos/zitadel/internal/iam/model"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
)

type IdpConfig struct {
	handler
}

const (
	idpConfigTable = "adminapi.idp_configs"
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
		AggregateTypeFilter(model.IamAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *IdpConfig) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.IamAggregate:
		err = m.processIdpConfig(event)
	}
	return err
}

func (m *IdpConfig) processIdpConfig(event *models.Event) (err error) {
	idp := new(iam_view_model.IdpConfigView)
	switch event.Type {
	case model.IdpConfigAdded:
		err = idp.AppendEvent(iam_model.IdpProviderTypeSystem, event)
	case model.IdpConfigChanged,
		model.OidcIdpConfigAdded,
		model.OidcIdpConfigChanged:
		err = idp.SetData(event)
		if err != nil {
			return err
		}
		idp, err = m.view.IdpConfigByID(idp.IdpConfigID)
		if err != nil {
			return err
		}
		err = idp.AppendEvent(iam_model.IdpProviderTypeSystem, event)
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
	logging.LogWithFields("SPOOL-Mslo9", "id", event.AggregateID).WithError(err).Warn("something went wrong in idp config handler")
	return spooler.HandleError(event, err, m.view.GetLatestIdpConfigFailedEvent, m.view.ProcessedIdpConfigFailedEvent, m.view.ProcessedIdpConfigSequence, m.errorCountUntilSkip)
}
