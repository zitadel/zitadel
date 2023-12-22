package notification

import (
	"context"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/handlers"
	_ "github.com/zitadel/zitadel/internal/notification/statik"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/query/projection"
)

func Start(
	ctx context.Context,
	userHandlerCustomConfig, quotaHandlerCustomConfig, telemetryHandlerCustomConfig projection.CustomConfig,
	telemetryCfg handlers.TelemetryPusherConfig,
	externalDomain string,
	externalPort uint16,
	externalSecure bool,
	commands *command.Commands,
	queries *query.Queries,
	es *eventstore.Eventstore,
	otpEmailTmpl string,
	fileSystemPath string,
	userEncryption, smtpEncryption, smsEncryption crypto.EncryptionAlgorithm,
) {
	q := handlers.NewNotificationQueries(queries, es, externalDomain, externalPort, externalSecure, fileSystemPath, userEncryption, smtpEncryption, smsEncryption)
	c := newChannels(q)
	handlers.NewUserNotifier(ctx, projection.ApplyCustomConfig(userHandlerCustomConfig), commands, q, c, otpEmailTmpl).Start(ctx)
	handlers.NewQuotaNotifier(ctx, projection.ApplyCustomConfig(quotaHandlerCustomConfig), commands, q, c).Start(ctx)
	if telemetryCfg.Enabled {
		handlers.NewTelemetryPusher(ctx, telemetryCfg, projection.ApplyCustomConfig(telemetryHandlerCustomConfig), commands, q, c).Start(ctx)
	}
}
