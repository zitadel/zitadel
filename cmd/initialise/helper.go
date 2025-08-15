package initialise

import (
	"context"
	"errors"
	"slices"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/zitadel/zitadel/internal/database"
)

func exec(ctx context.Context, db database.ContextExecuter, stmt string, possibleErrCodes []string, args ...any) error {
	_, err := db.ExecContext(ctx, stmt, args...)
	pgErr := new(pgconn.PgError)
	if errors.As(err, &pgErr) {
		if slices.Contains(possibleErrCodes, pgErr.Code) {
			return nil
		}
	}
	return err
}
