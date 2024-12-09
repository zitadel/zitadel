package command

import (
	"context"
	"strings"
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/activity"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/repository/samlrequest"
	"github.com/zitadel/zitadel/internal/repository/samlsession"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SAMLSession struct {
	SessionID         string
	SAMLResponseID    string
	EntityID          string
	UserID            string
	Audience          []string
	Expiration        time.Time
	AuthMethods       []domain.UserAuthMethodType
	AuthTime          time.Time
	PreferredLanguage *language.Tag
	UserAgent         *domain.UserAgent
}

type SAMLRequestComplianceChecker func(context.Context, *SAMLRequestWriteModel) error

func (c *Commands) CreateSAMLSessionFromSAMLRequest(ctx context.Context, samlReqId string, complianceCheck SAMLRequestComplianceChecker, samlResponseID string, samlResponseLifetime time.Duration) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if samlReqId == "" {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-Sf3g2", "Errors.SAMLRequest.InvalidCode")
	}

	samlReqModel, err := c.getSAMLRequestWriteModel(ctx, samlReqId)
	if err != nil {
		return err
	}

	instanceID := authz.GetInstance(ctx).InstanceID()
	sessionModel := NewSessionWriteModel(samlReqModel.SessionID, instanceID)
	err = c.eventstore.FilterToQueryReducer(ctx, sessionModel)
	if err != nil {
		return err
	}
	if err = sessionModel.CheckIsActive(); err != nil {
		return err
	}

	cmd, err := c.newSAMLSessionAddEvents(ctx, sessionModel.UserID, sessionModel.UserResourceOwner)
	if err != nil {
		return err
	}
	if err = complianceCheck(ctx, samlReqModel); err != nil {
		return err
	}

	cmd.AddSession(ctx,
		sessionModel.UserID,
		sessionModel.UserResourceOwner,
		sessionModel.AggregateID,
		samlReqModel.Issuer,
		[]string{samlReqModel.Issuer},
		samlReqModel.AuthMethods,
		samlReqModel.AuthTime,
		sessionModel.PreferredLanguage,
		sessionModel.UserAgent,
	)

	if err = cmd.AddSAMLResponse(ctx, samlResponseID, samlResponseLifetime); err != nil {
		return err
	}
	cmd.SetSAMLRequestSuccessful(ctx, samlReqModel.aggregate)
	_, err = cmd.PushEvents(ctx)
	return err
}

func (c *Commands) CreateSAMLSession(ctx context.Context,
	userID,
	resourceOwner,
	clientID,
	backChannelLogoutURI string,
	scope,
	audience []string,
	authMethods []domain.UserAuthMethodType,
	authTime time.Time,
	nonce string,
	preferredLanguage *language.Tag,
	userAgent *domain.UserAgent,
	reason domain.TokenReason,
	actor *domain.TokenActor,
	needRefreshToken bool,
	sessionID string,
	responseType domain.OIDCResponseType,
) (session *OIDCSession, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	cmd, err := c.newOIDCSessionAddEvents(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if reason == domain.TokenReasonImpersonation {
		if err := c.checkPermission(ctx, "impersonation", resourceOwner, userID); err != nil {
			return nil, err
		}
		cmd.UserImpersonated(ctx, userID, resourceOwner, clientID, actor)
	}

	cmd.AddSession(ctx, userID, resourceOwner, sessionID, clientID, audience, scope, authMethods, authTime, nonce, preferredLanguage, userAgent)
	cmd.RegisterLogout(ctx, sessionID, userID, clientID, backChannelLogoutURI)
	if responseType != domain.OIDCResponseTypeIDToken {
		if err = cmd.AddAccessToken(ctx, scope, userID, resourceOwner, reason, actor); err != nil {
			return nil, err
		}
	}
	if needRefreshToken {
		if err = cmd.AddRefreshToken(ctx, userID); err != nil {
			return nil, err
		}
	}
	postCommit, err := cmd.SetMilestones(ctx, clientID, sessionID != "")
	if err != nil {
		return nil, err
	}
	if session, err = cmd.PushEvents(ctx); err != nil {
		return nil, err
	}
	postCommit(ctx)
	return session, nil
}

func samlSessionTokenIDsFromToken(token string) (oidcSessionID, refreshTokenID, accessTokenID string, err error) {
	split := strings.Split(token, TokenDelimiter)
	if len(split) != 2 {
		return "", "", "", zerrors.ThrowPreconditionFailed(nil, "OIDCS-S87kl", "Errors.OIDCSession.Token.Invalid")
	}
	if strings.HasPrefix(split[1], RefreshTokenPrefix) {
		return split[0], split[1], "", nil
	}
	if strings.HasPrefix(split[1], AccessTokenPrefix) {
		return split[0], "", split[1], nil
	}
	return "", "", "", zerrors.ThrowPreconditionFailed(nil, "OIDCS-S87kl", "Errors.OIDCSession.Token.Invalid")
}

func (c *Commands) newSAMLSessionAddEvents(ctx context.Context, userID, resourceOwner string, pending ...eventstore.Command) (*SAMLSessionEvents, error) {
	userStateModel, err := c.userStateWriteModel(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !userStateModel.UserState.IsEnabled() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "SAML-kj3g2", "Errors.User.NotActive")
	}
	sessionID, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	sessionID = IDPrefixV2 + sessionID
	return &SAMLSessionEvents{
		commands:              c,
		idGenerator:           c.idGenerator,
		encryptionAlg:         c.keyAlgorithm,
		events:                pending,
		samlSessionWriteModel: NewSAMLSessionWriteModel(sessionID, resourceOwner),
		userStateModel:        userStateModel,
	}, nil
}

type SAMLSessionEvents struct {
	commands              *Commands
	idGenerator           id.Generator
	encryptionAlg         crypto.EncryptionAlgorithm
	events                []eventstore.Command
	samlSessionWriteModel *SAMLSessionWriteModel
	userStateModel        *UserV2WriteModel

	// samlResponseID is set by the command
	samlResponseID string
}

func (c *SAMLSessionEvents) AddSession(
	ctx context.Context,
	userID,
	userResourceOwner,
	sessionID,
	entityID string,
	audience []string,
	authMethods []domain.UserAuthMethodType,
	authTime time.Time,
	preferredLanguage *language.Tag,
	userAgent *domain.UserAgent,
) {
	c.events = append(c.events, samlsession.NewAddedEvent(
		ctx,
		c.samlSessionWriteModel.aggregate,
		userID,
		userResourceOwner,
		sessionID,
		entityID,
		audience,
		authMethods,
		authTime,
		preferredLanguage,
		userAgent,
	))
}

func (c *SAMLSessionEvents) SetSAMLRequestSuccessful(ctx context.Context, samlRequestAggregate *eventstore.Aggregate) {
	c.events = append(c.events, samlrequest.NewSucceededEvent(ctx, samlRequestAggregate))
}

func (c *SAMLSessionEvents) SetSAMLRequestFailed(ctx context.Context, authRequestAggregate *eventstore.Aggregate, err error) {
	c.events = append(c.events, samlrequest.NewFailedEvent(ctx, authRequestAggregate, domain.SAMLErrorReasonFromError(err)))
}

func (c *SAMLSessionEvents) AddSAMLResponse(ctx context.Context, id string, lifetime time.Duration) error {
	c.events = append(c.events, samlsession.NewSAMLResponseAddedEvent(ctx, c.samlSessionWriteModel.aggregate, id, lifetime))
	return nil
}

func (c *SAMLSessionEvents) PushEvents(ctx context.Context) (*SAMLSession, error) {
	pushedEvents, err := c.commands.eventstore.Push(ctx, c.events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(c.samlSessionWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	session := &SAMLSession{
		SessionID:         c.samlSessionWriteModel.SessionID,
		EntityID:          c.samlSessionWriteModel.EntityID,
		UserID:            c.samlSessionWriteModel.UserID,
		Audience:          c.samlSessionWriteModel.Audience,
		Expiration:        c.samlSessionWriteModel.SAMLResponseExpiration,
		AuthMethods:       c.samlSessionWriteModel.AuthMethods,
		AuthTime:          c.samlSessionWriteModel.AuthTime,
		PreferredLanguage: c.samlSessionWriteModel.PreferredLanguage,
		UserAgent:         c.samlSessionWriteModel.UserAgent,
		SAMLResponseID:    c.samlSessionWriteModel.SAMLResponseID,
	}
	activity.Trigger(ctx, c.samlSessionWriteModel.UserResourceOwner, c.samlSessionWriteModel.UserID, activity.SAMLResponse, c.commands.eventstore.FilterToQueryReducer)
	return session, nil
}
