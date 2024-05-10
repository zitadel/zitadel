package repository

import (
	"encoding/json"
	"unsafe"
)

type Field[T any] struct {
	addr *Field[T]

	didChange bool
	value     T
	previous  T
}

func (f *Field[T]) UnmarshalJSON(data []byte) error {
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	f.Set(value)
	return nil
}

func (f *Field[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.value)
}

func (f *Field[T]) Set(value T) {
	f.copyCheck()

	f.didChange = true
	f.previous = f.value
	f.value = value
}

func (f *Field[T]) Value() T {
	f.copyCheck()
	return f.value
}

func (f *Field[T]) DidChange() bool {
	f.copyCheck()
	return f.didChange
}

func (f *Field[T]) copyCheck() {
	if f.addr == nil {
		// This hack works around a failing of Go's escape analysis
		// that was causing b to escape and be heap allocated.
		// See issue 23382.
		// TODO: once issue 7921 is fixed, this should be reverted to
		// just "stmt.addr = stmt".
		f.addr = (*Field[T])(noescape(unsafe.Pointer(f)))
		// TODO: condition must know if it's args are named parameters or not
		// stmt.namedArgs = make(map[placeholder]string)
	} else if f.addr != f {
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
