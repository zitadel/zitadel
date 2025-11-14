package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
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
		name:          projection.DebugNotificationProviderTable,
		instanceIDCol: projection.DebugNotificationProviderInstanceIDCol,
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
	NotificationProviderColumnInstanceID = Column{
		name:  projection.DebugNotificationProviderInstanceIDCol,
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

func (q *Queries) NotificationProviderByIDAndType(ctx context.Context, aggID string, providerType domain.NotificationProviderType) (provider *DebugNotificationProvider, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareDebugNotificationProviderQuery()
	stmt, args, err := query.Where(
		sq.And{
			sq.Eq{NotificationProviderColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()},
			sq.Or{
				sq.Eq{
					NotificationProviderColumnAggID.identifier(): aggID,
					NotificationProviderColumnType.identifier():  providerType,
				},
			},
		}).
		Limit(1).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-f9jSf", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		provider, err = scan(row)
		return err
	}, stmt, args...)
	return provider, err
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
		).From(notificationProviderTable.identifier()).
			PlaceholderFormat(sq.Dollar),
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
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-s9ujf", "Errors.NotificationProvider.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-2liu0", "Errors.Internal")
			}
			return p, nil
		}
}
