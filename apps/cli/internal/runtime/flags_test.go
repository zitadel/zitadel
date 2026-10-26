package runtime

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"

	// Blank imports to register proto descriptors in protoregistry.GlobalFiles.
	_ "github.com/zitadel/zitadel/pkg/grpc/application/v2/applicationconnect"
	_ "github.com/zitadel/zitadel/pkg/grpc/user/v2/userconnect"
)

// helper: resolve a method and return the request message descriptor.
func mustResolveReq(t *testing.T, fullMethod string) protoreflect.MessageDescriptor {
	t.Helper()
	_, reqDesc, err := resolveMethod(fullMethod)
	if err != nil {
		t.Fatalf("resolveMethod(%q): %v", fullMethod, err)
	}
	return reqDesc
}

// helper: register flags for a method and return the flagset + values + aliases.
func setupFlags(t *testing.T, fullMethod string) (*pflag.FlagSet, map[string]interface{}, map[string][]string) {
	t.Helper()
	reqDesc := mustResolveReq(t, fullMethod)
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	values := make(map[string]interface{})
	aliases := make(map[string][]string)
	registerFlags(flags, reqDesc, "", values, nil, 0, true, aliases)
	return flags, values, aliases
}

// helper: create a cobra.Command wrapping the given flagset for applyFlags.
func cobraCmd(flags *pflag.FlagSet) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Flags().AddFlagSet(flags)
	return cmd
}

// --- Recursive flag expansion ---

func TestRecursiveExpansion(t *testing.T) {
	// CreateUser.Human has profile.given_name etc. which should be recursively expanded.
	// The Human message is inside CreateUserRequest as a oneof, so we resolve the
	// Human sub-message directly via its own descriptor.
	reqDesc := mustResolveReq(t, "zitadel.user.v2.UserService/CreateUser")

	// The request has a oneof user_type with Human/Machine. Find Human's message.
	humanField := reqDesc.Fields().ByName("human")
	if humanField == nil {
		t.Fatal("expected 'human' field on CreateUserRequest")
	}
	humanDesc := humanField.Message()

	// Register flags for the Human sub-message — this is what buildOneofCommand does.
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	values := make(map[string]interface{})
	aliases := make(map[string][]string)
	registerFlags(flags, humanDesc, "", values, nil, 0, true, aliases)

	// Should have expanded profile.given_name → --profile-given-name
	if flags.Lookup("profile-given-name") == nil {
		t.Error("expected --profile-given-name flag from recursive expansion")
	}
	if flags.Lookup("profile-family-name") == nil {
		t.Error("expected --profile-family-name flag from recursive expansion")
	}
}

// --- Deduplication ---

func TestDeduplication(t *testing.T) {
	reqDesc := mustResolveReq(t, "zitadel.user.v2.UserService/CreateUser")
	humanField := reqDesc.Fields().ByName("human")
	humanDesc := humanField.Message()

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	values := make(map[string]interface{})
	registerFlags(flags, humanDesc, "", values, nil, 0, true, nil)

	// email.email should be deduplicated to just --email, not --email-email
	if flags.Lookup("email") == nil {
		t.Error("expected --email flag (deduplicated from email.email)")
	}
	if flags.Lookup("email-email") != nil {
		t.Error("should NOT have --email-email (should be deduplicated)")
	}
}

// --- Alias registration ---

func TestAliasRegistration(t *testing.T) {
	reqDesc := mustResolveReq(t, "zitadel.user.v2.UserService/CreateUser")
	humanField := reqDesc.Fields().ByName("human")
	humanDesc := humanField.Message()

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	values := make(map[string]interface{})
	aliases := make(map[string][]string)
	registerFlags(flags, humanDesc, "", values, nil, 0, true, aliases)

	// "given-name" should map to ["profile-given-name"] in the alias tracker
	canonicals, ok := aliases["given-name"]
	if !ok {
		t.Fatal("expected 'given-name' in alias tracker")
	}
	if len(canonicals) != 1 || canonicals[0] != "profile-given-name" {
		t.Errorf("expected alias given-name → profile-given-name, got %v", canonicals)
	}
}

// --- Repeated enum as StringSlice ---

func TestRepeatedEnumRegistration(t *testing.T) {
	// CreateApplication's OIDCConfiguration has response_types (repeated enum).
	reqDesc := mustResolveReq(t, "zitadel.application.v2.ApplicationService/CreateApplication")
	oidcField := reqDesc.Fields().ByName("oidc_configuration")
	if oidcField == nil {
		t.Fatal("expected 'oidc_configuration' field")
	}
	oidcDesc := oidcField.Message()

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	values := make(map[string]interface{})
	registerFlags(flags, oidcDesc, "", values, nil, 0, true, nil)

	// response_types is a repeated enum — should be registered as StringSlice, not String.
	rtFlag := flags.Lookup("response-types")
	if rtFlag == nil {
		t.Fatal("expected --response-types flag")
	}
	if rtFlag.Value.Type() != "stringSlice" {
		t.Errorf("expected --response-types type=stringSlice, got %s", rtFlag.Value.Type())
	}

	// grant_types likewise
	gtFlag := flags.Lookup("grant-types")
	if gtFlag == nil {
		t.Fatal("expected --grant-types flag")
	}
	if gtFlag.Value.Type() != "stringSlice" {
		t.Errorf("expected --grant-types type=stringSlice, got %s", gtFlag.Value.Type())
	}
}

// --- Depth limit ---

func TestDepthLimit(t *testing.T) {
	// At maxFlagDepth=3, deeply nested fields beyond that should not appear.
	// We just verify that the function doesn't panic and produces some flags.
	reqDesc := mustResolveReq(t, "zitadel.user.v2.UserService/CreateUser")
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	values := make(map[string]interface{})
	registerFlags(flags, reqDesc, "", values, nil, 0, true, nil)

	// Should have at least some flags (organization_id, username, etc.)
	count := 0
	flags.VisitAll(func(_ *pflag.Flag) { count++ })
	if count == 0 {
		t.Error("expected some flags to be registered")
	}
}

// --- Required field markers ---

func TestRequiredMarkers(t *testing.T) {
	reqDesc := mustResolveReq(t, "zitadel.user.v2.UserService/CreateUser")
	humanField := reqDesc.Fields().ByName("human")
	humanDesc := humanField.Message()

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	values := make(map[string]interface{})
	registerFlags(flags, humanDesc, "", values, nil, 0, true, nil)

	// given_name has (google.api.field_behavior) = REQUIRED on SetHumanProfile,
	// and profile has REQUIRED on Human → should show [REQUIRED].
	givenName := flags.Lookup("profile-given-name")
	if givenName == nil {
		t.Fatal("expected --profile-given-name flag")
	}
	if !strings.Contains(givenName.Usage, "[REQUIRED]") {
		t.Errorf("expected [REQUIRED] in given-name usage, got: %s", givenName.Usage)
	}

	// email has REQUIRED on Human → should show [REQUIRED].
	email := flags.Lookup("email")
	if email == nil {
		t.Fatal("expected --email flag")
	}
	if !strings.Contains(email.Usage, "[REQUIRED]") {
		t.Errorf("expected [REQUIRED] in email usage, got: %s", email.Usage)
	}
}

// --- Required field propagation (no inheritance for optional parents) ---

func TestRequiredPropagation(t *testing.T) {
	reqDesc := mustResolveReq(t, "zitadel.user.v2.UserService/CreateUser")
	humanField := reqDesc.Fields().ByName("human")
	humanDesc := humanField.Message()

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	values := make(map[string]interface{})
	registerFlags(flags, humanDesc, "", values, nil, 0, true, nil)

	// phone is optional on Human, so phone (the string field inside SetHumanPhone)
	// should NOT show [REQUIRED] even though phone.phone has REQUIRED annotation
	// in the SetPhoneRequest context.
	phone := flags.Lookup("phone")
	if phone != nil && strings.Contains(phone.Usage, "[REQUIRED]") {
		t.Errorf("phone should NOT be [REQUIRED] (parent is optional), got: %s", phone.Usage)
	}

	// nick_name is optional on SetHumanProfile, should NOT show [REQUIRED]
	nickName := flags.Lookup("profile-nick-name")
	if nickName != nil && strings.Contains(nickName.Usage, "[REQUIRED]") {
		t.Errorf("nick-name should NOT be [REQUIRED], got: %s", nickName.Usage)
	}
}

// --- Apply single enum ---

func TestApplyEnum(t *testing.T) {
	reqDesc := mustResolveReq(t, "zitadel.application.v2.ApplicationService/CreateApplication")
	oidcField := reqDesc.Fields().ByName("oidc_configuration")
	oidcDesc := oidcField.Message()
	oidcMsg := dynamicpb.NewMessage(oidcDesc)

	// auth_method_type is a singular enum
	fd := oidcDesc.Fields().ByName("auth_method_type")
	if fd == nil {
		t.Fatal("expected auth_method_type field")
	}

	val := "OIDC_AUTH_METHOD_TYPE_NONE"
	applyFlagValue(oidcMsg, fd, &val)

	got := oidcMsg.Get(fd).Enum()
	// OIDC_AUTH_METHOD_TYPE_NONE = 2
	expected := fd.Enum().Values().ByName("OIDC_AUTH_METHOD_TYPE_NONE")
	if expected == nil {
		t.Fatal("expected to find OIDC_AUTH_METHOD_TYPE_NONE enum value")
	}
	if got != expected.Number() {
		t.Errorf("expected enum number %d, got %d", expected.Number(), got)
	}
}

// --- Apply repeated enum ---

func TestApplyRepeatedEnum(t *testing.T) {
	reqDesc := mustResolveReq(t, "zitadel.application.v2.ApplicationService/CreateApplication")
	oidcField := reqDesc.Fields().ByName("oidc_configuration")
	oidcDesc := oidcField.Message()
	oidcMsg := dynamicpb.NewMessage(oidcDesc)

	fd := oidcDesc.Fields().ByName("response_types")
	if fd == nil {
		t.Fatal("expected response_types field")
	}

	vals := []string{"OIDC_RESPONSE_TYPE_CODE", "OIDC_RESPONSE_TYPE_ID_TOKEN"}
	applyFlagValue(oidcMsg, fd, &vals)

	list := oidcMsg.Get(fd).List()
	if list.Len() != 2 {
		t.Fatalf("expected 2 enum values, got %d", list.Len())
	}
	// OIDC_RESPONSE_TYPE_CODE = 1
	if list.Get(0).Enum() != 1 {
		t.Errorf("expected first value = 1 (CODE), got %d", list.Get(0).Enum())
	}
	// OIDC_RESPONSE_TYPE_ID_TOKEN = 2
	if list.Get(1).Enum() != 2 {
		t.Errorf("expected second value = 2 (ID_TOKEN), got %d", list.Get(1).Enum())
	}
}

// --- Apply nested message via flags ---

func TestApplyNestedMessage(t *testing.T) {
	reqDesc := mustResolveReq(t, "zitadel.user.v2.UserService/CreateUser")
	humanField := reqDesc.Fields().ByName("human")
	humanDesc := humanField.Message()

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	values := make(map[string]interface{})
	registerFlags(flags, humanDesc, "", values, nil, 0, true, nil)

	// Simulate setting flags
	if err := flags.Set("profile-given-name", "Alice"); err != nil {
		t.Fatalf("set profile-given-name: %v", err)
	}
	if err := flags.Set("profile-family-name", "Smith"); err != nil {
		t.Fatalf("set profile-family-name: %v", err)
	}
	if err := flags.Set("email", "alice@example.com"); err != nil {
		t.Fatalf("set email: %v", err)
	}

	// Apply flags to a dynamic message
	humanMsg := dynamicpb.NewMessage(humanDesc)
	cmd := cobraCmd(flags)
	applyFlags(cmd, humanMsg, humanDesc, values)

	// Verify nested profile was built correctly
	profileFD := humanDesc.Fields().ByName("profile")
	if profileFD == nil {
		t.Fatal("expected profile field")
	}
	profileMsg := humanMsg.Get(profileFD).Message()
	givenNameFD := profileFD.Message().Fields().ByName("given_name")
	got := profileMsg.Get(givenNameFD).String()
	if got != "Alice" {
		t.Errorf("expected given_name='Alice', got %q", got)
	}
}

