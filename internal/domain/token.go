package domain

import (
	"strings"
	"time"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

type Token struct {
	es_models.ObjectRoot

	TokenID           string
	ApplicationID     string
	UserAgentID       string
	RefreshTokenID    string
	Audience          []string
	Expiration        time.Time
	Scopes            []string
	PreferredLanguage string
}

func AddAudScopeToAudience(audience, scopes []string) []string {
	for _, scope := range scopes {
		if strings.HasPrefix(scope, ProjectIDScope) && strings.HasSuffix(scope, AudSuffix) {
			audience = append(audience, strings.TrimSuffix(strings.TrimPrefix(scope, ProjectIDScope), AudSuffix))
		}
	}
	return audience
}
