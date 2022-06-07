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

var (
	passwordAgeTable = table{
		name: projection.PasswordAgeTable,
	}
	PasswordAgeColID = Column{
		name:  projection.AgePolicyIDCol,
		table: passwordAgeTable,
	}
	PasswordAgeColSequence = Column{
		name:  projection.AgePolicySequenceCol,
		table: passwordAgeTable,
	}
	PasswordAgeColCreationDate = Column{
		name:  projection.AgePolicyCreationDateCol,
		table: passwordAgeTable,
	}
	PasswordAgeColChangeDate = Column{
		name:  projection.AgePolicyChangeDateCol,
		table: passwordAgeTable,
	}
	PasswordAgeColResourceOwner = Column{
		name:  projection.AgePolicyResourceOwnerCol,
		table: passwordAgeTable,
	}
	PasswordAgeColWarnDays = Column{
		name:  projection.AgePolicyExpireWarnDaysCol,
		table: passwordAgeTable,
	}
	PasswordAgeColMaxAge = Column{
		name:  projection.AgePolicyMaxAgeDaysCol,
		table: passwordAgeTable,
	}
	PasswordAgeColIsDefault = Column{
		name:  projection.AgePolicyIsDefaultCol,
		table: passwordAgeTable,
	}
	PasswordAgeColState = Column{
		name:  projection.AgePolicyStateCol,
		table: passwordAgeTable,
	}
)

func (q *Queries) PasswordAgePolicyByOrg(ctx context.Context, shouldRealTime bool, orgID string) (*PasswordAgePolicy, error) {
	if shouldRealTime {
		projection.PasswordAgeProjection.TriggerBulk(ctx)
	}
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
		return nil, errors.ThrowInternal(err, "QUERY-SKR6X", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) DefaultPasswordAgePolicy(ctx context.Context, shouldRealTime bool) (*PasswordAgePolicy, error) {
	if shouldRealTime {
		projection.PasswordAgeProjection.TriggerBulk(ctx)
	}
	stmt, scan := preparePasswordAgePolicyQuery()
	query, args, err := stmt.Where(sq.Eq{
		PasswordAgeColID.identifier(): q.iamID,
	}).
		OrderBy(PasswordAgeColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-mN0Ci", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

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
					return nil, errors.ThrowNotFound(err, "QUERY-63mtI", "Errors.Org.PasswordComplexity.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-uulCZ", "Errors.Internal")
			}
			return policy, nil
		}
}
