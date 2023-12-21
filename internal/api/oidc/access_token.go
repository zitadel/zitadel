package oidc

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/user/model"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type accessToken struct {
	tokenID         string
	userID          string
	subject         string
	clientID        string
	audience        []string
	scope           []string
	tokenCreation   time.Time
	tokenExpiration time.Time
	isPAT           bool
}

func (s *Server) verifyAccessToken(ctx context.Context, tkn string) (*accessToken, error) {
	var tokenID, subject string

	if tokenIDSubject, err := s.Provider().Crypto().Decrypt(tkn); err == nil {
		split := strings.Split(tokenIDSubject, ":")
		if len(split) != 2 {
			return nil, errors.New("invalid token format")
		}
		tokenID, subject = split[0], split[1]
	} else {
		verifier := op.NewAccessTokenVerifier(op.IssuerFromContext(ctx), s.keySet)
		claims, err := op.VerifyAccessToken[*oidc.AccessTokenClaims](ctx, tkn, verifier)
		if err != nil {
			return nil, err
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
	return &accessToken{
		tokenID:         tokenID,
		userID:          token.UserID,
		subject:         subject,
		clientID:        token.ApplicationID,
		audience:        token.Audience,
		scope:           token.Scopes,
		tokenCreation:   token.CreationDate,
		tokenExpiration: token.Expiration,
		isPAT:           token.IsPAT,
	}
}

func accessTokenV2(tokenID, subject string, token *query.OIDCSessionAccessTokenReadModel) *accessToken {
	return &accessToken{
		tokenID:         tokenID,
		userID:          token.UserID,
		subject:         subject,
		clientID:        token.ClientID,
		audience:        token.Audience,
		scope:           token.Scope,
		tokenCreation:   token.AccessTokenCreation,
		tokenExpiration: token.AccessTokenExpiration,
	}
}

func (s *Server) assertClientScopesForPAT(ctx context.Context, token *accessToken, clientID, projectID string) error {
	token.audience = append(token.audience, clientID)
	projectIDQuery, err := query.NewProjectRoleProjectIDSearchQuery(projectID)
	if err != nil {
		return zerrors.ThrowInternal(err, "OIDC-Cyc78", "Errors.Internal")
	}
	roles, err := s.query.SearchProjectRoles(ctx, s.features.TriggerIntrospectionProjections, &query.ProjectRoleSearchQueries{Queries: []query.SearchQuery{projectIDQuery}})
	if err != nil {
		return err
	}
	for _, role := range roles.ProjectRoles {
		token.scope = append(token.scope, ScopeProjectRolePrefix+role.Key)
	}
	return nil
}
