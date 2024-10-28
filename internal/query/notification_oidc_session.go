package query

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type NotificationOIDCSession struct {
	domain.ObjectDetails

	SessionID            string
	UserID               string
	ClientID             string
	BackChannelLogoutURI string
}

var (
	//go:embed notification_oidc_session_by_session_id.sql
	notificationOIDCSessionsBySessionID string
)

func (q *Queries) NotificationOIDCSessions(ctx context.Context, sessionID string, triggerBulk bool) (out []NotificationOIDCSession, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	ctx, err = triggerNotificationOIDCSessionProjection(ctx, triggerBulk)
	if err != nil {
		return nil, err
	}

	err = q.client.QueryContext(ctx,
		func(rows *sql.Rows) error {
			for rows.Next() {
				var dst NotificationOIDCSession
				var backChannelLogoutURI sql.NullString
				err := rows.Scan(
					&dst.ID,
					&dst.CreationDate,
					&dst.EventDate,
					&dst.ResourceOwner,
					&dst.Sequence,
					&dst.SessionID,
					&dst.UserID,
					&dst.ClientID,
					&backChannelLogoutURI,
				)
				if err != nil {
					return err
				}
				dst.BackChannelLogoutURI = backChannelLogoutURI.String
				out = append(out, dst)
			}
			return rows.Err()
		},
		notificationOIDCSessionsBySessionID,
		authz.GetInstance(ctx).InstanceID(),
		sessionID,
	)

	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-nooZ2", "Errors.Internal")
	}
	return out, nil
}

func triggerNotificationOIDCSessionProjection(ctx context.Context, trigger bool) (_ context.Context, err error) {
	if trigger {
		ctx, span := tracing.NewSpan(ctx)
		defer func() { span.EndWithError(err) }()
		return projection.NotificationOIDCSessionProjection.Trigger(ctx, handler.WithAwaitRunning())
	}
	return ctx, nil
}
