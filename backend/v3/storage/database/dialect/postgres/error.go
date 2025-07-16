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
	// 23514: check_violation - A value violates a CHECK constraint.
	case "23514":
		return database.NewCheckError(pgxErr.TableName, pgxErr.ConstraintName, pgxErr)
	// 23505: unique_violation - A value violates a UNIQUE constraint.
	case "23505":
		return database.NewUniqueError(pgxErr.TableName, pgxErr.ConstraintName, pgxErr)
	// 23503: foreign_key_violation - A value violates a foreign key constraint.
	case "23503":
		return database.NewForeignKeyError(pgxErr.TableName, pgxErr.ConstraintName, pgxErr)
	// 23502: not_null_violation - A value violates a NOT NULL constraint.
	case "23502":
		return database.NewNotNullError(pgxErr.TableName, pgxErr.ConstraintName, pgxErr)
	}
	return database.NewUnknownError(err)
}
