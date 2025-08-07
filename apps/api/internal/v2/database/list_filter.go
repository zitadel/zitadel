package database

import "github.com/zitadel/logging"

type ListFilter[V value] struct {
	comp listCompare
	list []V
}

func NewListEquals[V value](list ...V) *ListFilter[V] {
	return newListFilter[V](listEqual, list)
}

func NewListContains[V value](list ...V) *ListFilter[V] {
	return newListFilter[V](listContain, list)
}

func NewListNotContains[V value](list ...V) *ListFilter[V] {
	return newListFilter[V](listNotContain, list)
}

func newListFilter[V value](comp listCompare, list []V) *ListFilter[V] {
	return &ListFilter[V]{
		comp: comp,
		list: list,
	}
}

func (f ListFilter[V]) Write(stmt *Statement, columnName string) {
	if len(f.list) == 0 {
		logging.WithFields("column", columnName).Debug("skip list filter because no entries defined")
		return
	}
	if f.comp == listNotContain {
		stmt.WriteString("NOT(")
	}
	stmt.WriteString(columnName)
	stmt.WriteString(" = ")
	if f.comp != listEqual {
		stmt.WriteString("ANY(")
	}
	stmt.WriteArg(f.list)
	if f.comp != listEqual {
		stmt.WriteString(")")
	}
	if f.comp == listNotContain {
		stmt.WriteRune(')')
	}
}

type listCompare uint8

const (
	listEqual listCompare = iota
	listContain
	listNotContain
)
