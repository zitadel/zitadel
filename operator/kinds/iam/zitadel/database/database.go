package database

import (
	"errors"
)

type Current struct {
	Host  string
	Port  string
	Users []string
}

func SetDatabaseInQueried(queried map[string]interface{}, current *Current) {
	queried["database"] = current
}

func GetDatabaseInQueried(queried map[string]interface{}) (*Current, error) {
	curr, ok := queried["database"].(*Current)
	if !ok {
		return nil, errors.New("database current not in supported format")
	}

	return curr, nil
}
