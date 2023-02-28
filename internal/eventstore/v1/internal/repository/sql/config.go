package sql

import (
	"github.com/zitadel/zitadel/internal/database"
)

func Start(client *database.DB) *SQL {
	return &SQL{
		client: client,
	}
}
