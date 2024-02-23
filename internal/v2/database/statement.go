package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type Query struct {
	Columns  []Column
	From     From
	Join     Join
	Where    Where
	Limit    uint32
	Offset   uint32
	GroupBy  GroupBy
	Having   Having
	Ordering Ordering
}

type Column interface {
	Write(stmt *Statement)
	WriteIdentifier(stmt *Statement)
}

type column struct {
	Name string
	As   string
}

type function struct {
	Name   string
	Column Column
}

func (f *function) Write(stmt *Statement) {
	stmt.Builder.WriteString(f.Name)
	stmt.Builder.WriteRune('(')
	f.Column.WriteIdentifier(stmt)
	stmt.Builder.WriteRune(')')
}

type From interface {
	Write(stmt *Statement)
}

type Join interface {
	Write(stmt *Statement)
}

type Where interface {
	Write(stmt *Statement)
}

type where struct{}

type condition struct{}

type Ordering interface {
	Write(stmt *Statement)
}

type orderBy struct {
	Fields []*orderByField
}

func (ob *orderBy) Write(stmt *Statement) {
	for i, field := range ob.Fields {
		field.Write(stmt)
		if i < len(ob.Fields)-1 {
			stmt.Builder.WriteString(", ")
		}
	}
}

type orderByField struct {
	Identifier string
	Desc       bool
}

func (f *orderByField) Write(stmt *Statement) {
	stmt.Builder.WriteString(f.Identifier)
	if f.Desc {
		stmt.Builder.WriteString(" DESC")
	}
}

type GroupBy interface {
	Write(stmt *Statement)
}

type groupBy struct {
	Fields []*groupByField
}

func (gb *groupBy) Write(stmt *Statement) {
	for i, field := range gb.Fields {
		field.Write(stmt)
		if i < len(gb.Fields)-1 {
			stmt.Builder.WriteString(", ")
		}
	}
}

type groupByField struct {
	Identifier string
}

func (f *groupByField) Write(stmt *Statement) {
	stmt.Builder.WriteString(f.Identifier)
}

type Having interface {
	Write(stmt *Statement)
}

type from interface {
	~string | *Query
}

type Statement struct {
	addr *Statement
	strings.Builder

	args []any
	// key is the name of the arg and value the placeholder
	namedArgs        map[string]string
	namedArgReplaces []string
}

func (stmt *Statement) Args() []any {
	stmt.copyCheck()
	return stmt.args
}

func (stmt *Statement) Reset() {
	stmt.Builder.Reset()
	stmt.addr = nil
	stmt.args = nil
}

// AppendArgs appends the args without writing it to [stmt.Builder]
func (stmt *Statement) AppendArgs(args ...any) {
	stmt.copyCheck()
	for _, arg := range args {
		if namedArg, ok := arg.(sql.NamedArg); ok {
			stmt.appendNamedArg(namedArg)
			continue
		}
		stmt.args = append(stmt.args, arg)
	}
}

// AppendArgs appends the args and adds the placeholders comma separated to [stmt.Builder]
func (stmt *Statement) WriteArgs(args ...any) {
	stmt.copyCheck()
	for i, arg := range args {
		stmt.AppendArg(arg)
		if i < len(args)-1 {
			stmt.Builder.WriteString(", ")
		}
	}
}

// AppendArg appends the arg and adds the placeholder to [stmt.Builder]
func (stmt *Statement) AppendArg(arg any) {
	stmt.copyCheck()
	if namedArg, ok := arg.(sql.NamedArg); ok {
		stmt.writeNamedPlaceholder(namedArg)
		return
	}
	stmt.args = append(stmt.args, arg)
	stmt.Builder.WriteString("$")
	stmt.Builder.WriteString(strconv.Itoa(len(stmt.args)))
}

// String builds the query and replaces placeholders starting with "@"
// with the corresponding named arg placeholder
func (stmt *Statement) String() string {
	query := stmt.Builder.String()
	if len(stmt.namedArgs) == 0 {
		return query
	}
	return strings.NewReplacer(stmt.namedArgReplaces...).Replace(query)
}

// Debug builds the statement and replaces the placeholders with the parameters
func (stmt *Statement) Debug() string {
	query := stmt.String()

	for i := len(stmt.args) - 1; i >= 0; i-- {
		var argText string
		switch arg := stmt.args[i].(type) {
		case time.Time:
			argText = "'" + arg.Format("2006-01-02 15:04:05Z07:00") + "'"
		case string:
			argText = "'" + arg + "'"
		case []string:
			argText = "ARRAY["
			for i, a := range arg {
				if i > 0 {
					argText += ", "
				}
				argText += "'" + a + "'"
			}
			argText += "]"
		default:
			argText = fmt.Sprint(arg)
		}
		query = strings.ReplaceAll(query, "$"+strconv.Itoa(i+1), argText)
	}

	return query
}

func (stmt *Statement) writeNamedPlaceholder(arg sql.NamedArg) {
	placeholder, ok := stmt.namedArgs[arg.Name]
	if !ok {
		placeholder = stmt.appendNamedArg(arg)
	}
	stmt.Builder.WriteString(placeholder)
}

func (stmt *Statement) appendNamedArg(arg sql.NamedArg) (placeholder string) {
	stmt.args = append(stmt.args, arg.Value)
	placeholder = fmt.Sprintf("$%d", len(stmt.args))
	stmt.namedArgs[arg.Name] = placeholder
	if !strings.HasPrefix(arg.Name, "@") {
		arg.Name = fmt.Sprintf("@%s", arg.Name)
	}
	stmt.namedArgReplaces = append(stmt.namedArgReplaces, arg.Name, placeholder)

	return placeholder
}

// copyCheck allows uninitialized usage of stmt
func (stmt *Statement) copyCheck() {
	if stmt.addr == nil {
		// This hack works around a failing of Go's escape analysis
		// that was causing b to escape and be heap allocated.
		// See issue 23382.
		// TODO: once issue 7921 is fixed, this should be reverted to
		// just "stmt.addr = stmt".
		stmt.addr = (*Statement)(noescape(unsafe.Pointer(stmt)))
		stmt.namedArgs = make(map[string]string)
	} else if stmt.addr != stmt {
		panic("statement: illegal use of non-zero Builder copied by value")
	}
}

// noescape hides a pointer from escape analysis. It is the identity function
// but escape analysis doesn't think the output depends on the input.
// noescape is inlined and currently compiles down to zero instructions.
// USE CAREFULLY!
// This was copied from the runtime; see issues 23382 and 7921.
//
//go:nosplit
//go:nocheckptr
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}
