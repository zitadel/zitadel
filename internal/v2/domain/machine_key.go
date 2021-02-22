package domain

import (
	"time"

	"github.com/caos/zitadel/internal/eventstore/models"
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
