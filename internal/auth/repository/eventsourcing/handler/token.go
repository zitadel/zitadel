package handler

import (
	"context"

	"github.com/muhlemmer/gu"
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
	case user.UserTokenAddedType:
		e, ok := event.(*user.UserTokenAddedEvent)
		if !ok {
			return nil, zerrors.ThrowInvalidArgumentf(nil, "MODEL-W4tnq", "reduce.wrong.event.type %s", user.UserTokenAddedType)
		}
		return handler.NewCreateStatement(event,
			[]handler.Column{
				handler.NewCol(instanceIDCol, event.Aggregate().InstanceID),
				handler.NewCol(userIDCol, event.Aggregate().ID),
				handler.NewCol(resourceOwnerCol, event.Aggregate().ResourceOwner),
				handler.NewCol("id", e.TokenID),
				handler.NewCol(creationDateCol, event.CreatedAt()),
				handler.NewCol(changeDateCol, event.CreatedAt()),
				handler.NewCol("application_id", e.ApplicationID),
				handler.NewCol(userAgentIDCol, e.UserAgentID),
				handler.NewCol("audience", e.Audience),
				handler.NewCol("scopes", e.Scopes),
				handler.NewCol("expiration", e.Expiration),
				handler.NewCol("preferred_language", e.PreferredLanguage),
				handler.NewCol("refresh_token_id", e.RefreshTokenID),
				handler.NewCol("actor", view_model.TokenActor{TokenActor: e.Actor}),
				handler.NewCol("is_pat", false),
			},
		), nil
	case user.PersonalAccessTokenAddedType:
		e, ok := event.(*user.PersonalAccessTokenAddedEvent)
		if !ok {
			return nil, zerrors.ThrowInvalidArgumentf(nil, "MODEL-zF3rb", "reduce.wrong.event.type %s", user.PersonalAccessTokenAddedType)
		}
		return handler.NewCreateStatement(event,
			[]handler.Column{
				handler.NewCol(instanceIDCol, event.Aggregate().InstanceID),
				handler.NewCol(userIDCol, event.Aggregate().ID),
				handler.NewCol(resourceOwnerCol, event.Aggregate().ResourceOwner),
				handler.NewCol("id", e.TokenID),
				handler.NewCol(creationDateCol, event.CreatedAt()),
				handler.NewCol(changeDateCol, event.CreatedAt()),
				handler.NewCol("scopes", e.Scopes),
				handler.NewCol("expiration", e.Expiration),
				handler.NewCol("is_pat", true),
			},
		), nil
	case user.UserV1ProfileChangedType,
		user.HumanProfileChangedType:
		e, ok := event.(*user.HumanProfileChangedEvent)
		if !ok {
			return nil, zerrors.ThrowInvalidArgumentf(nil, "MODEL-ASF2t", "reduce.wrong.event.type %s", user.HumanProfileChangedType)
		}
		if e.PreferredLanguage == nil {
			return handler.NewNoOpStatement(event), nil
		}
		return handler.NewUpdateStatement(event,
			[]handler.Column{
				handler.NewCol("preferred_language", gu.Value(e.PreferredLanguage).String()),
			},
			[]handler.Condition{
				handler.NewCond(instanceIDCol, e.Aggregate().InstanceID),
				handler.NewCond(userIDCol, e.Aggregate().ID),
			},
		), nil
	case user.UserV1SignedOutType,
		user.HumanSignedOutType:
		id, err := agentIDFromSession(event)
		if err != nil {
			return nil, err
		}
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(instanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond(userIDCol, event.Aggregate().ID),
				handler.NewCond(userAgentIDCol, id),
			},
		), nil
	case user.UserLockedType,
		user.UserDeactivatedType,
		user.UserRemovedType:
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(instanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond(userIDCol, event.Aggregate().ID),
			},
		), nil
	case user.UserTokenRemovedType,
		user.PersonalAccessTokenRemovedType:
		id, err := tokenIDFromRemovedEvent(event)
		if err != nil {
			return nil, err
		}
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(instanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond("id", id),
			},
		), nil
	case user.HumanRefreshTokenRemovedType:
		id, err := refreshTokenIDFromRemovedEvent(event)
		if err != nil {
			return nil, err
		}
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(instanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond("refresh_token_id", id),
			},
		), nil
	case project.ApplicationDeactivatedType,
		project.ApplicationRemovedType:
		application, err := applicationFromSession(event)
		if err != nil {
			return nil, err
		}
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(instanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond("application_id", application.AppID),
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

		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(instanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond("application_id", applicationIDs),
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
				handler.NewCond(resourceOwnerCol, event.Aggregate().ResourceOwner),
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
