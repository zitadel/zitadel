package view

import (
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
	"github.com/zitadel/zitadel/internal/iam/repository/view"
	"github.com/zitadel/zitadel/internal/iam/repository/view/model"
	global_view "github.com/zitadel/zitadel/internal/view/repository"
)

const (
	idpProviderTable = "auth.idp_providers"
)

func (v *View) IDPProviderByAggregateAndIDPConfigID(aggregateID, idpConfigID, instanceID string) (*model.IDPProviderView, error) {
	return view.GetIDPProviderByAggregateIDAndConfigID(v.Db, idpProviderTable, aggregateID, idpConfigID, instanceID)
}

func (v *View) IDPProvidersByIDPConfigID(idpConfigID, instanceID string) ([]*model.IDPProviderView, error) {
	return view.IDPProvidersByIdpConfigID(v.Db, idpProviderTable, idpConfigID, instanceID)
}

func (v *View) IDPProvidersByAggregateIDAndState(aggregateID, instanceID string, idpConfigState iam_model.IDPConfigState) ([]*model.IDPProviderView, error) {
	return view.IDPProvidersByAggregateIDAndState(v.Db, idpProviderTable, aggregateID, instanceID, idpConfigState)
}

func (v *View) SearchIDPProviders(request *iam_model.IDPProviderSearchRequest) ([]*model.IDPProviderView, uint64, error) {
	return view.SearchIDPProviders(v.Db, idpProviderTable, request)
}

func (v *View) PutIDPProvider(provider *model.IDPProviderView, event *models.Event) error {
	err := view.PutIDPProvider(v.Db, idpProviderTable, provider)
	if err != nil {
		return err
	}
	return v.ProcessedIDPProviderSequence(event)
}

func (v *View) PutIDPProviders(event *models.Event, providers ...*model.IDPProviderView) error {
	err := view.PutIDPProviders(v.Db, idpProviderTable, providers...)
	if err != nil {
		return err
	}
	return v.ProcessedIDPProviderSequence(event)
}

func (v *View) DeleteIDPProvider(aggregateID, idpConfigID, instanceID string, event *models.Event) error {
	err := view.DeleteIDPProvider(v.Db, idpProviderTable, aggregateID, idpConfigID, instanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedIDPProviderSequence(event)
}

func (v *View) DeleteIDPProvidersByAggregateID(aggregateID, instanceID string, event *models.Event) error {
	err := view.DeleteIDPProvidersByAggregateID(v.Db, idpProviderTable, aggregateID, instanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedIDPProviderSequence(event)
}

func (v *View) GetLatestIDPProviderSequence(instanceID string) (*global_view.CurrentSequence, error) {
	return v.latestSequence(idpProviderTable, instanceID)
}

func (v *View) GetLatestIDPProviderSequences() ([]*global_view.CurrentSequence, error) {
	return v.latestSequences(idpProviderTable)
}

func (v *View) ProcessedIDPProviderSequence(event *models.Event) error {
	return v.saveCurrentSequence(idpProviderTable, event)
}

func (v *View) UpdateIDPProviderSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(idpProviderTable)
}

func (v *View) GetLatestIDPProviderFailedEvent(sequence uint64, instanceID string) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(idpProviderTable, instanceID, sequence)
}

func (v *View) ProcessedIDPProviderFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
