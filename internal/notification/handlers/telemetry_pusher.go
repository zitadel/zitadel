package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	_ "github.com/zitadel/zitadel/internal/notification/statik"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/pseudo"
)

const (
	TelemetryProjectionTable = "projections.telemetry"
)

type TelemetryPusherConfig struct {
	Enabled   bool
	Endpoints []string
}

type telemetryPusher struct {
	crdb.StatementHandler
	commands                       *command.Commands
	queries                        *NotificationQueries
	metricSuccessfulDeliveriesJSON string
	metricFailedDeliveriesJSON     string
	endpoints                      []string
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
	handlerCfg.Reducers = []handler.AggregateReducer{{}}
	if telemetryCfg.Enabled {
		handlerCfg.Reducers = p.reducers()
	}
	p.endpoints = telemetryCfg.Endpoints
	p.StatementHandler = crdb.NewStatementHandler(ctx, handlerCfg)
	p.commands = commands
	p.queries = queries
	p.metricSuccessfulDeliveriesJSON = metricSuccessfulDeliveriesJSON
	p.metricFailedDeliveriesJSON = metricFailedDeliveriesJSON
	projection.TelemetryPusherProjection = p
	return p
}

func (t *telemetryPusher) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{{
		Aggregate: pseudo.AggregateType,
		EventRedusers: []handler.EventReducer{{
			Event:  pseudo.ScheduledEventType,
			Reduce: t.pushMilestones,
		}},
	}}
}

func printEvent(event eventstore.Event) {
	bytes, err := json.MarshalIndent(event, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(event.Type(), string(bytes))
}

func (t *telemetryPusher) pushMilestones(event eventstore.Event) (*handler.Statement, error) {
	printEvent(event)
	ctx := call.WithTimestamp(context.Background())
	scheduledEvent, ok := event.(*pseudo.ScheduledEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-lDTs5", "reduce.wrong.event.type %s", event.Type())
	}

	isReached, err := query.NewNotNullQuery(query.MilestoneReachedDateColID)
	if err != nil {
		return nil, err
	}
	isNotPushed, err := query.NewIsNullQuery(query.MilestonePushedDateColID)
	if err != nil {
		return nil, err
	}
	hasPrimaryDomain, err := query.NewNotNullQuery(query.MilestonePrimaryDomainColID)
	if err != nil {
		return nil, err
	}
	unpushedMilestones, err := t.queries.Queries.SearchMilestones(ctx, scheduledEvent.InstanceIDs, &query.MilestonesSearchQueries{
		SearchRequest: query.SearchRequest{
			Limit:         100,
			SortingColumn: query.MilestoneReachedDateColID,
			Asc:           true,
		},
		Queries: []query.SearchQuery{isReached, isNotPushed, hasPrimaryDomain},
	})
	if err != nil {
		return nil, err
	}
	var errs int
	for _, ms := range unpushedMilestones.Milestones {
		if err = t.pushMilestone(ctx, scheduledEvent, ms); err != nil {
			errs++
			logging.Warnf("pushing milestone %+v failed: %s", *ms, err.Error())
		}
	}
	if errs > 0 {
		return nil, fmt.Errorf("pushing %d of %d milestones failed", errs, unpushedMilestones.Count)
	}

	return crdb.NewNoOpStatement(scheduledEvent), nil
}

func (t *telemetryPusher) pushMilestone(ctx context.Context, event *pseudo.ScheduledEvent, ms *query.Milestone) error {
	ctx = authz.WithInstanceID(ctx, ms.InstanceID)
	for _, endpoint := range t.endpoints {
		if err := types.SendJSON(
			ctx,
			webhook.Config{
				CallURL: endpoint,
				Method:  http.MethodPost,
			},
			t.queries.GetFileSystemProvider,
			t.queries.GetLogProvider,
			ms,
			event,
			t.metricSuccessfulDeliveriesJSON,
			t.metricFailedDeliveriesJSON,
		).WithoutTemplate(); err != nil {
			return err
		}
	}
	return t.commands.MilestonePushed(ctx, ms.Type, t.endpoints, ms.PrimaryDomain)
}
