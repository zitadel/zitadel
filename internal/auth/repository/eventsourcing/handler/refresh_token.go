package handler

import (
	"context"

	"github.com/zitadel/logging"

	auth_view "github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
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
	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		switch event.Type() {
		case user.HumanRefreshTokenAddedType:
			token := new(view_model.RefreshTokenView)
			err := token.AppendEvent(event)
			if err != nil {
				return err
			}

			return t.view.PutRefreshToken(token)
		case user.HumanRefreshTokenRenewedType:
			e := new(user.HumanRefreshTokenRenewedEvent)
			if err := event.Unmarshal(e); err != nil {
				logging.WithError(err).Error("could not unmarshal event data")
				return caos_errs.ThrowInternal(nil, "MODEL-BHn75", "could not unmarshal data")
			}
			token, err := t.view.RefreshTokenByID(e.TokenID, event.Aggregate().InstanceID)
			if err != nil {
				return err
			}
			err = token.AppendEvent(event)
			if err != nil {
				return err
			}
			return t.view.PutRefreshToken(token)
		case user.HumanRefreshTokenRemovedType:
			e := new(user.HumanRefreshTokenRemovedEvent)
			if err := event.Unmarshal(e); err != nil {
				logging.WithError(err).Error("could not unmarshal event data")
				return caos_errs.ThrowInternal(nil, "MODEL-Bz653", "could not unmarshal data")
			}
			return t.view.DeleteRefreshToken(e.TokenID, event.Aggregate().InstanceID)
		case user.UserLockedType,
			user.UserDeactivatedType,
			user.UserRemovedType:

			return t.view.DeleteUserRefreshTokens(event.Aggregate().ID, event.Aggregate().InstanceID)
		case instance.InstanceRemovedEventType:

			return t.view.DeleteInstanceRefreshTokens(event.Aggregate().InstanceID)
		case org.OrgRemovedEventType:
			return t.view.DeleteOrgRefreshTokens(event)
		default:
			return nil
		}
	}), nil
}
