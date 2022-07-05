package initialise

import (
	"database/sql"
)

func exists(query string, args ...interface{}) func(*sql.DB) (exists bool, err error) {
	return func(db *sql.DB) (exists bool, err error) {
		row := db.QueryRow("SELECT EXISTS("+query+")", args...)
		err = row.Scan(&exists)
		return exists, err
	}
}

func exec(stmt string, args ...interface{}) func(*sql.DB) error {
	return func(db *sql.DB) error {
		_, err := db.Exec(stmt, args...)
		return err
	}
}

func verify(db *sql.DB, checkExists func(*sql.DB) (bool, error), create func(*sql.DB) error) error {
	exists, err := checkExists(db)
	if exists || err != nil {
		return err
	}
	return create(db)
}
