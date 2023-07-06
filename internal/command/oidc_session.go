package command

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/oidc/amr"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/repository/authrequest"
	"github.com/zitadel/zitadel/internal/repository/oidcsession"
)

// AddOIDCSessionAccessToken creates a new OIDC Session, creates an access token and returns its id and expiration
// if the underlying [AuthRequest] is a OIDC Auth Code Flow, it will set the code as exchange
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

// AddOIDCSessionRefreshAndAccessToken creates a new OIDC Session, creates an access token and refresh token
// it returns the access token id and expiration and the refresh token
// if the underlying [AuthRequest] is a OIDC Auth Code Flow, it will set the code as exchange
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

// ExchangeOIDCSessionRefreshAndAccessToken updates an existing OIDC Session, creates a new access and refresh token
// it returns the access token id and expiration and the new refresh token
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

// OIDCSessionByRefreshToken computes the current state of an existing OIDCSession by a refresh_token (to start a Refresh Token Grant)
// if either the session is not active, the token is invalid or expired (incl. idle expiration) an invalid refresh token error will be returned
func (c *Commands) OIDCSessionByRefreshToken(ctx context.Context, refreshToken string) (*OIDCSessionWriteModel, error) {
	split := strings.Split(refreshToken, ":")
	if len(split) != 2 {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "OIDCS-JOI23", "Errors.OIDCSession.RefreshTokenInvalid")
	}
	writeModel := NewOIDCSessionWriteModel(split[0], "")
	err := c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(err, "OIDCS-SAF31", "Errors.OIDCSession.RefreshTokenInvalid")
	}
	if err = writeModel.CheckRefreshToken(split[1]); err != nil {
		return nil, err
	}
	return writeModel, nil
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
		oidcSessionWriteModel:    NewOIDCSessionWriteModel(sessionID, authz.GetInstance(ctx).InstanceID()), //TODO: ro?
		sessionWriteModel:        sessionWriteModel,
		authRequestWriteModel:    authRequestWriteModel,
		accessTokenLifetime:      accessTokenLifetime,
		refreshTokenLifeTime:     refreshTokenLifeTime,
		refreshTokenIdleLifetime: refreshTokenIdleLifetime,
	}, nil
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
	split := strings.Split(decrypted, ":")
	if len(split) != 2 {
		return "", caos_errs.ThrowPreconditionFailed(nil, "OIDCS-Sj3lk", "Errors.OIDCSession.RefreshTokenInvalid")
	}
	return split[1], nil
}

func (c *Commands) newOIDCSessionUpdateEvents(ctx context.Context, oidcSessionID, refreshToken string) (*OIDCSessionEvents, error) {
	refreshTokenID, err := c.decryptRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	sessionWriteModel := NewOIDCSessionWriteModel(oidcSessionID, authz.GetInstance(ctx).InstanceID()) //TODO: ro?
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
		amr.List(c.sessionWriteModel),
		c.sessionWriteModel.AuthenticationTime(),
	))
}

func (c *OIDCSessionEvents) SetAuthRequestSuccessful(ctx context.Context) {
	c.events = append(c.events, authrequest.NewSucceededEvent(ctx, c.authRequestWriteModel.aggregate))
}

func (c *OIDCSessionEvents) AddAccessToken(ctx context.Context, scope []string) (err error) {
	c.accessTokenID, err = c.idGenerator.Next()
	if err != nil {
		return err
	}
	c.events = append(c.events, oidcsession.NewAccessTokenAddedEvent(ctx, c.oidcSessionWriteModel.aggregate, c.accessTokenID, scope, c.accessTokenLifetime))
	return nil
}

func (c *OIDCSessionEvents) AddRefreshToken(ctx context.Context) (err error) {
	var refreshTokenID string
	refreshTokenID, c.refreshToken, err = c.generateRefreshToken()
	if err != nil {
		return err
	}
	c.events = append(c.events, oidcsession.NewRefreshTokenAddedEvent(ctx, c.oidcSessionWriteModel.aggregate, refreshTokenID, c.refreshTokenLifeTime, c.refreshTokenIdleLifetime))
	return nil
}

func (c *OIDCSessionEvents) RenewRefreshToken(ctx context.Context) (err error) {
	var refreshTokenID string
	refreshTokenID, c.refreshToken, err = c.generateRefreshToken()
	if err != nil {
		return err
	}
	c.events = append(c.events, oidcsession.NewRefreshTokenRenewedEvent(ctx, c.oidcSessionWriteModel.aggregate, refreshTokenID, c.refreshTokenIdleLifetime))
	return nil
}

func (c *OIDCSessionEvents) generateRefreshToken() (refreshTokenID, refreshToken string, err error) {
	refreshTokenID, err = c.idGenerator.Next()
	if err != nil {
		return "", "", err
	}
	token, err := c.encryptionAlg.Encrypt([]byte(c.oidcSessionWriteModel.AggregateID + ":" + refreshTokenID))
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
	return c.oidcSessionWriteModel.AggregateID + "-" + c.accessTokenID, c.refreshToken, c.oidcSessionWriteModel.AccessTokenExpiration, nil
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
