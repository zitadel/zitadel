package postgres

import (
	"database/sql"
	"strconv"
	"strings"
	"unsafe"
)

type statement struct {
	addr    *statement
	builder strings.Builder
	args    []any
}

func (stmt *statement) reset() {
	stmt.builder.Reset()
	stmt.addr = nil
	stmt.args = nil
}

func (stmt *statement) appendArgs(args ...any) {
	stmt.copyCheck()
	stmt.args = append(stmt.args, args...)
}

func (stmt *statement) writeArgs(args ...any) {
	stmt.copyCheck()
	for i, arg := range args {
		stmt.args = append(stmt.args, arg)
		stmt.writePlaceholder(arg)
		if i < len(args)-1 {
			stmt.builder.WriteString(", ")
		}
	}
}

func (stmt *statement) writePlaceholder(arg any) {
	if namedArg, ok := arg.(sql.NamedArg); ok {
		stmt.builder.WriteRune('@')
		stmt.builder.WriteString(namedArg.Name)
		return
	}
	stmt.builder.WriteString("$")
	stmt.builder.WriteString(strconv.Itoa(len(stmt.args)))
}

// copyCheck allows uninitialized usage of stmt
func (stmt *statement) copyCheck() {
	if stmt.addr == nil {
		// This hack works around a failing of Go's escape analysis
		// that was causing b to escape and be heap allocated.
		// See issue 23382.
		// TODO: once issue 7921 is fixed, this should be reverted to
		// just "stmt.addr = stmt".
		stmt.addr = (*statement)(noescape(unsafe.Pointer(stmt)))
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
