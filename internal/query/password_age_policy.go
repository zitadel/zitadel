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
		name:          projection.PasswordAgeTable,
		instanceIDCol: projection.AgePolicyInstanceIDCol,
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
	PasswordAgeColInstanceID = Column{
		name:  projection.AgePolicyInstanceIDCol,
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
	PasswordAgeColOwnerRemoved = Column{
		name:  projection.AgePolicyOwnerRemovedCol,
		table: passwordAgeTable,
	}
)

func (q *Queries) PasswordAgePolicyByOrg(ctx context.Context, shouldTriggerBulk bool, orgID string, withOwnerRemoved bool) (_ *PasswordAgePolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		projection.PasswordAgeProjection.Trigger(ctx)
	}
	eq := sq.Eq{PasswordAgeColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	if !withOwnerRemoved {
		eq[PasswordAgeColOwnerRemoved.identifier()] = false
	}
	stmt, scan := preparePasswordAgePolicyQuery(ctx, q.client)
	query, args, err := stmt.Where(
		sq.And{
			eq,
			sq.Or{
				sq.Eq{PasswordAgeColID.identifier(): orgID},
				sq.Eq{PasswordAgeColID.identifier(): authz.GetInstance(ctx).InstanceID()},
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

func (q *Queries) DefaultPasswordAgePolicy(ctx context.Context, shouldTriggerBulk bool) (_ *PasswordAgePolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		projection.PasswordAgeProjection.Trigger(ctx)
	}

	stmt, scan := preparePasswordAgePolicyQuery(ctx, q.client)
	query, args, err := stmt.Where(sq.Eq{
		PasswordAgeColID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}).
		OrderBy(PasswordAgeColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-mN0Ci", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func preparePasswordAgePolicyQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*PasswordAgePolicy, error)) {
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
			From(passwordAgeTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
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
