package model

import (
	"time"

	"github.com/caos/zitadel/internal/eventstore/models"
	key_model "github.com/caos/zitadel/internal/key/model"
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
	Type           key_model.AuthNKeyType
	ExpirationDate time.Time
	PrivateKey     []byte
}
