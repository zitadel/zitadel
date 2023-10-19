package handlers

import (
	"context"
	"net/http"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	_ "github.com/zitadel/zitadel/internal/notification/statik"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

const (
	QuotaNotificationsProjectionTable = "projections.notifications_quota"
)

type quotaNotifier struct {
	commands *command.Commands
	queries  *NotificationQueries
	channels types.ChannelChains
}

func NewQuotaNotifier(
	ctx context.Context,
	config handler.Config,
	commands *command.Commands,
	queries *NotificationQueries,
	channels types.ChannelChains,
) *handler.Handler {
	return handler.NewHandler(ctx, &config, &quotaNotifier{
		commands: commands,
		queries:  queries,
		channels: channels,
	})

}

func (*quotaNotifier) Name() string {
	return QuotaNotificationsProjectionTable
}

func (u *quotaNotifier) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: quota.AggregateType,
			EventReducers: []handler.EventReducer{
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

	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		ctx := HandlerContext(event.Aggregate())
		alreadyHandled, err := u.queries.IsAlreadyHandled(ctx, event, map[string]interface{}{"dueEventID": e.ID}, quota.AggregateType, quota.NotifiedEventType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}
		err = types.SendJSON(ctx, webhook.Config{CallURL: e.CallURL, Method: http.MethodPost}, u.channels, e, e).WithoutTemplate()
		if err != nil {
			return err
		}
		return u.commands.UsageNotificationSent(ctx, e)
	}), nil
}
