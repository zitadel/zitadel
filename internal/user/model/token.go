package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"time"
)

type Token struct {
	es_models.ObjectRoot

	TokenID           string
	ApplicationID     string
	UserAgentID       string
	Audience          []string
	Expiration        time.Time
	Scopes            []string
	PreferredLanguage string
}
