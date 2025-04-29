package repository

import (
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type query struct{ database.Querier }

func Query(querier database.Querier) *query {
	return &query{Querier: querier}
}

type executor struct{ database.Executor }

func Execute(exec database.Executor) *executor {
	return &executor{Executor: exec}
}
