package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/cache/config"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	"github.com/caos/zitadel/internal/id"
	agent_model "github.com/caos/zitadel/internal/user_agent/model"
	"github.com/caos/zitadel/internal/user_agent/repository/eventsourcing/model"
)

type UserAgentEventstore struct {
	es_int.Eventstore
	agentCache  *UserAgentCache
	idGenerator id.Generator
}

type UserAgentConfig struct {
	Eventstore es_int.Eventstore
	Cache      *config.CacheConfig
}

type aggregateFunc func(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent) es_sdk.AggregateFunc
type pwCheckAggregateFn func(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent, userSessionID string) func(ctx context.Context) (*es_models.Aggregate, error)
type mfaCheckAggregateFn func(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent, userSessionID string, mfaType int32) func(ctx context.Context) (*es_models.Aggregate, error)

func StartUserAgent(conf UserAgentConfig) (*UserAgentEventstore, error) {
	agentCache, err := StartCache(conf.Cache)
	if err != nil {
		return nil, err
	}
	return &UserAgentEventstore{
		Eventstore: conf.Eventstore,
		agentCache: agentCache,
	}, nil
}

func (es *UserAgentEventstore) UserAgentByID(ctx context.Context, id string) (*agent_model.UserAgent, error) {
	userAgent, _ := es.agentCache.getUserAgent(id)

	query, err := UserAgentByIDQuery(userAgent.AggregateID, userAgent.Sequence)
	if err != nil {
		return nil, err
	}
	err = es_sdk.Filter(ctx, es.FilterEvents, userAgent.AppendEvents, query)
	if err != nil && !(caos_errs.IsNotFound(err) && userAgent.Sequence != 0) {
		return nil, err
	}
	es.agentCache.cacheUserAgent(userAgent)
	return model.UserAgentToModel(userAgent), nil
}

func (es *UserAgentEventstore) CreateUserAgent(ctx context.Context, userAgent *agent_model.UserAgent) (*agent_model.UserAgent, error) {
	if !userAgent.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-sdf32", "agent not valid")
	}
	id, err := es.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	userAgent.AggregateID = id
	userAgent.State = agent_model.UserAgentStateActive
	repoUserAgent := model.UserAgentFromModel(userAgent)

	createAggregate := UserAgentCreateAggregate(es.AggregateCreator(), repoUserAgent)
	err = es_sdk.Push(ctx, es.PushAggregates, repoUserAgent.AppendEvents, createAggregate)
	if err != nil {
		return nil, err
	}

	es.agentCache.cacheUserAgent(repoUserAgent)
	return model.UserAgentToModel(repoUserAgent), nil
}

func (es *UserAgentEventstore) RevokeUserAgent(ctx context.Context, id string) (*agent_model.UserAgent, error) {
	userAgent, err := es.UserAgentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !userAgent.IsActive() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-6a3fa", "user agent must be active")
	}
	repoUserAgent := model.UserAgentFromModel(userAgent)

	revocationAggregate := UserAgentRevocationAggregate(es.AggregateCreator(), repoUserAgent)
	err = es_sdk.Push(ctx, es.PushAggregates, repoUserAgent.AppendEvents, revocationAggregate)
	if err != nil {
		return nil, err
	}

	es.agentCache.cacheUserAgent(repoUserAgent)
	return model.UserAgentToModel(repoUserAgent), nil
}

func (es *UserAgentEventstore) AddUserSession(ctx context.Context, agentID, authSessionID string, userSession *agent_model.UserSession) (*agent_model.AuthSession, error) {
	repoUserAgent, err := es.userAgentByID(ctx, agentID)
	if err != nil {
		return nil, err
	}
	repoUserSession := model.UserSessionFromModel(userSession)

	addedAggregate := UserSessionAddedAggregate(es.AggregateCreator(), repoUserAgent, repoUserSession)
	err = es_sdk.Push(ctx, es.PushAggregates, repoUserAgent.AppendEvents, addedAggregate)
	if err != nil {
		return nil, err
	}

	es.agentCache.cacheUserAgent(repoUserAgent)
	return authSessionFromUserAgent(repoUserAgent, userSession.SessionID, authSessionID)
}

func (es *UserAgentEventstore) TerminateUserSession(ctx context.Context, userAgentID, userSessionID string) error {
	repoUserAgent, err := es.userAgentByID(ctx, userAgentID)
	if err != nil {
		return err
	}
	addedAggregate := UserSessionTerminatedAggregate(es.AggregateCreator(), repoUserAgent, userSessionID)
	err = es_sdk.Push(ctx, es.PushAggregates, repoUserAgent.AppendEvents, addedAggregate)
	if err != nil {
		return err
	}
	es.agentCache.cacheUserAgent(repoUserAgent)
	return nil
}

//
//func (es *UserAgentEventstore) PasswordCheckSucceeded(ctx context.Context, agentID, userSessionID, authSessionID string) (*agent_model.AuthSession, error) {
//	return es.passwordCheck(ctx, agentID, userSessionID, authSessionID, PasswordCheckSucceededAggregate)
//}
//
//func (es *UserAgentEventstore) PasswordCheckFailed(ctx context.Context, agentID, userSessionID, authSessionID string) (*agent_model.AuthSession, error) {
//	return es.passwordCheck(ctx, agentID, userSessionID, authSessionID, PasswordCheckFailedAggregate)
//}
//
//func (es *UserAgentEventstore) passwordCheck(ctx context.Context, agentID, userSessionID, authSessionID string, checkAgg pwCheckAggregateFn) (*agent_model.AuthSession, error) {
//	checkAggregate := func(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent) es_sdk.AggregateFunc {
//		return checkAgg(aggCreator, userAgent, userSessionID)
//	}
//	return es.authSessionEventIDs(ctx, agentID, userSessionID, authSessionID, checkAggregate)
//}
//
//func (es *UserAgentEventstore) MfaCheckSucceeded(ctx context.Context, agentID, userSessionID, authSessionID string, mfaType int32) (*agent_model.AuthSession, error) {
//	return es.mfaCheck(ctx, agentID, userSessionID, authSessionID, mfaType, MfaCheckSucceededAggregate)
//}
//
//func (es *UserAgentEventstore) MfaCheckFailed(ctx context.Context, agentID, userSessionID, authSessionID string, mfaType int32) (*agent_model.AuthSession, error) {
//	return es.mfaCheck(ctx, agentID, userSessionID, authSessionID, mfaType, MfaCheckFailedAggregate)
//}
//
//func (es *UserAgentEventstore) mfaCheck(ctx context.Context, agentID, userSessionID, authSessionID string, mfaType int32, checkAgg mfaCheckAggregateFn) (*agent_model.AuthSession, error) {
//	checkAggregate := func(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent) es_sdk.AggregateFunc {
//		return checkAgg(aggCreator, userAgent, userSessionID, mfaType)
//	}
//	return es.authSessionEventIDs(ctx, agentID, userSessionID, authSessionID, checkAggregate)
//}

func (es *UserAgentEventstore) ReAuthenticationRequested(ctx context.Context, userAgentID, userSessionID, authSessionID string) (*agent_model.AuthSession, error) {
	reauthAggregate := func(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent) es_sdk.AggregateFunc {
		return ReAuthRequestAggregate(aggCreator, userAgent, userSessionID)
	}
	return es.authSessionEventIDs(ctx, userAgentID, userSessionID, authSessionID, reauthAggregate)
}

func (es *UserAgentEventstore) AuthSessionByIDs(ctx context.Context, agentID, userSessionID, authSessionID string) (*agent_model.AuthSession, error) {
	agent, err := es.userAgentByID(ctx, agentID)
	if err != nil {
		return nil, err
	}
	return authSessionFromUserAgent(agent, userSessionID, authSessionID)
}

func (es *UserAgentEventstore) AuthSessionAdded(ctx context.Context, authSession *agent_model.AuthSession) (*agent_model.AuthSession, error) {
	repoUserAgent, err := es.userAgentByID(ctx, authSession.AggregateID)
	if err != nil {
		return nil, err
	}
	repoAuthSession := model.AuthSessionFromModel(authSession)

	addedAggregate := AuthSessionAddedAggregate(es.AggregateCreator(), repoUserAgent, repoAuthSession)
	err = es_sdk.Push(ctx, es.PushAggregates, repoUserAgent.AppendEvents, addedAggregate)
	if err != nil {
		return nil, err
	}

	es.agentCache.cacheUserAgent(repoUserAgent)
	return model.AuthSessionToModel(repoAuthSession), nil
}

func (es *UserAgentEventstore) AuthSessionSetUserSession(ctx context.Context, userAgentID, userSessionID, authSessionID string) (*agent_model.AuthSession, error) {
	setAggregate := func(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent) es_sdk.AggregateFunc {
		return AuthSessionSetUserSessionAggregate(aggCreator, userAgent, userSessionID, authSessionID)
	}
	return es.authSessionEventIDs(ctx, userAgentID, userSessionID, authSessionID, setAggregate)
}

func (es *UserAgentEventstore) AuthSessionSetNewUserSession(ctx context.Context, userAgentID, userSessionID, authSessionID string) (*agent_model.AuthSession, error) {
	setAggregate := func(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent) es_sdk.AggregateFunc {
		return AuthSessionSetUserSessionAggregate(aggCreator, userAgent, userSessionID, authSessionID)
	}
	return es.authSessionEventIDs(ctx, userAgentID, userSessionID, authSessionID, setAggregate)
}

func (es *UserAgentEventstore) CreateToken(ctx context.Context, userAgentID, userSessionID, authSessionID string) (*agent_model.Token, error) {

	repoUserAgent, err := es.userAgentByID(ctx, userAgentID)
	if err != nil {
		return nil, err
	}
	_, userSession := model.GetUserSession(repoUserAgent.UserSessions, userSessionID)
	if userSession == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9sdml", "user session not found")
	}
	if _, authSession := model.GetAuthSession(userSession.AuthSessions, authSessionID); authSession == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-s3dGk", "auth session not found")
	}
	err = es_sdk.Push(ctx, es.PushAggregates, repoUserAgent.AppendEvents, TokenAddedAggregate(es.AggregateCreator(), repoUserAgent, token))
	if err != nil {
		return nil, err
	}
	es.agentCache.cacheUserAgent(repoUserAgent)

}

//
//func (es *UserAgentEventstore) userSessionEvent(ctx context.Context, userSession *model.UserSession, aggregateFunc aggregateFunc) (*agent_model.UserSession, error) {
//	return es.userSessionEventIDs(ctx, userSession.AggregateID, userSession.SessionID, aggregateFunc)
//}
func (es *UserAgentEventstore) authSessionEvent(ctx context.Context, authSession *model.AuthSession, aggregateFunc aggregateFunc) (*agent_model.AuthSession, error) {
	return es.authSessionEventIDs(ctx, authSession.AggregateID, authSession.UserSession.SessionID, authSession.SessionID, aggregateFunc)
}

//func (es *UserAgentEventstore) userSessionEventIDs(ctx context.Context, agentID, sessionID string, aggregateFunc aggregateFunc) (*agent_model.UserSession, error) {
//	repoUserAgent, err := es.userAgentByID(ctx, agentID)
//	if err != nil {
//		return nil, err
//	}
//	if _, s := model.GetUserSession(repoUserAgent.UserSessions, sessionID); s == nil {
//		return  nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9sdml", "user session not found")
//	}
//	err = es_sdk.Push(ctx, es.PushAggregates, repoUserAgent.AppendEvents, aggregateFunc(es.AggregateCreator(), repoUserAgent))
//	if err != nil {
//		return nil, err
//	}
//	es.agentCache.cacheUserAgent(repoUserAgent)
//
//	if _, s := model.GetUserSession(repoUserAgent.UserSessions, sessionID); s != nil {
//		return model.UserSessionToModel(s), nil
//	}
//	return nil, caos_errs.ThrowInternal(nil, "EVENT-sk3t5", "Could not find user session in list")
//}

func (es *UserAgentEventstore) authSessionEventIDs(ctx context.Context, agentID, userSessionID, authSessionID string, aggregateFunc aggregateFunc) (*agent_model.AuthSession, error) {
	repoUserAgent, err := es.userAgentByID(ctx, agentID)
	if err != nil {
		return nil, err
	}
	_, userSession := model.GetUserSession(repoUserAgent.UserSessions, userSessionID)
	if userSession == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9sdml", "user session not found")
	}
	if _, authSession := model.GetAuthSession(userSession.AuthSessions, authSessionID); authSession == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-s3dGk", "auth session not found")
	}
	err = es_sdk.Push(ctx, es.PushAggregates, repoUserAgent.AppendEvents, aggregateFunc(es.AggregateCreator(), repoUserAgent))
	if err != nil {
		return nil, err
	}
	es.agentCache.cacheUserAgent(repoUserAgent)
	return authSessionFromUserAgent(repoUserAgent, userSessionID, authSessionID)
}

func (es *UserAgentEventstore) userAgentByID(ctx context.Context, id string) (*model.UserAgent, error) {
	userAgent, err := es.UserAgentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !userAgent.IsActive() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-s3hbD", "user agent must be active")
	}
	return model.UserAgentFromModel(userAgent), nil
}

func authSessionFromUserAgent(repoUserAgent *model.UserAgent, userSessionID, authSessionID string) (*agent_model.AuthSession, error) {
	if _, userSession := model.GetUserSession(repoUserAgent.UserSessions, userSessionID); userSession != nil {
		if _, authSession := model.GetAuthSession(userSession.AuthSessions, authSessionID); authSession != nil {
			return model.AuthSessionToModel(authSession), nil
		}
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-sk3t5", "Could not find user session in list")
}
