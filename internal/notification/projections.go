package notification

import (
	"context"

	statik_fs "github.com/rakyll/statik/fs"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/handlers"
	_ "github.com/zitadel/zitadel/internal/notification/statik"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/metrics"
)

const (
	metricSuccessfulDeliveriesEmail = "successful_deliveries_email"
	metricFailedDeliveriesEmail     = "failed_deliveries_email"
	metricSuccessfulDeliveriesSMS   = "successful_deliveries_sms"
	metricFailedDeliveriesSMS       = "failed_deliveries_sms"
	metricSuccessfulDeliveriesJSON  = "successful_deliveries_json"
	metricFailedDeliveriesJSON      = "failed_deliveries_json"
)

func Start(
	ctx context.Context,
	userHandlerCustomConfig projection.CustomConfig,
	quotaHandlerCustomConfig projection.CustomConfig,
	telemetryHandlerCustomConfig projection.CustomConfig,
	telemetryCfg handlers.TelemetryPusherConfig,
	externalDomain string,
	externalPort uint16,
	externalSecure bool,
	commands *command.Commands,
	queries *query.Queries,
	es *eventstore.Eventstore,
	assetsPrefix func(context.Context) string,
	fileSystemPath string,
	userEncryption,
	smtpEncryption,
	smsEncryption crypto.EncryptionAlgorithm,
) {
	statikFS, err := statik_fs.NewWithNamespace("notification")
	logging.OnError(err).Panic("unable to start listener")
	err = metrics.RegisterCounter(metricSuccessfulDeliveriesEmail, "Successfully delivered emails")
	logging.WithFields("metric", metricSuccessfulDeliveriesEmail).OnError(err).Panic("unable to register counter")
	err = metrics.RegisterCounter(metricFailedDeliveriesEmail, "Failed email deliveries")
	logging.WithFields("metric", metricFailedDeliveriesEmail).OnError(err).Panic("unable to register counter")
	err = metrics.RegisterCounter(metricSuccessfulDeliveriesSMS, "Successfully delivered SMS")
	logging.WithFields("metric", metricSuccessfulDeliveriesSMS).OnError(err).Panic("unable to register counter")
	err = metrics.RegisterCounter(metricFailedDeliveriesSMS, "Failed SMS deliveries")
	logging.WithFields("metric", metricFailedDeliveriesSMS).OnError(err).Panic("unable to register counter")
	err = metrics.RegisterCounter(metricSuccessfulDeliveriesJSON, "Successfully delivered JSON messages")
	logging.WithFields("metric", metricSuccessfulDeliveriesJSON).OnError(err).Panic("unable to register counter")
	err = metrics.RegisterCounter(metricFailedDeliveriesJSON, "Failed JSON message deliveries")
	logging.WithFields("metric", metricFailedDeliveriesJSON).OnError(err).Panic("unable to register counter")
	q := handlers.NewNotificationQueries(queries, es, externalDomain, externalPort, externalSecure, fileSystemPath, userEncryption, smtpEncryption, smsEncryption, statikFS)
	handlers.NewUserNotifier(
		ctx,
		projection.ApplyCustomConfig(userHandlerCustomConfig),
		commands,
		q,
		assetsPrefix,
		metricSuccessfulDeliveriesEmail,
		metricFailedDeliveriesEmail,
		metricSuccessfulDeliveriesSMS,
		metricFailedDeliveriesSMS,
	).Start()
	handlers.NewQuotaNotifier(
		ctx,
		projection.ApplyCustomConfig(quotaHandlerCustomConfig),
		commands,
		q,
		metricSuccessfulDeliveriesJSON,
		metricFailedDeliveriesJSON,
	).Start()
	if telemetryCfg.Enabled {
		handlers.NewTelemetryPusher(
			ctx,
			telemetryCfg,
			projection.ApplyCustomConfig(telemetryHandlerCustomConfig),
			commands,
			q,
			metricSuccessfulDeliveriesJSON,
			metricFailedDeliveriesJSON,
		).Start()
	}
}
