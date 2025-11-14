package postgres

import (
	"errors"
	"strings"

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
	if errors.Is(err, pgx.ErrTooManyRows) {
		return database.NewMultipleRowsFoundError(err)
	}

	var pgxErr *pgconn.PgError
	if errors.As(err, &pgxErr) {
		return wrapPgError(pgxErr)
	}

	// scany only exports its errors as strings
	if strings.HasPrefix(err.Error(), "scany: expected 1 row, got: ") {
		return database.NewMultipleRowsFoundError(err)
	}
	if strings.HasPrefix(err.Error(), "scany:") || strings.HasPrefix(err.Error(), "scanning:") {
		return database.NewScanError(err)
	}

	return database.NewUnknownError(err)
}

func wrapPgError(err *pgconn.PgError) error {
	switch err.Code {
	// 23514: check_violation - A value violates a CHECK constraint.
	case "23514":
		return database.NewCheckError(err.TableName, err.ConstraintName, err)
	// 23505: unique_violation - A value violates a UNIQUE constraint.
	case "23505":
		return database.NewUniqueError(err.TableName, err.ConstraintName, err)
	// 23503: foreign_key_violation - A value violates a foreign key constraint.
	case "23503":
		return database.NewForeignKeyError(err.TableName, err.ConstraintName, err)
	// 23502: not_null_violation - A value violates a NOT NULL constraint.
	case "23502":
		return database.NewNotNullError(err.TableName, err.ConstraintName, err)
	default:
		return database.NewUnknownError(err)
	}
}
