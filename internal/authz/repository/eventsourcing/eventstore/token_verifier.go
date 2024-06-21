package eventstore

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/muhlemmer/gu"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/authz/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	usr_model "github.com/zitadel/zitadel/internal/user/model"
	usr_view "github.com/zitadel/zitadel/internal/user/repository/view"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type TokenVerifierRepo struct {
	TokenVerificationKey crypto.EncryptionAlgorithm
	Eventstore           *eventstore.Eventstore
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

	// always load the latest sequence first, so in case the token was not found by id,
	// the sequence will be equal or lower than the actual projection and no events are lost
	sequence, err := repo.View.GetLatestState(ctx)
	logging.WithFields("instanceID", instanceID, "userID", userID, "tokenID", tokenID).
		OnError(err).
		Errorf("could not get current sequence for token check")

	token, viewErr := repo.View.TokenByIDs(tokenID, userID, instanceID)
	if viewErr != nil && !zerrors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if zerrors.IsNotFound(viewErr) {
		token = new(model.TokenView)
		token.ID = tokenID
		token.UserID = userID
		if sequence != nil {
			token.ChangeDate = sequence.EventCreatedAt
		}
	}

	events, esErr := repo.getUserEvents(ctx, userID, instanceID, token.ChangeDate, token.GetRelevantEventTypes())
	if zerrors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, zerrors.ThrowNotFound(nil, "EVENT-4T90g", "Errors.Token.NotFound")
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
		return nil, zerrors.ThrowNotFound(nil, "EVENT-5Bm9s", "Errors.Token.NotFound")
	}
	return model.TokenViewToModel(token), nil
}

func (repo *TokenVerifierRepo) VerifyAccessToken(ctx context.Context, tokenString, verifierClientID, projectID string) (userID string, agentID string, clientID, prefLang, resourceOwner string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	tokenID, subject, ok := repo.getTokenIDAndSubject(ctx, tokenString)
	if !ok {
		return "", "", "", "", "", zerrors.ThrowUnauthenticated(nil, "APP-Reb32", "invalid token")
	}
	if strings.HasPrefix(tokenID, command.IDPrefixV2) {
		return repo.verifyAccessTokenV2(ctx, tokenID, verifierClientID, projectID)
	}
	if sessionID, ok := strings.CutPrefix(tokenID, authz.SessionTokenPrefix); ok {
		userID, clientID, resourceOwner, err = repo.verifySessionToken(ctx, sessionID, tokenString)
		return
	}
	return repo.verifyAccessTokenV1(ctx, tokenID, subject, verifierClientID, projectID)
}

func (repo *TokenVerifierRepo) verifyAccessTokenV1(ctx context.Context, tokenID, subject, verifierClientID, projectID string) (userID, agentID, clientID, prefLang, resourceOwner string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	_, tokenSpan := tracing.NewNamedSpan(ctx, "tokenByID")
	token, err := repo.tokenByID(ctx, tokenID, subject)
	tokenSpan.EndWithError(err)
	if err != nil {
		return "", "", "", "", "", zerrors.ThrowUnauthenticated(err, "APP-BxUSiL", "invalid token")
	}
	if token.Actor != nil {
		return "", "", "", "", "", zerrors.ThrowPermissionDenied(nil, "APP-wai8O", "Errors.TokenExchange.Token.NotForAPI")
	}
	if !token.Expiration.After(time.Now().UTC()) {
		return "", "", "", "", "", zerrors.ThrowUnauthenticated(err, "APP-k9KS0", "invalid token")
	}
	if token.IsPAT {
		return token.UserID, "", "", "", token.ResourceOwner, nil
	}
	if err = verifyAudience(token.Audience, verifierClientID, projectID); err != nil {
		return "", "", "", "", "", err
	}
	return token.UserID, token.UserAgentID, token.ApplicationID, token.PreferredLanguage, token.ResourceOwner, nil
}

func (repo *TokenVerifierRepo) verifyAccessTokenV2(ctx context.Context, token, verifierClientID, projectID string) (userID, agentID, clientID, prefLang, resourceOwner string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	activeToken, err := repo.Query.ActiveAccessTokenByToken(ctx, token)
	if err != nil {
		return "", "", "", "", "", err
	}
	if activeToken.Actor != nil {
		return "", "", "", "", "", zerrors.ThrowPermissionDenied(nil, "APP-Shi0J", "Errors.TokenExchange.Token.NotForAPI")
	}
	if err = verifyAudience(activeToken.Audience, verifierClientID, projectID); err != nil {
		return "", "", "", "", "", err
	}
	if err = repo.checkAuthentication(ctx, activeToken.AuthMethods, activeToken.UserID); err != nil {
		return "", "", "", "", "", err
	}
	prefLang = gu.Value(activeToken.PreferredLanguage).String()
	agentID = gu.Value(gu.Value(activeToken.UserAgent).FingerprintID)

	return activeToken.UserID, agentID, activeToken.ClientID, prefLang, activeToken.ResourceOwner, nil
}

func (repo *TokenVerifierRepo) verifySessionToken(ctx context.Context, sessionID, token string) (userID, clientID, resourceOwner string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	session, err := repo.Query.SessionByID(ctx, true, sessionID, token)
	if err != nil {
		return "", "", "", err
	}
	if !session.Expiration.IsZero() && session.Expiration.Before(time.Now()) {
		return "", "", "", zerrors.ThrowPermissionDenied(nil, "AUTHZ-EGDo3", "session expired")
	}
	if err = repo.checkAuthentication(ctx, authMethodsFromSession(session), session.UserFactor.UserID); err != nil {
		return "", "", "", err
	}
	return session.UserFactor.UserID, "", session.UserFactor.ResourceOwner, nil
}

// checkAuthentication ensures the session or token was authenticated (at least a single [domain.UserAuthMethodType]).
// It will also check if there was a multi factor authentication, if either MFA is forced by the login policy or if the user has set up any second factor
func (repo *TokenVerifierRepo) checkAuthentication(ctx context.Context, authMethods []domain.UserAuthMethodType, userID string) error {
	if len(authMethods) == 0 {
		return zerrors.ThrowPermissionDenied(nil, "AUTHZ-Kl3p0", "authentication required")
	}
	if domain.HasMFA(authMethods) {
		return nil
	}
	requirements, err := repo.Query.ListUserAuthMethodTypesRequired(setCallerCtx(ctx, userID), userID)
	if err != nil {
		return err
	}
	if requirements.UserType == domain.UserTypeMachine {
		return nil
	}
	if domain.RequiresMFA(
		requirements.ForceMFA,
		requirements.ForceMFALocalOnly,
		!hasIDPAuthentication(authMethods),
	) {
		return zerrors.ThrowPermissionDenied(nil, "AUTHZ-Kl3p0", "mfa required")
	}
	return nil
}

func hasIDPAuthentication(authMethods []domain.UserAuthMethodType) bool {
	for _, method := range authMethods {
		if method == domain.UserAuthMethodTypeIDP {
			return true
		}
	}
	return false
}

func authMethodsFromSession(session *query.Session) []domain.UserAuthMethodType {
	types := make([]domain.UserAuthMethodType, 0, domain.UserAuthMethodTypeIDP)
	if !session.PasswordFactor.PasswordCheckedAt.IsZero() {
		types = append(types, domain.UserAuthMethodTypePassword)
	}
	if !session.WebAuthNFactor.WebAuthNCheckedAt.IsZero() {
		if session.WebAuthNFactor.UserVerified {
			types = append(types, domain.UserAuthMethodTypePasswordless)
		} else {
			types = append(types, domain.UserAuthMethodTypeU2F)
		}
	}
	if !session.IntentFactor.IntentCheckedAt.IsZero() {
		types = append(types, domain.UserAuthMethodTypeIDP)
	}
	if !session.TOTPFactor.TOTPCheckedAt.IsZero() {
		types = append(types, domain.UserAuthMethodTypeTOTP)
	}
	if !session.OTPSMSFactor.OTPCheckedAt.IsZero() {
		types = append(types, domain.UserAuthMethodTypeOTPSMS)
	}
	if !session.OTPEmailFactor.OTPCheckedAt.IsZero() {
		types = append(types, domain.UserAuthMethodTypeOTPEmail)
	}
	return types
}

func setCallerCtx(ctx context.Context, userID string) context.Context {
	ctxData := authz.GetCtxData(ctx)
	ctxData.UserID = userID
	return authz.SetCtxData(ctx, ctxData)
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

func (repo *TokenVerifierRepo) getUserEvents(ctx context.Context, userID, instanceID string, changeDate time.Time, eventTypes []eventstore.EventType) (_ []eventstore.Event, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	query, err := usr_view.UserByIDQuery(userID, instanceID, changeDate, eventTypes)
	if err != nil {
		return nil, err
	}
	return repo.Eventstore.Filter(ctx, query)
}

// getTokenIDAndSubject returns the TokenID and Subject of both opaque tokens and JWTs
func (repo *TokenVerifierRepo) getTokenIDAndSubject(ctx context.Context, accessToken string) (tokenID string, subject string, valid bool) {
	// accessToken can be either opaque or JWT
	// let's try opaque first:
	tokenIDSubject, err := repo.decryptAccessToken(accessToken)
	if err != nil {
		logging.WithError(err).Warn("token verifier repo: decrypt access token")
		// if decryption did not work, it might be a JWT
		accessTokenClaims, err := op.VerifyAccessToken[*oidc.AccessTokenClaims](ctx, accessToken, repo.jwtTokenVerifier(ctx))
		if err != nil {
			logging.WithError(err).Warn("token verifier repo: verify JWT access token")
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

func (repo *TokenVerifierRepo) jwtTokenVerifier(ctx context.Context) *op.AccessTokenVerifier {
	keySet := &openIDKeySet{repo.Query}
	issuer := http_util.BuildOrigin(authz.GetInstance(ctx).RequestedHost(), repo.ExternalSecure)
	return op.NewAccessTokenVerifier(issuer, keySet)
}

func (repo *TokenVerifierRepo) decryptAccessToken(token string) (string, error) {
	tokenData, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return "", zerrors.ThrowUnauthenticated(nil, "APP-ASdgg", "invalid token")
	}
	tokenIDSubject, err := repo.TokenVerificationKey.DecryptString(tokenData, repo.TokenVerificationKey.EncryptionKeyID())
	if err != nil {
		return "", zerrors.ThrowUnauthenticated(nil, "APP-8EF0zZ", "invalid token")
	}
	return tokenIDSubject, nil
}

func verifyAudience(audience []string, verifierClientID, projectID string) error {
	for _, aud := range audience {
		if verifierClientID == aud || projectID == aud {
			return nil
		}
	}
	return zerrors.ThrowUnauthenticated(nil, "APP-Zxfako", "invalid audience")
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
