package setup

import (
	"embed"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 09/cockroach/index.sql
	//go:embed 09/postgres/index.sql
	stmts09 embed.FS
)

func New09(db *database.DB) *EventstoreIndexesNew {
	return &EventstoreIndexesNew{
		dbClient: db,
		name:     "09_optimise_indexes",
		step:     "09",
		fileName: "index.sql",
		stmts:    stmts09,
	}
}
