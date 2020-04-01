package sql

import (
	// postgres dialect
	"database/sql"

	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/errors"
	_ "github.com/lib/pq"
)

type Config struct {
	types.SQL
}

func Start(conf Config) (*SQL, error) {
	client, err := sql.Open("postgres", conf.ConnectionString())
	if err != nil {
		return nil, errors.ThrowPreconditionFailed(err, "SQL-9qBtr", "unable to open database connection")
	}
	return &SQL{
		client: client,
	}, nil
}
