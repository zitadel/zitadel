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

type NotificationPolicy struct {
	ID            string
	Sequence      uint64
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain.PolicyState

	PasswordChange bool

	IsDefault bool
}

var (
	notificationPolicyTable = table{
		name:          projection.NotificationPolicyProjectionTable,
		instanceIDCol: projection.NotificationPolicyColumnInstanceID,
	}
	NotificationPolicyColID = Column{
		name:  projection.NotificationPolicyColumnID,
		table: notificationPolicyTable,
	}
	NotificationPolicyColSequence = Column{
		name:  projection.NotificationPolicyColumnSequence,
		table: notificationPolicyTable,
	}
	NotificationPolicyColCreationDate = Column{
		name:  projection.NotificationPolicyColumnCreationDate,
		table: notificationPolicyTable,
	}
	NotificationPolicyColChangeDate = Column{
		name:  projection.NotificationPolicyColumnChangeDate,
		table: notificationPolicyTable,
	}
	NotificationPolicyColResourceOwner = Column{
		name:  projection.NotificationPolicyColumnResourceOwner,
		table: notificationPolicyTable,
	}
	NotificationPolicyColInstanceID = Column{
		name:  projection.NotificationPolicyColumnInstanceID,
		table: notificationPolicyTable,
	}
	NotificationPolicyColPasswordChange = Column{
		name:  projection.NotificationPolicyColumnPasswordChange,
		table: notificationPolicyTable,
	}
	NotificationPolicyColIsDefault = Column{
		name:  projection.NotificationPolicyColumnIsDefault,
		table: notificationPolicyTable,
	}
	NotificationPolicyColState = Column{
		name:  projection.NotificationPolicyColumnStateCol,
		table: notificationPolicyTable,
	}
	NotificationPolicyColOwnerRemoved = Column{
		name:  projection.NotificationPolicyColumnOwnerRemoved,
		table: notificationPolicyTable,
	}
)

func (q *Queries) NotificationPolicyByOrg(ctx context.Context, shouldTriggerBulk bool, orgID string, withOwnerRemoved bool) (_ *NotificationPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		if err := projection.NotificationPolicyProjection.Trigger(ctx); err != nil {
			return nil, err
		}
	}
	eq := sq.Eq{NotificationPolicyColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	if !withOwnerRemoved {
		eq[NotificationPolicyColOwnerRemoved.identifier()] = false
	}
	stmt, scan := prepareNotificationPolicyQuery(ctx, q.client)
	query, args, err := stmt.Where(
		sq.And{
			eq,
			sq.Or{
				sq.Eq{NotificationPolicyColID.identifier(): orgID},
				sq.Eq{NotificationPolicyColID.identifier(): authz.GetInstance(ctx).InstanceID()},
			},
		}).
		OrderBy(NotificationPolicyColIsDefault.identifier()).Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Xuoapqm", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) DefaultNotificationPolicy(ctx context.Context, shouldTriggerBulk bool) (_ *NotificationPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		if err := projection.NotificationPolicyProjection.Trigger(ctx); err != nil {
			return nil, err
		}
	}

	stmt, scan := prepareNotificationPolicyQuery(ctx, q.client)
	query, args, err := stmt.Where(sq.Eq{
		NotificationPolicyColID.identifier():         authz.GetInstance(ctx).InstanceID(),
		NotificationPolicyColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}).
		OrderBy(NotificationPolicyColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-xlqp209", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func prepareNotificationPolicyQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*NotificationPolicy, error)) {
	return sq.Select(
			NotificationPolicyColID.identifier(),
			NotificationPolicyColSequence.identifier(),
			NotificationPolicyColCreationDate.identifier(),
			NotificationPolicyColChangeDate.identifier(),
			NotificationPolicyColResourceOwner.identifier(),
			NotificationPolicyColPasswordChange.identifier(),
			NotificationPolicyColIsDefault.identifier(),
			NotificationPolicyColState.identifier(),
		).
			From(notificationPolicyTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*NotificationPolicy, error) {
			policy := new(NotificationPolicy)
			err := row.Scan(
				&policy.ID,
				&policy.Sequence,
				&policy.CreationDate,
				&policy.ChangeDate,
				&policy.ResourceOwner,
				&policy.PasswordChange,
				&policy.IsDefault,
				&policy.State,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-x0so2p", "Errors.NotificationPolicy.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-Zixoooq", "Errors.Internal")
			}
			return policy, nil
		}
}
