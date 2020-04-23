package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user_agent/repository/eventsourcing/model"
)

func UserAgentAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent) (*es_models.Aggregate, error) {
	if userAgent == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-sd823", "user agent must not be nil")
	}
	return aggCreator.NewAggregate(ctx, userAgent.AggregateID, model.UserAgentAggregate, model.UserAgentVersion, userAgent.Sequence)
}

func UserAgentCreateAggregate(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAgentAggregate(ctx, aggCreator, userAgent)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.UserAgentAdded, userAgent)
	}
}

func UserAgentRevocationAggregate(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAgentAggregate(ctx, aggCreator, userAgent)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.UserAgentRevoked, nil)
	}
}

func UserSessionAddedAggregate(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent, userSession *model.UserSession) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if userSession == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-32jkF", "user session must not be nil")
		}
		agg, err := UserAgentAggregate(ctx, aggCreator, userAgent)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.UserSessionAdded, userSession)
	}
}

func UserSessionTerminatedAggregate(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent, userSessionID string) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAgentAggregate(ctx, aggCreator, userAgent)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.UserSessionTerminated, &model.UserSessionID{UserSessionID: userSessionID})
	}
}

//func UsernameCheckSucceededAggregate(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent, userSessionID string) func(ctx context.Context) (*es_models.Aggregate, error) {
//	return passwordCheckAggregate(aggCreator, userAgent, userSessionID, model.UserNameCheckSucceeded)
//}
//func UsernameCheckFailedAggregate(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent, userSessionID string) func(ctx context.Context) (*es_models.Aggregate, error) {
//	return passwordCheckAggregate(aggCreator, userAgent, userSessionID, model.UserNameCheckFailed)
//}
//
//func usernameCheckAggregate(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent, userSessionID string, check models.EventType) func(ctx context.Context) (*es_models.Aggregate, error) {
//	return func(ctx context.Context) (*es_models.Aggregate, error) {
//		agg, err := UserAgentAggregate(ctx, aggCreator, userAgent)
//		if err != nil {
//			return nil, err
//		}
//		for _, session := range userAgent.UserSessions {
//			if session.SessionID == userSessionID {
//				agg.AppendEvent(check, userSessionID)
//				return agg, nil
//			}
//		}
//		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-89KNd", "user session does not exist")
//	}
//}

func PasswordCheckSucceededAggregate(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent, userSessionID string) func(ctx context.Context) (*es_models.Aggregate, error) {
	return passwordCheckAggregate(aggCreator, userAgent, userSessionID, model.PasswordCheckSucceeded)
}

func PasswordCheckFailedAggregate(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent, userSessionID string) func(ctx context.Context) (*es_models.Aggregate, error) {
	return passwordCheckAggregate(aggCreator, userAgent, userSessionID, model.PasswordCheckFailed)
}

func passwordCheckAggregate(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent, userSessionID string, check es_models.EventType) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAgentAggregate(ctx, aggCreator, userAgent)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(check, &model.UserSessionID{UserSessionID: userSessionID})
	}
}

func MfaCheckSucceededAggregate(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent, userSessionID string, mfaType int32) func(ctx context.Context) (*es_models.Aggregate, error) {
	return mfaCheckAggregate(aggCreator, userAgent, userSessionID, mfaType, model.MfaCheckSucceeded)
}

func MfaCheckFailedAggregate(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent, userSessionID string, mfaType int32) func(ctx context.Context) (*es_models.Aggregate, error) {
	return mfaCheckAggregate(aggCreator, userAgent, userSessionID, mfaType, model.MfaCheckFailed)
}

func mfaCheckAggregate(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent, userSessionID string, mfaType int32, check es_models.EventType) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAgentAggregate(ctx, aggCreator, userAgent)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(check, &model.MfaUserSession{UserSessionID: userSessionID, MfaType: mfaType})
	}
}

func ReAuthRequestAggregate(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent, userSessionID string) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAgentAggregate(ctx, aggCreator, userAgent)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.ReAuthRequested, &model.UserSessionID{UserSessionID: userSessionID})
	}
}

func AuthSessionAddedAggregate(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent, authSession *model.AuthSession) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAgentAggregate(ctx, aggCreator, userAgent)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.AuthSessionAdded, authSession)
	}
}

func AuthSessionSetUserSessionAggregate(aggCreator *es_models.AggregateCreator, userAgent *model.UserAgent, authSessionID, userSessionID string) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserAgentAggregate(ctx, aggCreator, userAgent)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserSessionSet, &model.SetUserSession{UserSessionID: userSessionID, AuthSessionID: authSessionID})
	}
}
