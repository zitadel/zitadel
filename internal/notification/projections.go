package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/notification/handlers"
	_ "github.com/zitadel/zitadel/internal/notification/statik"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/queue"
)

var (
	projections []*handler.Handler
)

func Register(
	ctx context.Context,
	userHandlerCustomConfig, quotaHandlerCustomConfig, telemetryHandlerCustomConfig, backChannelLogoutHandlerCustomConfig projection.CustomConfig,
	notificationWorkerConfig handlers.WorkerConfig,
	telemetryCfg handlers.TelemetryPusherConfig,
	externalDomain string,
	externalPort uint16,
	externalSecure bool,
	commands *command.Commands,
	queries *query.Queries,
	es *eventstore.Eventstore,
	otpEmailTmpl, fileSystemPath string,
	userEncryption, smtpEncryption, smsEncryption, keysEncryptionAlg crypto.EncryptionAlgorithm,
	tokenLifetime time.Duration,
	queue *queue.Queue,
) {
	if !notificationWorkerConfig.LegacyEnabled {
		queue.ShouldStart()
	}

	// make sure the slice does not contain old values
	projections = nil

	q := handlers.NewNotificationQueries(queries, es, externalDomain, externalPort, externalSecure, fileSystemPath, userEncryption, smtpEncryption, smsEncryption)
	c := newChannels(q)
	projections = append(projections, handlers.NewUserNotifier(ctx, projection.ApplyCustomConfig(userHandlerCustomConfig), commands, q, c, otpEmailTmpl, notificationWorkerConfig, queue))
	projections = append(projections, handlers.NewQuotaNotifier(ctx, projection.ApplyCustomConfig(quotaHandlerCustomConfig), commands, q, c))
	projections = append(projections, handlers.NewBackChannelLogoutNotifier(
		ctx,
		projection.ApplyCustomConfig(backChannelLogoutHandlerCustomConfig),
		commands,
		q,
		es,
		keysEncryptionAlg,
		c,
		tokenLifetime,
	))
	if telemetryCfg.Enabled {
		projections = append(projections, handlers.NewTelemetryPusher(ctx, telemetryCfg, projection.ApplyCustomConfig(telemetryHandlerCustomConfig), commands, q, c))
	}
	if !notificationWorkerConfig.LegacyEnabled {
		queue.AddWorkers(handlers.NewNotificationWorker(notificationWorkerConfig, commands, q, c))
	}
}

func Start(ctx context.Context) {
	for _, projection := range projections {
		projection.Start(ctx)
	}
}

func SetCurrentState(ctx context.Context, es *eventstore.Eventstore) error {
	if len(projections) == 0 {
		return nil
	}
	position, err := es.LatestPosition(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsMaxPosition).InstanceID(authz.GetInstance(ctx).InstanceID()).OrderDesc().Limit(1))
	if err != nil {
		return err
	}

	for i, projection := range projections {
		logging.WithFields("name", projection.ProjectionName(), "instance", authz.GetInstance(ctx).InstanceID(), "index", fmt.Sprintf("%d/%d", i, len(projections))).Info("set current state of notification projection")
		_, err = projection.Trigger(ctx, handler.WithMinPosition(position))
		if err != nil {
			return err
		}
		logging.WithFields("name", projection.ProjectionName(), "instance", authz.GetInstance(ctx).InstanceID(), "index", fmt.Sprintf("%d/%d", i, len(projections))).Info("current state of notification projection set")
	}
	return nil
}

func ProjectInstance(ctx context.Context) error {
	for i, projection := range projections {
		logging.WithFields("name", projection.ProjectionName(), "instance", authz.GetInstance(ctx).InstanceID(), "index", fmt.Sprintf("%d/%d", i, len(projections))).Info("starting notification projection")
		_, err := projection.Trigger(ctx)
		if err != nil {
			return err
		}
		logging.WithFields("name", projection.ProjectionName(), "instance", authz.GetInstance(ctx).InstanceID(), "index", fmt.Sprintf("%d/%d", i, len(projections))).Info("notification projection done")
	}
	return nil
}

func Projections() []*handler.Handler {
	return projections
}
