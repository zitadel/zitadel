package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
)

type LockoutPolicy struct {
	ID            string
	Sequence      uint64
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string

	MaxPasswordAttempts uint64
	ShowFailures        bool

	IsDefault bool
}

var (
	lockoutTable = table{
		name: projection.LockoutPolicyTable,
	}
	LockoutColID = Column{
		name: projection.LockoutPolicyIDCol,
	}
	LockoutColSequence = Column{
		name: projection.LockoutPolicySequenceCol,
	}
	LockoutColCreationDate = Column{
		name: projection.LockoutPolicyCreationDateCol,
	}
	LockoutColChangeDate = Column{
		name: projection.LockoutPolicyChangeDateCol,
	}
	LockoutColResourceOwner = Column{
		name: projection.LockoutPolicyResourceOwnerCol,
	}
	LockoutColShowFailures = Column{
		name: projection.LockoutPolicyShowLockOutFailuresCol,
	}
	LockoutColMaxPasswordAttempts = Column{
		name: projection.LockoutPolicyMaxPasswordAttemptsCol,
	}
	LockoutColIsDefault = Column{
		name: projection.LockoutPolicyIsDefaultCol,
	}
)

func (q *Queries) LockoutPolicyByOrg(ctx context.Context, orgID string) (*LockoutPolicy, error) {
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
