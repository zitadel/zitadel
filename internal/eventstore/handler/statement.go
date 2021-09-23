package handler

import (
	"database/sql"
	"errors"

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

type Column struct {
	Name  string
	Value interface{}
}

func NewCol(name string, value interface{}) Column {
	return Column{
		Name:  name,
		Value: value,
	}
}

type Condition Column

func NewCond(name string, value interface{}) Condition {
	return Condition{
		Name:  name,
		Value: value,
	}
}
