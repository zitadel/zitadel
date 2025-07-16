package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

func wrapError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return database.NewNoRowFoundError(err)
	}
	var pgxErr *pgconn.PgError
	if !errors.As(err, &pgxErr) {
		return database.NewUnknownError(err)
	}
	switch pgxErr.Code {
	case "23514": // check_violation
		return database.NewCheckError(pgxErr.TableName, pgxErr.ConstraintName, pgxErr)
	case "23505": // unique_violation
		return database.NewUniqueError(pgxErr.TableName, pgxErr.ConstraintName, pgxErr)
	case "23503": // foreign_key_violation
		return database.NewForeignKeyError(pgxErr.TableName, pgxErr.ConstraintName, pgxErr)
	case "23502": // not_null_violation
		return database.NewNotNullError(pgxErr.TableName, pgxErr.ConstraintName, pgxErr)
	}
	return database.NewUnknownError(err)
}
