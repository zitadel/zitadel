package view

import (
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
	"github.com/zitadel/zitadel/internal/iam/repository/view"
	iam_es_model "github.com/zitadel/zitadel/internal/iam/repository/view/model"
	global_view "github.com/zitadel/zitadel/internal/view/repository"
)

const (
	idpConfigTable = "auth.idp_configs"
)

func (v *View) IDPConfigByID(idpID, instanceID string) (*iam_es_model.IDPConfigView, error) {
	return view.IDPByID(v.Db, idpConfigTable, idpID, instanceID)
}

func (v *View) GetIDPConfigsByAggregateID(aggregateID, instanceID string) ([]*iam_es_model.IDPConfigView, error) {
	return view.GetIDPConfigsByAggregateID(v.Db, idpConfigTable, aggregateID, instanceID)
}

func (v *View) SearchIDPConfigs(request *iam_model.IDPConfigSearchRequest) ([]*iam_es_model.IDPConfigView, uint64, error) {
	return view.SearchIDPs(v.Db, idpConfigTable, request)
}

func (v *View) PutIDPConfig(idp *iam_es_model.IDPConfigView, event *models.Event) error {
	err := view.PutIDP(v.Db, idpConfigTable, idp)
	if err != nil {
		return err
	}
	return v.ProcessedIDPConfigSequence(event)
}

func (v *View) DeleteIDPConfig(idpID string, event *models.Event) error {
	err := view.DeleteIDP(v.Db, idpConfigTable, idpID, event.InstanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedIDPConfigSequence(event)
}

func (v *View) GetLatestIDPConfigSequence(instanceID string) (*global_view.CurrentSequence, error) {
	return v.latestSequence(idpConfigTable, instanceID)
}

func (v *View) GetLatestIDPConfigSequences() ([]*global_view.CurrentSequence, error) {
	return v.latestSequences(idpConfigTable)
}

func (v *View) ProcessedIDPConfigSequence(event *models.Event) error {
	return v.saveCurrentSequence(idpConfigTable, event)
}

func (v *View) UpdateIDPConfigSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(idpConfigTable)
}

func (v *View) GetLatestIDPConfigFailedEvent(sequence uint64, instanceID string) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(idpConfigTable, instanceID, sequence)
}

func (v *View) ProcessedIDPConfigFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
