package sql

import (
	"github.com/zitadel/zitadel/internal/database"
)

func Start(client *database.DB, allowOrderByCreationDate bool) *SQL {
	return &SQL{
		client:                   client,
		allowOrderByCreationDate: allowOrderByCreationDate,
	}
}
