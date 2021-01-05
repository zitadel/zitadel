package domain

import "github.com/caos/zitadel/internal/eventstore/models"

type Machine struct {
	models.ObjectRoot

	Name        string
	Description string
}

func (sa *Machine) IsValid() bool {
	return sa.Name != ""
}
