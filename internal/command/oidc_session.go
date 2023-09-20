package command

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/repository/authrequest"
	"github.com/zitadel/zitadel/internal/repository/oidcsession"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	TokenDelimiter            = "-"
	AccessTokenPrefix         = "at_"
	RefreshTokenPrefix        = "rt_"
	oidcTokenSubjectDelimiter = ":"
	oidcTokenFormat           = "%s" + oidcTokenSubjectDelimiter + "%s"
)

// AddOIDCSessionAccessToken creates a new OIDC Session, creates an access token and returns its id and expiration.
// If the underlying [AuthRequest] is a OIDC Auth Code Flow, it will set the code as exchanged.
func (c *Commands) AddOIDCSessionAccessToken(ctx context.Context, authRequestID string) (string, time.Time, error) {
	cmd, err := c.newOIDCSessionAddEvents(ctx, authRequestID)
	if err != nil {
		return "", time.Time{}, err
	}
	cmd.AddSession(ctx)
	if err = cmd.AddAccessToken(ctx, cmd.authRequestWriteModel.Scope); err != nil {
		return "", time.Time{}, err
	}
	cmd.SetAuthRequestSuccessful(ctx)
	accessTokenID, _, accessTokenExpiration, err := cmd.PushEvents(ctx)
	return accessTokenID, accessTokenExpiration, err
}

// AddOIDCSessionRefreshAndAccessToken creates a new OIDC Session, creates an access token and refresh token.
// It returns the access token id, expiration and the refresh token.
// If the underlying [AuthRequest] is a OIDC Auth Code Flow, it will set the code as exchanged.
func (c *Commands) AddOIDCSessionRefreshAndAccessToken(ctx context.Context, authRequestID string) (tokenID, refreshToken string, tokenExpiration time.Time, err error) {
	cmd, err := c.newOIDCSessionAddEvents(ctx, authRequestID)
	if err != nil {
		return "", "", time.Time{}, err
	}
	cmd.AddSession(ctx)
	if err = cmd.AddAccessToken(ctx, cmd.authRequestWriteModel.Scope); err != nil {
		return "", "", time.Time{}, err
	}
	if err = cmd.AddRefreshToken(ctx); err != nil {
		return "", "", time.Time{}, err
	}
	cmd.SetAuthRequestSuccessful(ctx)
	return cmd.PushEvents(ctx)
}

// ExchangeOIDCSessionRefreshAndAccessToken updates an existing OIDC Session, creates a new access and refresh token.
// It returns the access token id and expiration and the new refresh token.
func (c *Commands) ExchangeOIDCSessionRefreshAndAccessToken(ctx context.Context, oidcSessionID, refreshToken string, scope []string) (tokenID, newRefreshToken string, tokenExpiration time.Time, err error) {
	cmd, err := c.newOIDCSessionUpdateEvents(ctx, oidcSessionID, refreshToken)
	if err != nil {
		return "", "", time.Time{}, err
	}
	if err = cmd.AddAccessToken(ctx, scope); err != nil {
		return "", "", time.Time{}, err
	}
	if err = cmd.RenewRefreshToken(ctx); err != nil {
		return "", "", time.Time{}, err
	}
	return cmd.PushEvents(ctx)
}

// OIDCSessionByRefreshToken computes the current state of an existing OIDCSession by a refresh_token (to start a Refresh Token Grant).
// If either the session is not active, the token is invalid or expired (incl. idle expiration) an invalid refresh token error will be returned.
func (c *Commands) OIDCSessionByRefreshToken(ctx context.Context, refreshToken string) (*OIDCSessionWriteModel, error) {
	oidcSessionID, refreshTokenID, err := parseRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	writeModel := NewOIDCSessionWriteModel(oidcSessionID, "")
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(err, "OIDCS-SAF31", "Errors.OIDCSession.RefreshTokenInvalid")
	}
	if err = writeModel.CheckRefreshToken(refreshTokenID); err != nil {
		return nil, err
	}
	return writeModel, nil
}

func oidcSessionTokenIDsFromToken(token string) (oidcSessionID, refreshTokenID, accessTokenID string, err error) {
	split := strings.Split(token, TokenDelimiter)
	if len(split) != 2 {
		return "", "", "", caos_errs.ThrowPreconditionFailed(nil, "OIDCS-S87kl", "Errors.OIDCSession.Token.Invalid")
	}
	if strings.HasPrefix(split[1], RefreshTokenPrefix) {
		return split[0], split[1], "", nil
	}
	if strings.HasPrefix(split[1], AccessTokenPrefix) {
		return split[0], "", split[1], nil
	}
	return "", "", "", caos_errs.ThrowPreconditionFailed(nil, "OIDCS-S87kl", "Errors.OIDCSession.Token.Invalid")
}

// RevokeOIDCSessionToken revokes an access_token or refresh_token
// if the OIDCSession cannot be retrieved by the provided token, is not active or if the token is already revoked,
// then no error will be returned.
// The only possible error (except db connection or other internal errors) occurs if a client tries to revoke a token,
// which was not part of the audience.
func (c *Commands) RevokeOIDCSessionToken(ctx context.Context, token, clientID string) (err error) {
	oidcSessionID, refreshTokenID, accessTokenID, err := oidcSessionTokenIDsFromToken(token)
	if err != nil {
		logging.WithError(err).Info("token revocation with invalid token format")
		return nil
	}
	writeModel := NewOIDCSessionWriteModel(oidcSessionID, "")
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return caos_errs.ThrowInternal(err, "OIDCS-NB3t2", "Errors.Internal")
	}
	if err = writeModel.CheckClient(clientID); err != nil {
		return err
	}
	if refreshTokenID != "" {
		if err = writeModel.CheckRefreshToken(refreshTokenID); err != nil {
			logging.WithFields("oidcSessionID", oidcSessionID, "refreshTokenID", refreshTokenID).WithError(err).
				Info("refresh token revocation with invalid token")
			return nil
		}
		return c.pushAppendAndReduce(ctx, writeModel, oidcsession.NewRefreshTokenRevokedEvent(ctx, writeModel.aggregate))
	}
	if err = writeModel.CheckAccessToken(accessTokenID); err != nil {
		logging.WithFields("oidcSessionID", oidcSessionID, "accessTokenID", accessTokenID).WithError(err).
			Info("access token revocation with invalid token")
		return nil
	}
	return c.pushAppendAndReduce(ctx, writeModel, oidcsession.NewAccessTokenRevokedEvent(ctx, writeModel.aggregate))
}

func (c *Commands) newOIDCSessionAddEvents(ctx context.Context, authRequestID string) (*OIDCSessionEvents, error) {
	authRequestWriteModel, err := c.getAuthRequestWriteModel(ctx, authRequestID)
	if err != nil {
		return nil, err
	}
	if err = authRequestWriteModel.CheckAuthenticated(); err != nil {
		return nil, err
	}
	sessionWriteModel := NewSessionWriteModel(authRequestWriteModel.SessionID, authz.GetCtxData(ctx).OrgID)
	err = c.eventstore.FilterToQueryReducer(ctx, sessionWriteModel)
	if err != nil {
		return nil, err
	}
	if sessionWriteModel.State != domain.SessionStateActive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "OIDCS-sjkl3", "Errors.Session.Terminated")
	}
	resourceOwner, err := c.getResourceOwnerOfSessionUser(ctx, sessionWriteModel.UserID, sessionWriteModel.InstanceID)
	if err != nil {
		return nil, err
	}
	accessTokenLifetime, refreshTokenLifeTime, refreshTokenIdleLifetime, err := c.tokenTokenLifetimes(ctx)
	if err != nil {
		return nil, err
	}
	sessionID, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	sessionID = IDPrefixV2 + sessionID
	return &OIDCSessionEvents{
		eventstore:               c.eventstore,
		idGenerator:              c.idGenerator,
		encryptionAlg:            c.keyAlgorithm,
		oidcSessionWriteModel:    NewOIDCSessionWriteModel(sessionID, resourceOwner),
		sessionWriteModel:        sessionWriteModel,
		authRequestWriteModel:    authRequestWriteModel,
		accessTokenLifetime:      accessTokenLifetime,
		refreshTokenLifeTime:     refreshTokenLifeTime,
		refreshTokenIdleLifetime: refreshTokenIdleLifetime,
	}, nil
}

func (c *Commands) getResourceOwnerOfSessionUser(ctx context.Context, userID, instanceID string) (string, error) {
	events, err := c.eventstore.Filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(instanceID).
		AllowTimeTravel().
		OrderAsc().
		Limit(1).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(userID).
		Builder())
	if err != nil || len(events) != 1 {
		return "", caos_errs.ThrowInternal(err, "OIDCS-sferh", "Errors.Internal")
	}
	return events[0].Aggregate().ResourceOwner, nil
}

func (c *Commands) decryptRefreshToken(refreshToken string) (refreshTokenID string, err error) {
	decoded, err := base64.RawURLEncoding.DecodeString(refreshToken)
	if err != nil {
		return "", err
	}
	decrypted, err := c.keyAlgorithm.DecryptString(decoded, c.keyAlgorithm.EncryptionKeyID())
	if err != nil {
		return "", err
	}
	_, refreshTokenID, err = parseRefreshToken(decrypted)
	return refreshTokenID, err
}

func parseRefreshToken(refreshToken string) (oidcSessionID, refreshTokenID string, err error) {
	split := strings.Split(refreshToken, TokenDelimiter)
	if len(split) < 2 || !strings.HasPrefix(split[1], RefreshTokenPrefix) {
		return "", "", caos_errs.ThrowPreconditionFailed(nil, "OIDCS-JOI23", "Errors.OIDCSession.RefreshTokenInvalid")
	}
	// the oidc library requires that every token has the format of <tokenID>:<userID>
	// the V2 tokens don't use the userID anymore, so let's just remove it
	return split[0], strings.Split(split[1], oidcTokenSubjectDelimiter)[0], nil
}

func (c *Commands) newOIDCSessionUpdateEvents(ctx context.Context, oidcSessionID, refreshToken string) (*OIDCSessionEvents, error) {
	refreshTokenID, err := c.decryptRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	sessionWriteModel := NewOIDCSessionWriteModel(oidcSessionID, "")
	if err = c.eventstore.FilterToQueryReducer(ctx, sessionWriteModel); err != nil {
		return nil, err
	}
	if err = sessionWriteModel.CheckRefreshToken(refreshTokenID); err != nil {
		return nil, err
	}
	accessTokenLifetime, refreshTokenLifeTime, refreshTokenIdleLifetime, err := c.tokenTokenLifetimes(ctx)
	if err != nil {
		return nil, err
	}
	return &OIDCSessionEvents{
		eventstore:               c.eventstore,
		idGenerator:              c.idGenerator,
		encryptionAlg:            c.keyAlgorithm,
		oidcSessionWriteModel:    sessionWriteModel,
		accessTokenLifetime:      accessTokenLifetime,
		refreshTokenLifeTime:     refreshTokenLifeTime,
		refreshTokenIdleLifetime: refreshTokenIdleLifetime,
	}, nil
}

type OIDCSessionEvents struct {
	eventstore               *eventstore.Eventstore
	idGenerator              id.Generator
	encryptionAlg            crypto.EncryptionAlgorithm
	events                   []eventstore.Command
	oidcSessionWriteModel    *OIDCSessionWriteModel
	sessionWriteModel        *SessionWriteModel
	authRequestWriteModel    *AuthRequestWriteModel
	accessTokenLifetime      time.Duration
	refreshTokenLifeTime     time.Duration
	refreshTokenIdleLifetime time.Duration

	// accessTokenID is set by the command
	accessTokenID string

	// refreshToken is set by the command
	refreshToken string
}

func (c *OIDCSessionEvents) AddSession(ctx context.Context) {
	c.events = append(c.events, oidcsession.NewAddedEvent(
		ctx,
		c.oidcSessionWriteModel.aggregate,
		c.sessionWriteModel.UserID,
		c.sessionWriteModel.AggregateID,
		c.authRequestWriteModel.ClientID,
		c.authRequestWriteModel.Audience,
		c.authRequestWriteModel.Scope,
		c.sessionWriteModel.AuthMethodTypes(),
		c.sessionWriteModel.AuthenticationTime(),
	))
}

func (c *OIDCSessionEvents) SetAuthRequestSuccessful(ctx context.Context) {
	c.events = append(c.events, authrequest.NewSucceededEvent(ctx, c.authRequestWriteModel.aggregate))
}

func (c *OIDCSessionEvents) AddAccessToken(ctx context.Context, scope []string) error {
	accessTokenID, err := c.idGenerator.Next()
	if err != nil {
		return err
	}
	c.accessTokenID = AccessTokenPrefix + accessTokenID
	c.events = append(c.events, oidcsession.NewAccessTokenAddedEvent(ctx, c.oidcSessionWriteModel.aggregate, c.accessTokenID, scope, c.accessTokenLifetime))
	return nil
}

func (c *OIDCSessionEvents) AddRefreshToken(ctx context.Context) (err error) {
	var refreshTokenID string
	refreshTokenID, c.refreshToken, err = c.generateRefreshToken(c.sessionWriteModel.UserID)
	if err != nil {
		return err
	}
	c.events = append(c.events, oidcsession.NewRefreshTokenAddedEvent(ctx, c.oidcSessionWriteModel.aggregate, refreshTokenID, c.refreshTokenLifeTime, c.refreshTokenIdleLifetime))
	return nil
}

func (c *OIDCSessionEvents) RenewRefreshToken(ctx context.Context) (err error) {
	var refreshTokenID string
	refreshTokenID, c.refreshToken, err = c.generateRefreshToken(c.oidcSessionWriteModel.UserID)
	if err != nil {
		return err
	}
	c.events = append(c.events, oidcsession.NewRefreshTokenRenewedEvent(ctx, c.oidcSessionWriteModel.aggregate, refreshTokenID, c.refreshTokenIdleLifetime))
	return nil
}

func (c *OIDCSessionEvents) generateRefreshToken(userID string) (refreshTokenID, refreshToken string, err error) {
	refreshTokenID, err = c.idGenerator.Next()
	if err != nil {
		return "", "", err
	}
	refreshTokenID = RefreshTokenPrefix + refreshTokenID
	token, err := c.encryptionAlg.Encrypt([]byte(fmt.Sprintf(oidcTokenFormat, c.oidcSessionWriteModel.OIDCRefreshTokenID(refreshTokenID), userID)))
	if err != nil {
		return "", "", err
	}
	return refreshTokenID, base64.RawURLEncoding.EncodeToString(token), nil
}

func (c *OIDCSessionEvents) PushEvents(ctx context.Context) (accessTokenID string, refreshToken string, accessTokenExpiration time.Time, err error) {
	pushedEvents, err := c.eventstore.Push(ctx, c.events...)
	if err != nil {
		return "", "", time.Time{}, err
	}
	err = AppendAndReduce(c.oidcSessionWriteModel, pushedEvents...)
	if err != nil {
		return "", "", time.Time{}, err
	}
	// prefix the returned id with the oidcSessionID so that we can retrieve it later on
	// we need to use `-` as a delimiter because the OIDC library uses `:` and will check for a length of 2 parts
	return c.oidcSessionWriteModel.AggregateID + TokenDelimiter + c.accessTokenID, c.refreshToken, c.oidcSessionWriteModel.AccessTokenExpiration, nil
}

func (c *Commands) tokenTokenLifetimes(ctx context.Context) (accessTokenLifetime time.Duration, refreshTokenLifetime time.Duration, refreshTokenIdleLifetime time.Duration, err error) {
	oidcSettings := NewInstanceOIDCSettingsWriteModel(ctx)
	err = c.eventstore.FilterToQueryReducer(ctx, oidcSettings)
	if err != nil {
		return 0, 0, 0, err
	}
	accessTokenLifetime = c.defaultAccessTokenLifetime
	refreshTokenLifetime = c.defaultRefreshTokenLifetime
	refreshTokenIdleLifetime = c.defaultRefreshTokenIdleLifetime
	if oidcSettings.AccessTokenLifetime > 0 {
		accessTokenLifetime = oidcSettings.AccessTokenLifetime
	}
	if oidcSettings.RefreshTokenExpiration > 0 {
		refreshTokenLifetime = oidcSettings.RefreshTokenExpiration
	}
	if oidcSettings.RefreshTokenIdleExpiration > 0 {
		refreshTokenIdleLifetime = oidcSettings.RefreshTokenIdleExpiration
	}
	return accessTokenLifetime, refreshTokenLifetime, refreshTokenIdleLifetime, nil
}
