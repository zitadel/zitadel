package domain

import (
	"math/rand/v2"
	"strconv"

	"github.com/zitadel/zitadel/backend/v3/storage/cache"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/crypto"
)

var (
	pool              database.Pool
	userCodeAlgorithm crypto.EncryptionAlgorithm
	tracer            tracing.Tracer

	// userRepo     func(database.QueryExecutor) UserRepository
	instanceRepo func(database.QueryExecutor) InstanceRepository
	cryptoRepo   func(database.QueryExecutor) CryptoRepository
	orgRepo      func(database.QueryExecutor) OrgRepository

	instanceCache cache.Cache[string, string, *Instance]

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

// func SetUserRepository(repo func(database.QueryExecutor) UserRepository) {
// 	userRepo = repo
// }

func SetInstanceRepository(repo func(database.QueryExecutor) InstanceRepository) {
	instanceRepo = repo
}

func SetCryptoRepository(repo func(database.QueryExecutor) CryptoRepository) {
	cryptoRepo = repo
}
