package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	_ "github.com/zitadel/zitadel/internal/notification/statik"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/milestone"
	"github.com/zitadel/zitadel/internal/repository/pseudo"
)

const (
	TelemetryProjectionTable = "projections.telemetry"
)

type TelemetryPusherConfig struct {
	Enabled   bool
	Endpoints []string
	Headers   http.Header
	Limit     uint64
}

type telemetryPusher struct {
	cfg                            TelemetryPusherConfig
	commands                       *command.Commands
	queries                        *NotificationQueries
	metricSuccessfulDeliveriesJSON string
	metricFailedDeliveriesJSON     string
}

func NewTelemetryPusher(
	ctx context.Context,
	telemetryCfg TelemetryPusherConfig,
	handlerCfg handler.Config,
	commands *command.Commands,
	queries *NotificationQueries,
	metricSuccessfulDeliveriesJSON,
	metricFailedDeliveriesJSON string,
) *handler.Handler {
	pusher := &telemetryPusher{
		cfg:                            telemetryCfg,
		commands:                       commands,
		queries:                        queries,
		metricSuccessfulDeliveriesJSON: metricFailedDeliveriesJSON,
		metricFailedDeliveriesJSON:     metricFailedDeliveriesJSON,
	}
	handlerCfg.TriggerWithoutEvents = pusher.pushMilestones
	return handler.NewHandler(
		ctx,
		&handlerCfg,
		pusher,
	)
}

func (u *telemetryPusher) Name() string {
	return TelemetryProjectionTable
}

func (t *telemetryPusher) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{{
		Aggregate: pseudo.AggregateType,
		EventRedusers: []handler.EventReducer{{
			Event:  pseudo.ScheduledEventType,
			Reduce: t.pushMilestones,
		}},
	}}
}

func (t *telemetryPusher) pushMilestones(event eventstore.Event) (*handler.Statement, error) {
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
			Limit:         t.cfg.Limit,
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

	return handler.NewNoOpStatement(scheduledEvent), nil
}

func (t *telemetryPusher) pushMilestone(ctx context.Context, event *pseudo.ScheduledEvent, ms *query.Milestone) error {
	ctx = authz.WithInstanceID(ctx, ms.InstanceID)
	alreadyHandled, err := t.queries.IsAlreadyHandled(ctx, event, map[string]interface{}{"type": ms.Type.String()}, milestone.AggregateType, milestone.PushedEventType)
	if err != nil {
		return err
	}
	if alreadyHandled {
		return nil
	}

	for _, endpoint := range t.cfg.Endpoints {
		if err := types.SendJSON(
			ctx,
			webhook.Config{
				CallURL: endpoint,
				Method:  http.MethodPost,
				Headers: t.cfg.Headers,
			},
			t.queries.GetFileSystemProvider,
			t.queries.GetLogProvider,
			&struct {
				InstanceID     string         `json:"instanceId"`
				ExternalDomain string         `json:"externalDomain"`
				PrimaryDomain  string         `json:"primaryDomain"`
				Type           milestone.Type `json:"type"`
				ReachedDate    time.Time      `json:"reached"`
			}{
				InstanceID:     ms.InstanceID,
				ExternalDomain: t.queries.externalDomain,
				PrimaryDomain:  ms.PrimaryDomain,
				Type:           ms.Type,
				ReachedDate:    ms.ReachedDate,
			},
			event,
			t.metricSuccessfulDeliveriesJSON,
			t.metricFailedDeliveriesJSON,
		).WithoutTemplate(); err != nil {
			return err
		}
	}
	return t.commands.MilestonePushed(ctx, ms.Type, t.cfg.Endpoints, ms.PrimaryDomain)
}
