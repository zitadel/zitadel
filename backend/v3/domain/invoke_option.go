package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/eventstore"
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

func WithIDPRepo(repo IDProviderRepository) InvokeOpt {
	return func(opts *InvokeOpts) {
		opts.idpRepo = repo
	}
}

func WithLockoutSettingsRepo(repo LockoutSettingsRepository) InvokeOpt {
	return func(opts *InvokeOpts) {
		opts.lockoutSettingRepo = repo
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

// InvokeOpts are passed to each command
type InvokeOpts struct {
	// db is the database client.
	// [Executor]s MUST NOT access this field directly, use [InvokeOpts.DB] to access it.
	//
	// [Invoker]s may manipulate this field for example changing it to a transaction.
	// Its their responsibility to restore it after ending the transaction.
	db                     database.QueryExecutor
	legacyEventstore       eventstore.LegacyEventstore
	Invoker                Invoker
	Permissions            PermissionChecker
	organizationRepo       OrganizationRepository
	organizationDomainRepo OrganizationDomainRepository
	projectRepo            ProjectRepository
	instanceRepo           InstanceRepository
	instanceDomainRepo     InstanceDomainRepository
	sessionRepo            SessionRepository
	userRepo               UserRepository
	idpRepo                IDProviderRepository

	// Settings repos
	lockoutSettingRepo LockoutSettingsRepository
}

func (o *InvokeOpts) DB() database.QueryExecutor {
	if o.db != nil {
		return o.db
	}
	o.db = pool
	return o.db
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
		Invoker:     invoker,
		Permissions: &noopPermissionChecker{}, // prevent panics for now
	}
}
