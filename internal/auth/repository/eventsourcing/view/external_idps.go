package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/user/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	externalIDPTable = "auth.user_external_idps"
)

func (v *View) ExternalIDPByExternalUserIDAndIDPConfigID(externalUserID, idpConfigID, instanceID string) (*model.ExternalIDPView, error) {
	return view.ExternalIDPByExternalUserIDAndIDPConfigID(v.Db, externalIDPTable, externalUserID, idpConfigID, instanceID)
}

func (v *View) ExternalIDPByExternalUserIDAndIDPConfigIDAndResourceOwner(externalUserID, idpConfigID, resourceOwner, instanceID string) (*model.ExternalIDPView, error) {
	return view.ExternalIDPByExternalUserIDAndIDPConfigIDAndResourceOwner(v.Db, externalIDPTable, externalUserID, idpConfigID, resourceOwner, instanceID)
}

func (v *View) ExternalIDPsByIDPConfigID(idpConfigID, instanceID string) ([]*model.ExternalIDPView, error) {
	return view.ExternalIDPsByIDPConfigID(v.Db, externalIDPTable, idpConfigID, instanceID)
}

func (v *View) PutExternalIDP(externalIDP *model.ExternalIDPView, event *models.Event) error {
	err := view.PutExternalIDP(v.Db, externalIDPTable, externalIDP)
	if err != nil {
		return err
	}
	return v.ProcessedExternalIDPSequence(event)
}

func (v *View) PutExternalIDPs(event *models.Event, externalIDPs ...*model.ExternalIDPView) error {
	err := view.PutExternalIDPs(v.Db, externalIDPTable, externalIDPs...)
	if err != nil {
		return err
	}
	return v.ProcessedExternalIDPSequence(event)
}

func (v *View) DeleteExternalIDP(externalUserID, idpConfigID, instanceID string, event *models.Event) error {
	err := view.DeleteExternalIDP(v.Db, externalIDPTable, externalUserID, idpConfigID, instanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedExternalIDPSequence(event)
}

func (v *View) DeleteExternalIDPsByUserID(userID, instanceID string, event *models.Event) error {
	err := view.DeleteExternalIDPsByUserID(v.Db, externalIDPTable, userID, instanceID)
	if err != nil {
		return err
	}
	return v.ProcessedExternalIDPSequence(event)
}

func (v *View) GetLatestExternalIDPSequence(instanceID string) (*global_view.CurrentSequence, error) {
	return v.latestSequence(externalIDPTable, instanceID)
}

func (v *View) GetLatestExternalIDPSequences() ([]*global_view.CurrentSequence, error) {
	return v.latestSequences(externalIDPTable)
}

func (v *View) ProcessedExternalIDPSequence(event *models.Event) error {
	return v.saveCurrentSequence(externalIDPTable, event)
}

func (v *View) UpdateExternalIDPSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(externalIDPTable)
}

func (v *View) GetLatestExternalIDPFailedEvent(sequence uint64, instanceID string) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(externalIDPTable, instanceID, sequence)
}

func (v *View) ProcessedExternalIDPFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
