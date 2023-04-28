package handler

import (
	"context"
	"encoding/json"

	"github.com/zitadel/logging"

	auth_view "github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	handler2 "github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

const (
	refreshTokenTable = "auth.refresh_tokens"
)

var _ handler2.Projection = (*RefreshToken)(nil)

type RefreshToken struct {
	view *auth_view.View
}

func newRefreshToken(
	ctx context.Context,
	config handler2.Config,
	view *auth_view.View,
) *handler2.Handler {
	return handler2.NewHandler(
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
func (t *RefreshToken) Reducers() []handler2.AggregateReducer {
	return []handler2.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventRedusers: []handler2.EventReducer{
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
		// TODO: which events should be processed?
		{
			Aggregate:     project.AggregateType,
			EventRedusers: []handler2.EventReducer{},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler2.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: t.Reduce,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler2.EventReducer{
				{
					Event:  org.OrgRemovedEventType,
					Reduce: t.Reduce,
				},
			},
		},
	}
}

func (t *RefreshToken) Reduce(event eventstore.Event) (_ *handler2.Statement, err error) {
	switch eventstore.EventType(event.Type()) {
	case user.HumanRefreshTokenAddedType:
		token := new(view_model.RefreshTokenView)
		err := token.AppendEvent(event)
		if err != nil {
			return nil, err
		}
		return handler2.NewStatement(event,
			func(ex handler2.Executer, projectionName string) error {
				return t.view.PutRefreshToken(token)
			}), nil
	case user.HumanRefreshTokenRenewedType:
		e := new(user.HumanRefreshTokenRenewedEvent)
		if err := json.Unmarshal(event.DataAsBytes(), e); err != nil {
			logging.WithError(err).Error("could not unmarshal event data")
			return nil, caos_errs.ThrowInternal(nil, "MODEL-BHn75", "could not unmarshal data")
		}
		token, err := t.view.RefreshTokenByID(e.TokenID, event.Aggregate().InstanceID)
		if err != nil {
			return nil, err
		}
		err = token.AppendEvent(event)
		if err != nil {
			return nil, err
		}
		return handler2.NewStatement(event,
			func(ex handler2.Executer, projectionName string) error {
				return t.view.PutRefreshToken(token)
			}), nil
	case user.HumanRefreshTokenRemovedType:
		e := new(user.HumanRefreshTokenRemovedEvent)
		if err := json.Unmarshal(event.DataAsBytes(), e); err != nil {
			logging.WithError(err).Error("could not unmarshal event data")
			return nil, caos_errs.ThrowInternal(nil, "MODEL-Bz653", "could not unmarshal data")
		}
		return handler2.NewStatement(event,
			func(ex handler2.Executer, projectionName string) error {
				return t.view.DeleteRefreshToken(e.TokenID, event.Aggregate().InstanceID)
			}), nil
	case user.UserLockedType,
		user.UserDeactivatedType,
		user.UserRemovedType:
		return handler2.NewStatement(event,
			func(ex handler2.Executer, projectionName string) error {
				return t.view.DeleteUserRefreshTokens(event.Aggregate().ID, event.Aggregate().InstanceID)
			}), nil
	case instance.InstanceRemovedEventType:
		return handler2.NewStatement(event,
			func(ex handler2.Executer, projectionName string) error {
				return t.view.DeleteInstanceRefreshTokens(event.Aggregate().InstanceID)
			}), nil
	case org.OrgRemovedEventType:
		return handler2.NewStatement(event,
			func(ex handler2.Executer, projectionName string) error {
				return t.view.DeleteOrgRefreshTokens(event)
			}), nil
	default:
		return handler2.NewNoOpStatement(event), nil
	}
}
