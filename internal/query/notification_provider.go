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

type DebugNotificationProvider struct {
	AggregateID   string
	CreationDate  time.Time
	ChangeDate    time.Time
	Sequence      uint64
	ResourceOwner string
	State         domain.NotificationProviderState
	Type          domain.NotificationProviderType
	Compact       bool
}

var (
	notificationProviderTable = table{
		name: projection.DebugNotificationProviderTable,
	}
	NotificationProviderColumnAggID = Column{
		name:  projection.DebugNotificationProviderAggIDCol,
		table: notificationProviderTable,
	}
	NotificationProviderColumnCreationDate = Column{
		name:  projection.DebugNotificationProviderCreationDateCol,
		table: notificationProviderTable,
	}
	NotificationProviderColumnChangeDate = Column{
		name:  projection.DebugNotificationProviderChangeDateCol,
		table: notificationProviderTable,
	}
	NotificationProviderColumnSequence = Column{
		name:  projection.DebugNotificationProviderSequenceCol,
		table: notificationProviderTable,
	}
	NotificationProviderColumnResourceOwner = Column{
		name:  projection.DebugNotificationProviderResourceOwnerCol,
		table: notificationProviderTable,
	}
	NotificationProviderColumnState = Column{
		name:  projection.DebugNotificationProviderStateCol,
		table: notificationProviderTable,
	}
	NotificationProviderColumnType = Column{
		name:  projection.DebugNotificationProviderTypeCol,
		table: notificationProviderTable,
	}
	NotificationProviderColumnCompact = Column{
		name:  projection.DebugNotificationProviderCompactCol,
		table: notificationProviderTable,
	}
)

func (q *Queries) NotificationProviderByIDAndType(ctx context.Context, aggID string, providerType domain.NotificationProviderType) (*DebugNotificationProvider, error) {
	query, scan := prepareDebugNotificationProviderQuery()
	stmt, args, err := query.Where(
		sq.Or{
			sq.Eq{
				LoginPolicyColumnOrgID.identifier(): aggID,
			},
		}).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-f9jSf", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	return scan(row)
}

func prepareDebugNotificationProviderQuery() (sq.SelectBuilder, func(*sql.Row) (*DebugNotificationProvider, error)) {
	return sq.Select(
			NotificationProviderColumnAggID.identifier(),
			NotificationProviderColumnCreationDate.identifier(),
			NotificationProviderColumnChangeDate.identifier(),
			NotificationProviderColumnSequence.identifier(),
			NotificationProviderColumnResourceOwner.identifier(),
			NotificationProviderColumnState.identifier(),
			NotificationProviderColumnType.identifier(),
			NotificationProviderColumnCompact.identifier(),
		).From(notificationProviderTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*DebugNotificationProvider, error) {
			p := new(DebugNotificationProvider)
			err := row.Scan(
				&p.AggregateID,
				&p.CreationDate,
				&p.ChangeDate,
				&p.Sequence,
				&p.ResourceOwner,
				&p.State,
				&p.Type,
				&p.Compact,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-s9ujf", "Errors.NotificationProvider.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-2liu0", "Errors.Internal")
			}
			return p, nil
		}
}
