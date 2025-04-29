package stmt

import (
	"fmt"
	"strings"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type statement[T any] struct {
	builder strings.Builder
	client  database.QueryExecutor

	columns []Column[T]

	schema string
	table  string
	alias  string

	condition Condition[T]

	limit  uint32
	offset uint32
	// order by fieldname and sort direction false for asc true for desc
	// orderBy SortingColumns[C]
	args         []any
	existingArgs map[any]string
}

func (s *statement[T]) scanners(t *T) []any {
	scanners := make([]any, len(s.columns))
	for i, column := range s.columns {
		scanners[i] = column.scanner(t)
	}
	return scanners
}

func (s *statement[T]) query() string {
	s.builder.WriteString(`SELECT `)
	for i, column := range s.columns {
		if i > 0 {
			s.builder.WriteString(", ")
		}
		column.Apply(s)
	}
	s.builder.WriteString(` FROM `)
	s.builder.WriteString(s.schema)
	s.builder.WriteRune('.')
	s.builder.WriteString(s.table)
	if s.alias != "" {
		s.builder.WriteString(" AS ")
		s.builder.WriteString(s.alias)
	}

	s.builder.WriteString(` WHERE `)

	s.condition.Apply(s)

	if s.limit > 0 {
		s.builder.WriteString(` LIMIT `)
		s.builder.WriteString(s.appendArg(s.limit))
	}
	if s.offset > 0 {
		s.builder.WriteString(` OFFSET `)
		s.builder.WriteString(s.appendArg(s.offset))
	}

	return s.builder.String()
}

// func (s *statement[T]) Where(condition Condition[T]) *statement[T] {
// 	s.condition = condition
// 	return s
// }

// func (s *statement[T]) Limit(limit uint32) *statement[T] {
// 	s.limit = limit
// 	return s
// }

// func (s *statement[T]) Offset(offset uint32) *statement[T] {
// 	s.offset = offset
// 	return s
// }

func (s *statement[T]) columnPrefix() string {
	if s.alias != "" {
		return s.alias + "."
	}
	return s.schema + "." + s.table + "."
}

func (s *statement[T]) appendArg(arg any) string {
	if s.existingArgs == nil {
		s.existingArgs = make(map[any]string)
	}
	if existing, ok := s.existingArgs[arg]; ok {
		return existing
	}
	s.args = append(s.args, arg)
	placeholder := fmt.Sprintf("$%d", len(s.args))
	s.existingArgs[arg] = placeholder
	return placeholder
}
