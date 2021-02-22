package domain

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"strings"
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

func AddAudScopeToAudience(audience, scopes []string) []string {
	for _, scope := range scopes {
		if strings.HasPrefix(scope, ProjectIDScope) && strings.HasSuffix(scope, AudSuffix) {
			audience = append(audience, strings.TrimSuffix(strings.TrimPrefix(scope, ProjectIDScope), AudSuffix))
		}
	}
	return audience
}
