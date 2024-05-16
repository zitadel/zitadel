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

	refreshTokenIDCol            = "id"
	refreshTokenUserAgentIDCol   = "user_agent_id"
	refreshTokenUserIDCol        = "user_id"
	refreshTokenInstanceIDCol    = "instance_id"
	refreshTokenCreationDateCol  = "creation_date"
	refreshTokenChangeDateCol    = "change_date"
	refreshTokenResourceOwnerCol = "resource_owner"
	refreshTokenSequenceOwnerCol = "sequence"
	refreshTokenActorCol         = "actor"
	refreshTokenAMR              = "amr"
	refreshTokenAuthTime         = "auth_time"
	refreshTokenAudience         = "audience"
	refreshTokenClientID         = "client_id"
	refreshTokenExpiration       = "expiration"
	refreshTokenIdleExpiration   = "idle_expiration"
	refreshTokenScopes           = "scopes"
	refreshTokenToken            = "token"
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
		columns := []handler.Column{
			handler.NewCol(refreshTokenClientID, e.ClientID),
			handler.NewCol(refreshTokenUserAgentIDCol, e.UserAgentID),
			handler.NewCol(refreshTokenUserIDCol, e.Aggregate().ID),
			handler.NewCol(refreshTokenInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCol(refreshTokenIDCol, e.TokenID),
			handler.NewCol(refreshTokenResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCol(refreshTokenCreationDateCol, event.CreatedAt()),
			handler.NewCol(refreshTokenChangeDateCol, event.CreatedAt()),
			handler.NewCol(refreshTokenSequenceOwnerCol, event.Sequence()),
			handler.NewCol(refreshTokenAMR, e.AuthMethodsReferences),
			handler.NewCol(refreshTokenAuthTime, e.AuthTime),
			handler.NewCol(refreshTokenAudience, e.Audience),
			handler.NewCol(refreshTokenExpiration, event.CreatedAt().Add(e.Expiration)),
			handler.NewCol(refreshTokenIdleExpiration, event.CreatedAt().Add(e.IdleExpiration)),
			handler.NewCol(refreshTokenScopes, e.Scopes),
			handler.NewCol(refreshTokenToken, e.TokenID),
			handler.NewCol(refreshTokenActorCol, view_model.TokenActor{TokenActor: e.Actor}),
		}
		return handler.NewUpsertStatement(event, columns[0:3], columns), nil
	case user.HumanRefreshTokenRenewedType:
		e, ok := event.(*user.HumanRefreshTokenRenewedEvent)
		if !ok {
			return nil, zerrors.ThrowInvalidArgumentf(nil, "MODEL-AG43hq", "reduce.wrong.event.type %s", user.HumanRefreshTokenRenewedType)
		}
		return handler.NewUpdateStatement(event,
			[]handler.Column{
				handler.NewCol(refreshTokenIdleExpiration, event.CreatedAt().Add(e.IdleExpiration)),
				handler.NewCol(refreshTokenToken, e.RefreshToken),
				handler.NewCol(refreshTokenChangeDateCol, e.CreatedAt()),
				handler.NewCol(refreshTokenSequenceOwnerCol, event.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(refreshTokenIDCol, e.TokenID),
				handler.NewCond(refreshTokenInstanceIDCol, e.Aggregate().InstanceID),
			},
		), nil
	case user.HumanRefreshTokenRemovedType:
		e, ok := event.(*user.HumanRefreshTokenRemovedEvent)
		if !ok {
			return nil, zerrors.ThrowInvalidArgumentf(nil, "MODEL-SFF3t", "reduce.wrong.event.type %s", user.HumanRefreshTokenRemovedType)
		}
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(refreshTokenInstanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond(refreshTokenIDCol, e.TokenID),
			},
		), nil
	case user.UserLockedType,
		user.UserDeactivatedType,
		user.UserRemovedType:
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(refreshTokenInstanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond(refreshTokenUserIDCol, event.Aggregate().ID),
			},
		), nil
	case instance.InstanceRemovedEventType:
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(refreshTokenInstanceIDCol, event.Aggregate().InstanceID),
			},
		), nil
	case org.OrgRemovedEventType:
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(refreshTokenInstanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond(refreshTokenResourceOwnerCol, event.Aggregate().ResourceOwner),
			},
		), nil
	default:
		return handler.NewNoOpStatement(event), nil
	}
}
