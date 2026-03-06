package output

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const describeKeyWidth = 28

// Describe writes a human-readable key:value representation of msg to stdout.
// Each field is shown on its own line; nested messages are indented.
// Timestamps are formatted as RFC3339. Unset (zero-value) fields are omitted.
func Describe(msg proto.Message) {
	describeMessage(os.Stdout, msg.ProtoReflect(), 0)
}

func describeMessage(w io.Writer, msg protoreflect.Message, depth int) {
	indent := strings.Repeat("  ", depth)
	msg.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		label := describeLabel(fd)
		switch {
		case fd.IsList():
			list := v.List()
			if list.Len() == 0 {
				return true
			}
			if fd.Kind() == protoreflect.MessageKind {
				fmt.Fprintf(w, "%s%s:\n", indent, label)
				for i := 0; i < list.Len(); i++ {
					fmt.Fprintf(w, "%s  -\n", indent)
					describeMessage(w, list.Get(i).Message(), depth+2)
				}
			} else {
				items := make([]string, list.Len())
				for i := 0; i < list.Len(); i++ {
					items[i] = describeScalar(fd, list.Get(i))
				}
				fmt.Fprintf(w, "%s%-*s %s\n", indent, describeKeyWidth, label+":", strings.Join(items, ", "))
			}
		case fd.Kind() == protoreflect.MessageKind:
			if fd.Message().FullName() == "google.protobuf.Timestamp" {
				secs := v.Message().Get(fd.Message().Fields().ByName("seconds")).Int()
				nanos := v.Message().Get(fd.Message().Fields().ByName("nanos")).Int()
				t := time.Unix(secs, nanos).UTC()
				fmt.Fprintf(w, "%s%-*s %s\n", indent, describeKeyWidth, label+":", t.Format(time.RFC3339))
			} else {
				fmt.Fprintf(w, "%s%s:\n", indent, label)
				describeMessage(w, v.Message(), depth+1)
			}
		default:
			s := describeScalar(fd, v)
			if s != "" {
				fmt.Fprintf(w, "%s%-*s %s\n", indent, describeKeyWidth, label+":", s)
			}
		}
		return true
	})
}

func describeLabel(fd protoreflect.FieldDescriptor) string {
	return strings.ToUpper(strings.ReplaceAll(string(fd.Name()), "_", " "))
}

func describeScalar(fd protoreflect.FieldDescriptor, v protoreflect.Value) string {
	if fd.Kind() == protoreflect.EnumKind {
		if enumVal := fd.Enum().Values().ByNumber(v.Enum()); enumVal != nil {
			return string(enumVal.Name())
		}
	}
	return fmt.Sprint(v.Interface())
}
