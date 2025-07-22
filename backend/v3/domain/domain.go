package domain

import "github.com/zitadel/zitadel/backend/v3/storage/database"

type DomainValidationType string

const (
	DomainValidationTypeDNS  DomainValidationType = "dns"
	DomainValidationTypeHTTP DomainValidationType = "http"
)

type domainColumns interface {
	// InstanceIDColumn returns the column for the instance id field.
	// `qualified` indicates if the column should be qualified with the table name.
	InstanceIDColumn(qualified bool) database.Column
	// DomainColumn returns the column for the domain field.
	// `qualified` indicates if the column should be qualified with the table name.
	DomainColumn(qualified bool) database.Column
	// IsVerifiedColumn returns the column for the is verified field.
	// `qualified` indicates if the column should be qualified with the table name.
	IsVerifiedColumn(qualified bool) database.Column
	// IsPrimaryColumn returns the column for the is primary field.
	// `qualified` indicates if the column should be qualified with the table name.
	IsPrimaryColumn(qualified bool) database.Column
	// ValidationTypeColumn returns the column for the verification type field.
	// `qualified` indicates if the column should be qualified with the table name.
	ValidationTypeColumn(qualified bool) database.Column
	// CreatedAtColumn returns the column for the created at field.
	// `qualified` indicates if the column should be qualified with the table name.
	CreatedAtColumn(qualified bool) database.Column
	// UpdatedAtColumn returns the column for the updated at field.
	// `qualified` indicates if the column should be qualified with the table name.
	UpdatedAtColumn(qualified bool) database.Column
}

type domainConditions interface {
	// InstanceIDCondition returns a filter on the instance id field.
	InstanceIDCondition(instanceID string) database.Condition
	// DomainCondition returns a filter on the domain field.
	DomainCondition(op database.TextOperation, domain string) database.Condition
	// IsPrimaryCondition returns a filter on the is primary field.
	IsPrimaryCondition(isPrimary bool) database.Condition
	// IsVerifiedCondition returns a filter on the is verified field.
	IsVerifiedCondition(isVerified bool) database.Condition
}

type domainChanges interface {
	// SetVerified sets the is verified column to true.
	SetVerified() database.Change
	// SetPrimary sets a domain as primary based on the condition.
	// All other domains will be set to non-primary.
	//
	// An error is returned if:
	// - The condition identifies multiple domains.
	// - The condition does not identify any domain.
	//
	// This is a no-op if:
	// - The domain is already primary.
	// - No domain matches the condition.
	SetPrimary() database.Change
	// SetValidationType sets the verification type column.
	// If the domain is already verified, this is a no-op.
	SetValidationType(verificationType DomainValidationType) database.Change
}

// import (
// 	"math/rand/v2"
// 	"strconv"

// 	"github.com/zitadel/zitadel/backend/v3/storage/cache"
// 	"github.com/zitadel/zitadel/backend/v3/storage/database"

// 	// "github.com/zitadel/zitadel/backend/v3/telemetry/logging"
// 	"github.com/zitadel/zitadel/backend/v3/telemetry/tracing"
// 	"github.com/zitadel/zitadel/internal/crypto"
// )

// // The variables could also be moved to a struct.
// // I just started with the singleton pattern and kept it like this.
// var (
// 	pool              database.Pool
// 	userCodeAlgorithm crypto.EncryptionAlgorithm
// 	tracer            tracing.Tracer
// 	// logger            logging.Logger

// 	userRepo func(database.QueryExecutor) UserRepository
// 	// instanceRepo func(database.QueryExecutor) InstanceRepository
// 	cryptoRepo func(database.QueryExecutor) CryptoRepository
// 	orgRepo    func(database.QueryExecutor) OrgRepository

// 	// instanceCache cache.Cache[instanceCacheIndex, string, *Instance]
// 	orgCache cache.Cache[orgCacheIndex, string, *Org]

// 	generateID func() (string, error) = func() (string, error) {
// 		return strconv.FormatUint(rand.Uint64(), 10), nil
// 	}
// )

// func SetPool(p database.Pool) {
// 	pool = p
// }

// func SetUserCodeAlgorithm(algorithm crypto.EncryptionAlgorithm) {
// 	userCodeAlgorithm = algorithm
// }

// func SetTracer(t tracing.Tracer) {
// 	tracer = t
// }

// // func SetLogger(l logging.Logger) {
// // 	logger = l
// // }

// func SetUserRepository(repo func(database.QueryExecutor) UserRepository) {
// 	userRepo = repo
// }

// func SetOrgRepository(repo func(database.QueryExecutor) OrgRepository) {
// 	orgRepo = repo
// }

// // func SetInstanceRepository(repo func(database.QueryExecutor) InstanceRepository) {
// // 	instanceRepo = repo
// // }

// func SetCryptoRepository(repo func(database.QueryExecutor) CryptoRepository) {
// 	cryptoRepo = repo
// }
