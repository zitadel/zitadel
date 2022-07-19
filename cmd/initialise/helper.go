package initialise

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgconn"
)

func exec(db *sql.DB, stmt string, possibleErrCodes []string, args ...interface{}) error {
	// s, err := db.Prepare(stmt)
	// if err != nil {
	// 	return err
	// }
	// _, err = s.Exec(args...)
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
