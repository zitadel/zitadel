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

	tokenInstanceIDCol        = "instance_id"
	tokenUserIDCol            = "user_id"
	tokenResourceOwnerCol     = "resource_owner"
	tokenCreationDateCol      = "creation_date"
	tokenChangeDateCol        = "change_date"
	tokenApplicationIDCol     = "application_id"
	tokenUserAgentIDCol       = "user_agent_id"
	tokenSequencerCol         = "sequence"
	tokenActorCol             = "actor"
	tokenIDCol                = "id"
	tokenAudienceCol          = "audience"
	tokenPreferredLanguageCol = "preferred_language"
	tokenExpirationCol        = "expiration"
	tokenRefreshTokenIDCol    = "refresh_token_id"
	tokenScopesCol            = "scopes"
	tokenIsPatCol             = "is_pat"
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
				handler.NewCol(tokenInstanceIDCol, event.Aggregate().InstanceID),
				handler.NewCol(tokenUserIDCol, event.Aggregate().ID),
				handler.NewCol(tokenResourceOwnerCol, event.Aggregate().ResourceOwner),
				handler.NewCol(tokenIDCol, e.TokenID),
				handler.NewCol(tokenCreationDateCol, event.CreatedAt()),
				handler.NewCol(tokenChangeDateCol, event.CreatedAt()),
				handler.NewCol(tokenSequencerCol, event.Sequence()),
				handler.NewCol(tokenApplicationIDCol, e.ApplicationID),
				handler.NewCol(tokenUserAgentIDCol, e.UserAgentID),
				handler.NewCol(tokenAudienceCol, e.Audience),
				handler.NewCol(tokenScopesCol, e.Scopes),
				handler.NewCol(tokenExpirationCol, e.Expiration),
				handler.NewCol(tokenPreferredLanguageCol, e.PreferredLanguage),
				handler.NewCol(tokenRefreshTokenIDCol, e.RefreshTokenID),
				handler.NewCol(tokenActorCol, view_model.TokenActor{TokenActor: e.Actor}),
				handler.NewCol(tokenIsPatCol, false),
			},
		), nil
	case user.PersonalAccessTokenAddedType:
		e, ok := event.(*user.PersonalAccessTokenAddedEvent)
		if !ok {
			return nil, zerrors.ThrowInvalidArgumentf(nil, "MODEL-zF3rb", "reduce.wrong.event.type %s", user.PersonalAccessTokenAddedType)
		}
		return handler.NewCreateStatement(event,
			[]handler.Column{
				handler.NewCol(tokenInstanceIDCol, event.Aggregate().InstanceID),
				handler.NewCol(tokenUserIDCol, event.Aggregate().ID),
				handler.NewCol(tokenResourceOwnerCol, event.Aggregate().ResourceOwner),
				handler.NewCol(tokenIDCol, e.TokenID),
				handler.NewCol(tokenCreationDateCol, event.CreatedAt()),
				handler.NewCol(tokenChangeDateCol, event.CreatedAt()),
				handler.NewCol(tokenSequencerCol, event.Sequence()),
				handler.NewCol(tokenScopesCol, e.Scopes),
				handler.NewCol(tokenExpirationCol, e.Expiration),
				handler.NewCol(tokenIsPatCol, true),
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
				handler.NewCol(tokenPreferredLanguageCol, gu.Value(e.PreferredLanguage).String()),
				handler.NewCol(tokenChangeDateCol, event.CreatedAt()),
				handler.NewCol(tokenSequencerCol, event.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(tokenInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCond(tokenUserIDCol, e.Aggregate().ID),
			},
		), nil
	case user.UserV1SignedOutType,
		user.HumanSignedOutType:
		e, ok := event.(*user.HumanSignedOutEvent)
		if !ok {
			return nil, zerrors.ThrowInvalidArgumentf(nil, "MODEL-Wtn2q", "reduce.wrong.event.type %s", user.HumanSignedOutType)
		}
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(tokenInstanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond(tokenUserIDCol, event.Aggregate().ID),
				handler.NewCond(tokenUserAgentIDCol, e.UserAgentID),
			},
		), nil
	case user.UserLockedType,
		user.UserDeactivatedType,
		user.UserRemovedType:
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(tokenInstanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond(tokenUserIDCol, event.Aggregate().ID),
			},
		), nil
	case user.UserTokenRemovedType,
		user.PersonalAccessTokenRemovedType:
		var tokenID string
		switch e := event.(type) {
		case *user.UserTokenRemovedEvent:
			tokenID = e.TokenID
		case *user.PersonalAccessTokenRemovedEvent:
			tokenID = e.TokenID
		default:
			return nil, zerrors.ThrowInvalidArgumentf(nil, "MODEL-SF3ga", "reduce.wrong.event.type %s", user.UserTokenRemovedType)
		}
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(tokenInstanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond(tokenIDCol, tokenID),
			},
		), nil
	case user.HumanRefreshTokenRemovedType:
		e, ok := event.(*user.HumanRefreshTokenRemovedEvent)
		if !ok {
			return nil, zerrors.ThrowInvalidArgumentf(nil, "MODEL-Sfe11", "reduce.wrong.event.type %s", user.HumanRefreshTokenRemovedType)
		}
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(tokenInstanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond(tokenRefreshTokenIDCol, e.TokenID),
			},
		), nil
	case project.ApplicationDeactivatedType,
		project.ApplicationRemovedType:
		var applicationID string
		switch e := event.(type) {
		case *project.ApplicationDeactivatedEvent:
			applicationID = e.AppID
		case *project.ApplicationRemovedEvent:
			applicationID = e.AppID
		default:
			return nil, zerrors.ThrowInvalidArgumentf(nil, "MODEL-SF3fq", "reduce.wrong.event.type  %v", []eventstore.EventType{project.ApplicationDeactivatedType, project.ApplicationRemovedType})
		}
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(tokenInstanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond(tokenApplicationIDCol, applicationID),
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
				handler.NewCond(tokenInstanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond(tokenApplicationIDCol, applicationIDs),
			},
		), nil
	case instance.InstanceRemovedEventType:
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(tokenInstanceIDCol, event.Aggregate().InstanceID),
			},
		), nil
	case org.OrgRemovedEventType:
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(tokenInstanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond(tokenResourceOwnerCol, event.Aggregate().ResourceOwner),
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
