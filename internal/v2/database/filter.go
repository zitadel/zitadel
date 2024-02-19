package database

type Filter[C compare, V value] struct {
	comp  C
	value V
}

func (f *Filter[C, V]) Write(stmt *Statement, columnName string) {
	stmt.Builder.WriteString(columnName)
	stmt.Builder.WriteRune(' ')
	stmt.Builder.WriteString(f.comp.String())
	stmt.Builder.WriteRune(' ')
	stmt.AppendArg(f.value)
}

type compare interface {
	numberCompare | textCompare | listCompare
	String() string
}

type value interface {
	number | text
}
