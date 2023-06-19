package handlers

import (
	"context"
	"net/http"

	"github.com/zitadel/zitadel/internal/repository/milestone"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	_ "github.com/zitadel/zitadel/internal/notification/statik"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/query/projection"
)

const (
	TelemetryProjectionTable = "projections.telemetry"
)

type TelemetryPusherConfig struct {
	Enabled   bool
	Endpoints []string
}

type telemetryPusher struct {
	cfg TelemetryPusherConfig
	crdb.StatementHandler
	commands                       *command.Commands
	queries                        *NotificationQueries
	metricSuccessfulDeliveriesJSON string
	metricFailedDeliveriesJSON     string
}

func NewTelemetryPusher(
	ctx context.Context,
	telemetryCfg TelemetryPusherConfig,
	handlerCfg crdb.StatementHandlerConfig,
	commands *command.Commands,
	queries *NotificationQueries,
	metricSuccessfulDeliveriesJSON,
	metricFailedDeliveriesJSON string,
) *telemetryPusher {
	p := new(telemetryPusher)
	handlerCfg.ProjectionName = TelemetryProjectionTable
	handlerCfg.Reducers = p.reducers()
	p.cfg = telemetryCfg
	p.StatementHandler = crdb.NewStatementHandler(ctx, handlerCfg)
	p.commands = commands
	p.queries = queries
	p.metricSuccessfulDeliveriesJSON = metricSuccessfulDeliveriesJSON
	p.metricFailedDeliveriesJSON = metricFailedDeliveriesJSON
	projection.TelemetryPusherProjection = p
	return p
}

func (t *telemetryPusher) reducers() []handler.AggregateReducer {
	if !t.cfg.Enabled {
		return nil
	}
	return []handler.AggregateReducer{
		{
			Aggregate: milestone.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  milestone.ReachedEventType,
					Reduce: t.reduceTelemetryPushDue,
				},
			},
		},
	}
}

func (t *telemetryPusher) reduceTelemetryPushDue(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*milestone.ReachedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-UjA3E", "reduce.wrong.event.type %s", milestone.ReachedEventType)
	}
	ctx := HandlerContext(event.Aggregate())
	alreadyHandled, err := t.queries.IsAlreadyHandled(ctx, event, nil, milestone.AggregateType, milestone.PushedEventType)
	if err != nil {
		return nil, err
	}
	if alreadyHandled {
		return crdb.NewNoOpStatement(e), nil
	}
	for _, endpoint := range t.cfg.Endpoints {
		if err = types.SendJSON(
			ctx,
			webhook.Config{
				CallURL: endpoint,
				Method:  http.MethodPost,
			},
			t.queries.GetFileSystemProvider,
			t.queries.GetLogProvider,
			e,
			e,
			t.metricSuccessfulDeliveriesJSON,
			t.metricFailedDeliveriesJSON,
		).WithoutTemplate(); err != nil {
			return nil, err
		}
	}

	err = t.commands.TelemetryPushed(ctx, e, t.cfg.Endpoints)
	if err != nil {
		return nil, err
	}
	return crdb.NewNoOpStatement(e), nil
}
