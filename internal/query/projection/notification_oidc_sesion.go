package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/oidcsession"
)

const (
	NotificationOIDCSessionProjectionTable = "projections.notification_oidc_sessions"

	NotificationOIDCSessionColumnID            = "id"
	NotificationOIDCSessionColumnCreationDate  = "creation_date"
	NotificationOIDCSessionColumnChangeDate    = "change_date"
	NotificationOIDCSessionColumnResourceOwner = "resource_owner"
	NotificationOIDCSessionColumnInstanceID    = "instance_id"
	NotificationOIDCSessionColumnSequence      = "sequence"
	NotificationOIDCSessionColumnSessionID     = "session_id"
	NotificationOIDCSessionColumnClientID      = "client_id"
	NotificationOIDCSessionColumnUserID        = "user_id"
)

type notificationOIDCSessionProjection struct{}

func (*notificationOIDCSessionProjection) Name() string {
	return NotificationOIDCSessionProjectionTable
}

func newNotificationOIDCSessionProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(notificationOIDCSessionProjection))
}

// Init implements [handler.initializer]
func (p *notificationOIDCSessionProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(NotificationOIDCSessionColumnID, handler.ColumnTypeText),
			handler.NewColumn(NotificationOIDCSessionColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(NotificationOIDCSessionColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(NotificationOIDCSessionColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(NotificationOIDCSessionColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(NotificationOIDCSessionColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(NotificationOIDCSessionColumnSessionID, handler.ColumnTypeText),
			handler.NewColumn(NotificationOIDCSessionColumnClientID, handler.ColumnTypeText),
			handler.NewColumn(NotificationOIDCSessionColumnUserID, handler.ColumnTypeText),
		},
			handler.NewPrimaryKey(NotificationOIDCSessionColumnInstanceID, NotificationOIDCSessionColumnID),
		),
	)
}

func (p *notificationOIDCSessionProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: oidcsession.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  oidcsession.AddedType,
					Reduce: p.reduceOIDCSessionAdded,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(NotificationOIDCSessionColumnInstanceID),
				},
			},
		},
	}
}

func (p *notificationOIDCSessionProjection) reduceOIDCSessionAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*oidcsession.AddedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewCreateStatement(
			e,
			[]handler.Column{
				handler.NewCol(NotificationOIDCSessionColumnID, e.Aggregate().ID),
				handler.NewCol(NotificationOIDCSessionColumnCreationDate, e.CreationDate()),
				handler.NewCol(NotificationOIDCSessionColumnChangeDate, e.CreationDate()),
				handler.NewCol(NotificationOIDCSessionColumnResourceOwner, e.Aggregate().ResourceOwner),
				handler.NewCol(NotificationOIDCSessionColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(NotificationOIDCSessionColumnSequence, e.Sequence()),
				handler.NewCol(NotificationOIDCSessionColumnSessionID, e.SessionID),
				handler.NewCol(NotificationOIDCSessionColumnClientID, e.ClientID),
				handler.NewCol(NotificationOIDCSessionColumnUserID, e.UserID),
			}),
		nil
}
