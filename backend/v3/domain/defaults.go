package domain

import (
	"log/slog"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/eventstore"
	"github.com/zitadel/zitadel/backend/v3/telemetry/logging"
	"github.com/zitadel/zitadel/backend/v3/telemetry/tracing"
)

// The variables could also be moved to a struct.
// I just started with the singleton pattern and kept it like this.
var (
	pool database.Pool
	// userCodeAlgorithm crypto.EncryptionAlgorithm
	tracer           tracing.Tracer
	logger           logging.Logger = *logging.NewLogger(slog.Default())
	legacyEventstore eventstore.LegacyEventstore

	// instanceCache cache.Cache[instanceCacheIndex, string, *Instance]
	// orgCache cache.Cache[orgCacheIndex, string, *Org]

	// generateID func() (string, error) = func() (string, error) {
	// 	return strconv.FormatUint(rand.Uint64(), 10), nil
	// }
)

func SetPool(p database.Pool) {
	pool = p
}

// func SetUserCodeAlgorithm(algorithm crypto.EncryptionAlgorithm) {
// 	userCodeAlgorithm = algorithm
// }

func SetTracer(t tracing.Tracer) {
	tracer = t
}

func SetLogger(l logging.Logger) {
	logger = l
}

func SetLegacyEventstore(es eventstore.LegacyEventstore) {
	legacyEventstore = es
}

// func SetUserRepository(repo func(database.QueryExecutor) UserRepository) {
// 	userRepo = repo
// }

// func SetInstanceRepository(repo func(database.QueryExecutor) InstanceRepository) {
// 	instanceRepo = repo
// }

// func SetCryptoRepository(repo func(database.QueryExecutor) CryptoRepository) {
// 	cryptoRepo = repo
// }
