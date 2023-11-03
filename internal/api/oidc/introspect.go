package oidc

import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"github.com/zitadel/zitadel/internal/command"
	errz "github.com/zitadel/zitadel/internal/errors"
)

func (s *Server) Introspect(ctx context.Context, r *op.Request[op.IntrospectionRequest]) (_ *op.Response, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	clientChan := make(chan *instrospectionClientResult)
	go s.instrospectionClientAuth(ctx, r.Data.ClientCredentials, clientChan)

	tokenChan := make(chan *introspectionTokenResult)
	go s.introspectionToken(ctx, r.Data.Token, tokenChan)

	var (
		client *instrospectionClientResult
		token  *introspectionTokenResult
	)

	// make sure both channels are always read,
	// and cancel the context on first error
	for i := 0; i < 2; i++ {
		var resErr error

		select {
		case client = <-clientChan:
			resErr = client.err
		case token = <-tokenChan:
			resErr = token.err
		}

		if resErr == nil {
			continue
		}
		cancel()

		// we only care for the first error that occured,
		// as the next error is most probably a context error.
		if err == nil {
			err = resErr
		}
	}

	// only client auth errors should be returned
	var target *oidc.Error
	if errors.As(err, &target) && target.ErrorType == oidc.UnauthorizedClient {
		return nil, err
	}

	// all other errors should result in a response with active: false.
	response := new(oidc.IntrospectionResponse)
	if err != nil {
		// TODO: log error
		return op.NewResponse(response), nil
	}
	if err = validateIntrospectionAudience(token.audience, client.clientID, client.projectID); err != nil {
		// TODO: log error
		return op.NewResponse(response), nil
	}
	userInfo, err := s.storage.query.GetOIDCUserinfo(ctx, token.userID, token.scope, []string{client.projectID})
	if err != nil {
		// TODO: log error
		return op.NewResponse(response), nil
	}
	response.SetUserInfo(userinfoToOIDC(userInfo, token.scope))
	response.Scope = token.scope
	response.ClientID = token.clientID
	response.TokenType = oidc.BearerToken
	response.Expiration = oidc.FromTime(token.tokenExpiration)
	response.IssuedAt = oidc.FromTime(token.tokenCreation)
	response.NotBefore = oidc.FromTime(token.tokenCreation)
	response.Audience = token.audience
	response.Issuer = op.IssuerFromContext(ctx)
	response.JWTID = token.tokenID
	response.Active = true
	return op.NewResponse(response), nil
}

type instrospectionClientResult struct {
	clientID  string
	projectID string
	err       error
}

func (s *Server) instrospectionClientAuth(ctx context.Context, cc *op.ClientCredentials, rc chan<- *instrospectionClientResult) {
	clientID := cc.ClientID

	if cc.ClientAssertion != "" {
		verifier := op.NewJWTProfileVerifier(s.storage, op.IssuerFromContext(ctx), 1*time.Hour, time.Second)
		profile, err := op.VerifyJWTAssertion(ctx, cc.ClientAssertion, verifier)
		if err != nil {
			rc <- &instrospectionClientResult{
				err: oidc.ErrUnauthorizedClient().WithParent(err),
			}
			return
		}
		clientID = profile.Issuer
	} else {
		if err := s.storage.AuthorizeClientIDSecret(ctx, cc.ClientID, cc.ClientSecret); err != nil {
			if err != nil {
				rc <- &instrospectionClientResult{
					err: oidc.ErrUnauthorizedClient().WithParent(err),
				}
				return
			}
		}

	}

	// TODO: give clients their own aggregate, so we can skip this query
	projectID, err := s.storage.query.ProjectIDFromClientID(ctx, clientID, false)
	if err != nil {
		rc <- &instrospectionClientResult{err: err}
		return
	}

	rc <- &instrospectionClientResult{
		clientID:  clientID,
		projectID: projectID,
	}
}

type introspectionTokenResult struct {
	tokenID         string
	userID          string
	subject         string
	clientID        string
	audience        []string
	scope           []string
	tokenCreation   time.Time
	tokenExpiration time.Time
	isPAT           bool

	err error
}

func (s *Server) introspectionToken(ctx context.Context, accessToken string, rc chan<- *introspectionTokenResult) {
	var tokenID, subject string

	if tokenIDSubject, err := s.Provider().Crypto().Decrypt(accessToken); err == nil {
		split := strings.Split(tokenIDSubject, ":")
		if len(split) != 2 {
			rc <- &introspectionTokenResult{err: errors.New("invalid token format")}
			return
		}
		tokenID, subject = split[0], split[1]
	} else {
		verifier := op.NewAccessTokenVerifier(op.IssuerFromContext(ctx), s.storage.keySet)
		claims, err := op.VerifyAccessToken[*oidc.AccessTokenClaims](ctx, accessToken, verifier)
		if err != nil {
			rc <- &introspectionTokenResult{err: err}
			return
		}
		tokenID, subject = claims.JWTID, claims.Subject
	}

	if strings.HasPrefix(tokenID, command.IDPrefixV2) {
		token, err := s.storage.query.ActiveAccessTokenByToken(ctx, tokenID)
		if err != nil {
			rc <- &introspectionTokenResult{err: err}
			return
		}
		rc <- &introspectionTokenResult{
			tokenID:         tokenID,
			userID:          token.UserID,
			subject:         subject,
			clientID:        token.ClientID,
			audience:        token.Audience,
			scope:           token.Scope,
			tokenCreation:   token.AccessTokenCreation,
			tokenExpiration: token.AccessTokenExpiration,
		}
		return
	}

	token, err := s.storage.repo.TokenByIDs(ctx, subject, tokenID)
	if err != nil {
		rc <- &introspectionTokenResult{
			err: errz.ThrowPermissionDenied(err, "OIDC-Dsfb2", "token is not valid or has expired"),
		}
		return
	}
	rc <- &introspectionTokenResult{
		tokenID:         tokenID,
		userID:          token.UserID,
		subject:         subject,
		clientID:        token.ApplicationID, // check correctness?
		audience:        token.Audience,
		scope:           token.Scopes,
		tokenCreation:   token.CreationDate,
		tokenExpiration: token.Expiration,
		isPAT:           token.IsPAT,
	}
}

func validateIntrospectionAudience(audience []string, clientID, projectID string) error {
	if slices.ContainsFunc(audience, func(entry string) bool {
		return entry == clientID || entry == projectID
	}) {
		return nil
	}

	return errz.ThrowPermissionDenied(nil, "OIDC-sdg3G", "token is not valid for this client")
}
