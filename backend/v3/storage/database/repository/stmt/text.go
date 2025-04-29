package stmt

type Text interface {
	~string | ~[]byte
}

type TextCondition[V Text, T any] struct {
	condition[V, T, TextOperation]
}

func (tc *TextCondition[V, T]) Apply(stmt *statement[T]) {
	placeholder := stmt.appendArg(tc.value)

	switch tc.op {
	case TextOperationEqual, TextOperationNotEqual:
		tc.field.Apply(stmt)
		operation[T, TextOperation]{tc.op}.Apply(stmt)
		stmt.builder.WriteString(placeholder)
	case TextOperationEqualIgnoreCase:
		if desc, ok := tc.field.(ignoreCaseColumnDescriptor[T]); ok {
			desc.ApplyIgnoreCase(stmt)
		} else {
			stmt.builder.WriteString("LOWER(")
			tc.field.Apply(stmt)
			stmt.builder.WriteString(")")
		}
		operation[T, TextOperation]{tc.op}.Apply(stmt)
		stmt.builder.WriteString("LOWER(")
		stmt.builder.WriteString(placeholder)
		stmt.builder.WriteString(")")
	case TextOperationStartsWith:
		tc.field.Apply(stmt)
		operation[T, TextOperation]{tc.op}.Apply(stmt)
		stmt.builder.WriteString(placeholder)
		stmt.builder.WriteString("|| '%'")
	case TextOperationStartsWithIgnoreCase:
		if desc, ok := tc.field.(ignoreCaseColumnDescriptor[T]); ok {
			desc.ApplyIgnoreCase(stmt)
		} else {
			stmt.builder.WriteString("LOWER(")
			tc.field.Apply(stmt)
			stmt.builder.WriteString(")")
		}
		operation[T, TextOperation]{tc.op}.Apply(stmt)
		stmt.builder.WriteString("LOWER(")
		stmt.builder.WriteString(placeholder)
		stmt.builder.WriteString(")")
		stmt.builder.WriteString("|| '%'")
	}
}

type TextOperation uint8

const (
	TextOperationEqual TextOperation = iota + 1
	TextOperationEqualIgnoreCase
	TextOperationNotEqual
	TextOperationStartsWith
	TextOperationStartsWithIgnoreCase
)

var textOperations = map[TextOperation]string{
	TextOperationEqual:                " = ",
	TextOperationEqualIgnoreCase:      " = ",
	TextOperationNotEqual:             " <> ",
	TextOperationStartsWith:           " LIKE ",
	TextOperationStartsWithIgnoreCase: " LIKE ",
}

func (to TextOperation) String() string {
	return textOperations[to]
}
