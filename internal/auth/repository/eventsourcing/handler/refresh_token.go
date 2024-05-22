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
	// in case anything needs to be change here check if appendEvent function needs the change as well
	switch event.Type() {
	case user.HumanRefreshTokenAddedType:
		e, ok := event.(*user.HumanRefreshTokenAddedEvent)
		if !ok {
			return nil, zerrors.ThrowInvalidArgumentf(nil, "MODEL-IoF6j", "reduce.wrong.event.type %s", user.HumanRefreshTokenAddedType)
		}
		columns := []handler.Column{
			handler.NewCol(view_model.RefreshTokenKeyClientID, e.ClientID),
			handler.NewCol(view_model.RefreshTokenKeyUserAgentID, e.UserAgentID),
			handler.NewCol(view_model.RefreshTokenKeyUserID, e.Aggregate().ID),
			handler.NewCol(view_model.RefreshTokenKeyInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(view_model.RefreshTokenKeyTokenID, e.TokenID),
			handler.NewCol(view_model.RefreshTokenKeyResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(view_model.RefreshTokenKeyCreationDate, event.CreatedAt()),
			handler.NewCol(view_model.RefreshTokenKeyChangeDate, event.CreatedAt()),
			handler.NewCol(view_model.RefreshTokenKeySequence, event.Sequence()),
			handler.NewCol(view_model.RefreshTokenKeyAMR, e.AuthMethodsReferences),
			handler.NewCol(view_model.RefreshTokenKeyAuthTime, e.AuthTime),
			handler.NewCol(view_model.RefreshTokenKeyAudience, e.Audience),
			handler.NewCol(view_model.RefreshTokenKeyExpiration, event.CreatedAt().Add(e.Expiration)),
			handler.NewCol(view_model.RefreshTokenKeyIdleExpiration, event.CreatedAt().Add(e.IdleExpiration)),
			handler.NewCol(view_model.RefreshTokenKeyScopes, e.Scopes),
			handler.NewCol(view_model.RefreshTokenKeyToken, e.TokenID),
			handler.NewCol(view_model.RefreshTokenKeyActor, view_model.TokenActor{TokenActor: e.Actor}),
		}
		return handler.NewUpsertStatement(event, columns[0:3], columns), nil
	case user.HumanRefreshTokenRenewedType:
		e, ok := event.(*user.HumanRefreshTokenRenewedEvent)
		if !ok {
			return nil, zerrors.ThrowInvalidArgumentf(nil, "MODEL-AG43hq", "reduce.wrong.event.type %s", user.HumanRefreshTokenRenewedType)
		}
		return handler.NewUpdateStatement(event,
			[]handler.Column{
				handler.NewCol(view_model.RefreshTokenKeyIdleExpiration, event.CreatedAt().Add(e.IdleExpiration)),
				handler.NewCol(view_model.RefreshTokenKeyToken, e.RefreshToken),
				handler.NewCol(view_model.RefreshTokenKeyChangeDate, e.CreatedAt()),
				handler.NewCol(view_model.RefreshTokenKeySequence, event.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(view_model.RefreshTokenKeyTokenID, e.TokenID),
				handler.NewCond(view_model.RefreshTokenKeyInstanceID, e.Aggregate().InstanceID),
			},
		), nil
	case user.HumanRefreshTokenRemovedType:
		e, ok := event.(*user.HumanRefreshTokenRemovedEvent)
		if !ok {
			return nil, zerrors.ThrowInvalidArgumentf(nil, "MODEL-SFF3t", "reduce.wrong.event.type %s", user.HumanRefreshTokenRemovedType)
		}
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(view_model.RefreshTokenKeyInstanceID, event.Aggregate().InstanceID),
				handler.NewCond(view_model.RefreshTokenKeyTokenID, e.TokenID),
			},
		), nil
	case user.UserLockedType,
		user.UserDeactivatedType,
		user.UserRemovedType:
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(view_model.RefreshTokenKeyInstanceID, event.Aggregate().InstanceID),
				handler.NewCond(view_model.RefreshTokenKeyUserID, event.Aggregate().ID),
			},
		), nil
	case instance.InstanceRemovedEventType:
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(view_model.RefreshTokenKeyInstanceID, event.Aggregate().InstanceID),
			},
		), nil
	case org.OrgRemovedEventType:
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(view_model.RefreshTokenKeyInstanceID, event.Aggregate().InstanceID),
				handler.NewCond(view_model.RefreshTokenKeyResourceOwner, event.Aggregate().ResourceOwner),
			},
		), nil
	default:
		return handler.NewNoOpStatement(event), nil
	}
}
