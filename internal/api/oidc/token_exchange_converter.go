package oidc

import (
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

type exchangeToken struct {
	tokenType         oidc.TokenType
	userID            string
	issuer            string
	resourceOwner     string
	authTime          time.Time
	authMethods       []domain.UserAuthMethodType
	actor             *domain.TokenActor
	audience          []string
	scopes            []string
	preferredLanguage *language.Tag
}

func (et *exchangeToken) nestedActor() *domain.TokenActor {
	return &domain.TokenActor{
		Actor:  et.actor,
		UserID: et.userID,
		Issuer: et.issuer,
	}
}

func accessToExchangeToken(token *accessToken, issuer string) *exchangeToken {
	return &exchangeToken{
		tokenType:         oidc.AccessTokenType,
		userID:            token.userID,
		issuer:            issuer,
		resourceOwner:     token.resourceOwner,
		authMethods:       token.authMethods,
		actor:             token.actor,
		audience:          token.audience,
		scopes:            token.scope,
		preferredLanguage: token.preferredLanguage,
	}
}

func idTokenClaimsToExchangeToken(claims *oidc.IDTokenClaims, resourceOwner string) *exchangeToken {
	var preferredLanguage *language.Tag
	if tag := claims.Locale.Tag(); !tag.IsRoot() {
		preferredLanguage = &tag
	}
	return &exchangeToken{
		tokenType:         oidc.IDTokenType,
		userID:            claims.Subject,
		issuer:            claims.Issuer,
		resourceOwner:     resourceOwner,
		authTime:          claims.GetAuthTime(),
		authMethods:       AMRToAuthMethodTypes(claims.AuthenticationMethodsReferences),
		actor:             actorClaimsToDomain(claims.Actor),
		audience:          claims.Audience,
		preferredLanguage: preferredLanguage,
	}
}

func actorClaimsToDomain(actor *oidc.ActorClaims) *domain.TokenActor {
	if actor == nil {
		return nil
	}
	return &domain.TokenActor{
		Actor:  actorClaimsToDomain(actor.Actor),
		UserID: actor.Subject,
		Issuer: actor.Issuer,
	}
}

func actorDomainToClaims(actor *domain.TokenActor) *oidc.ActorClaims {
	if actor == nil {
		return nil
	}
	return &oidc.ActorClaims{
		Actor:   actorDomainToClaims(actor.Actor),
		Subject: actor.UserID,
		Issuer:  actor.Issuer,
	}
}

func jwtToExchangeToken(jwt *oidc.JWTTokenRequest, resourceOwner string, preferredLanguage *language.Tag) *exchangeToken {
	return &exchangeToken{
		tokenType:     oidc.JWTTokenType,
		userID:        jwt.Subject,
		issuer:        jwt.Issuer,
		resourceOwner: resourceOwner,
		scopes:        jwt.Scopes,
		authTime:      jwt.IssuedAt.AsTime(),
		// audience omitted as we don't thrust audiences not signed by us
		preferredLanguage: preferredLanguage,
	}
}

func userToExchangeToken(user *query.User) *exchangeToken {
	return &exchangeToken{
		tokenType:     UserIDTokenType,
		userID:        user.ID,
		resourceOwner: user.ResourceOwner,
	}
}
