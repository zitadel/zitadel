package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	limitSettingsTable = table{
		name:          projection.LimitsProjectionTable,
		instanceIDCol: projection.LimitsColumnInstanceID,
	}
	LimitsColumnAggregateID = Column{
		name:  projection.LimitsColumnAggregateID,
		table: limitSettingsTable,
	}
	LimitsColumnCreationDate = Column{
		name:  projection.LimitsColumnCreationDate,
		table: limitSettingsTable,
	}
	LimitsColumnChangeDate = Column{
		name:  projection.LimitsColumnChangeDate,
		table: limitSettingsTable,
	}
	LimitsColumnResourceOwner = Column{
		name:  projection.LimitsColumnResourceOwner,
		table: limitSettingsTable,
	}
	LimitsColumnInstanceID = Column{
		name:  projection.LimitsColumnInstanceID,
		table: limitSettingsTable,
	}
	LimitsColumnSequence = Column{
		name:  projection.LimitsColumnSequence,
		table: limitSettingsTable,
	}
	LimitsColumnAuditLogRetention = Column{
		name:  projection.LimitsColumnAuditLogRetention,
		table: limitSettingsTable,
	}
)

type Limits struct {
	AggregateID   string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64

	AuditLogRetention *time.Duration
}

func (q *Queries) Limits(ctx context.Context, resourceOwner string) (limits *Limits, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareLimitsQuery(ctx, q.client)
	query, args, err := stmt.Where(sq.Eq{
		LimitsColumnInstanceID.identifier():    authz.GetInstance(ctx).InstanceID(),
		LimitsColumnResourceOwner.identifier(): resourceOwner,
	}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-jJe80", "Errors.Query.SQLStatment")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		limits, err = scan(row)
		return err
	}, query, args...)
	return limits, err
}

func prepareLimitsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*Limits, error)) {
	return sq.Select(
			LimitsColumnAggregateID.identifier(),
			LimitsColumnCreationDate.identifier(),
			LimitsColumnChangeDate.identifier(),
			LimitsColumnResourceOwner.identifier(),
			LimitsColumnSequence.identifier(),
			LimitsColumnAuditLogRetention.identifier(),
		).
			From(limitSettingsTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Limits, error) {
			var (
				limits            = new(Limits)
				auditLogRetention database.NullDuration
			)
			err := row.Scan(
				&limits.AggregateID,
				&limits.CreationDate,
				&limits.ChangeDate,
				&limits.ResourceOwner,
				&limits.Sequence,
				&auditLogRetention,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-GU1em", "Errors.Limits.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-00jgy", "Errors.Internal")
			}
			if auditLogRetention.Valid {
				limits.AuditLogRetention = &auditLogRetention.Duration
			}
			return limits, nil
		}
}
