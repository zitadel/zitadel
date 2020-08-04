package view

import (
	"github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	idpConfigTable = "adminapi.idp_configs"
)

func (v *View) IdpConfigByID(idpID string) (*model.IdpConfigView, error) {
	return view.IdpByID(v.Db, idpConfigTable, idpID)
}

func (v *View) SearchIdpConfigs(request *iam_model.IdpConfigSearchRequest) ([]*model.IdpConfigView, int, error) {
	return view.SearchIdps(v.Db, idpConfigTable, request)
}

func (v *View) PutIdpConfig(idp *model.IdpConfigView, sequence uint64) error {
	err := view.PutIdp(v.Db, idpConfigTable, idp)
	if err != nil {
		return err
	}
	return v.ProcessedIdpConfigSequence(sequence)
}

func (v *View) DeleteIdpConfig(idpID string, eventSequence uint64) error {
	err := view.DeleteIdp(v.Db, idpConfigTable, idpID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedIdpConfigSequence(eventSequence)
}

func (v *View) GetLatestIdpConfigSequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(idpConfigTable)
}

func (v *View) ProcessedIdpConfigSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(idpConfigTable, eventSequence)
}

func (v *View) GetLatestIdpConfigFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(idpConfigTable, sequence)
}

func (v *View) ProcessedIdpConfigFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
