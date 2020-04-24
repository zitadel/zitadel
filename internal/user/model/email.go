package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type Email struct {
	es_models.ObjectRoot

	EmailAddress    string
	IsEmailVerified bool
}

func (e *Email) IsValid() bool {
	return e.EmailAddress != ""
}
