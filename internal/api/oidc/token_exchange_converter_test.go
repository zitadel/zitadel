package oidc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
)

func Test_accessToExchangeToken(t *testing.T) {
	authTime := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	preferredLanguage := language.German
	actor := &domain.TokenActor{UserID: "actor", Issuer: "https://issuer.example"}

	got := accessToExchangeToken(&accessToken{
		userID:            "user",
		resourceOwner:     "org",
		audience:          []string{"client", "project"},
		scope:             []string{oidc.ScopeOpenID, oidc.ScopeProfile},
		authMethods:       []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
		authTime:          authTime,
		actor:             actor,
		preferredLanguage: &preferredLanguage,
	}, "https://issuer.example")

	assert.Equal(t, &exchangeToken{
		tokenType:         oidc.AccessTokenType,
		userID:            "user",
		issuer:            "https://issuer.example",
		resourceOwner:     "org",
		authTime:          authTime,
		authMethods:       []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
		actor:             actor,
		audience:          []string{"client", "project"},
		scopes:            []string{oidc.ScopeOpenID, oidc.ScopeProfile},
		preferredLanguage: &preferredLanguage,
	}, got)
}
