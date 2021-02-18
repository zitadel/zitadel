package domain

import (
	auth_req_model "github.com/caos/zitadel/internal/auth_request/model"
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
		if strings.HasPrefix(scope, auth_req_model.ProjectIDScope) && strings.HasSuffix(scope, auth_req_model.AudSuffix) {
			audience = append(audience, strings.TrimSuffix(strings.TrimPrefix(scope, auth_req_model.ProjectIDScope), auth_req_model.AudSuffix))
		}
	}
	return audience
}
