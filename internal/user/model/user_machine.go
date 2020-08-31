package model

import (
	"time"

	"github.com/caos/zitadel/internal/eventstore/models"
)

type Machine struct {
	models.ObjectRoot

	Name        string
	Description string
}

func (sa *Machine) IsValid() bool {
	return sa.Name != ""
}

type MachineKey struct {
	models.ObjectRoot

	KeyID          string
	Type           MachineKeyType
	ExpirationDate time.Time
	PrivateKey     []byte
}

type MachineKeyType int32

const (
	MachineKeyTypeNONE = iota
	MachineKeyTypeJSON
)
