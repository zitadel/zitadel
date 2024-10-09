package eventstore

import (
	"context"
	_ "embed"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed snapshot_get.sql
	setSnapshotQuery string
	//go:embed snapshot_set.sql
	getSnapshotQuery string
)

func (es *Eventstore) SetSnapshot(ctx context.Context, snapshot *eventstore.SnapshotData) error {
	_, err := es.client.Pool.Exec(ctx, setSnapshotQuery,
		snapshot.InstanceID,
		snapshot.SnapshotType,
		snapshot.AggregateID,
		snapshot.Position,
		snapshot.ChangeDate,
		snapshot.Payload,
	)
	return err
}

func (es *Eventstore) GetSnapshot(ctx context.Context, instanceID string, typ eventstore.SnapshotType, aggregateID string) (*eventstore.SnapshotData, error) {
	rows, _ := es.client.Pool.Query(ctx, getSnapshotQuery, instanceID, typ, aggregateID)
	s, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[snapshotScanner])
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	return &eventstore.SnapshotData{
		SnapshotBase: eventstore.SnapshotBase{
			InstanceID:   instanceID,
			SnapshotType: typ,
			AggregateID:  aggregateID,
			Position:     s.Position,
			ChangeDate:   s.ChangeDate,
		},
		Payload: s.Payload,
	}, nil

}

type snapshotScanner struct {
	Position   float64   `db:"position"`
	ChangeDate time.Time `db:"change_date"`
	Payload    []byte    `db:"payload"`
}
