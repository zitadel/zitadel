package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	zitade_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

var (
	restrictionsTable = table{
		name:          projection.RestrictionsProjectionTable,
		instanceIDCol: projection.RestrictionsColumnInstanceID,
	}
	RestrictionsColumnAggregateID = Column{
		name:  projection.RestrictionsColumnAggregateID,
		table: restrictionsTable,
	}
	RestrictionsColumnCreationDate = Column{
		name:  projection.RestrictionsColumnCreationDate,
		table: restrictionsTable,
	}
	RestrictionsColumnChangeDate = Column{
		name:  projection.RestrictionsColumnChangeDate,
		table: restrictionsTable,
	}
	RestrictionsColumnResourceOwner = Column{
		name:  projection.RestrictionsColumnResourceOwner,
		table: restrictionsTable,
	}
	RestrictionsColumnInstanceID = Column{
		name:  projection.RestrictionsColumnInstanceID,
		table: restrictionsTable,
	}
	RestrictionsColumnSequence = Column{
		name:  projection.RestrictionsColumnSequence,
		table: restrictionsTable,
	}
	RestrictionsColumnDisallowPublicOrgRegistrations = Column{
		name:  projection.RestrictionsColumnDisallowPublicOrgRegistration,
		table: restrictionsTable,
	}
)

type Restrictions struct {
	AggregateID   string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64

	DisallowPublicOrgRegistration bool
}

func (q *Queries) GetInstanceRestrictions(ctx context.Context) (restrictions Restrictions, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareRestrictionsQuery(ctx, q.client)
	instanceID := authz.GetInstance(ctx).InstanceID()
	query, args, err := stmt.Where(sq.Eq{
		RestrictionsColumnInstanceID.identifier():    instanceID,
		RestrictionsColumnResourceOwner.identifier(): instanceID,
	}).ToSql()
	if err != nil {
		return restrictions, zitade_errors.ThrowInternal(err, "QUERY-XnLMQ", "Errors.Query.SQLStatment")
	}
	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		restrictions, err = scan(row)
		return err
	}, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		// not found is not an error
		err = nil
	}
	return restrictions, err
}

func prepareRestrictionsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (Restrictions, error)) {
	return sq.Select(
			RestrictionsColumnAggregateID.identifier(),
			RestrictionsColumnCreationDate.identifier(),
			RestrictionsColumnChangeDate.identifier(),
			RestrictionsColumnResourceOwner.identifier(),
			RestrictionsColumnSequence.identifier(),
			RestrictionsColumnDisallowPublicOrgRegistrations.identifier(),
		).
			From(restrictionsTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (restrictions Restrictions, err error) {
			return restrictions, row.Scan(
				&restrictions.AggregateID,
				&restrictions.CreationDate,
				&restrictions.ChangeDate,
				&restrictions.ResourceOwner,
				&restrictions.Sequence,
				&restrictions.DisallowPublicOrgRegistration,
			)
		}
}
