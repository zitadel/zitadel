package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	idpProviderTable = "adminapi.idp_providers"
)

func (v *View) IDPProviderByAggregateAndIdpConfigID(aggregateID, idpConfigID string) (*model.IDPProviderView, error) {
	return view.GetIDPProviderByAggregateIDAndConfigID(v.Db, idpProviderTable, aggregateID, idpConfigID)
}

func (v *View) IDPProvidersByIDPConfigID(idpConfigID string) ([]*model.IDPProviderView, error) {
	return view.IDPProvidersByIdpConfigID(v.Db, idpProviderTable, idpConfigID)
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

func (v *View) DeleteIDPProvider(aggregateID, idpConfigID string, event *models.Event) error {
	err := view.DeleteIDPProvider(v.Db, idpProviderTable, aggregateID, idpConfigID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedIDPProviderSequence(event)
}

func (v *View) GetLatestIDPProviderSequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(idpProviderTable)
}

func (v *View) ProcessedIDPProviderSequence(event *models.Event) error {
	return v.saveCurrentSequence(idpProviderTable, event)
}

func (v *View) UpdateIDPProviderSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(idpProviderTable)
}

func (v *View) GetLatestIDPProviderFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(idpProviderTable, sequence)
}

func (v *View) ProcessedIDPProviderFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
