package model

import (
	"time"

	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type RefreshToken struct {
	es_models.ObjectRoot

	TokenID           string
	ApplicationID     string
	UserAgentID       string
	Audience          []string
	Expiration        time.Time
	Scopes            []string
	PreferredLanguage string
}
