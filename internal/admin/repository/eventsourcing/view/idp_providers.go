package view

import (
	"github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	idpProviderTable = "adminapi.idp_providers"
)

func (v *View) IdpProviderByAggregateAndIdpConfigID(aggregateID, idpConfigID string) (*model.IdpProviderView, error) {
	return view.GetIdpProviderByAggregateIDAndConfigID(v.Db, idpProviderTable, aggregateID, idpConfigID)
}

func (v *View) IdpProvidersByIdpConfigID(aggregateID, idpConfigID string) ([]*model.IdpProviderView, error) {
	return view.IdpProvidersByIdpConfigID(v.Db, idpProviderTable, idpConfigID)
}

func (v *View) SearchIdpProviders(request *iam_model.IdpProviderSearchRequest) ([]*model.IdpProviderView, uint64, error) {
	return view.SearchIdpProviders(v.Db, idpProviderTable, request)
}

func (v *View) PutIdpProvider(provider *model.IdpProviderView, sequence uint64) error {
	err := view.PutIdpProvider(v.Db, idpProviderTable, provider)
	if err != nil {
		return err
	}
	return v.ProcessedIdpProviderSequence(sequence)
}

func (v *View) PutIdpProviders(sequence uint64, providers ...*model.IdpProviderView) error {
	err := view.PutIdpProviders(v.Db, idpProviderTable, providers...)
	if err != nil {
		return err
	}
	return v.ProcessedIdpProviderSequence(sequence)
}

func (v *View) DeleteIdpProvider(aggregateID, idpConfigID string, eventSequence uint64) error {
	err := view.DeleteIdpProvider(v.Db, idpProviderTable, aggregateID, idpConfigID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedIdpProviderSequence(eventSequence)
}

func (v *View) GetLatestIdpProviderSequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(idpProviderTable)
}

func (v *View) ProcessedIdpProviderSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(idpProviderTable, eventSequence)
}

func (v *View) GetLatestIdpProviderFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(idpProviderTable, sequence)
}

func (v *View) ProcessedIdpProviderFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
