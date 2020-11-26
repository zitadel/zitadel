package view

import (
	"github.com/caos/zitadel/internal/errors"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
	"time"
)

const (
	externalIDPTable = "auth.user_external_idps"
)

func (v *View) ExternalIDPByExternalUserIDAndIDPConfigID(externalUserID, idpConfigID string) (*model.ExternalIDPView, error) {
	return view.ExternalIDPByExternalUserIDAndIDPConfigID(v.Db, externalIDPTable, externalUserID, idpConfigID)
}

func (v *View) ExternalIDPByExternalUserIDAndIDPConfigIDAndResourceOwner(externalUserID, idpConfigID, resourceOwner string) (*model.ExternalIDPView, error) {
	return view.ExternalIDPByExternalUserIDAndIDPConfigIDAndResourceOwner(v.Db, externalIDPTable, externalUserID, idpConfigID, resourceOwner)
}

func (v *View) ExternalIDPsByIDPConfigID(idpConfigID string) ([]*model.ExternalIDPView, error) {
	return view.ExternalIDPsByIDPConfigID(v.Db, externalIDPTable, idpConfigID)
}

func (v *View) ExternalIDPsByUserID(userID string) ([]*model.ExternalIDPView, error) {
	return view.ExternalIDPsByUserID(v.Db, externalIDPTable, userID)
}

func (v *View) SearchExternalIDPs(request *usr_model.ExternalIDPSearchRequest) ([]*model.ExternalIDPView, uint64, error) {
	return view.SearchExternalIDPs(v.Db, externalIDPTable, request)
}

func (v *View) PutExternalIDP(externalIDP *model.ExternalIDPView, sequence uint64, eventTimestamp time.Time) error {
	err := view.PutExternalIDP(v.Db, externalIDPTable, externalIDP)
	if err != nil {
		return err
	}
	return v.ProcessedExternalIDPSequence(sequence, eventTimestamp)
}

func (v *View) PutExternalIDPs(sequence uint64, eventTimestamp time.Time, externalIDPs ...*model.ExternalIDPView) error {
	err := view.PutExternalIDPs(v.Db, externalIDPTable, externalIDPs...)
	if err != nil {
		return err
	}
	return v.ProcessedExternalIDPSequence(sequence, eventTimestamp)
}

func (v *View) DeleteExternalIDP(externalUserID, idpConfigID string, eventSequence uint64, eventTimestamp time.Time) error {
	err := view.DeleteExternalIDP(v.Db, externalIDPTable, externalUserID, idpConfigID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedExternalIDPSequence(eventSequence, eventTimestamp)
}

func (v *View) DeleteExternalIDPsByUserID(userID string, eventSequence uint64, eventTimestamp time.Time) error {
	err := view.DeleteExternalIDPsByUserID(v.Db, externalIDPTable, userID)
	if err != nil {
		return err
	}
	return v.ProcessedExternalIDPSequence(eventSequence, eventTimestamp)
}

func (v *View) GetLatestExternalIDPSequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(externalIDPTable)
}

func (v *View) ProcessedExternalIDPSequence(eventSequence uint64, eventTimestamp time.Time) error {
	return v.saveCurrentSequence(externalIDPTable, eventSequence, eventTimestamp)
}

func (v *View) UpdateExternalIDPSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(externalIDPTable)
}

func (v *View) GetLatestExternalIDPFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(externalIDPTable, sequence)
}

func (v *View) ProcessedExternalIDPFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
