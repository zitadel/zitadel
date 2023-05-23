package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type LockoutPolicy struct {
	ID            string
	Sequence      uint64
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain.PolicyState

	MaxPasswordAttempts uint64
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
	LockoutColIsDefault = Column{
		name:  projection.LockoutPolicyIsDefaultCol,
		table: lockoutTable,
	}
	LockoutColState = Column{
		name:  projection.LockoutPolicyStateCol,
		table: lockoutTable,
	}
	LockoutPolicyOwnerRemoved = Column{
		name:  projection.LockoutPolicyOwnerRemovedCol,
		table: lockoutTable,
	}
)

func (q *Queries) LockoutPolicyByOrg(ctx context.Context, shouldTriggerBulk bool, orgID string, withOwnerRemoved bool) (_ *LockoutPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		projection.LockoutPolicyProjection.Trigger(ctx)
	}
	eq := sq.Eq{
		LockoutColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[LockoutPolicyOwnerRemoved.identifier()] = false
	}

	stmt, scan := prepareLockoutPolicyQuery(ctx, q.client)
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
		return nil, errors.ThrowInternal(err, "QUERY-SKR6X", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) DefaultLockoutPolicy(ctx context.Context) (_ *LockoutPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareLockoutPolicyQuery(ctx, q.client)
	query, args, err := stmt.Where(sq.Eq{
		LockoutColID.identifier():         authz.GetInstance(ctx).InstanceID(),
		LockoutColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}).
		OrderBy(LockoutColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-mN0Ci", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func prepareLockoutPolicyQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*LockoutPolicy, error)) {
	return sq.Select(
			LockoutColID.identifier(),
			LockoutColSequence.identifier(),
			LockoutColCreationDate.identifier(),
			LockoutColChangeDate.identifier(),
			LockoutColResourceOwner.identifier(),
			LockoutColShowFailures.identifier(),
			LockoutColMaxPasswordAttempts.identifier(),
			LockoutColIsDefault.identifier(),
			LockoutColState.identifier(),
		).
			From(lockoutTable.identifier() + db.Timetravel(call.Took(ctx))).
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
				&policy.IsDefault,
				&policy.State,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-63mtI", "Errors.PasswordComplexityPolicy.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-uulCZ", "Errors.Internal")
			}
			return policy, nil
		}
}
