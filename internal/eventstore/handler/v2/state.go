package handler

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"time"

	"github.com/jackc/pgconn"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type state struct {
	instanceID     string
	aggregateType  eventstore.AggregateType
	aggregateID    string
	eventTimestamp time.Time
	eventSequence  uint64
}

var (
	//go:embed state_get.sql
	currentStateStmt string
	//go:embed state_set.sql
	updateStateStmt string
	//go:embed state_lock.sql
	lockStateStmt string
)

func (h *Handler) currentState(ctx context.Context, tx *sql.Tx) (currentState *state, shouldSkip bool, err error) {
	currentState = &state{
		instanceID: authz.GetInstance(ctx).InstanceID(),
	}

	timestamp := new(sql.NullTime)
	aggregateType := new(sql.NullString)
	aggregateID := new(sql.NullString)
	sequence := new(sql.NullInt64)
	row := tx.QueryRowContext(ctx, currentStateStmt, currentState.instanceID, h.projection.Name())
	err = row.Scan(
		timestamp,
		aggregateType,
		aggregateID,
		sequence,
	)
	pgErr := new(pgconn.PgError)
	if errors.Is(err, sql.ErrNoRows) {
		err = h.lockState(ctx, tx, currentState.instanceID)
	} else if errors.As(err, &pgErr) {
		// error returned if the row is currently locked by another connection
		if pgErr.Code == "55P03" {
			return nil, true, nil
		}
	}
	if err != nil {
		logging.WithError(err).Debug("unable to query current state")
		return nil, false, err
	}

	currentState.eventTimestamp = timestamp.Time
	currentState.aggregateType = eventstore.AggregateType(aggregateType.String)
	currentState.aggregateID = aggregateID.String
	currentState.eventSequence = uint64(sequence.Int64)
	return currentState, false, nil
}

func (h *Handler) setState(ctx context.Context, tx *sql.Tx, updatedState *state) error {
	res, err := tx.ExecContext(ctx, updateStateStmt,
		h.projection.Name(),
		updatedState.instanceID,
		updatedState.eventTimestamp,
		updatedState.aggregateType,
		updatedState.aggregateID,
		updatedState.eventSequence,
	)
	if err != nil {
		logging.WithError(err).Debug("unable to update state")
		return err
	}
	if affected, err := res.RowsAffected(); affected == 0 {
		logging.OnError(err).Error("unable to check if states are updated")
		return errs.ThrowInternal(err, "V2-lpiK0", "unable to update state")
	}
	return nil
}

func (h *Handler) lockState(ctx context.Context, tx *sql.Tx, instanceID string) error {
	res, err := tx.ExecContext(ctx, lockStateStmt,
		h.projection.Name(),
		instanceID,
	)
	if err != nil {
		logging.WithError(err).Debug("unable to lock state")
		return err
	}
	if affected, err := res.RowsAffected(); affected == 0 {
		logging.OnError(err).Error("projection is already locked")
		return errs.ThrowInternal(err, "V2-lpiK0", "projection already locked")
	}
	return nil
}
