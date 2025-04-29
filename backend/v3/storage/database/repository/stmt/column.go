package stmt

import "fmt"

type Column[T any] interface {
	fmt.Stringer
	statementApplier[T]
	scanner(t *T) any
}

type columnDescriptor[T any] struct {
	name string
	scan func(*T) any
}

func (cd columnDescriptor[T]) scanner(t *T) any {
	return cd.scan(t)
}

// Apply implements [Column].
func (f columnDescriptor[T]) Apply(stmt *statement[T]) {
	stmt.builder.WriteString(stmt.columnPrefix())
	stmt.builder.WriteString(f.String())
}

// String implements [Column].
func (f columnDescriptor[T]) String() string {
	return f.name
}

var _ Column[any] = (*columnDescriptor[any])(nil)

type ignoreCaseColumnDescriptor[T any] struct {
	columnDescriptor[T]
	fieldNameSuffix string
}

func (f ignoreCaseColumnDescriptor[T]) ApplyIgnoreCase(stmt *statement[T]) {
	stmt.builder.WriteString(f.String())
	stmt.builder.WriteString(f.fieldNameSuffix)
}

var _ Column[any] = (*ignoreCaseColumnDescriptor[any])(nil)
