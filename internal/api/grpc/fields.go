package grpc

import (
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// AllFieldsSet recusively checks if all values in a message
// have a non-zero value.
func AllFieldsSet(t testing.TB, msg protoreflect.Message, ignoreTypes ...protoreflect.FullName) {
	ignore := make(map[protoreflect.FullName]bool, len(ignoreTypes))
	for _, name := range ignoreTypes {
		ignore[name] = true
	}

	md := msg.Descriptor()
	name := md.FullName()
	if ignore[name] {
		return
	}

	fields := md.Fields()

	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		if !msg.Has(fd) {
			t.Errorf("not all fields set in %q, missing %q", name, fd.Name())
			continue
		}

		if fd.Kind() == protoreflect.MessageKind {
			if m, ok := msg.Get(fd).Interface().(protoreflect.Message); ok {
				AllFieldsSet(t, m, ignoreTypes...)
			}
		}
	}
}
