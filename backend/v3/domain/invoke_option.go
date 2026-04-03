package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type InvokeOpt func(*InvokeOpts)

func WithOrganizationRepo(repo OrganizationRepository) InvokeOpt {
	return func(opts *InvokeOpts) {
		opts.organizationRepo = repo
	}
}

func WithOrganizationDomainRepo(repo OrganizationDomainRepository) InvokeOpt {
	return func(opts *InvokeOpts) {
		opts.organizationDomainRepo = repo
	}
}

func WithProjectRepo(repo ProjectRepository) InvokeOpt {
	return func(opts *InvokeOpts) {
		opts.projectRepo = repo
	}
}

func WithInstanceRepo(repo InstanceRepository) InvokeOpt {
	return func(opts *InvokeOpts) {
		opts.instanceRepo = repo
	}
}

func WithInstanceDomainRepo(repo InstanceDomainRepository) InvokeOpt {
	return func(opts *InvokeOpts) {
		opts.instanceDomainRepo = repo
	}
}

func WithSessionRepo(repo SessionRepository) InvokeOpt {
	return func(opts *InvokeOpts) {
		opts.sessionRepo = repo
	}
}

func WithUserRepo(repo UserRepository) InvokeOpt {
	return func(opts *InvokeOpts) {
		opts.userRepo = repo
	}
}

func WithLockoutSettingsRepo(repo LockoutSettingsRepository) InvokeOpt {
	return func(opts *InvokeOpts) {
		opts.lockoutSettingRepo = repo
	}
}

func WithSecretGeneratorSettingsRepo(repo SecretGeneratorSettingsRepository) InvokeOpt {
	return func(opts *InvokeOpts) {
		opts.secretGeneratorSettingsRepo = repo
	}
}

func WithPermissionChecker(checker PermissionChecker) InvokeOpt {
	return func(opts *InvokeOpts) {
		opts.Permissions = checker
	}
}

func WithIDPIntentRepo(repo IDPIntentRepository) InvokeOpt {
	return func(opts *InvokeOpts) {
		opts.idpIntentRepo = repo
	}
}

// WithQueryExecutor sets the database client to be used by the command.
// If not set, the default pool will be used.
// This is mainly used for testing.
func WithQueryExecutor(executor database.QueryExecutor) InvokeOpt {
	return func(opts *InvokeOpts) {
		opts.db = executor
	}
}

// WithLegacyEventstore sets the eventstore to be used by the command.
// If not set, the default one will be used.
// This is mainly used for testing.
func WithLegacyEventstore(es eventstore.LegacyEventstore) InvokeOpt {
	return func(opts *InvokeOpts) {
		opts.legacyEventstore = es
	}
}

// WithSessionTokenDecryptor sets the decryptor for session tokens used by the commands.
// If not set, the default one will be used.
// This is mainly used for testing
func WithSessionTokenDecryptor(decryptor SessionTokenDecryptor) InvokeOpt {
	return func(opts *InvokeOpts) {
		opts.sessionTokenDecryptor = decryptor
	}
}

// InvokeOpts are passed to each command
type InvokeOpts struct {
	// db is the database client.
	// [Executor]s MUST NOT access this field directly, use [InvokeOpts.DB] to access it.
	//
	// [Invoker]s may manipulate this field for example changing it to a transaction.
	// It's their responsibility to restore it after ending the transaction.
	db                     database.QueryExecutor
	legacyEventstore       eventstore.LegacyEventstore
	Invoker                Invoker
	Permissions            PermissionChecker
	sessionTokenDecryptor  SessionTokenDecryptor
	organizationRepo       OrganizationRepository
	organizationDomainRepo OrganizationDomainRepository
	projectRepo            ProjectRepository
	instanceRepo           InstanceRepository
	instanceDomainRepo     InstanceDomainRepository
	sessionRepo            SessionRepository
	userRepo               UserRepository
	idpIntentRepo          IDPIntentRepository

	// Settings repos
	lockoutSettingRepo          LockoutSettingsRepository
	secretGeneratorSettingsRepo SecretGeneratorSettingsRepository
}

func (o *InvokeOpts) DB() database.QueryExecutor {
	if o.db != nil {
		return o.db
	}
	o.db = pool
	return o.db
}

// StartTransactionFromDB returns a [database.Transaction] from the input [database.QueryExecutor].
// Optionally, the caller can pass [database.TransactionOptions] for a customised transaction type.
//
// If db doesn't implement [database.Beginner] an internal error is returned.
//
// If the transaction [database.Beginner.Begin] call fails, an internal error is returned.
//
// The caller is in charge of calling [database.Transaction.End], [database.Transaction.Commit]
// or [database.Transaction.Rollback] as they see fit.
func (o *InvokeOpts) StartTransactionFromDB(ctx context.Context, db database.QueryExecutor, opts *database.TransactionOptions) (database.Transaction, error) {
	beginner, ok := db.(database.Beginner)
	if !ok {
		return nil, zerrors.CreateZitadelError(zerrors.KindInternal, nil, "DOM-LqxZbk", "database doesn't implement database.Beginner", 1)
	}

	tx, txErr := beginner.Begin(ctx, opts)
	if txErr != nil {
		return nil, zerrors.CreateZitadelError(zerrors.KindInternal, txErr, "DOM-sAAd3V", "failed starting transaction", 1)
	}

	return tx, nil
}

// StartTransaction is the same as [domain.StartTransactionFromDB] but uses the DB provided by [InvokeOpts]
func (o *InvokeOpts) StartTransaction(ctx context.Context, opts *database.TransactionOptions) (database.Transaction, error) {
	return o.StartTransactionFromDB(ctx, o.DB(), opts)
}

func (o *InvokeOpts) LegacyEventstore() eventstore.LegacyEventstore {
	if o.legacyEventstore != nil {
		return o.legacyEventstore
	}
	o.legacyEventstore = legacyEventstore
	return o.legacyEventstore
}

func (o *InvokeOpts) Invoke(ctx context.Context, executor Executor) error {
	if o.Invoker == nil {
		return executor.Execute(ctx, o)
	}
	return o.Invoker.Invoke(ctx, executor, o)
}

func DefaultOpts(invoker Invoker) *InvokeOpts {
	if invoker == nil {
		invoker = &noopInvoker{}
	}
	return &InvokeOpts{
		Invoker:               invoker,
		Permissions:           &noopPermissionChecker{}, // prevent panics for now
		sessionTokenDecryptor: sessionTokenDecryptor,
	}
}
