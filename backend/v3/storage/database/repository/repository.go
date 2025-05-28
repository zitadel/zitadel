package repository

import "github.com/zitadel/zitadel/backend/v3/storage/database"

type repository struct {
	// we can't reuse builder after it's been used already, I think we should remove it
	builder database.StatementBuilder
	client  database.QueryExecutor
}
