package repository

import "github.com/zitadel/zitadel/backend/storage/database"

type executor struct {
	client database.Executor
}

func execute(client database.Executor) *executor {
	return &executor{client: client}
}

type querier struct {
	client database.Querier
}

func query(client database.Querier) *querier {
	return &querier{client: client}
}
