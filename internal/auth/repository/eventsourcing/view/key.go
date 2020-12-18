package view

import (
	"github.com/caos/zitadel/internal/eventstore/models"
	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/key/repository/view"
	"github.com/caos/zitadel/internal/key/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	keyTable = "auth.keys"
)

func (v *View) KeyByIDAndType(keyID string, private bool) (*model.KeyView, error) {
	return view.KeyByIDAndType(v.Db, keyTable, keyID, private)
}

func (v *View) GetSigningKey() (*key_model.SigningKey, error) {
	key, err := view.GetSigningKey(v.Db, keyTable)
	if err != nil {
		return nil, err
	}
	return key_model.SigningKeyFromKeyView(model.KeyViewToModel(key), v.keyAlgorithm)
}

func (v *View) GetActiveKeySet() ([]*key_model.PublicKey, error) {
	keys, err := view.GetActivePublicKeys(v.Db, keyTable)
	if err != nil {
		return nil, err
	}
	return key_model.PublicKeysFromKeyView(model.KeyViewsToModel(keys), v.keyAlgorithm)
}

func (v *View) PutKeys(privateKey, publicKey *model.KeyView, event *models.Event) error {
	err := view.PutKeys(v.Db, keyTable, privateKey, publicKey)
	if err != nil {
		return err
	}
	return v.ProcessedKeySequence(event)
}

func (v *View) DeleteKey(keyID string, private bool, event *models.Event) error {
	err := view.DeleteKey(v.Db, keyTable, keyID, private)
	if err != nil {
		return nil
	}
	return v.ProcessedKeySequence(event)
}

func (v *View) DeleteKeyPair(keyID string, event *models.Event) error {
	err := view.DeleteKeyPair(v.Db, keyTable, keyID)
	if err != nil {
		return nil
	}
	return v.ProcessedKeySequence(event)
}

func (v *View) GetLatestKeySequence(aggregateType string) (*repository.CurrentSequence, error) {
	return v.latestSequence(keyTable, aggregateType)
}

func (v *View) ProcessedKeySequence(event *models.Event) error {
	return v.saveCurrentSequence(keyTable, event)
}

func (v *View) UpdateKeySpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(keyTable)
}

func (v *View) GetLatestKeyFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(keyTable, sequence)
}

func (v *View) ProcessedKeyFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
