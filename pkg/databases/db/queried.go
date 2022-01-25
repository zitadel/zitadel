package db

import (
	"errors"

	"github.com/caos/orbos/pkg/tree"
)

const queriedName = "database"

var ErrNoCurrentState = errors.New("no current state for database found")

func ParseQueriedForDatabase(queried map[string]interface{}) (Client, error) {
	queriedDB, ok := queried[queriedName]
	if !ok {
		return nil, ErrNoCurrentState
	}
	currentDBTree, ok := queriedDB.(*tree.Tree)
	if !ok {
		return nil, errors.New("current state does not fullfil interface")
	}
	currentDB, ok := currentDBTree.Parsed.(Client)
	if !ok {
		return nil, errors.New("current state does not fullfil interface")
	}

	return currentDB, nil
}

func SetQueriedForDatabase(queried map[string]interface{}, databaseCurrent *tree.Tree) {
	queried[queriedName] = databaseCurrent
}
