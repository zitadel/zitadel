package handler

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
)

var (
	ErrNoProjection    = errors.New("no projection")
	ErrNoValues        = errors.New("no values")
	ErrNoCondition     = errors.New("no condition")
	ErrSomeStmtsFailed = errors.New("some statements failed")
)

type Statement struct {
	AggregateType    eventstore.AggregateType
	Sequence         uint64
	PreviousSequence uint64

	Execute func(ex Executer, projectionName string) error
}

func (s *Statement) IsNoop() bool {
	return s.Execute == nil
}

type Executer interface {
	Exec(string, ...interface{}) (sql.Result, error)
}

func NewCol(name string, value interface{}) Column {
	return Column{
		Name:  name,
		Value: value,
	}
}

func NewJSONCol(name string, value interface{}) Column {
	marshalled, err := json.Marshal(value)
	if err != nil {
		logging.LogWithFields("HANDL-oFvsl", "column", name).WithError(err).Panic("unable to marshal column")
	}

	return NewCol(name, marshalled)
}

type Column struct {
	Name  string
	Value interface{}
}
