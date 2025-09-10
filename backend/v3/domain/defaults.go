package domain

import (
	"log/slog"
	"math/rand/v2"
	"strconv"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	// "github.com/zitadel/zitadel/backend/v3/telemetry/logging"
	"github.com/zitadel/zitadel/backend/v3/telemetry/logging"
	"github.com/zitadel/zitadel/backend/v3/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
)

// The variables could also be moved to a struct.
// I just started with the singleton pattern and kept it like this.
var (
	pool              database.Pool
	userCodeAlgorithm crypto.EncryptionAlgorithm
	tracer            tracing.Tracer
	logger            logging.Logger = *logging.NewLogger(slog.Default())
	legacyEventstore  *eventstore.Eventstore

	// userRepo func(database.QueryExecutor) UserRepository
	// instanceRepo func(database.QueryExecutor) InstanceRepository
	// cryptoRepo func(database.QueryExecutor) CryptoRepository
	// orgRepo    func(database.QueryExecutor) OrgRepository

	// instanceCache cache.Cache[instanceCacheIndex, string, *Instance]
	// orgCache cache.Cache[orgCacheIndex, string, *Org]

	generateID func() (string, error) = func() (string, error) {
		return strconv.FormatUint(rand.Uint64(), 10), nil
	}
)

func SetPool(p database.Pool) {
	pool = p
}

func SetUserCodeAlgorithm(algorithm crypto.EncryptionAlgorithm) {
	userCodeAlgorithm = algorithm
}

func SetTracer(t tracing.Tracer) {
	tracer = t
}

func SetLogger(l logging.Logger) {
	logger = l
}

func SetLegacyEventstore(es *eventstore.Eventstore) {
	legacyEventstore = es
}

// func SetUserRepository(repo func(database.QueryExecutor) UserRepository) {
// 	userRepo = repo
// }

// func SetOrgRepository(repo func(database.QueryExecutor) OrgRepository) {
// 	orgRepo = repo
// }

// func SetInstanceRepository(repo func(database.QueryExecutor) InstanceRepository) {
// 	instanceRepo = repo
// }

// func SetCryptoRepository(repo func(database.QueryExecutor) CryptoRepository) {
// 	cryptoRepo = repo
// }
