package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/key/repository/view"
	"github.com/caos/zitadel/internal/key/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	authNKeyTable = "management.authn_keys"
)

func (v *View) AuthNKeyByIDs(objectID, keyID string) (*model.AuthNKeyView, error) {
	return view.AuthNKeyByIDs(v.Db, authNKeyTable, objectID, keyID)
}

func (v *View) AuthNKeysByObjectID(objectID string) ([]*model.AuthNKeyView, error) {
	return view.AuthNKeysByObjectID(v.Db, authNKeyTable, objectID)
}

func (v *View) AuthNKeyByID(keyID string) (*model.AuthNKeyView, error) {
	return view.AuthNKeyByID(v.Db, authNKeyTable, keyID)
}

func (v *View) SearchAuthNKeys(request *key_model.AuthNKeySearchRequest) ([]*model.AuthNKeyView, uint64, error) {
	return view.SearchAuthNKeys(v.Db, authNKeyTable, request)
}

func (v *View) PutAuthNKey(key *model.AuthNKeyView, event *models.Event) error {
	err := view.PutAuthNKey(v.Db, authNKeyTable, key)
	if err != nil {
		return err
	}
	return v.ProcessedAuthNKeySequence(event)
}

func (v *View) DeleteAuthNKey(keyID string, event *models.Event) error {
	err := view.DeleteAuthNKey(v.Db, authNKeyTable, keyID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedAuthNKeySequence(event)
}

func (v *View) DeleteAuthNKeysByObjectID(userID string, event *models.Event) error {
	err := view.DeleteAuthNKey(v.Db, authNKeyTable, userID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedAuthNKeySequence(event)
}

func (v *View) GetLatestAuthNKeySequence() (*repository.CurrentSequence, error) {
	return v.latestSequence(authNKeyTable)
}

func (v *View) ProcessedAuthNKeySequence(event *models.Event) error {
	return v.saveCurrentSequence(authNKeyTable, event)
}

func (v *View) UpdateAuthNKeySpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(authNKeyTable)
}

func (v *View) GetLatestAuthNKeyFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(authNKeyTable, sequence)
}

func (v *View) ProcessedAuthNKeyFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
