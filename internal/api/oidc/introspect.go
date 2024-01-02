package oidc

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (s *Server) Introspect(ctx context.Context, r *op.Request[op.IntrospectionRequest]) (resp *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if s.features.LegacyIntrospection {
		return s.LegacyServer.Introspect(ctx, r)
	}
	if s.features.TriggerIntrospectionProjections {
		// Execute all triggers in one concurrent sweep.
		query.TriggerIntrospectionProjections(ctx)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	clientChan := make(chan *introspectionClientResult)
	go s.introspectionClientAuth(ctx, r.Data.ClientCredentials, clientChan)

	tokenChan := make(chan *introspectionTokenResult)
	go s.introspectionToken(ctx, r.Data.Token, tokenChan)

	var (
		client *introspectionClientResult
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

		// we only care for the first error that occurred,
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

	// remaining errors shouldn't be returned to the client,
	// so we catch errors here, log them and return the response
	// with active: false
	defer func() {
		if err != nil {
			s.getLogger(ctx).ErrorContext(ctx, "oidc introspection", "err", err)
			resp, err = op.NewResponse(new(oidc.IntrospectionResponse)), nil
		}
	}()

	if err != nil {
		return nil, err
	}

	// TODO: can we get rid of this separate query?
	if token.isPAT {
		if err = s.assertClientScopesForPAT(ctx, token.accessToken, client.clientID, client.projectID); err != nil {
			return nil, err
		}
	}

	if err = validateIntrospectionAudience(token.audience, client.clientID, client.projectID); err != nil {
		return nil, err
	}
	userInfo, err := s.userInfo(ctx, token.userID, client.projectID, token.scope, []string{client.projectID})
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

type introspectionClientResult struct {
	clientID  string
	projectID string
	err       error
}

var errNoClientSecret = errors.New("client has no configured secret")

func (s *Server) introspectionClientAuth(ctx context.Context, cc *op.ClientCredentials, rc chan<- *introspectionClientResult) {
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
			return client.ClientID, client.ProjectID, nil

		}
		if client.ClientSecret != nil {
			if err := crypto.CompareHash(client.ClientSecret, []byte(cc.ClientSecret), s.hashAlg); err != nil {
				return "", "", oidc.ErrUnauthorizedClient().WithParent(err)
			}
			return client.ClientID, client.ProjectID, nil
		}
		return "", "", oidc.ErrUnauthorizedClient().WithParent(errNoClientSecret)
	}()

	span.EndWithError(err)

	rc <- &introspectionClientResult{
		clientID:  clientID,
		projectID: projectID,
		err:       err,
	}
}

// clientFromCredentials parses the client ID early,
// and makes a single query for the client for either auth methods.
func (s *Server) clientFromCredentials(ctx context.Context, cc *op.ClientCredentials) (client *query.IntrospectionClient, err error) {
	clientID, assertion, err := clientIDFromCredentials(cc)
	if err != nil {
		return nil, err
	}
	client, err = s.query.GetIntrospectionClientByID(ctx, clientID, assertion)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, oidc.ErrUnauthorizedClient().WithParent(err)
	}
	// any other error is regarded internal and should not be reported back to the client.
	return client, err
}

type introspectionTokenResult struct {
	*accessToken
	err error
}

func (s *Server) introspectionToken(ctx context.Context, tkn string, rc chan<- *introspectionTokenResult) {
	ctx, span := tracing.NewSpan(ctx)
	token, err := s.verifyAccessToken(ctx, tkn)
	span.EndWithError(err)

	rc <- &introspectionTokenResult{
		accessToken: token,
		err:         err,
	}
}

func validateIntrospectionAudience(audience []string, clientID, projectID string) error {
	if slices.ContainsFunc(audience, func(entry string) bool {
		return entry == clientID || entry == projectID
	}) {
		return nil
	}

	return zerrors.ThrowPermissionDenied(nil, "OIDC-sdg3G", "token is not valid for this client")
}
