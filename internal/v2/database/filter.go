package database

import "database/sql"

type Condition interface {
	Write(stmt *Statement, columnName string)
}

type Filter[C compare, V value] struct {
	comp  C
	value V
}

func (f Filter[C, V]) Write(stmt *Statement, columnName string) {
	prepareWrite(stmt, columnName, f.comp)
	stmt.AppendArg(f.value)
}

func (f Filter[C, V]) WriteNamed(stmt *Statement, columnName string) {
	prepareWrite(stmt, columnName, f.comp)
	stmt.AppendArg(sql.Named(columnName, f.value))
}

func prepareWrite[C compare](stmt *Statement, columnName string, comp C) {
	stmt.WriteString(columnName)
	stmt.WriteRune(' ')
	stmt.WriteString(comp.String())
	stmt.WriteRune(' ')
}

type compare interface {
	numberCompare | textCompare | listCompare
	String() string
}

type value interface {
	number | text
}
