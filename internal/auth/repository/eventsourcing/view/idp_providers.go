package view

import (
	"github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
	"time"
)

const (
	idpProviderTable = "auth.idp_providers"
)

func (v *View) IDPProviderByAggregateAndIDPConfigID(aggregateID, idpConfigID string) (*model.IDPProviderView, error) {
	return view.GetIDPProviderByAggregateIDAndConfigID(v.Db, idpProviderTable, aggregateID, idpConfigID)
}

func (v *View) IDPProvidersByIDPConfigID(idpConfigID string) ([]*model.IDPProviderView, error) {
	return view.IDPProvidersByIdpConfigID(v.Db, idpProviderTable, idpConfigID)
}

func (v *View) IDPProvidersByAggregateIDAndState(aggregateID string, idpConfigState iam_model.IDPConfigState) ([]*model.IDPProviderView, error) {
	return view.IDPProvidersByAggregateIDAndState(v.Db, idpProviderTable, aggregateID, idpConfigState)
}

func (v *View) SearchIDPProviders(request *iam_model.IDPProviderSearchRequest) ([]*model.IDPProviderView, uint64, error) {
	return view.SearchIDPProviders(v.Db, idpProviderTable, request)
}

func (v *View) PutIDPProvider(provider *model.IDPProviderView, sequence uint64, eventTimestamp time.Time) error {
	err := view.PutIDPProvider(v.Db, idpProviderTable, provider)
	if err != nil {
		return err
	}
	return v.ProcessedIDPProviderSequence(sequence, eventTimestamp)
}

func (v *View) PutIDPProviders(sequence uint64, eventTimestamp time.Time, providers ...*model.IDPProviderView) error {
	err := view.PutIDPProviders(v.Db, idpProviderTable, providers...)
	if err != nil {
		return err
	}
	return v.ProcessedIDPProviderSequence(sequence, eventTimestamp)
}

func (v *View) DeleteIDPProvider(aggregateID, idpConfigID string, eventSequence uint64, eventTimestamp time.Time) error {
	err := view.DeleteIDPProvider(v.Db, idpProviderTable, aggregateID, idpConfigID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedIDPProviderSequence(eventSequence, eventTimestamp)
}

func (v *View) DeleteIDPProvidersByAggregateID(aggregateID string, eventSequence uint64, eventTimestamp time.Time) error {
	err := view.DeleteIDPProvidersByAggregateID(v.Db, idpProviderTable, aggregateID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedIDPProviderSequence(eventSequence, eventTimestamp)
}

func (v *View) GetLatestIDPProviderSequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(idpProviderTable)
}

func (v *View) ProcessedIDPProviderSequence(eventSequence uint64, eventTimestamp time.Time) error {
	return v.saveCurrentSequence(idpProviderTable, eventSequence, eventTimestamp)
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
