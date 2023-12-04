package setup

import (
	"context"
	"database/sql"
	_ "embed"
)

var (
	//go:embed 05.sql
	lastFailedStmts string
)

type LastFailed struct {
	dbClient *sql.DB
}

func (mig *LastFailed) Execute(ctx context.Context) error {
	_, err := mig.dbClient.ExecContext(ctx, lastFailedStmts)
	return err
}

func (mig *LastFailed) String() string {
	return "05_last_failed"
}

func (mig *LastFailed) ShouldSkip() bool {
	return false
}
