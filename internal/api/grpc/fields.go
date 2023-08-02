package grpc

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"
)

var CustomMappers = map[protoreflect.FullName]func(testing.TB, protoreflect.ProtoMessage) any{
	"google.protobuf.Struct": func(t testing.TB, msg protoreflect.ProtoMessage) any {
		e, ok := msg.(*structpb.Struct)
		require.True(t, ok)
		return e.AsMap()
	},
}

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

func AllFieldsEqual(t testing.TB, expected, actual protoreflect.Message, customMappers map[protoreflect.FullName]func(testing.TB, protoreflect.ProtoMessage) any) {
	md := expected.Descriptor()
	name := md.FullName()
	if mapper := customMappers[name]; mapper != nil {
		require.Equal(t, mapper(t, expected.Interface()), mapper(t, actual.Interface()))
		return
	}

	fields := md.Fields()

	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)

		if fd.Kind() == protoreflect.MessageKind {
			AllFieldsEqual(t, expected.Get(fd).Message(), actual.Get(fd).Message(), customMappers)
		} else {
			require.Equal(t, expected.Get(fd).Interface(), actual.Get(fd).Interface())
		}
	}
}
