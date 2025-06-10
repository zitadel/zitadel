package postgres

import (
	"github.com/jackc/pgx/v5"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ database.Rows = (*Rows)(nil)

type Rows struct{ pgx.Rows }
