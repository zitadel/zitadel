package main

import (
	"slices"
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// buildPlugin constructs a *protogen.Plugin from a synthetic FileDescriptorProto.
// The file is automatically marked for generation.
func buildPlugin(t *testing.T, fd *descriptorpb.FileDescriptorProto) *protogen.Plugin {
	t.Helper()
	// protogen requires a go_package option to resolve the Go import path.
	if fd.Options == nil {
		fd.Options = &descriptorpb.FileOptions{}
	}
	if fd.Options.GoPackage == nil {
		fd.Options.GoPackage = proto.String("example.com/test;test")
	}
	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{fd.GetName()},
		ProtoFile:      []*descriptorpb.FileDescriptorProto{fd},
	}
	plugin, err := protogen.Options{}.New(req)
	if err != nil {
		t.Fatalf("protogen.Options{}.New: %v", err)
	}
	return plugin
}

// findMessage returns the named message from the first file in the plugin.
func findMessage(t *testing.T, plugin *protogen.Plugin, name string) *protogen.Message {
	t.Helper()
	for _, f := range plugin.Files {
		for _, m := range f.Messages {
			if string(m.Desc.Name()) == name {
				return m
			}
		}
	}
	t.Fatalf("message %q not found in plugin", name)
	return nil
}

// detailsMessage returns a *descriptorpb.DescriptorProto representing a minimal
// "Details" message with a resource_owner string field.
func detailsDescriptor() *descriptorpb.DescriptorProto {
	return &descriptorpb.DescriptorProto{
		Name: proto.String("Details"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:     proto.String("resource_owner"),
				Number:   proto.Int32(3),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				JsonName: proto.String("resourceOwner"),
			},
		},
	}
}

// stringField is a convenience builder for a simple string FieldDescriptorProto.
func stringField(name string, number int32) *descriptorpb.FieldDescriptorProto {
	return &descriptorpb.FieldDescriptorProto{
		Name:     proto.String(name),
		Number:   proto.Int32(number),
		Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
		Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
		JsonName: proto.String(name),
	}
}

// messageField returns a FieldDescriptorProto for a message-typed field.
func messageField(name string, number int32, typeName string) *descriptorpb.FieldDescriptorProto {
	return &descriptorpb.FieldDescriptorProto{
		Name:     proto.String(name),
		Number:   proto.Int32(number),
		Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
		Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
		TypeName: proto.String(typeName),
		JsonName: proto.String(name),
	}
}

// columnHeaders returns just the header strings from a []columnDef slice.
func columnHeaders(cols []columnDef) []string {
	headers := make([]string, len(cols))
	for i, c := range cols {
		headers[i] = c.Header
	}
	return headers
}

// ---- Tests for extractMessageColumns ----

// TestExtractMessageColumns_IDPLike verifies that a resource with a bare "id" field
// and a "details" message containing resource_owner shows both ID and ORGANIZATION ID.
// This is the regression test for Bug 1 (the suppression rule that incorrectly dropped
// ORGANIZATION ID from IDP-like resources).
func TestExtractMessageColumns_IDPLike(t *testing.T) {
	fd := &descriptorpb.FileDescriptorProto{
		Name:    proto.String("test.proto"),
		Syntax:  proto.String("proto3"),
		Package: proto.String("test"),
		MessageType: []*descriptorpb.DescriptorProto{
			detailsDescriptor(),
			{
				Name: proto.String("IDP"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField("id", 1),
					messageField("details", 2, ".test.Details"),
					stringField("name", 3),
				},
			},
		},
	}

	plugin := buildPlugin(t, fd)
	msg := findMessage(t, plugin, "IDP")
	cols := extractMessageColumns(msg)
	headers := columnHeaders(cols)

	if !slices.Contains(headers, "ID") {
		t.Errorf("expected ID column, got %v", headers)
	}
	if !slices.Contains(headers, "ORGANIZATION ID") {
		t.Errorf("expected ORGANIZATION ID column, got %v", headers)
	}
	// ID must come before ORGANIZATION ID
	idIdx := slices.Index(headers, "ID")
	orgIdx := slices.Index(headers, "ORGANIZATION ID")
	if idIdx > orgIdx {
		t.Errorf("expected ID before ORGANIZATION ID, got %v", headers)
	}
}

// TestExtractMessageColumns_OrganizationLike verifies that when the message is named
// "Organization" (so primaryIDHeader == "ORGANIZATION ID"), details.resource_owner is
// still placed in the ORGANIZATION ID column and does NOT overwrite the primary id.
// This is the regression test for Bug 2 (switch-case ordering: "ORGANIZATION ID" must
// be checked before primaryIDHeader).
func TestExtractMessageColumns_OrganizationLike(t *testing.T) {
	fd := &descriptorpb.FileDescriptorProto{
		Name:    proto.String("test.proto"),
		Syntax:  proto.String("proto3"),
		Package: proto.String("test"),
		MessageType: []*descriptorpb.DescriptorProto{
			detailsDescriptor(),
			{
				Name: proto.String("Organization"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField("id", 1),
					messageField("details", 2, ".test.Details"),
					stringField("name", 3),
					stringField("primary_domain", 4),
				},
			},
		},
	}

	plugin := buildPlugin(t, fd)
	msg := findMessage(t, plugin, "Organization")
	cols := extractMessageColumns(msg)
	headers := columnHeaders(cols)

	if !slices.Contains(headers, "ID") {
		t.Errorf("expected ID column, got %v", headers)
	}
	if !slices.Contains(headers, "ORGANIZATION ID") {
		t.Errorf("expected ORGANIZATION ID column, got %v", headers)
	}
	// Both ID (from the bare id field) and ORGANIZATION ID (from details.resource_owner)
	// must be present and separate.
	if len(headers) < 2 {
		t.Fatalf("expected at least 2 columns, got %v", headers)
	}
	// The first column must be ID (the resource's own key).
	if headers[0] != "ID" {
		t.Errorf("expected ID as first column, got %v", headers)
	}
	// ORGANIZATION ID must come second.
	if headers[1] != "ORGANIZATION ID" {
		t.Errorf("expected ORGANIZATION ID as second column, got %v", headers)
	}
	// ID accessor must use GetId(), not GetDetails().GetResourceOwner().
	if cols[0].GoAccessor != "GetId()" {
		t.Errorf("expected ID accessor to be GetId(), got %q", cols[0].GoAccessor)
	}
}

// TestExtractMessageColumns_UserLike verifies that a resource with a "user_id" field
// (not bare "id") has it renamed to "ID", and details.resource_owner becomes ORGANIZATION ID.
func TestExtractMessageColumns_UserLike(t *testing.T) {
	fd := &descriptorpb.FileDescriptorProto{
		Name:    proto.String("test.proto"),
		Syntax:  proto.String("proto3"),
		Package: proto.String("test"),
		MessageType: []*descriptorpb.DescriptorProto{
			detailsDescriptor(),
			{
				Name: proto.String("User"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField("user_id", 1),
					messageField("details", 2, ".test.Details"),
					stringField("username", 3),
				},
			},
		},
	}

	plugin := buildPlugin(t, fd)
	msg := findMessage(t, plugin, "User")
	cols := extractMessageColumns(msg)
	headers := columnHeaders(cols)

	// user_id should be renamed to plain "ID".
	if !slices.Contains(headers, "ID") {
		t.Errorf("expected ID column (renamed from user_id), got %v", headers)
	}
	if slices.Contains(headers, "USER ID") {
		t.Errorf("expected user_id to be renamed to ID, still found USER ID in %v", headers)
	}
	if !slices.Contains(headers, "ORGANIZATION ID") {
		t.Errorf("expected ORGANIZATION ID column from details.resource_owner, got %v", headers)
	}
	// ID must come before ORGANIZATION ID
	idIdx := slices.Index(headers, "ID")
	orgIdx := slices.Index(headers, "ORGANIZATION ID")
	if idIdx > orgIdx {
		t.Errorf("expected ID before ORGANIZATION ID, got %v", headers)
	}
}

// TestExtractMessageColumns_NoOrgID verifies that a resource with no details field and
// no organization_id field does not show an ORGANIZATION ID column.
func TestExtractMessageColumns_NoOrgID(t *testing.T) {
	fd := &descriptorpb.FileDescriptorProto{
		Name:    proto.String("test.proto"),
		Syntax:  proto.String("proto3"),
		Package: proto.String("test"),
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: proto.String("Application"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField("application_id", 1),
					stringField("project_id", 2),
					stringField("name", 3),
				},
			},
		},
	}

	plugin := buildPlugin(t, fd)
	msg := findMessage(t, plugin, "Application")
	cols := extractMessageColumns(msg)
	headers := columnHeaders(cols)

	if slices.Contains(headers, "ORGANIZATION ID") {
		t.Errorf("expected no ORGANIZATION ID column for Application-like resource, got %v", headers)
	}
	// application_id should be renamed to ID
	if !slices.Contains(headers, "ID") {
		t.Errorf("expected ID column (renamed from application_id), got %v", headers)
	}
	// PROJECT ID should be present as a foreign ID
	if !slices.Contains(headers, "PROJECT ID") {
		t.Errorf("expected PROJECT ID column, got %v", headers)
	}
}

// ---- Tests for toKebab ----

func TestToKebab(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Organization", "organization"},
		{"OrganizationService", "organization-service"},
		// Acronyms at the end of a word are kept together, not split char-by-char.
		{"UserID", "user-id"},
		{"GetUserByID", "get-user-by-id"},
		// All-uppercase strings have no dash insertions.
		{"IDP", "idp"},
		{"ListIDPs", "list-id-ps"},
		{"OIDCConfiguration", "oidc-configuration"},
		{"PrimaryDomain", "primary-domain"},
		{"snake_case", "snake-case"},
		{"already-kebab", "already-kebab"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toKebab(tt.input)
			if got != tt.want {
				t.Errorf("toKebab(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// ---- Tests for enumShortValues ----

func TestEnumShortValues(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "strips common prefix",
			input: []string{"USER_STATE_UNSPECIFIED", "USER_STATE_ACTIVE", "USER_STATE_INACTIVE"},
			want:  []string{"ACTIVE", "INACTIVE"},
		},
		{
			name:  "strips longer prefix",
			input: []string{"ORGANIZATION_FIELD_NAME_UNSPECIFIED", "ORGANIZATION_FIELD_NAME_NAME", "ORGANIZATION_FIELD_NAME_CREATION_DATE"},
			want:  []string{"NAME", "CREATION_DATE"},
		},
		{
			name:  "no common prefix returns original values",
			input: []string{"FOO", "BAR", "BAZ"},
			want:  []string{"FOO", "BAR", "BAZ"},
		},
		{
			name:  "empty input",
			input: nil,
			want:  nil,
		},
		{
			// With a single value, the function strips the full value as the common prefix,
			// leaving an empty string which is not filtered (only UNSPECIFIED is filtered).
			name:  "single value strips to empty",
			input: []string{"ORGANIZATION_STATE_ACTIVE"},
			want:  []string{""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := enumShortValues(tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("enumShortValues(%v) = %v, want %v", tt.input, got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("enumShortValues(%v)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}

// ---- Tests for rpcNameToVerbAndSuffix ----

func TestRPCNameToVerbAndSuffix(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		wantVerb    string
		wantSuffix  string
	}{
		{"ListOrganizations", "OrganizationService", "list", ""},
		{"CreateOrganization", "OrganizationService", "create", ""},
		// "ByID" suffix — toKebab keeps "ID" together as an acronym → "by-id"
		{"GetOrganizationByID", "OrganizationService", "get", "by-id"},
		{"ListOrganizationDomains", "OrganizationService", "list", "domains"},
		{"AddOrganizationMember", "OrganizationService", "create", "member"},
		{"DeleteOrganization", "OrganizationService", "delete", ""},
		{"ListUsers", "UserService", "list", ""},
		{"SetUserMetadata", "UserService", "set", "metadata"},
		{"RemoveUserPasskey", "UserService", "remove", "passkey"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verb, suffix := rpcNameToVerbAndSuffix(tt.name, tt.serviceName)
			if verb != tt.wantVerb {
				t.Errorf("rpcNameToVerbAndSuffix(%q, %q) verb = %q, want %q", tt.name, tt.serviceName, verb, tt.wantVerb)
			}
			if suffix != tt.wantSuffix {
				t.Errorf("rpcNameToVerbAndSuffix(%q, %q) suffix = %q, want %q", tt.name, tt.serviceName, suffix, tt.wantSuffix)
			}
		})
	}
}

// ---- Tests for extractTopLevelColumns and extractResponseColumns ----

// TestExtractTopLevelColumns_NestedResourceUnwrapped verifies that a response wrapping
// a nested resource message (like GetUserByIDResponse) returns the resource's columns,
// not just a fallback "CHANGE DATE".  This is the regression test for the get-by-id bug
// where only "CHANGE DATE" was shown instead of the full resource fields.
func TestExtractTopLevelColumns_NestedResourceUnwrapped(t *testing.T) {
resource := &descriptorpb.DescriptorProto{
Name: proto.String("Resource"),
Field: []*descriptorpb.FieldDescriptorProto{
stringField("id", 1),
stringField("name", 2),
},
}
response := &descriptorpb.DescriptorProto{
Name: proto.String("GetResourceByIDResponse"),
Field: []*descriptorpb.FieldDescriptorProto{
messageField("details", 1, ".Details"),
messageField("resource", 2, ".Resource"),
},
}
fd := &descriptorpb.FileDescriptorProto{
Name:        proto.String("test.proto"),
Syntax:      proto.String("proto3"),
MessageType: []*descriptorpb.DescriptorProto{detailsDescriptor(), resource, response},
}
plugin := buildPlugin(t, fd)
msg := findMessage(t, plugin, "GetResourceByIDResponse")

// extractTopLevelColumns should return empty (nested message → hasNestedMessage)
topLevel := extractTopLevelColumns(msg)
if len(topLevel) != 0 {
t.Errorf("extractTopLevelColumns: want empty (nested message present), got %v", columnHeaders(topLevel))
}

// extractResponseColumns should unwrap the nested resource → ID + NAME columns
_, _, cols := extractResponseColumns(msg, "get")
headers := columnHeaders(cols)
if !slices.Contains(headers, "ID") {
t.Errorf("extractResponseColumns: expected ID column in %v", headers)
}
if !slices.Contains(headers, "NAME") {
t.Errorf("extractResponseColumns: expected NAME column in %v", headers)
}
if len(headers) == 1 && headers[0] == "CHANGE DATE" {
t.Errorf("extractResponseColumns: got only CHANGE DATE, expected full resource columns")
}
}

// TestExtractTopLevelColumns_ActionResponseChangeDateFallback verifies that a simple
// "action" response (only a details field with no nested resource message) still shows
// CHANGE DATE as a confirmation that the operation completed.
func TestExtractTopLevelColumns_ActionResponseChangeDateFallback(t *testing.T) {
details := &descriptorpb.DescriptorProto{
Name: proto.String("Details"),
Field: []*descriptorpb.FieldDescriptorProto{
{
Name:     proto.String("change_date"),
Number:   proto.Int32(1),
Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
TypeName: proto.String(".Timestamp"),
JsonName: proto.String("changeDate"),
},
},
}
// Minimal Timestamp message (stands in for google.protobuf.Timestamp in the
// synthetic descriptor — full-name check won't match, but the fallback path
// does not require that; it just checks the field name "change_date").
timestamp := &descriptorpb.DescriptorProto{
Name:  proto.String("Timestamp"),
Field: []*descriptorpb.FieldDescriptorProto{},
}
response := &descriptorpb.DescriptorProto{
Name: proto.String("DeleteResourceResponse"),
Field: []*descriptorpb.FieldDescriptorProto{
messageField("details", 1, ".Details"),
},
}
fd := &descriptorpb.FileDescriptorProto{
Name:        proto.String("test.proto"),
Syntax:      proto.String("proto3"),
MessageType: []*descriptorpb.DescriptorProto{details, timestamp, response},
}
plugin := buildPlugin(t, fd)
msg := findMessage(t, plugin, "DeleteResourceResponse")

cols := extractTopLevelColumns(msg)
headers := columnHeaders(cols)
if !slices.Contains(headers, "CHANGE DATE") {
t.Errorf("extractTopLevelColumns: want CHANGE DATE fallback for action response, got %v", headers)
}
}
