package database

import "github.com/zitadel/logging"

type ListFilter[V value] struct {
	comp listCompare
	list []V
}

func NewListEquals[V value](list []V) *ListFilter[V] {
	return newListFilter[V](listEqual, list)
}

func NewListContains[V value](list []V) *ListFilter[V] {
	return newListFilter[V](listContain, list)
}

func newListFilter[V value](comp listCompare, list []V) *ListFilter[V] {
	return &ListFilter[V]{
		comp: comp,
		list: list,
	}
}

func (f ListFilter[V]) Write(stmt *Statement, columnName string) {
	stmt.Builder.WriteString(columnName)
	stmt.Builder.WriteRune(' ')
	stmt.Builder.WriteString(f.comp.String())
	stmt.Builder.WriteString(" ANY(")
	stmt.AppendArg(f.list)
	stmt.Builder.WriteString(")")
}

type listCompare uint8

const (
	listEqual listCompare = iota
	listContain
)

func (c listCompare) String() string {
	switch c {
	case listEqual, listContain:
		return "="
	default:
		logging.WithFields("compare", c).Panic("comparison type not implemented")
		return ""
	}
}
