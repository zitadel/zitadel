package repository

import "github.com/zitadel/zitadel/backend/v3/storage/database"

type repository struct {
	// builder database.StatementBuilder
	client database.QueryExecutor
}
