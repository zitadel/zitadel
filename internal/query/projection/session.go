package projection

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/session"
)

const (
	SessionsProjectionTable = "projections.sessions"

	SessionColumnID                = "id"
	SessionColumnCreationDate      = "creation_date"
	SessionColumnChangeDate        = "change_date"
	SessionColumnSequence          = "sequence"
	SessionColumnState             = "state"
	SessionColumnResourceOwner     = "resource_owner"
	SessionColumnInstanceID        = "instance_id"
	SessionColumnCreator           = "creator" //TODO: client?
	SessionColumnUserID            = "user_id"
	SessionColumnUserCheckedAt     = "user_checked_at"
	SessionColumnPasswordCheckedAt = "password_checked_at"
	SessionColumnMetadata          = "metadata"
)

type sessionProjection struct {
	crdb.StatementHandler
}

func newSessionProjection(ctx context.Context, config crdb.StatementHandlerConfig) *sessionProjection {
	p := new(sessionProjection)
	config.ProjectionName = SessionsProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewMultiTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(SessionColumnID, crdb.ColumnTypeText),
			crdb.NewColumn(SessionColumnCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(SessionColumnChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(SessionColumnSequence, crdb.ColumnTypeInt64),
			crdb.NewColumn(SessionColumnState, crdb.ColumnTypeEnum),
			crdb.NewColumn(SessionColumnResourceOwner, crdb.ColumnTypeText),
			crdb.NewColumn(SessionColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(SessionColumnCreator, crdb.ColumnTypeText),
			crdb.NewColumn(SessionColumnUserID, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(SessionColumnUserCheckedAt, crdb.ColumnTypeTimestamp, crdb.Nullable()),
			crdb.NewColumn(SessionColumnPasswordCheckedAt, crdb.ColumnTypeTimestamp, crdb.Nullable()),
			crdb.NewColumn(SessionColumnMetadata, crdb.ColumnTypeBytes, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(SessionColumnInstanceID, SessionColumnID),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *sessionProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: session.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  session.SetType,
					Reduce: p.reduceSessionSet,
				},
				{
					Event:  session.TerminateType,
					Reduce: p.reduceSessionTerminated,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(SMSColumnInstanceID),
				},
			},
		},
	}
}

func (p *sessionProjection) reduceSessionSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*session.SetEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-SFegf", "reduce.wrong.event.type %s", session.SetType)
	}

	columns := []handler.Column{
		handler.NewCol(SessionColumnID, e.Aggregate().ID),
		handler.NewCol(SessionColumnInstanceID, e.Aggregate().InstanceID),
		handler.NewCol(SessionColumnCreationDate, e.CreationDate()), // TODO: overwrites every time
		handler.NewCol(SessionColumnChangeDate, e.CreationDate()),
		handler.NewCol(SessionColumnResourceOwner, e.Aggregate().ResourceOwner),
		handler.NewCol(SessionColumnState, domain.SessionStateActive),
		handler.NewCol(SessionColumnSequence, e.Sequence()),
		handler.NewCol(SessionColumnCreator, e.User),
	}
	if e.UserID != nil {
		columns = append(columns, handler.NewCol(SessionColumnUserID, *e.UserID))
	}
	if e.UserCheckedAt != nil {
		columns = append(columns, handler.NewCol(SessionColumnUserCheckedAt, *e.UserCheckedAt))
	}
	if e.PasswordCheckedAt != nil {
		columns = append(columns, handler.NewCol(SessionColumnPasswordCheckedAt, *e.PasswordCheckedAt))
	}
	if len(e.Metadata) != 0 {
		m, _ := json.Marshal(e.Metadata)
		columns = append(columns, handler.NewCol(SessionColumnMetadata, m))
	}

	return crdb.NewUpsertStatement(
		e,
		[]handler.Column{
			handler.NewCol(SessionColumnID, e.Aggregate().ID),
			handler.NewCol(SessionColumnInstanceID, e.Aggregate().InstanceID),
		},
		columns,
	), nil
}

func (p *sessionProjection) reduceSessionTerminated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*session.TerminateEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-SAftn", "reduce.wrong.event.type %s", session.TerminateType)
	}

	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(SessionColumnID, e.Aggregate().ID),
			handler.NewCond(SessionColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}
