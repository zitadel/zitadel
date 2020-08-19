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
		AggregateTypeFilter(model.OrgAggregate, iam_es_model.IamAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *IdpConfig) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate:
		err = m.processIdpConfig(iam_model.IdpProviderTypeOrg, event)
	case iam_es_model.IamAggregate:
		err = m.processIdpConfig(iam_model.IdpProviderTypeOrg, event)
	}
	return err
}

func (m *IdpConfig) processIdpConfig(providerType iam_model.IdpProviderType, event *models.Event) (err error) {
	idp := new(iam_view_model.IdpConfigView)
	switch event.Type {
	case model.IdpConfigAdded,
		iam_es_model.IdpConfigAdded:
		err = idp.AppendEvent(providerType, event)
	case model.IdpConfigChanged, iam_es_model.IdpConfigChanged,
		model.OidcIdpConfigAdded, iam_es_model.OidcIdpConfigAdded,
		model.OidcIdpConfigChanged, iam_es_model.OidcIdpConfigChanged:
		err = idp.SetData(event)
		if err != nil {
			return err
		}
		idp, err = m.view.IdpConfigByID(idp.IdpConfigID)
		if err != nil {
			return err
		}
		err = idp.AppendEvent(providerType, event)
	case model.IdpConfigRemoved, iam_es_model.IdpConfigRemoved:
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
