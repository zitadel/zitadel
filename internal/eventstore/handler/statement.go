package handler

import (
	"database/sql"
	"errors"

	"github.com/caos/zitadel/internal/eventstore"
)

var (
	ErrNoAggregateType = errors.New("no aggregate")
	ErrNoTable         = errors.New("no table")
	ErrPrevSeqGtSeq    = errors.New("prev seq >= seq")
	ErrNoValues        = errors.New("no values")
	ErrNoCondition     = errors.New("no condition")
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

type Column struct {
	Name  string
	Value interface{}
}
