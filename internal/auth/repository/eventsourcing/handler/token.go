package handler

import (
	"context"

	"github.com/zitadel/logging"

	auth_view "github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	proj_model "github.com/zitadel/zitadel/internal/project/model"
	project_es_model "github.com/zitadel/zitadel/internal/project/repository/eventsourcing/model"
	proj_view "github.com/zitadel/zitadel/internal/project/repository/view"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	tokenTable = "auth.tokens"
)

var _ handler.Projection = (*Token)(nil)

type Token struct {
	view *auth_view.View
	es   handler.EventStore
}

func newToken(
	ctx context.Context,
	config handler.Config,
	view *auth_view.View,
) *handler.Handler {
	return handler.NewHandler(
		ctx,
		&config,
		&Token{
			view: view,
			es:   config.Eventstore,
		},
	)
}

// Name implements [handler.Projection]
func (*Token) Name() string {
	return tokenTable
}

// Reducers implements [handler.Projection]
func (t *Token) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  user.PersonalAccessTokenAddedType,
					Reduce: t.Reduce,
				},
				{
					Event:  user.UserTokenAddedType,
					Reduce: t.Reduce,
				},
				{
					Event:  user.UserV1ProfileChangedType,
					Reduce: t.Reduce,
				},
				{
					Event:  user.HumanProfileChangedType,
					Reduce: t.Reduce,
				},
				{
					Event:  user.UserV1SignedOutType,
					Reduce: t.Reduce,
				},
				{
					Event:  user.HumanSignedOutType,
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
				{
					Event:  user.UserTokenRemovedType,
					Reduce: t.Reduce,
				},
				{
					Event:  user.PersonalAccessTokenRemovedType,
					Reduce: t.Reduce,
				},
				{
					Event:  user.HumanRefreshTokenRemovedType,
					Reduce: t.Reduce,
				},
			},
		},
		{
			Aggregate: project.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  project.ApplicationDeactivatedType,
					Reduce: t.Reduce,
				},
				{
					Event:  project.ApplicationRemovedType,
					Reduce: t.Reduce,
				},
				{
					Event:  project.ProjectDeactivatedType,
					Reduce: t.Reduce,
				},
				{
					Event:  project.ProjectRemovedType,
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
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: t.Reduce,
				},
			},
		},
	}
}

func (t *Token) Reduce(event eventstore.Event) (_ *handler.Statement, err error) { //nolint:gocognit
	switch event.Type() {
	case user.UserTokenAddedType,
		user.PersonalAccessTokenAddedType:
		token := new(view_model.TokenView)
		err := token.AppendEvent(event)
		if err != nil {
			return nil, err
		}
		return handler.NewCreateStatement(
			event,
			append(token.PKColumns(), token.Changes()...),
		), nil
	case user.UserV1ProfileChangedType,
		user.HumanProfileChangedType:
		user := new(view_model.UserView)
		err := user.AppendEvent(event)
		if err != nil {
			return nil, err
		}
		if !user.HumanView.Value().PreferredLanguage.DidChange() {
			return handler.NewNoOpStatement(event), nil
		}

		return handler.NewUpdateStatement(
			event,
			[]handler.Column{
				handler.NewCol(view_model.TokenKeyPreferredLanguage, user.HumanView.Value().PreferredLanguage.Value()),
			},
			[]handler.Condition{
				handler.NewCond(view_model.TokenKeyInstanceID, event.Aggregate().InstanceID),
				handler.NewCond(view_model.TokenKeyUserID, event.Aggregate().ID),
				handler.NewCond(view_model.TokenKeyResourceOwner, event.Aggregate().ResourceOwner),
			},
		), nil
	case user.UserV1SignedOutType,
		user.HumanSignedOutType:
		id, err := agentIDFromSession(event)
		if err != nil {
			return nil, err
		}

		return handler.NewDeleteStatement(
			event,
			[]handler.Condition{
				handler.NewCond(view_model.TokenKeyUserAgentID, id),
				handler.NewCond(view_model.TokenKeyUserID, event.Aggregate().ID),
				handler.NewCond(view_model.TokenKeyInstanceID, event.Aggregate().InstanceID),
			},
		), nil

	case user.UserLockedType,
		user.UserDeactivatedType,
		user.UserRemovedType:

		return handler.NewDeleteStatement(
			event,
			[]handler.Condition{
				handler.NewCond(view_model.TokenKeyUserID, event.Aggregate().ID),
				handler.NewCond(view_model.TokenKeyInstanceID, event.Aggregate().InstanceID),
			},
		), nil
	case user.UserTokenRemovedType,
		user.PersonalAccessTokenRemovedType:
		id, err := tokenIDFromRemovedEvent(event)
		if err != nil {
			return nil, err
		}
		return handler.NewDeleteStatement(
			event,
			[]handler.Condition{
				handler.NewCond(view_model.TokenKeyTokenID, id),
				handler.NewCond(view_model.TokenKeyInstanceID, event.Aggregate().InstanceID),
			},
		), nil
	case user.HumanRefreshTokenRemovedType:
		id, err := refreshTokenIDFromRemovedEvent(event)
		if err != nil {
			return nil, err
		}

		return handler.NewDeleteStatement(
			event,
			[]handler.Condition{
				handler.NewCond(view_model.TokenKeyRefreshTokenID, id),
				handler.NewCond(view_model.TokenKeyInstanceID, event.Aggregate().InstanceID),
			},
		), nil
	case project.ApplicationDeactivatedType,
		project.ApplicationRemovedType:
		application, err := applicationFromSession(event)
		if err != nil {
			return nil, err
		}

		return handler.NewDeleteStatement(
			event,
			[]handler.Condition{
				handler.NewCond(view_model.TokenKeyApplicationID, application.AppID),
				handler.NewCond(view_model.TokenKeyInstanceID, event.Aggregate().InstanceID),
			},
		), nil
	case project.ProjectDeactivatedType,
		project.ProjectRemovedType:
		project, err := t.getProjectByID(context.Background(), event.Aggregate().ID, event.Aggregate().InstanceID)
		if err != nil {
			return nil, err
		}
		applicationIDs := make([]string, 0, len(project.Applications))
		for _, app := range project.Applications {
			if app.OIDCConfig != nil && app.OIDCConfig.ClientID != "" {
				applicationIDs = append(applicationIDs, app.OIDCConfig.ClientID)
			}
		}

		deletes := make([]func(eventstore.Event) handler.Exec, len(applicationIDs))
		for i, applicationID := range applicationIDs {
			deletes[i] = handler.AddDeleteStatement([]handler.Condition{
				handler.NewCond(view_model.TokenKeyApplicationID, applicationID),
				handler.NewCond(view_model.TokenKeyInstanceID, event.Aggregate().InstanceID),
			})
		}

		return handler.NewMultiStatement(event, deletes...), nil
	case instance.InstanceRemovedEventType:
		return handler.NewDeleteStatement(
			event,
			[]handler.Condition{
				handler.NewCond(view_model.TokenKeyInstanceID, event.Aggregate().InstanceID),
			},
		), nil
	case org.OrgRemovedEventType:
		// deletes all tokens including PATs, which is expected for now
		// if there is an undo of the org deletion in the future,
		// we will need to have a look on how to handle the deleted PATs
		return handler.NewDeleteStatement(
			event,
			[]handler.Condition{
				handler.NewCond(view_model.TokenKeyResourceOwner, event.Aggregate().ID),
			},
		), nil
	default:
		return handler.NewNoOpStatement(event), nil
	}
}

func agentIDFromSession(event eventstore.Event) (string, error) {
	session := make(map[string]interface{})
	if err := event.Unmarshal(&session); err != nil {
		logging.WithError(err).Error("could not unmarshal event data")
		return "", zerrors.ThrowInternal(nil, "MODEL-sd325", "could not unmarshal data")
	}
	agentID, _ := session["userAgentID"].(string)
	return agentID, nil
}

func applicationFromSession(event eventstore.Event) (*project_es_model.Application, error) {
	application := new(project_es_model.Application)
	if err := event.Unmarshal(application); err != nil {
		logging.WithError(err).Error("could not unmarshal event data")
		return nil, zerrors.ThrowInternal(nil, "MODEL-Hrw1q", "could not unmarshal data")
	}
	return application, nil
}

func tokenIDFromRemovedEvent(event eventstore.Event) (string, error) {
	removed := make(map[string]interface{})
	if err := event.Unmarshal(&removed); err != nil {
		logging.WithError(err).Error("could not unmarshal event data")
		return "", zerrors.ThrowInternal(nil, "MODEL-Sff32", "could not unmarshal data")
	}
	return removed["tokenId"].(string), nil
}

func refreshTokenIDFromRemovedEvent(event eventstore.Event) (string, error) {
	removed := make(map[string]interface{})
	if err := event.Unmarshal(&removed); err != nil {
		logging.WithError(err).Error("could not unmarshal event data")
		return "", zerrors.ThrowInternal(nil, "MODEL-Dfb3w", "could not unmarshal data")
	}
	return removed["tokenId"].(string), nil
}

func (t *Token) getProjectByID(ctx context.Context, projID, instanceID string) (*proj_model.Project, error) {
	query, err := proj_view.ProjectByIDQuery(projID, instanceID, 0)
	if err != nil {
		return nil, err
	}
	esProject := &project_es_model.Project{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: projID,
		},
	}
	events, err := t.es.Filter(ctx, query)
	if err != nil {
		return nil, err
	}
	if err = esProject.AppendEvents(events...); err != nil {
		return nil, err
	}

	if esProject.Sequence == 0 {
		return nil, zerrors.ThrowNotFound(nil, "EVENT-Dsdw2", "Errors.Project.NotFound")
	}
	return project_es_model.ProjectToModel(esProject), nil
}
