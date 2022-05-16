package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
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
		name: projection.LockoutPolicyTable,
	}
	LockoutColID = Column{
		name:  projection.LockoutPolicyIDCol,
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
)

func (q *Queries) LockoutPolicyByOrg(ctx context.Context, shouldRealTime bool, orgID string) (*LockoutPolicy, error) {
	if shouldRealTime {
		projection.LockoutPolicyProjection.TriggerBulk(ctx)
	}
	stmt, scan := prepareLockoutPolicyQuery()
	query, args, err := stmt.Where(
		sq.Or{
			sq.Eq{
				LockoutColID.identifier(): orgID,
			},
			sq.Eq{
				LockoutColID.identifier(): q.iamID,
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

func (q *Queries) DefaultLockoutPolicy(ctx context.Context) (*LockoutPolicy, error) {
	stmt, scan := prepareLockoutPolicyQuery()
	query, args, err := stmt.Where(sq.Eq{
		LockoutColID.identifier(): q.iamID,
	}).
		OrderBy(LockoutColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-mN0Ci", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
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
			LockoutColIsDefault.identifier(),
			LockoutColState.identifier(),
		).
			From(lockoutTable.identifier()).PlaceholderFormat(sq.Dollar),
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
