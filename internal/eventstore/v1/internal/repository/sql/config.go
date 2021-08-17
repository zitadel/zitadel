package sql

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/errors"
)

type Config struct {
	SQL types.SQL
}

func Start(conf Config) (*SQL, *sql.DB, error) {
	client, err := conf.SQL.Start()
	if err != nil {
		return nil, nil, errors.ThrowPreconditionFailed(err, "SQL-9qBtr", "unable to open database connection")
	}
	// as we open many sql clients we set the max
	// open cons deep. now 3(maxconn) * 8(clients) = max 24 conns per pod
	client.SetMaxOpenConns(5)
	client.SetConnMaxLifetime(5 * time.Minute)

	return &SQL{
		client: client,
	}, client, nil
}
