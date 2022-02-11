package sql

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/caos/zitadel/internal/config/types"
)

type Config struct {
	SQL types.SQL
}

func Start(client *sql.DB) *SQL {
	return &SQL{
		client: client,
	}
}
