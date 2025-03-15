package sql

import (
	"github.com/zitadel/zitadel/backend/storage/database"
)

type executor[C database.Executor] struct {
	client C
}

func Execute[C database.Executor](client C) *executor[C] {
	return &executor[C]{client: client}
}

type querier[C database.Querier] struct {
	client C
}

func Query[C database.Querier](client C) *querier[C] {
	return &querier[C]{client: client}
}
