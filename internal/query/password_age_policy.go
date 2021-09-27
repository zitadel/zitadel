package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
)

type PasswordAgePolicy struct {
	ID            string
	Sequence      uint64
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain.PolicyState

	ExpireWarnDays uint64
	MaxAgeDays     uint64

	IsDefault bool
}

func (q *Queries) MyPasswordAgePolicy(ctx context.Context, orgID string) (*PasswordAgePolicy, error) {
	stmt, scan := preparePasswordAgePolicyQuery()
	query, args, err := stmt.Where(
		sq.Or{
			sq.Eq{
				PasswordAgeColID.identifier(): orgID,
			},
			sq.Eq{
				PasswordAgeColID.identifier(): q.iamID,
			},
		}).
		OrderBy(PasswordAgeColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-SKR6X", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) DefaultPasswordAgePolicy(ctx context.Context) (*PasswordAgePolicy, error) {
	stmt, scan := preparePasswordAgePolicyQuery()
	query, args, err := stmt.Where(sq.Eq{
		PasswordAgeColID.identifier(): q.iamID,
	}).
		OrderBy(PasswordAgeColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-mN0Ci", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

var (
	passwordAgeTable = table{
		name: projection.PasswordAgeTable,
	}
	PasswordAgeColID = Column{
		name: projection.AgePolicyIDCol,
	}
	PasswordAgeColSequence = Column{
		name: projection.AgePolicySequenceCol,
	}
	PasswordAgeColCreationDate = Column{
		name: projection.AgePolicyCreationDateCol,
	}
	PasswordAgeColChangeDate = Column{
		name: projection.AgePolicyChangeDateCol,
	}
	PasswordAgeColResourceOwner = Column{
		name: projection.AgePolicyResourceOwnerCol,
	}
	PasswordAgeColWarnDays = Column{
		name: projection.AgePolicyExpireWarnDaysCol,
	}
	PasswordAgeColMaxAge = Column{
		name: projection.AgePolicyMaxAgeDaysCol,
	}
	PasswordAgeColIsDefault = Column{
		name: projection.AgePolicyIsDefaultCol,
	}
	PasswordAgeColState = Column{
		name: projection.AgePolicyStateCol,
	}
)

func preparePasswordAgePolicyQuery() (sq.SelectBuilder, func(*sql.Row) (*PasswordAgePolicy, error)) {
	return sq.Select(
			PasswordAgeColID.identifier(),
			PasswordAgeColSequence.identifier(),
			PasswordAgeColCreationDate.identifier(),
			PasswordAgeColChangeDate.identifier(),
			PasswordAgeColResourceOwner.identifier(),
			PasswordAgeColWarnDays.identifier(),
			PasswordAgeColMaxAge.identifier(),
			PasswordAgeColIsDefault.identifier(),
			PasswordAgeColState.identifier(),
		).
			From(passwordAgeTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*PasswordAgePolicy, error) {
			policy := new(PasswordAgePolicy)
			err := row.Scan(
				&policy.ID,
				&policy.Sequence,
				&policy.CreationDate,
				&policy.ChangeDate,
				&policy.ResourceOwner,
				&policy.ExpireWarnDays,
				&policy.MaxAgeDays,
				&policy.IsDefault,
				&policy.State,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-63mtI", "errors.policy.password.complexity.not_found")
				}
				return nil, errors.ThrowInternal(err, "QUERY-uulCZ", "errors.internal")
			}
			return policy, nil
		}
}
