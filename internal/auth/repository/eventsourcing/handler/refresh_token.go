package handler

import (
	"context"

	auth_view "github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	refreshTokenTable = "auth.refresh_tokens"
)

var _ handler.Projection = (*RefreshToken)(nil)

type RefreshToken struct {
	view *auth_view.View
}

func newRefreshToken(
	ctx context.Context,
	config handler.Config,
	view *auth_view.View,
) *handler.Handler {
	return handler.NewHandler(
		ctx,
		&config,
		&RefreshToken{
			view: view,
		},
	)
}

// Name implements [handler.Projection]
func (*RefreshToken) Name() string {
	return refreshTokenTable
}

// Reducers implements [handler.Projection]
func (t *RefreshToken) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  user.HumanRefreshTokenAddedType,
					Reduce: t.Reduce,
				},
				{
					Event:  user.HumanRefreshTokenRenewedType,
					Reduce: t.Reduce,
				},
				{
					Event:  user.HumanRefreshTokenRemovedType,
					Reduce: t.Reduce,
				},
				{
					Event:  user.UserLockedType,
					Reduce: t.Reduce,
				},
				{
					Event:  user.UserDeactivatedType,
					Reduce: t.Reduce,
				},
				{
					Event:  user.UserRemovedType,
					Reduce: t.Reduce,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: t.Reduce,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.OrgRemovedEventType,
					Reduce: t.Reduce,
				},
			},
		},
	}
}

func (t *RefreshToken) Reduce(event eventstore.Event) (_ *handler.Statement, err error) {
	switch event.Type() {
	case user.HumanRefreshTokenAddedType:
		e, ok := event.(*user.HumanRefreshTokenAddedEvent)
		if !ok {
			return nil, zerrors.ThrowInvalidArgumentf(nil, "MODEL-IoF6j", "reduce.wrong.event.type %s", user.HumanRefreshTokenAddedType)
		}
		return handler.NewCreateStatement(event,
			[]handler.Column{
				handler.NewCol("id", e.TokenID),
				handler.NewCol(userIDCol, e.Aggregate().ID),
				handler.NewCol(resourceOwnerCol, e.Aggregate().ResourceOwner),
				handler.NewCol(instanceIDCol, e.Aggregate().InstanceID),
				handler.NewCol(creationDateCol, event.CreatedAt()),
				handler.NewCol("amr", e.AuthMethodsReferences),
				handler.NewCol("auth_time", e.AuthTime),
				handler.NewCol("audience", e.Audience),
				handler.NewCol("client_id", e.ClientID),
				handler.NewCol("expiration", event.CreatedAt().Add(e.Expiration)),
				handler.NewCol("idle_expiration", event.CreatedAt().Add(e.IdleExpiration)),
				handler.NewCol("scopes", e.Scopes),
				handler.NewCol("token", e.TokenID),
				handler.NewCol(userAgentIDCol, e.UserAgentID),
				handler.NewCol("actor", view_model.TokenActor{TokenActor: e.Actor}),
			}), nil
	case user.HumanRefreshTokenRenewedType:
		e, ok := event.(*user.HumanRefreshTokenRenewedEvent)
		if !ok {
			return nil, zerrors.ThrowInvalidArgumentf(nil, "MODEL-AG43hq", "reduce.wrong.event.type %s", user.HumanRefreshTokenRenewedType)
		}
		//e := new(user.HumanRefreshTokenRenewedEvent)
		//if err := event.Unmarshal(e); err != nil {
		//	logging.WithError(err).Error("could not unmarshal event data")
		//	return nil, zerrors.ThrowInternal(nil, "MODEL-BHn75", "could not unmarshal data")
		//}
		return handler.NewUpdateStatement(event,
			[]handler.Column{
				handler.NewCol("idle_expiration", event.CreatedAt().Add(e.IdleExpiration)),
				handler.NewCol("token", e.RefreshToken),
			},
			[]handler.Condition{
				handler.NewCond("id", e.TokenID),
				handler.NewCond(instanceIDCol, e.Aggregate().InstanceID),
			},
		), nil
	case user.HumanRefreshTokenRemovedType:
		e, ok := event.(*user.HumanRefreshTokenRemovedEvent)
		if !ok {
			return nil, zerrors.ThrowInvalidArgumentf(nil, "MODEL-SFF3t", "reduce.wrong.event.type %s", user.HumanRefreshTokenRemovedType)
		}
		//e := new(user.HumanRefreshTokenRemovedEvent)
		//if err := event.Unmarshal(e); err != nil {
		//	logging.WithError(err).Error("could not unmarshal event data")
		//	return nil, zerrors.ThrowInternal(nil, "MODEL-Bz653", "could not unmarshal data")
		//}
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(instanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond("id", e.TokenID),
			},
		), nil
	case user.UserLockedType,
		user.UserDeactivatedType,
		user.UserRemovedType:
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(instanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond(userIDCol, event.Aggregate().InstanceID),
			},
		), nil
	case instance.InstanceRemovedEventType:
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(instanceIDCol, event.Aggregate().InstanceID),
			},
		), nil
	case org.OrgRemovedEventType:
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(instanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond(resourceOwnerCol, event.Aggregate().InstanceID),
			},
		), nil
	default:
		return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
			return nil
		}), nil
	}
}
