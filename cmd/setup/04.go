package setup

import (
	"embed"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 04/cockroach/index.sql
	//go:embed 04/postgres/index.sql
	stmts04 embed.FS
)

func New04(db *database.DB) *EventstoreIndexesNew {
	return &EventstoreIndexesNew{
		dbClient: db,
		name:     "04_eventstore_indexes",
		step:     "04",
		fileName: "index.sql",
		stmts:    stmts04,
	}
}
