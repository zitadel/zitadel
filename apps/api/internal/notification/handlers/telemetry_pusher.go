package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	_ "github.com/zitadel/zitadel/internal/notification/statik"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/repository/milestone"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	TelemetryProjectionTable = "projections.telemetry"
)

type TelemetryPusherConfig struct {
	Enabled   bool
	Endpoints []string
	Headers   http.Header
}

type telemetryPusher struct {
	cfg      TelemetryPusherConfig
	commands *command.Commands
	queries  *NotificationQueries
	channels types.ChannelChains
}

func NewTelemetryPusher(
	ctx context.Context,
	telemetryCfg TelemetryPusherConfig,
	handlerCfg handler.Config,
	commands *command.Commands,
	queries *NotificationQueries,
	channels types.ChannelChains,
) *handler.Handler {
	pusher := &telemetryPusher{
		cfg:      telemetryCfg,
		commands: commands,
		queries:  queries,
		channels: channels,
	}
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
		Aggregate: milestone.AggregateType,
		EventReducers: []handler.EventReducer{{
			Event:  milestone.ReachedEventType,
			Reduce: t.pushMilestones,
		}},
	}}
}

func (t *telemetryPusher) pushMilestones(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*milestone.ReachedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-lDTs5", "reduce.wrong.event.type %s", event.Type())
	}
	return handler.NewStatement(event, func(ctx context.Context, _ handler.Executer, _ string) error {
		// Do not push the milestone again if this was a migration event.
		if e.ReachedDate != nil {
			return nil
		}
		return t.pushMilestone(ctx, e)
	}), nil
}

func (t *telemetryPusher) pushMilestone(ctx context.Context, e *milestone.ReachedEvent) error {
	for _, endpoint := range t.cfg.Endpoints {
		if err := types.SendJSON(
			ctx,
			webhook.Config{
				CallURL: endpoint,
				Method:  http.MethodPost,
				Headers: t.cfg.Headers,
			},
			t.channels,
			&struct {
				InstanceID     string         `json:"instanceId"`
				ExternalDomain string         `json:"externalDomain"`
				Type           milestone.Type `json:"type"`
				ReachedDate    time.Time      `json:"reached"`
			}{
				InstanceID:     e.Agg.InstanceID,
				ExternalDomain: t.queries.externalDomain,
				Type:           e.MilestoneType,
				ReachedDate:    e.GetReachedDate(),
			},
			e.EventType,
		).WithoutTemplate(); err != nil {
			return err
		}
	}
	return t.commands.MilestonePushed(ctx, e.Agg.InstanceID, e.MilestoneType, t.cfg.Endpoints)
}
