package projection

import (
	"time"
	"unsafe"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type projection struct {
	addr *projection

	reducers map[string]map[string]eventstore.ReduceEvent

	Position eventstore.GlobalPosition
}

func (p *projection) ShouldReduce(event *eventstore.StorageEvent) bool {
	return p.Position.IsLess(event.Position)
}

func (p *projection) reduce(event *eventstore.StorageEvent, reduce eventstore.ReduceEvent) error {
	p.copyCheck()

	if reduce == nil {
		return nil
	}

	if err := reduce(event); err != nil {
		return err
	}

	p.set(event)
	return nil
}

func (p *projection) set(event *eventstore.StorageEvent) {
	p.copyCheck()

	p.Position = event.Position
}

// TODO: condition must know if it's args are named parameters or not
// func (stmt *projection) writeNamedPlaceholder(arg placeholder) {
// 	placeholder, ok := stmt.namedArgs[arg]
// 	if !ok {
// 		logging.WithFields("named_placeholder", arg).Fatal("named placeholder not defined")
// 	}
// 	stmt.Builder.WriteString(placeholder)
// }

// copyCheck allows uninitialized usage of stmt
func (p *projection) copyCheck() {
	if p.addr == nil {
		// This hack works around a failing of Go's escape analysis
		// that was causing b to escape and be heap allocated.
		// See issue 23382.
		// TODO: once issue 7921 is fixed, this should be reverted to
		// just "stmt.addr = stmt".
		p.addr = (*projection)(noescape(unsafe.Pointer(p)))
		// TODO: condition must know if it's args are named parameters or not
		// stmt.namedArgs = make(map[placeholder]string)
	} else if p.addr != p {
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

type objectMetadata struct {
	addr *objectMetadata
	projection

	Sequence  uint32
	CreatedAt time.Time
	ChangedAt time.Time
}

func (om *objectMetadata) reduce(event *eventstore.StorageEvent, reduce eventstore.ReduceEvent) error {
	om.copyCheck()

	if reduce == nil {
		return nil
	}

	if err := reduce(event); err != nil {
		return err
	}

	om.set(event)
	return nil
}

func (om *objectMetadata) set(event *eventstore.StorageEvent) {
	om.copyCheck()
	om.projection.set(event)
	om.Position = event.Position
}

// TODO: condition must know if it's args are named parameters or not
// func (stmt *projection) writeNamedPlaceholder(arg placeholder) {
// 	placeholder, ok := stmt.namedArgs[arg]
// 	if !ok {
// 		logging.WithFields("named_placeholder", arg).Fatal("named placeholder not defined")
// 	}
// 	stmt.Builder.WriteString(placeholder)
// }

// copyCheck allows uninitialized usage of stmt
func (om *objectMetadata) copyCheck() {
	if om.addr == nil {
		// This hack works around a failing of Go's escape analysis
		// that was causing b to escape and be heap allocated.
		// See issue 23382.
		// TODO: once issue 7921 is fixed, this should be reverted to
		// just "stmt.addr = stmt".
		om.addr = (*objectMetadata)(noescape(unsafe.Pointer(om)))
		// TODO: condition must know if it's args are named parameters or not
		// stmt.namedArgs = make(map[placeholder]string)
	} else if om.addr != om {
		panic("statement: illegal use of non-zero Builder copied by value")
	}
}
