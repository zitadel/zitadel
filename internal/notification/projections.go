package notification

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/notification/handlers"
	_ "github.com/zitadel/zitadel/internal/notification/statik"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/query/projection"
)

var (
	projections []*handler.Handler
	worker      *handlers.NotificationWorker
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
	client *database.DB,
) {
	q := handlers.NewNotificationQueries(queries, es, externalDomain, externalPort, externalSecure, fileSystemPath, userEncryption, smtpEncryption, smsEncryption)
	c := newChannels(q)
	projections = append(projections, handlers.NewUserNotifier(ctx, projection.ApplyCustomConfig(userHandlerCustomConfig), commands, q, c, otpEmailTmpl, notificationWorkerConfig.LegacyEnabled))
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
	worker = handlers.NewNotificationWorker(notificationWorkerConfig, commands, q, es, client, c)
}

func Start(ctx context.Context) {
	for _, projection := range projections {
		projection.Start(ctx)
	}
	worker.Start(ctx)
}

func ProjectInstance(ctx context.Context) error {
	for _, projection := range projections {
		_, err := projection.Trigger(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func Projections() []*handler.Handler {
	return projections
}
