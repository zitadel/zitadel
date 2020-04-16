package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type Token struct {
	es_models.ObjectRoot
}

func (t *Token) IsValid() bool {
	return true
}
