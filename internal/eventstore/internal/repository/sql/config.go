package sql

import (
	_ "github.com/lib/pq"

	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/errors"
)

type Config struct {
	SQL types.SQL
}

func Start(conf Config) (*SQL, error) {
	client, err := conf.SQL.Start()
	if err != nil {
		return nil, errors.ThrowPreconditionFailed(err, "SQL-9qBtr", "unable to open database connection")
	}
	return &SQL{
		client: client,
	}, nil
}
