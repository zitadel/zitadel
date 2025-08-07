package database

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/zitadel/logging"
)

type Statement struct {
	addr    *Statement
	builder strings.Builder

	args []any
	// key is the name of the arg and value is the placeholder
	// TODO: condition must know if it's args are named parameters or not
	// namedArgs map[placeholder]string
}

func (stmt *Statement) Args() []any {
	if stmt == nil {
		return nil
	}
	return stmt.args
}

func (stmt *Statement) Reset() {
	stmt.builder.Reset()
	stmt.addr = nil
	stmt.args = nil
}

// TODO: condition must know if it's args are named parameters or not
// SetNamedArg sets the arg and makes it available for query construction
// func (stmt *Statement) SetNamedArg(name placeholder, value any) (placeholder string) {
// 	stmt.copyCheck()
// 	stmt.args = append(stmt.args, value)
// 	placeholder = fmt.Sprintf("$%d", len(stmt.args))
// 	if !strings.HasPrefix(name.string, "@") {
// 		name.string = "@" + name.string
// 	}
// 	stmt.namedArgs[name] = placeholder
// 	return placeholder
// }

// AppendArgs appends the args without writing it to Builder
// if any arg is a [placeholder] it's replaced with the placeholders parameter
func (stmt *Statement) AppendArgs(args ...any) {
	stmt.copyCheck()
	stmt.args = slices.Grow(stmt.args, len(args))
	for _, arg := range args {
		stmt.AppendArg(arg)
	}
}

// AppendArg appends the arg without writing it to Builder
// if the arg is a [placeholder] it's replaced with the placeholders parameter
func (stmt *Statement) AppendArg(arg any) int {
	stmt.copyCheck()

	// TODO: condition must know if it's args are named parameters or not
	// if namedArg, ok := arg.(sql.NamedArg); ok {
	// 	stmt.SetNamedArg(placeholder{namedArg.Name}, namedArg.Value)
	// 	return
	// }
	stmt.args = append(stmt.args, arg)
	return len(stmt.args)
}

// TODO: condition must know if it's args are named parameters or not
// func Placeholder(name string) placeholder {
// 	return placeholder{name}
// }

// TODO: condition must know if it's args are named parameters or not
// type placeholder struct {
// 	string
// }

// WriteArgs appends the args and adds the placeholders comma separated to [stmt.Builder]
// if any arg is a [placeholder] it's replaced with the placeholders parameter
func (stmt *Statement) WriteArgs(args ...any) {
	stmt.copyCheck()
	stmt.args = slices.Grow(stmt.args, len(args))
	for i, arg := range args {
		if i > 0 {
			stmt.WriteString(", ")
		}
		stmt.WriteArg(arg)
	}
}

// WriteArg appends the arg and adds the placeholder to [stmt.Builder]
// if the arg is a [placeholder] it's replaced with the placeholders parameter
func (stmt *Statement) WriteArg(arg any) {
	stmt.copyCheck()
	// TODO: condition must know if it's args are named parameters or not
	// if namedPlaceholder, ok := arg.(placeholder); ok {
	// 	stmt.writeNamedPlaceholder(namedPlaceholder)
	// 	return
	// }
	placeholder := stmt.AppendArg(arg)
	stmt.WriteString("$")
	stmt.WriteString(strconv.Itoa(placeholder))
}

// WriteString extends [strings.Builder.WriteString]
// it replaces named args with the previously provided named args
func (stmt *Statement) WriteString(s string) {
	// TODO: condition must know if it's args are named parameters or not
	// for name, placeholder := range stmt.namedArgs {
	// 	s = strings.ReplaceAll(s, name.string, placeholder)
	// }
	stmt.builder.WriteString(s)
}

// WriteRune extends [strings.Builder.WriteRune]
func (stmt *Statement) WriteRune(r rune) {
	// TODO: condition must know if it's args are named parameters or not
	// for name, placeholder := range stmt.namedArgs {
	// 	s = strings.ReplaceAll(s, name.string, placeholder)
	// }
	stmt.builder.WriteRune(r)
}

// WriteByte extends [strings.Builder.WriteByte]
func (stmt *Statement) WriteByte(b byte) {
	// TODO: condition must know if it's args are named parameters or not
	// for name, placeholder := range stmt.namedArgs {
	// 	s = strings.ReplaceAll(s, name.string, placeholder)
	// }
	err := stmt.builder.WriteByte(b)
	logging.OnError(err).Warn("unable to write bytes")
}

// Write extends [strings.Builder.Write]
// it replaces named args with the previously provided named args
func (stmt *Statement) Write(b []byte) {
	// TODO: condition must know if it's args are named parameters or not
	// for name, placeholder := range stmt.namedArgs {
	// 	bytes.ReplaceAll(b, []byte(name.string), []byte(placeholder))
	// }
	stmt.builder.Write(b)
}

// String builds the query and replaces placeholders starting with "@"
// with the corresponding named arg placeholder
func (stmt *Statement) String() string {
	return stmt.builder.String()
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

// TODO: condition must know if it's args are named parameters or not
// func (stmt *Statement) writeNamedPlaceholder(arg placeholder) {
// 	placeholder, ok := stmt.namedArgs[arg]
// 	if !ok {
// 		logging.WithFields("named_placeholder", arg).Fatal("named placeholder not defined")
// 	}
// 	stmt.Builder.WriteString(placeholder)
// }

// copyCheck allows uninitialized usage of stmt
func (stmt *Statement) copyCheck() {
	if stmt.addr == nil {
		// This hack works around a failing of Go's escape analysis
		// that was causing b to escape and be heap allocated.
		// See issue 23382.
		// TODO: once issue 7921 is fixed, this should be reverted to
		// just "stmt.addr = stmt".
		stmt.addr = (*Statement)(noescape(unsafe.Pointer(stmt)))
		// TODO: condition must know if it's args are named parameters or not
		// stmt.namedArgs = make(map[placeholder]string)
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
	//nolint: staticcheck
	return unsafe.Pointer(x ^ 0)
}
