package repository

import "strings"

type statement struct {
	builder strings.Builder
	args    []any
}

func (s *statement) appendArg(arg any) (placeholder string) {
	s.args = append(s.args, arg)
	return "$" + string(len(s.args))
}

func (s *statement) appendArgs(args ...any) (placeholders []string) {
	placeholders = make([]string, len(args))
	for i, arg := range args {
		placeholders[i] = s.appendArg(arg)
	}
	return placeholders
}
