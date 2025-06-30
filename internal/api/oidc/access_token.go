package oidc

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/user/model"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type accessToken struct {
	tokenID           string
	userID            string
	resourceOwner     string
	subject           string
	preferredLanguage *language.Tag
	clientID          string
	audience          []string
	scope             []string
	authMethods       []domain.UserAuthMethodType
	authTime          time.Time
	tokenCreation     time.Time
	tokenExpiration   time.Time
	isPAT             bool
	actor             *domain.TokenActor
}

var ErrInvalidTokenFormat = errors.New("invalid token format")

func (s *Server) verifyAccessToken(ctx context.Context, tkn string) (_ *accessToken, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	var tokenID, subject string

	if tokenIDSubject, err := s.Provider().Crypto().Decrypt(tkn); err == nil {
		split := strings.Split(tokenIDSubject, ":")
		if len(split) != 2 {
			return nil, zerrors.ThrowPermissionDenied(ErrInvalidTokenFormat, "OIDC-rei1O", "token is not valid or has expired")
		}
		tokenID, subject = split[0], split[1]
	} else {
		verifier := op.NewAccessTokenVerifier(op.IssuerFromContext(ctx), s.accessTokenKeySet,
			op.WithSupportedAccessTokenSigningAlgorithms(supportedSigningAlgs()...),
		)
		claims, err := op.VerifyAccessToken[*oidc.AccessTokenClaims](ctx, tkn, verifier)
		if err != nil {
			return nil, zerrors.ThrowPermissionDenied(err, "OIDC-Eib8e", "token is not valid or has expired")
		}
		tokenID, subject = claims.JWTID, claims.Subject
	}

	if strings.HasPrefix(tokenID, command.IDPrefixV2) {
		token, err := s.query.ActiveAccessTokenByToken(ctx, tokenID)
		if err != nil {
			return nil, err
		}
		return accessTokenV2(tokenID, subject, token), nil
	}

	token, err := s.repo.TokenByIDs(ctx, subject, tokenID)
	if err != nil {
		return nil, zerrors.ThrowPermissionDenied(err, "OIDC-Dsfb2", "token is not valid or has expired")
	}
	return accessTokenV1(tokenID, subject, token), nil
}

func accessTokenV1(tokenID, subject string, token *model.TokenView) *accessToken {
	var preferredLanguage *language.Tag
	if token.PreferredLanguage != "" {
		preferredLanguage = gu.Ptr(language.Make(token.PreferredLanguage))
	}
	return &accessToken{
		tokenID:           tokenID,
		userID:            token.UserID,
		resourceOwner:     token.ResourceOwner,
		subject:           subject,
		preferredLanguage: preferredLanguage,
		clientID:          token.ApplicationID,
		audience:          token.Audience,
		scope:             token.Scopes,
		tokenCreation:     token.CreationDate,
		tokenExpiration:   token.Expiration,
		isPAT:             token.IsPAT,
		actor:             token.Actor,
	}
}

func accessTokenV2(tokenID, subject string, token *query.OIDCSessionAccessTokenReadModel) *accessToken {
	return &accessToken{
		tokenID:           tokenID,
		userID:            token.UserID,
		resourceOwner:     token.ResourceOwner,
		subject:           subject,
		preferredLanguage: token.PreferredLanguage,
		clientID:          token.ClientID,
		audience:          token.Audience,
		scope:             token.Scope,
		authMethods:       token.AuthMethods,
		authTime:          token.AuthTime,
		tokenCreation:     token.AccessTokenCreation,
		tokenExpiration:   token.AccessTokenExpiration,
		actor:             token.Actor,
	}
}

func (s *Server) assertClientScopesForPAT(ctx context.Context, token *accessToken, clientID, projectID string) error {
	token.audience = append(token.audience, clientID, projectID)
	projectIDQuery, err := query.NewProjectRoleProjectIDSearchQuery(projectID)
	if err != nil {
		return zerrors.ThrowInternal(err, "OIDC-Cyc78", "Errors.Internal")
	}
	roles, err := s.query.SearchProjectRoles(ctx, false, &query.ProjectRoleSearchQueries{Queries: []query.SearchQuery{projectIDQuery}}, nil)
	if err != nil {
		return err
	}
	for _, role := range roles.ProjectRoles {
		token.scope = append(token.scope, ScopeProjectRolePrefix+role.Key)
	}
	return nil
}
