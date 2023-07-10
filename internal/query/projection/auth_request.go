package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
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

type authRequestProjection struct {
	crdb.StatementHandler
}

func newAuthRequestProjection(ctx context.Context, config crdb.StatementHandlerConfig) *authRequestProjection {
	p := new(authRequestProjection)
	config.ProjectionName = AuthRequestsProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewMultiTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(AuthRequestColumnID, crdb.ColumnTypeText),
			crdb.NewColumn(AuthRequestColumnCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(AuthRequestColumnChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(AuthRequestColumnSequence, crdb.ColumnTypeInt64),
			crdb.NewColumn(AuthRequestColumnResourceOwner, crdb.ColumnTypeText),
			crdb.NewColumn(AuthRequestColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(AuthRequestColumnLoginClient, crdb.ColumnTypeText),
			crdb.NewColumn(AuthRequestColumnClientID, crdb.ColumnTypeText),
			crdb.NewColumn(AuthRequestColumnRedirectURI, crdb.ColumnTypeText),
			crdb.NewColumn(AuthRequestColumnScope, crdb.ColumnTypeTextArray),
			crdb.NewColumn(AuthRequestColumnPrompt, crdb.ColumnTypeEnumArray, crdb.Nullable()),
			crdb.NewColumn(AuthRequestColumnUILocales, crdb.ColumnTypeTextArray, crdb.Nullable()),
			crdb.NewColumn(AuthRequestColumnMaxAge, crdb.ColumnTypeInt64, crdb.Nullable()),
			crdb.NewColumn(AuthRequestColumnLoginHint, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(AuthRequestColumnHintUserID, crdb.ColumnTypeText, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(AuthRequestColumnInstanceID, AuthRequestColumnID),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *authRequestProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: authrequest.AggregateType,
			EventRedusers: []handler.EventReducer{
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
			EventRedusers: []handler.EventReducer{
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

	return crdb.NewCreateStatement(
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

	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(AuthRequestColumnID, event.Aggregate().ID),
			handler.NewCond(AuthRequestColumnInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}
