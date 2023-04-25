package handler

import (
	"context"
	"database/sql"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
)

type CurrentState struct {
	LastRun      time.Time
	CreationDate time.Time
	Sequence     uint64
	Projection   string
}

func QueryCurrentState(ctx context.Context, db *sql.DB, projectionName string) (*CurrentState, error) {
	// TODO
	return nil, nil
}

type FailedEvent struct {
	Projection string

	AggregateType eventstore.AggregateType
	AggregateID   string
	Sequence      uint64

	Count      uint8
	Error      string
	LastFailed time.Time
}

func QueryFailedEvents(ctx context.Context, db *sql.DB) ([]*FailedEvent, error) {
	// TODO
	return nil, nil
}

func RemoveFailedEvent(ctx context.Context, db *sql.DB) error {
	// TODO
	return nil
}
