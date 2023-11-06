package oidc

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	errz "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/user/model"
)

func (s *Server) Introspect(ctx context.Context, r *op.Request[op.IntrospectionRequest]) (resp *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
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

	// remaining errors shoudn't be returned to the client,
	// so we catch errors here, log them and return the response
	// with active: false
	defer func() {
		if err != nil {
			s.getLogger(ctx).ErrorContext(ctx, "oidc introspection", "err", err)
		}
		resp, err = op.NewResponse(new(oidc.IntrospectionResponse)), nil
	}()

	if err != nil {
		return nil, err
	}
	if err = validateIntrospectionAudience(token.audience, client.clientID, client.projectID); err != nil {
		return nil, err
	}
	userInfo, err := s.getUserInfoWithRoles(ctx, token.userID, client.projectID, token.scope, []string{client.projectID})
	if err != nil {
		return nil, err
	}
	introspectionResp := &oidc.IntrospectionResponse{
		Active:     true,
		Scope:      token.scope,
		ClientID:   token.clientID,
		TokenType:  oidc.BearerToken,
		Expiration: oidc.FromTime(token.tokenExpiration),
		IssuedAt:   oidc.FromTime(token.tokenCreation),
		NotBefore:  oidc.FromTime(token.tokenCreation),
		Audience:   token.audience,
		Issuer:     op.IssuerFromContext(ctx),
		JWTID:      token.tokenID,
	}
	introspectionResp.SetUserInfo(userInfo)
	return op.NewResponse(introspectionResp), nil
}

type instrospectionClientResult struct {
	clientID  string
	projectID string
	err       error
}

func (s *Server) instrospectionClientAuth(ctx context.Context, cc *op.ClientCredentials, rc chan<- *instrospectionClientResult) {
	ctx, span := tracing.NewSpan(ctx)

	clientID, projectID, err := func() (string, string, error) {
		client, err := s.clientFromCredentials(ctx, cc)
		if err != nil {
			return "", "", err
		}

		if cc.ClientAssertion != "" {
			verifier := op.NewJWTProfileVerifierKeySet(keySetMap(client.PublicKeys), op.IssuerFromContext(ctx), time.Hour, time.Second)
			if _, err := op.VerifyJWTAssertion(ctx, cc.ClientAssertion, verifier); err != nil {
				return "", "", oidc.ErrUnauthorizedClient().WithParent(err)
			}
		} else {
			if err := crypto.CompareHash(client.ClientSecret, []byte(cc.ClientSecret), s.hashAlg); err != nil {
				return "", "", oidc.ErrUnauthorizedClient().WithParent(err)
			}
		}

		return client.ClientID, client.ProjectID, nil
	}()

	span.EndWithError(err)

	rc <- &instrospectionClientResult{
		clientID:  clientID,
		projectID: projectID,
		err:       err,
	}
}

// clientFromCredentials parses the client ID early,
// and makes a single query for the client for either auth methods.
func (s *Server) clientFromCredentials(ctx context.Context, cc *op.ClientCredentials) (client *query.IntrospectionClient, err error) {
	if cc.ClientAssertion != "" {
		claims := new(oidc.JWTTokenRequest)
		if _, err := oidc.ParseToken(cc.ClientAssertion, claims); err != nil {
			return nil, oidc.ErrUnauthorizedClient().WithParent(err)
		}
		client, err = s.storage.query.GetIntrospectionClientByID(ctx, claims.Issuer, true)
	} else {
		client, err = s.storage.query.GetIntrospectionClientByID(ctx, cc.ClientID, false)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, oidc.ErrUnauthorizedClient().WithParent(err)
	}
	// any other error is regarded internal and should not be reported back to the client.
	return client, err
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
	ctx, span := tracing.NewSpan(ctx)

	result, err := func() (_ *introspectionTokenResult, err error) {
		var tokenID, subject string

		if tokenIDSubject, err := s.Provider().Crypto().Decrypt(accessToken); err == nil {
			split := strings.Split(tokenIDSubject, ":")
			if len(split) != 2 {
				return nil, errors.New("invalid token format")
			}
			tokenID, subject = split[0], split[1]
		} else {
			verifier := op.NewAccessTokenVerifier(op.IssuerFromContext(ctx), s.storage.keySet)
			claims, err := op.VerifyAccessToken[*oidc.AccessTokenClaims](ctx, accessToken, verifier)
			if err != nil {
				return nil, err
			}
			tokenID, subject = claims.JWTID, claims.Subject
		}

		if strings.HasPrefix(tokenID, command.IDPrefixV2) {
			token, err := s.storage.query.ActiveAccessTokenByToken(ctx, tokenID)
			if err != nil {
				rc <- &introspectionTokenResult{err: err}
				return nil, err
			}
			return introspectionTokenResultV2(tokenID, subject, token), nil
		}

		token, err := s.storage.repo.TokenByIDs(ctx, subject, tokenID)
		if err != nil {
			return nil, errz.ThrowPermissionDenied(err, "OIDC-Dsfb2", "token is not valid or has expired")
		}
		return introspectionTokenResultV1(tokenID, subject, token), nil
	}()

	span.EndWithError(err)

	if err != nil {
		rc <- &introspectionTokenResult{err: err}
		return
	}
	rc <- result
}

func introspectionTokenResultV1(tokenID, subject string, token *model.TokenView) *introspectionTokenResult {
	return &introspectionTokenResult{
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

func introspectionTokenResultV2(tokenID, subject string, token *query.OIDCSessionAccessTokenReadModel) *introspectionTokenResult {
	return &introspectionTokenResult{
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

func validateIntrospectionAudience(audience []string, clientID, projectID string) error {
	if slices.ContainsFunc(audience, func(entry string) bool {
		return entry == clientID || entry == projectID
	}) {
		return nil
	}

	return errz.ThrowPermissionDenied(nil, "OIDC-sdg3G", "token is not valid for this client")
}
