package runtime

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"

	// Register proto descriptors.
	_ "github.com/zitadel/zitadel/pkg/grpc/application/v2/applicationconnect"
	_ "github.com/zitadel/zitadel/pkg/grpc/user/v2/userconnect"
)

// FuzzApplyFlagValue throws random values at applyFlagValue across all registered
// field descriptors to find panics (type mismatches, missing kind handlers).
func FuzzApplyFlagValue(f *testing.F) {
	// Seed corpus: cover different field types, services, and edge cases.
	// -- String fields --
	f.Add("zitadel.user.v2.UserService/CreateUser", "username", "alice")
	f.Add("zitadel.user.v2.UserService/CreateUser", "organization_id", "12345")
	f.Add("zitadel.application.v2.ApplicationService/CreateApplication", "name", "my-app")
	// -- Boolean fields --
	f.Add("zitadel.user.v2.UserService/CreateUser", "is_verification_required", "true")
	f.Add("zitadel.user.v2.UserService/CreateUser", "is_verification_required", "false")
	// -- Enum fields --
	f.Add("zitadel.user.v2.UserService/CreateUser", "gender", "GENDER_MALE")
	f.Add("zitadel.user.v2.UserService/CreateUser", "gender", "INVALID_ENUM_VALUE")
	// -- Integer-like string values --
	f.Add("zitadel.user.v2.UserService/CreateUser", "username", "0")
	f.Add("zitadel.user.v2.UserService/CreateUser", "username", "-1")
	f.Add("zitadel.user.v2.UserService/CreateUser", "username", "999999999999999999")
	// -- Different services --
	f.Add("zitadel.session.v2.SessionService/CreateSession", "password", "s3cret!")
	f.Add("zitadel.org.v2.OrganizationService/CreateOrganization", "name", "acme-corp")
	f.Add("zitadel.project.v2.ProjectService/CreateProject", "name", "my-project")
	// -- Edge cases: empty & whitespace --
	f.Add("zitadel.user.v2.UserService/CreateUser", "username", "")
	f.Add("zitadel.user.v2.UserService/CreateUser", "username", "   ")
	f.Add("zitadel.user.v2.UserService/CreateUser", "username", "\t\n")
	// -- Edge cases: unicode & special chars --
	f.Add("zitadel.user.v2.UserService/CreateUser", "username", "用户名")
	f.Add("zitadel.user.v2.UserService/CreateUser", "username", "al'ice; DROP TABLE users;--")
	f.Add("zitadel.user.v2.UserService/CreateUser", "username", "alice\x00bob")
	f.Add("zitadel.user.v2.UserService/CreateUser", "username", "../../../etc/passwd")
	// -- Edge cases: long value --
	f.Add("zitadel.user.v2.UserService/CreateUser", "username", strings.Repeat("x", 10000))
	// -- Invalid method / field combos --
	f.Add("nonexistent.Service/Method", "field", "value")
	f.Add("zitadel.user.v2.UserService/CreateUser", "nonexistent_field", "value")

	f.Fuzz(func(t *testing.T, method, fieldName, value string) {
		_, reqDesc, err := resolveMethod(method)
		if err != nil {
			t.Skip() // invalid method, skip
		}

		// Try to find the field on the message
		fd := findFieldByName(reqDesc, protoreflect.Name(fieldName), 3)
		if fd == nil {
			t.Skip() // field not found, skip
		}

		msg := dynamicpb.NewMessage(fd.Parent().(protoreflect.MessageDescriptor))

		// Must not panic — errors are fine
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic on method=%q field=%q value=%q: %v", method, fieldName, value, r)
			}
		}()

		// Try applying as different types the runtime might see
		switch fd.Kind() {
		case protoreflect.StringKind:
			if fd.IsList() {
				vals := []string{value}
				applyFlagValue(msg, fd, &vals)
			} else {
				applyFlagValue(msg, fd, &value)
			}
		case protoreflect.BoolKind:
			b := value == "true"
			applyFlagValue(msg, fd, &b)
		case protoreflect.EnumKind:
			if fd.IsList() {
				vals := []string{value}
				applyFlagValue(msg, fd, &vals)
			} else {
				applyFlagValue(msg, fd, &value)
			}
		default:
			applyFlagValue(msg, fd, &value)
		}
	})
}

// findFieldByName recursively searches for a field by name in a message descriptor.
func findFieldByName(desc protoreflect.MessageDescriptor, name protoreflect.Name, depth int) protoreflect.FieldDescriptor {
	if depth <= 0 {
		return nil
	}
	for i := 0; i < desc.Fields().Len(); i++ {
		fd := desc.Fields().Get(i)
		if fd.Name() == name {
			return fd
		}
		if fd.Kind() == protoreflect.MessageKind && !fd.IsList() {
			if found := findFieldByName(fd.Message(), name, depth-1); found != nil {
				return found
			}
		}
	}
	return nil
}
