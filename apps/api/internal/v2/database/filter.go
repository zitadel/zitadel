package database

type Condition interface {
	Write(stmt *Statement, columnName string)
}

type Filter[C compare, V value] struct {
	comp  C
	value V
}

func (f Filter[C, V]) Write(stmt *Statement, columnName string) {
	prepareWrite(stmt, columnName, f.comp)
	stmt.WriteArg(f.value)
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
	// TODO: condition must know if it's args are named parameters or not
	// number | text | placeholder
}
