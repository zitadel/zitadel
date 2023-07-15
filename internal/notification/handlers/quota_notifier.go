package handlers

import (
	"context"
	"net/http"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	_ "github.com/zitadel/zitadel/internal/notification/statik"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

const (
	QuotaNotificationsProjectionTable = "projections.notifications_quota"
)

type quotaNotifier struct {
	crdb.StatementHandler
	commands                       *command.Commands
	queries                        *NotificationQueries
	metricSuccessfulDeliveriesJSON string
	metricFailedDeliveriesJSON     string
}

func NewQuotaNotifier(
	ctx context.Context,
	config crdb.StatementHandlerConfig,
	commands *command.Commands,
	queries *NotificationQueries,
	metricSuccessfulDeliveriesJSON,
	metricFailedDeliveriesJSON string,
) *quotaNotifier {
	p := new(quotaNotifier)
	config.ProjectionName = QuotaNotificationsProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	p.commands = commands
	p.queries = queries
	p.metricSuccessfulDeliveriesJSON = metricSuccessfulDeliveriesJSON
	p.metricFailedDeliveriesJSON = metricFailedDeliveriesJSON
	projection.NotificationsQuotaProjection = p
	return p
}

func (u *quotaNotifier) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: quota.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  quota.NotificationDueEventType,
					Reduce: u.reduceNotificationDue,
				},
			},
		},
	}
}

func (u *quotaNotifier) reduceNotificationDue(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*quota.NotificationDueEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-DLxdE", "reduce.wrong.event.type %s", quota.NotificationDueEventType)
	}
	ctx := HandlerContext(event.Aggregate())
	alreadyHandled, err := u.queries.IsAlreadyHandled(ctx, event, map[string]interface{}{"dueEventID": e.ID}, quota.AggregateType, quota.NotifiedEventType)
	if err != nil {
		return nil, err
	}
	if alreadyHandled {
		return crdb.NewNoOpStatement(e), nil
	}
	err = types.SendJSON(
		ctx,
		webhook.Config{
			CallURL: e.CallURL,
			Method:  http.MethodPost,
		},
		u.queries.GetFileSystemProvider,
		u.queries.GetLogProvider,
		e,
		e,
		u.metricSuccessfulDeliveriesJSON,
		u.metricFailedDeliveriesJSON,
	).WithoutTemplate()
	if err != nil {
		return nil, err
	}
	err = u.commands.UsageNotificationSent(ctx, e)
	if err != nil {
		return nil, err
	}
	return crdb.NewNoOpStatement(e), nil
}
