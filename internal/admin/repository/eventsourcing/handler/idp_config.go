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

type IDPConfig struct {
	handler
}

const (
	idpConfigTable = "adminapi.idp_configs"
)

func (i *IDPConfig) ViewModel() string {
	return idpConfigTable
}

func (i *IDPConfig) EventQuery() (*models.SearchQuery, error) {
	sequence, err := i.view.GetLatestIDPConfigSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.IAMAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (i *IDPConfig) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.IAMAggregate:
		err = i.processIDPConfig(event)
	}
	return err
}

func (i *IDPConfig) processIDPConfig(event *models.Event) (err error) {
	idp := new(iam_view_model.IDPConfigView)
	switch event.Type {
	case model.IDPConfigAdded:
		err = idp.AppendEvent(iam_model.IDPProviderTypeSystem, event)
	case model.IDPConfigChanged,
		model.OIDCIDPConfigAdded,
		model.OIDCIDPConfigChanged:
		err = idp.SetData(event)
		if err != nil {
			return err
		}
		idp, err = i.view.IDPConfigByID(idp.IDPConfigID)
		if err != nil {
			return err
		}
		err = idp.AppendEvent(iam_model.IDPProviderTypeSystem, event)
	case model.IDPConfigRemoved:
		err = idp.SetData(event)
		if err != nil {
			return err
		}
		return i.view.DeleteIDPConfig(idp.IDPConfigID, event.Sequence, event.CreationDate)
	default:
		return i.view.ProcessedIDPConfigSequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return i.view.PutIDPConfig(idp, idp.Sequence, event.CreationDate)
}

func (i *IDPConfig) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-Mslo9", "id", event.AggregateID).WithError(err).Warn("something went wrong in idp config handler")
	return spooler.HandleError(event, err, i.view.GetLatestIDPConfigFailedEvent, i.view.ProcessedIDPConfigFailedEvent, i.view.ProcessedIDPConfigSequence, i.errorCountUntilSkip)
}

func (i *IDPConfig) OnSuccess() error {
	return spooler.HandleSuccess(i.view.UpdateIDPConfigSpoolerRunTimestamp)
}
