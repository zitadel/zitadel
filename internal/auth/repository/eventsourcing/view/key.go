package view

import (
	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/key/repository/view"
	"github.com/caos/zitadel/internal/key/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view"
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

func (v *View) PutKeys(privateKey, publicKey *model.KeyView, eventSequence uint64) error {
	err := view.PutKeys(v.Db, keyTable, privateKey, publicKey)
	if err != nil {
		return err
	}
	return v.ProcessedKeySequence(eventSequence)
}

func (v *View) DeleteKey(keyID string, private bool, eventSequence uint64) error {
	err := view.DeleteKey(v.Db, keyTable, keyID, private)
	if err != nil {
		return nil
	}
	return v.ProcessedKeySequence(eventSequence)
}

func (v *View) DeleteKeyPair(keyID string, eventSequence uint64) error {
	err := view.DeleteKeyPair(v.Db, keyTable, keyID)
	if err != nil {
		return nil
	}
	return v.ProcessedKeySequence(eventSequence)
}

func (v *View) GetLatestKeySequence() (uint64, error) {
	return v.latestSequence(keyTable)
}

func (v *View) ProcessedKeySequence(eventSequence uint64) error {
	return v.saveCurrentSequence(keyTable, eventSequence)
}

func (v *View) GetLatestKeyFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(keyTable, sequence)
}

func (v *View) ProcessedKeyFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
