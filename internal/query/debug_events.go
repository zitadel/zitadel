package query

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type DebugEventState struct {
	domain.ObjectDetails
	Blob string
}

var (
	//go:embed debug_events_state_by_id.sql
	debugEventsStateByIdQuery string
	//go:embed debug_events_states.sql
	debugEventsStatesQuery string
)

func (q *Queries) GetDebugEventsStateByID(ctx context.Context, id string, triggerBulk bool) (_ *DebugEventState, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	ctx, err = triggerDebugEventsProjection(ctx, triggerBulk)
	if err != nil {
		return nil, err
	}

	dst := new(DebugEventState)
	err = q.client.QueryRowContext(ctx,
		func(row *sql.Row) error {
			return row.Scan(
				&dst.ID,
				&dst.CreationDate,
				&dst.EventDate,
				&dst.ResourceOwner,
				&dst.Sequence,
				&dst.Blob,
			)
		},
		debugEventsStateByIdQuery,
		authz.GetInstance(ctx).InstanceID(),
		id,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, zerrors.ThrowNotFound(err, "QUERY-Eeth5", "debug event state not found")
	}
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-oe0Ae", "Errors.Internal")
	}
	return dst, err
}

func (q *Queries) ListDebugEventsStates(ctx context.Context, triggerBulk bool) (out []DebugEventState, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	ctx, err = triggerDebugEventsProjection(ctx, triggerBulk)
	if err != nil {
		return nil, err
	}

	err = q.client.QueryContext(ctx,
		func(rows *sql.Rows) error {
			for rows.Next() {
				var dst DebugEventState
				err := rows.Scan(
					&dst.ID,
					&dst.CreationDate,
					&dst.EventDate,
					&dst.ResourceOwner,
					&dst.Sequence,
					&dst.Blob,
				)
				if err != nil {
					return err
				}
				out = append(out, dst)
			}
			return rows.Err()
		},
		debugEventsStatesQuery,
		authz.GetInstance(ctx).InstanceID(),
	)

	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-nooZ2", "Errors.Internal")
	}
	return out, nil
}

func triggerDebugEventsProjection(ctx context.Context, trigger bool) (_ context.Context, err error) {
	if trigger {
		ctx, span := tracing.NewSpan(ctx)
		defer func() { span.EndWithError(err) }()
		return projection.DebugEventsProjection.Trigger(ctx, handler.WithAwaitRunning())
	}
	return ctx, nil
}
