package domain

import (
	"log/slog"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/eventstore"
	"github.com/zitadel/zitadel/backend/v3/telemetry/logging"
	"github.com/zitadel/zitadel/backend/v3/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
)

var (
	pool             database.Pool
	tracer           tracing.Tracer
	logger           logging.Logger = *logging.NewLogger(slog.Default())
	legacyEventstore eventstore.LegacyEventstore
	sysConfig        systemdefaults.SystemDefaults
)

func SetPool(p database.Pool) {
	pool = p
}

func SetTracer(t tracing.Tracer) {
	tracer = t
}

func SetLogger(l logging.Logger) {
	logger = l
}

func SetLegacyEventstore(es eventstore.LegacyEventstore) {
	legacyEventstore = es
}

func SetSystemConfig(cfg systemdefaults.SystemDefaults) {
	sysConfig = cfg
}
