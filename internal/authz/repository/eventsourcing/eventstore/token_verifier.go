package eventstore

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"github.com/zitadel/oidc/v2/pkg/op"
	"gopkg.in/square/go-jose.v2"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/authz/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/crypto"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	usr_model "github.com/zitadel/zitadel/internal/user/model"
	usr_view "github.com/zitadel/zitadel/internal/user/repository/view"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
)

type TokenVerifierRepo struct {
	TokenVerificationKey crypto.EncryptionAlgorithm
	IAMID                string
	Eventstore           v1.Eventstore
	View                 *view.View
	Query                *query.Queries
	ExternalSecure       bool
}

func (repo *TokenVerifierRepo) Health() error {
	return repo.View.Health()
}

func (repo *TokenVerifierRepo) tokenByID(ctx context.Context, tokenID, userID string) (_ *usr_model.TokenView, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	instanceID := authz.GetInstance(ctx).InstanceID()
	sequence, err := repo.View.GetLatestTokenSequence(instanceID)
	logging.WithFields("instanceID", instanceID, "userID", userID, "tokenID").
		OnError(err).
		Errorf("could not get current sequence for token check")

	token, viewErr := repo.View.TokenByIDs(tokenID, userID, instanceID)
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		token = new(model.TokenView)
		token.ID = tokenID
		token.UserID = userID
		if sequence != nil {
			token.Sequence = sequence.CurrentSequence
		}
	}

	events, esErr := repo.getUserEvents(ctx, userID, instanceID, token.Sequence)
	if caos_errs.IsNotFound(viewErr) && len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-4T90g", "Errors.Token.NotFound")
	}

	if esErr != nil {
		logging.WithError(viewErr).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("error retrieving new events")
		return model.TokenViewToModel(token), nil
	}
	viewToken := *token
	for _, event := range events {
		err := token.AppendEventIfMyToken(event)
		if err != nil {
			return model.TokenViewToModel(&viewToken), nil
		}
	}
	if !token.Expiration.After(time.Now().UTC()) || token.Deactivated {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-5Bm9s", "Errors.Token.NotFound")
	}
	return model.TokenViewToModel(token), nil
}

func (repo *TokenVerifierRepo) VerifyAccessToken(ctx context.Context, tokenString, verifierClientID, projectID string) (userID string, agentID string, clientID, prefLang, resourceOwner string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	tokenID, subject, ok := repo.getTokenIDAndSubject(ctx, tokenString)
	if !ok {
		return "", "", "", "", "", caos_errs.ThrowUnauthenticated(nil, "APP-Reb32", "invalid token")
	}
	_, tokenSpan := tracing.NewNamedSpan(ctx, "token")
	token, err := repo.tokenByID(ctx, tokenID, subject)
	tokenSpan.EndWithError(err)
	if err != nil {
		return "", "", "", "", "", caos_errs.ThrowUnauthenticated(err, "APP-BxUSiL", "invalid token")
	}
	if !token.Expiration.After(time.Now().UTC()) {
		return "", "", "", "", "", caos_errs.ThrowUnauthenticated(err, "APP-k9KS0", "invalid token")
	}
	if token.IsPAT {
		return token.UserID, "", "", "", token.ResourceOwner, nil
	}
	for _, aud := range token.Audience {
		if verifierClientID == aud || projectID == aud {
			return token.UserID, token.UserAgentID, token.ApplicationID, token.PreferredLanguage, token.ResourceOwner, nil
		}
	}
	return "", "", "", "", "", caos_errs.ThrowUnauthenticated(nil, "APP-Zxfako", "invalid audience")
}

func (repo *TokenVerifierRepo) ProjectIDAndOriginsByClientID(ctx context.Context, clientID string) (projectID string, origins []string, err error) {
	app, err := repo.View.ApplicationByOIDCClientID(ctx, clientID)
	if err != nil {
		return "", nil, err
	}
	return app.ProjectID, app.OIDCConfig.AllowedOrigins, nil
}

func (repo *TokenVerifierRepo) VerifierClientID(ctx context.Context, appName string) (clientID, projectID string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	app, err := repo.View.ApplicationByProjecIDAndAppName(ctx, authz.GetInstance(ctx).ProjectID(), appName)
	if err != nil {
		return "", "", err
	}
	if app.OIDCConfig != nil {
		clientID = app.OIDCConfig.ClientID
	} else if app.APIConfig != nil {
		clientID = app.APIConfig.ClientID
	}
	return clientID, app.ProjectID, nil
}

func (repo *TokenVerifierRepo) getUserEvents(ctx context.Context, userID, instanceID string, sequence uint64) (_ []*models.Event, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	query, err := usr_view.UserByIDQuery(userID, instanceID, sequence)
	if err != nil {
		return nil, err
	}
	return repo.Eventstore.FilterEvents(ctx, query)
}

// getTokenIDAndSubject returns the TokenID and Subject of both opaque tokens and JWTs
func (repo *TokenVerifierRepo) getTokenIDAndSubject(ctx context.Context, accessToken string) (tokenID string, subject string, valid bool) {
	// accessToken can be either opaque or JWT
	// let's try opaque first:
	tokenIDSubject, err := repo.decryptAccessToken(accessToken)
	if err != nil {
		// if decryption did not work, it might be a JWT
		accessTokenClaims, err := op.VerifyAccessToken[*oidc.AccessTokenClaims](ctx, accessToken, repo.jwtTokenVerifier(ctx))
		if err != nil {
			return "", "", false
		}
		return accessTokenClaims.JWTID, accessTokenClaims.Subject, true
	}
	splitToken := strings.Split(tokenIDSubject, ":")
	if len(splitToken) != 2 {
		return "", "", false
	}
	return splitToken[0], splitToken[1], true
}

func (repo *TokenVerifierRepo) jwtTokenVerifier(ctx context.Context) op.AccessTokenVerifier {
	keySet := &openIDKeySet{repo.Query}
	issuer := http_util.BuildOrigin(authz.GetInstance(ctx).RequestedHost(), repo.ExternalSecure)
	return op.NewAccessTokenVerifier(issuer, keySet)
}

func (repo *TokenVerifierRepo) decryptAccessToken(token string) (string, error) {
	tokenData, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return "", caos_errs.ThrowUnauthenticated(nil, "APP-ASdgg", "invalid token")
	}
	tokenIDSubject, err := repo.TokenVerificationKey.DecryptString(tokenData, repo.TokenVerificationKey.EncryptionKeyID())
	if err != nil {
		return "", caos_errs.ThrowUnauthenticated(nil, "APP-8EF0zZ", "invalid token")
	}
	return tokenIDSubject, nil
}

type openIDKeySet struct {
	*query.Queries
}

// VerifySignature implements the oidc.KeySet interface
// providing an implementation for the keys retrieved directly from Queries
func (o *openIDKeySet) VerifySignature(ctx context.Context, jws *jose.JSONWebSignature) ([]byte, error) {
	keySet, err := o.Queries.ActivePublicKeys(ctx, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error fetching keys: %w", err)
	}
	keyID, alg := oidc.GetKeyIDAndAlg(jws)
	key, err := oidc.FindMatchingKey(keyID, oidc.KeyUseSignature, alg, jsonWebKeys(keySet.Keys)...)
	if err != nil {
		return nil, fmt.Errorf("invalid signature: %w", err)
	}
	return jws.Verify(&key)
}

func jsonWebKeys(keys []query.PublicKey) []jose.JSONWebKey {
	webKeys := make([]jose.JSONWebKey, len(keys))
	for i, key := range keys {
		webKeys[i] = jose.JSONWebKey{
			KeyID:     key.ID(),
			Algorithm: key.Algorithm(),
			Use:       key.Use().String(),
			Key:       key.Key(),
		}
	}
	return webKeys
}
