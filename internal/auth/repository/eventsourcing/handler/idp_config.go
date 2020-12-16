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
	idpConfigTable = "auth.idp_configs"
)

func (i *IDPConfig) ViewModel() string {
	return idpConfigTable
}

func (_ *IDPConfig) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (i *IDPConfig) CurrentSequence() (uint64, error) {
	sequence, err := i.view.GetLatestIDPConfigSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (i *IDPConfig) EventQuery() (*models.SearchQuery, error) {
	sequence, err := i.view.GetLatestIDPConfigSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(i.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (i *IDPConfig) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate:
		err = i.processIdpConfig(iam_model.IDPProviderTypeOrg, event)
	case iam_es_model.IAMAggregate:
		err = i.processIdpConfig(iam_model.IDPProviderTypeSystem, event)
	}
	return err
}

func (i *IDPConfig) processIdpConfig(providerType iam_model.IDPProviderType, event *models.Event) (err error) {
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
		idp, err = i.view.IDPConfigByID(idp.IDPConfigID)
		if err != nil {
			return err
		}
		err = idp.AppendEvent(providerType, event)
	case model.IDPConfigRemoved, iam_es_model.IDPConfigRemoved:
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
	logging.LogWithFields("SPOOL-Ejf8s", "id", event.AggregateID).WithError(err).Warn("something went wrong in idp config handler")
	return spooler.HandleError(event, err, i.view.GetLatestIDPConfigFailedEvent, i.view.ProcessedIDPConfigFailedEvent, i.view.ProcessedIDPConfigSequence, i.errorCountUntilSkip)
}

func (i *IDPConfig) OnSuccess() error {
	return spooler.HandleSuccess(i.view.UpdateIDPConfigSpoolerRunTimestamp)
}
