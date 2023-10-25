package initialise

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/zitadel/zitadel/internal/database"
)

func exec(db *database.DB, stmt string, possibleErrCodes []string, args ...interface{}) error {
	_, err := db.Exec(stmt, args...)
	pgErr := new(pgconn.PgError)
	if errors.As(err, &pgErr) {
		for _, possibleCode := range possibleErrCodes {
			if possibleCode == pgErr.Code {
				return nil
			}
		}
	}
	return err
}
