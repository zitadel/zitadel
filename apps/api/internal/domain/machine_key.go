package domain

import (
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type MachineKey struct {
	models.ObjectRoot

	KeyID          string
	Type           AuthNKeyType
	ExpirationDate time.Time
	PrivateKey     []byte
	PublicKey      []byte
}

func (key *MachineKey) setPublicKey(publicKey []byte) {
	key.PublicKey = publicKey
}

func (key *MachineKey) setPrivateKey(privateKey []byte) {
	key.PrivateKey = privateKey
}

func (key *MachineKey) expirationDate() time.Time {
	return key.ExpirationDate
}

func (key *MachineKey) setExpirationDate(expiration time.Time) {
	key.ExpirationDate = expiration
}

func (key *MachineKey) Detail() ([]byte, error) {
	if key.Type == AuthNKeyTypeJSON {
		return key.MarshalJSON()
	}
	return nil, zerrors.ThrowPreconditionFailed(nil, "KEY-dsg52", "Errors.Internal")
}

func (key *MachineKey) MarshalJSON() ([]byte, error) {
	return MachineKeyMarshalJSON(key.KeyID, key.PrivateKey, key.ExpirationDate, key.AggregateID)
}

type MachineKeyState int32

const (
	MachineKeyStateUnspecified MachineKeyState = iota
	MachineKeyStateActive
	MachineKeyStateRemoved

	machineKeyStateCount
)

func (f MachineKeyState) Valid() bool {
	return f >= 0 && f < machineKeyStateCount
}

func MachineKeyMarshalJSON(keyID string, privateKey []byte, expirationDate time.Time, userID string) ([]byte, error) {
	return json.Marshal(struct {
		Type           string    `json:"type"`
		KeyID          string    `json:"keyId"`
		Key            string    `json:"key"`
		ExpirationDate time.Time `json:"expirationDate"`
		UserID         string    `json:"userId"`
	}{
		Type:           "serviceaccount",
		KeyID:          keyID,
		Key:            string(privateKey),
		ExpirationDate: expirationDate,
		UserID:         userID,
	})
}
