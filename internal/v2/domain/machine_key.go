package domain

import (
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/models"
	"time"
)

type MachineKey struct {
	models.ObjectRoot

	KeyID          string
	Type           MachineKeyType
	ExpirationDate time.Time
	PrivateKey     []byte
	PublicKey      []byte
}

type MachineKeyType int32

const (
	MachineKeyTypeNONE = iota
	MachineKeyTypeJSON

	keyCount
)

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

func (f MachineKeyType) Valid() bool {
	return f >= 0 && f < keyCount
}

func (key *MachineKey) GenerateNewMachineKeyPair(keySize int) error {
	privateKey, publicKey, err := crypto.GenerateKeyPair(keySize)
	if err != nil {
		return err
	}
	key.PublicKey, err = crypto.PublicKeyToBytes(publicKey)
	if err != nil {
		return err
	}
	key.PrivateKey = crypto.PrivateKeyToBytes(privateKey)
	return nil
}
