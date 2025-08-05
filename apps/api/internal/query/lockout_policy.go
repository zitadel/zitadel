package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type LockoutPolicy struct {
	ID            string
	Sequence      uint64
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain.PolicyState

	MaxPasswordAttempts uint64
	MaxOTPAttempts      uint64
	ShowFailures        bool

	IsDefault bool
}

var (
	lockoutTable = table{
		name:          projection.LockoutPolicyTable,
		instanceIDCol: projection.LockoutPolicyInstanceIDCol,
	}
	LockoutColID = Column{
		name:  projection.LockoutPolicyIDCol,
		table: lockoutTable,
	}
	LockoutColInstanceID = Column{
		name:  projection.LockoutPolicyInstanceIDCol,
		table: lockoutTable,
	}
	LockoutColSequence = Column{
		name:  projection.LockoutPolicySequenceCol,
		table: lockoutTable,
	}
	LockoutColCreationDate = Column{
		name:  projection.LockoutPolicyCreationDateCol,
		table: lockoutTable,
	}
	LockoutColChangeDate = Column{
		name:  projection.LockoutPolicyChangeDateCol,
		table: lockoutTable,
	}
	LockoutColResourceOwner = Column{
		name:  projection.LockoutPolicyResourceOwnerCol,
		table: lockoutTable,
	}
	LockoutColShowFailures = Column{
		name:  projection.LockoutPolicyShowLockOutFailuresCol,
		table: lockoutTable,
	}
	LockoutColMaxPasswordAttempts = Column{
		name:  projection.LockoutPolicyMaxPasswordAttemptsCol,
		table: lockoutTable,
	}
	LockoutColMaxOTPAttempts = Column{
		name:  projection.LockoutPolicyMaxOTPAttemptsCol,
		table: lockoutTable,
	}
	LockoutColIsDefault = Column{
		name:  projection.LockoutPolicyIsDefaultCol,
		table: lockoutTable,
	}
	LockoutColState = Column{
		name:  projection.LockoutPolicyStateCol,
		table: lockoutTable,
	}
)

func (q *Queries) LockoutPolicyByOrg(ctx context.Context, shouldTriggerBulk bool, orgID string) (policy *LockoutPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerLockoutPolicyProjection")
		ctx, err = projection.LockoutPolicyProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}
	eq := sq.Eq{
		LockoutColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}

	stmt, scan := prepareLockoutPolicyQuery()
	query, args, err := stmt.Where(
		sq.And{
			eq,
			sq.Or{
				sq.Eq{LockoutColID.identifier(): orgID},
				sq.Eq{LockoutColID.identifier(): authz.GetInstance(ctx).InstanceID()},
			},
		}).
		OrderBy(LockoutColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-SKR6X", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		policy, err = scan(row)
		return err
	}, query, args...)
	return policy, err
}

func (q *Queries) DefaultLockoutPolicy(ctx context.Context) (policy *LockoutPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareLockoutPolicyQuery()
	query, args, err := stmt.Where(sq.Eq{
		LockoutColID.identifier():         authz.GetInstance(ctx).InstanceID(),
		LockoutColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}).
		OrderBy(LockoutColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-mN0Ci", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		policy, err = scan(row)
		return err
	}, query, args...)
	return policy, err
}

func prepareLockoutPolicyQuery() (sq.SelectBuilder, func(*sql.Row) (*LockoutPolicy, error)) {
	return sq.Select(
			LockoutColID.identifier(),
			LockoutColSequence.identifier(),
			LockoutColCreationDate.identifier(),
			LockoutColChangeDate.identifier(),
			LockoutColResourceOwner.identifier(),
			LockoutColShowFailures.identifier(),
			LockoutColMaxPasswordAttempts.identifier(),
			LockoutColMaxOTPAttempts.identifier(),
			LockoutColIsDefault.identifier(),
			LockoutColState.identifier(),
		).
			From(lockoutTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*LockoutPolicy, error) {
			policy := new(LockoutPolicy)
			err := row.Scan(
				&policy.ID,
				&policy.Sequence,
				&policy.CreationDate,
				&policy.ChangeDate,
				&policy.ResourceOwner,
				&policy.ShowFailures,
				&policy.MaxPasswordAttempts,
				&policy.MaxOTPAttempts,
				&policy.IsDefault,
				&policy.State,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-38pZnUemLP", "Errors.IAM.PasswordLockoutPolicy.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-PJURxRUoYG", "Errors.Internal")
			}
			return policy, nil
		}
}
