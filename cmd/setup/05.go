package setup

import (
	"context"
	"database/sql"
)

const (
	adminAPI    = `ALTER TABLE adminapi.failed_events ADD COLUMN last_failed TIMESTAMPTZ;`
	auth        = `ALTER TABLE auth.failed_events ADD COLUMN last_failed TIMESTAMPTZ;`
	projections = `ALTER TABLE projections.failed_events ADD COLUMN last_failed TIMESTAMPTZ;`
)

type LastFailed struct {
	dbClient *sql.DB
}

func (mig *LastFailed) Execute(ctx context.Context) error {
	_, err := mig.dbClient.ExecContext(ctx, adminAPI+auth+projections)
	return err
}

func (mig *LastFailed) String() string {
	return "05_last_failed"
}
