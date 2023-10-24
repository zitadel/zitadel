package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/authrequest"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

const (
	AuthRequestsProjectionTable = "projections.auth_requests"

	AuthRequestColumnID            = "id"
	AuthRequestColumnCreationDate  = "creation_date"
	AuthRequestColumnChangeDate    = "change_date"
	AuthRequestColumnSequence      = "sequence"
	AuthRequestColumnResourceOwner = "resource_owner"
	AuthRequestColumnInstanceID    = "instance_id"
	AuthRequestColumnLoginClient   = "login_client"
	AuthRequestColumnClientID      = "client_id"
	AuthRequestColumnRedirectURI   = "redirect_uri"
	AuthRequestColumnScope         = "scope"
	AuthRequestColumnPrompt        = "prompt"
	AuthRequestColumnUILocales     = "ui_locales"
	AuthRequestColumnMaxAge        = "max_age"
	AuthRequestColumnLoginHint     = "login_hint"
	AuthRequestColumnHintUserID    = "hint_user_id"
)

type authRequestProjection struct{}

// Name implements handler.Projection.
func (*authRequestProjection) Name() string {
	return AuthRequestsProjectionTable
}

func newAuthRequestProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(authRequestProjection))
}

func (*authRequestProjection) Init() *old_handler.Check {
	return handler.NewMultiTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(AuthRequestColumnID, handler.ColumnTypeText),
			handler.NewColumn(AuthRequestColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(AuthRequestColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(AuthRequestColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(AuthRequestColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(AuthRequestColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(AuthRequestColumnLoginClient, handler.ColumnTypeText),
			handler.NewColumn(AuthRequestColumnClientID, handler.ColumnTypeText),
			handler.NewColumn(AuthRequestColumnRedirectURI, handler.ColumnTypeText),
			handler.NewColumn(AuthRequestColumnScope, handler.ColumnTypeTextArray),
			handler.NewColumn(AuthRequestColumnPrompt, handler.ColumnTypeEnumArray, handler.Nullable()),
			handler.NewColumn(AuthRequestColumnUILocales, handler.ColumnTypeTextArray, handler.Nullable()),
			handler.NewColumn(AuthRequestColumnMaxAge, handler.ColumnTypeInt64, handler.Nullable()),
			handler.NewColumn(AuthRequestColumnLoginHint, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(AuthRequestColumnHintUserID, handler.ColumnTypeText, handler.Nullable()),
		},
			handler.NewPrimaryKey(AuthRequestColumnInstanceID, AuthRequestColumnID),
		),
	)
}

func (p *authRequestProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: authrequest.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  authrequest.AddedType,
					Reduce: p.reduceAuthRequestAdded,
				},
				{
					Event:  authrequest.SucceededType,
					Reduce: p.reduceAuthRequestEnded,
				},
				{
					Event:  authrequest.FailedType,
					Reduce: p.reduceAuthRequestEnded,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(AuthRequestColumnInstanceID),
				},
			},
		},
	}
}

func (p *authRequestProjection) reduceAuthRequestAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*authrequest.AddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Sfwfa", "reduce.wrong.event.type %s", authrequest.AddedType)
	}

	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(AuthRequestColumnID, e.Aggregate().ID),
			handler.NewCol(AuthRequestColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(AuthRequestColumnCreationDate, e.CreationDate()),
			handler.NewCol(AuthRequestColumnChangeDate, e.CreationDate()),
			handler.NewCol(AuthRequestColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(AuthRequestColumnSequence, e.Sequence()),
			handler.NewCol(AuthRequestColumnLoginClient, e.LoginClient),
			handler.NewCol(AuthRequestColumnClientID, e.ClientID),
			handler.NewCol(AuthRequestColumnRedirectURI, e.RedirectURI),
			handler.NewCol(AuthRequestColumnScope, e.Scope),
			handler.NewCol(AuthRequestColumnPrompt, e.Prompt),
			handler.NewCol(AuthRequestColumnUILocales, e.UILocales),
			handler.NewCol(AuthRequestColumnMaxAge, e.MaxAge),
			handler.NewCol(AuthRequestColumnLoginHint, e.LoginHint),
			handler.NewCol(AuthRequestColumnHintUserID, e.HintUserID),
		},
	), nil
}

func (p *authRequestProjection) reduceAuthRequestEnded(event eventstore.Event) (*handler.Statement, error) {
	switch event.(type) {
	case *authrequest.SucceededEvent,
		*authrequest.FailedEvent:
		break
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ASF3h", "reduce.wrong.event.type %s", []eventstore.EventType{authrequest.SucceededType, authrequest.FailedType})
	}

	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(AuthRequestColumnID, event.Aggregate().ID),
			handler.NewCond(AuthRequestColumnInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}
