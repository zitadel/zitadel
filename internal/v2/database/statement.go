package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

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

func (stmt *Statement) WriteArgs(args ...any) {
	stmt.copyCheck()
	for i, arg := range args {
		stmt.WritePlaceholder(arg)
		if i < len(args)-1 {
			stmt.Builder.WriteString(", ")
		}
	}
}

func (stmt *Statement) WritePlaceholder(arg any) {
	if namedArg, ok := arg.(sql.NamedArg); ok {
		stmt.writeNamedPlaceholder(namedArg)
		return
	}
	stmt.args = append(stmt.args, arg)
	stmt.Builder.WriteString("$")
	stmt.Builder.WriteString(strconv.Itoa(len(stmt.args)))
}

// builds the query and replaces placeholders starting with "@"
// with the corresponding named arg placeholder
func (stmt *Statement) String() string {
	query := stmt.Builder.String()
	if len(stmt.namedArgs) == 0 {
		return query
	}
	return strings.NewReplacer(stmt.namedArgReplaces...).Replace(query)
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
