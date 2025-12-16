package command

import (
	"context"
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
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-0LxK6O31wH", "Errors.SAMLRequest.InvalidCode")
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
	postCommit, err := cmd.SetMilestones(ctx)
	if err != nil {
		return err
	}
	_, err = cmd.PushEvents(ctx)
	if err != nil {
		return err
	}
	postCommit(ctx)
	return err
}

func (c *Commands) newSAMLSessionAddEvents(ctx context.Context, userID, resourceOwner string, pending ...eventstore.Command) (*SAMLSessionEvents, error) {
	userStateModel, err := c.userStateWriteModel(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !userStateModel.UserState.IsEnabled() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "SAML-1768ZQpmcP", "Errors.User.NotActive")
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

func (c *SAMLSessionEvents) SetSAMLRequestFailed(ctx context.Context, samlRequestAggregate *eventstore.Aggregate, err domain.SAMLErrorReason) {
	c.events = append(c.events, samlrequest.NewFailedEvent(ctx, samlRequestAggregate, err))
}

func (c *SAMLSessionEvents) AddSAMLResponse(ctx context.Context, id string, lifetime time.Duration) error {
	c.samlResponseID = id
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
