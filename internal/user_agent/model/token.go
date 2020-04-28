package model

import (
	"time"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type Token struct {
	es_models.ObjectRoot

	TokenID       string
	AgentID       string
	AuthSessionID string
	UserSessionID string
	Expiry        time.Duration
}

func (t *Token) IsValid() bool {
	return true
}
