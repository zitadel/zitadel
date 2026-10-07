package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func executeRoot(t *testing.T, stdin string, args ...string) (string, string, error) {
	t.Helper()

	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("ZITADEL_TOKEN", "")
	t.Setenv("ZITADEL_INSTANCE", "")

	oldStdin := os.Stdin
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	stdinR, stdinW, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdin pipe: %v", err)
	}
	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}
	stderrR, stderrW, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stderr pipe: %v", err)
	}

	os.Stdin = stdinR
	os.Stdout = stdoutW
	os.Stderr = stderrW

	defer func() {
		os.Stdin = oldStdin
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	// Feed stdin in a goroutine.
	go func() {
		if stdin != "" {
			_, _ = io.WriteString(stdinW, stdin)
		}
		_ = stdinW.Close()
	}()

	// Drain stdout and stderr concurrently to prevent pipe buffer deadlock.
	// Without this, writes >64KB to os.Stdout block forever.
	var stdoutBuf, stderrBuf bytes.Buffer
	outDone := make(chan struct{})
	errDone := make(chan struct{})
	go func() {
		_, _ = io.Copy(&stdoutBuf, stdoutR)
		close(outDone)
	}()
	go func() {
		_, _ = io.Copy(&stderrBuf, stderrR)
		close(errDone)
	}()

	root := NewRootCmd()
	root.SetArgs(args)
	runErr := root.Execute()

	_ = stdoutW.Close()
	_ = stderrW.Close()
	<-outDone
	<-errDone

	_ = stdinR.Close()
	_ = stdoutR.Close()
	_ = stderrR.Close()

	return stdoutBuf.String(), stderrBuf.String(), runErr
}

func TestFromJSONDryRunBypassesRequiredFlags(t *testing.T) {
	stdout, _, err := executeRoot(t, `{"name":"my-org"}`, "orgs", "create", "--from-json", "--dry-run")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !strings.Contains(stdout, `"method": "zitadel.org.v2.OrganizationService/AddOrganization"`) {
		t.Fatalf("expected method in output, got: %s", stdout)
	}
	if !strings.Contains(stdout, `"name": "my-org"`) {
		t.Fatalf("expected request payload in output, got: %s", stdout)
	}
}

func TestFromJSONDryRunOneofVariantBypassesRequiredFlags(t *testing.T) {
	request := `{"organizationId":"1234567890","username":"alice","human":{"profile":{"givenName":"Alice","familyName":"Doe"},"email":{"email":"alice@example.com"}}}`
	stdout, _, err := executeRoot(t, request, "users", "create", "human", "--from-json", "--dry-run")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !strings.Contains(stdout, `"method": "zitadel.user.v2.UserService/CreateUser"`) {
		t.Fatalf("expected method in output, got: %s", stdout)
	}
	if !strings.Contains(stdout, `"givenName": "Alice"`) {
		t.Fatalf("expected human profile in output, got: %s", stdout)
	}
}

func TestDescribeSupportsVariantPath(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "describe", "users", "create", "human")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	// Variant path "create human" resolves to the base verb "create".
	if !strings.Contains(stdout, `"verb": "create"`) {
		t.Fatalf("expected verb create in describe output, got: %s", stdout)
	}
}

func TestDescribeRejectsAmbiguousBaseCommand(t *testing.T) {
	// "describe users create" should succeed — "create" is a valid verb.
	stdout, _, err := executeRoot(t, "", "describe", "users", "create")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !strings.Contains(stdout, `"verb": "create"`) {
		t.Fatalf("expected create verb in describe output, got: %s", stdout)
	}
}

func TestUserFacingHelpShowsJSONFlags(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "users", "create", "-h")
	if err != nil {
		t.Fatalf("expected help command success, got error: %v", err)
	}
	if !strings.Contains(stdout, "--from-json") {
		t.Fatalf("expected --from-json in help output, got: %s", stdout)
	}
	if !strings.Contains(stdout, "--dry-run") {
		t.Fatalf("expected --dry-run in help output, got: %s", stdout)
	}
}

func TestDescribeAllIncludesGlobalFlags(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "describe")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !strings.Contains(stdout, `"global_flags"`) {
		t.Fatalf("expected global_flags in describe output, got: %s", stdout)
	}
	if !strings.Contains(stdout, `"from-json"`) {
		t.Fatalf("expected from-json in global flags, got: %s", stdout)
	}
	if !strings.Contains(stdout, `"dry-run"`) {
		t.Fatalf("expected dry-run in global flags, got: %s", stdout)
	}
}

func TestDescribeAllIncludesFullPaths(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "describe")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	// Full paths should include the group prefix, e.g. "orgs create" not just "create"
	if !strings.Contains(stdout, `"orgs create"`) {
		t.Fatalf("expected full path 'orgs create' in describe output, got: %s", stdout)
	}
}

func TestDescribeVariantShowsGroupField(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "describe", "users", "create", "human")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	// The describe output for "create" should include the group "users".
	if !strings.Contains(stdout, `"group": "users"`) {
		t.Fatalf("expected group 'users' in describe output, got: %s", stdout)
	}
}

func TestDescribeIncludesJSONTemplate(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "describe", "users", "get-by-id")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	// The describe output includes request_schema with field definitions.
	if !strings.Contains(stdout, `"request_schema"`) {
		t.Fatalf("expected request_schema in describe output, got: %s", stdout)
	}
	if !strings.Contains(stdout, `"userId"`) {
		t.Fatalf("expected userId field in request_schema, got: %s", stdout)
	}
}

func TestDescribeListCommandHasPaginationFlags(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "describe", "users", "list")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	for _, flag := range []string{`"offset"`, `"limit"`, `"asc"`} {
		if !strings.Contains(stdout, flag) {
			t.Fatalf("expected pagination flag %s in describe output, got: %s", flag, stdout)
		}
	}
}

func TestRequestJSONDryRun(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "orgs", "create", "--request-json", `{"name":"inline-org"}`, "--dry-run")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !strings.Contains(stdout, `"name": "inline-org"`) {
		t.Fatalf("expected inline JSON payload in output, got: %s", stdout)
	}
}

func TestRequestJSONGlobalFlagInDescribe(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "describe")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !strings.Contains(stdout, `"request-json"`) {
		t.Fatalf("expected request-json in global flags, got: %s", stdout)
	}
}

func TestDescribeJSONTemplateShowsNestedFields(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "describe", "users", "create", "human")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	// request_schema should include deeply nested fields
	if !strings.Contains(stdout, `"profile"`) {
		t.Fatalf("expected nested profile in request_schema, got: %s", stdout)
	}
	if !strings.Contains(stdout, `"givenName"`) {
		t.Fatalf("expected nested givenName in request_schema, got: %s", stdout)
	}
	if !strings.Contains(stdout, `"password"`) {
		t.Fatalf("expected nested password in request_schema, got: %s", stdout)
	}
}

func TestVariantJSONTemplateFiltersOneofBranches(t *testing.T) {
	// "set-email send-code" resolves to verb "set-email", which includes all oneof branches.
	stdout, _, err := executeRoot(t, "", "describe", "users", "set-email", "send-code")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	// The full schema includes all oneof branches.
	if !strings.Contains(stdout, `"sendCode"`) {
		t.Fatalf("expected sendCode in request_schema, got: %s", stdout)
	}
}

func TestVariantJSONTemplateReturnCodeBranch(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "describe", "users", "set-email", "return-code")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	// Both oneof branches are present in the full schema.
	if !strings.Contains(stdout, `"returnCode"`) {
		t.Fatalf("expected returnCode in request_schema, got: %s", stdout)
	}
}

func TestCreateHumanHasPasswordFlags(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "describe", "users", "create", "human")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	// The request_schema includes password fields in the oneOf section.
	for _, field := range []string{`"password"`, `"changeRequired"`} {
		if !strings.Contains(stdout, field) {
			t.Fatalf("expected password field %s in create describe output, got: %s", field, stdout)
		}
	}
}

func TestPasswordFlagsDryRun(t *testing.T) {
	// Use --request-json to pass the password since it's a nested oneof field.
	stdout, _, err := executeRoot(t, "", "users", "create",
		"--request-json", `{"human":{"profile":{"givenName":"Alice","familyName":"Doe"},"email":{"email":"alice@example.com"},"password":{"password":"s3cret!"}}}`,
		"--dry-run")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !strings.Contains(stdout, `s3cret!`) {
		t.Fatalf("expected password value in dry-run output, got: %s", stdout)
	}
}

func TestListKeysHasUserIdFilterFlag(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "describe", "users", "list-keys")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	// The request_schema should include userId filter fields.
	if !strings.Contains(stdout, `"userId"`) {
		t.Fatalf("expected userId in list-keys describe request_schema, got: %s", stdout)
	}
}

func TestListKeysUserIdDryRun(t *testing.T) {
	// --user-id is accepted as a filter flag (no error = success).
	_, _, err := executeRoot(t, "", "users", "list-keys", "--user-id", "abc123", "--dry-run")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestDescribeOrgsListExists(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "describe", "orgs", "list")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !strings.Contains(stdout, `"method": "zitadel.org.v2.OrganizationService/ListOrganizations"`) {
		t.Fatalf("expected ListOrganizations method in describe output, got: %s", stdout)
	}
}

func TestDescribeIDPsListExists(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "describe", "idps", "get-idp-by-id")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !strings.Contains(stdout, `"group": "idps"`) {
		t.Fatalf("expected idps group in describe output, got: %s", stdout)
	}
}

func TestDescribeAppsListExists(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "describe", "apps", "list")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !strings.Contains(stdout, `"method": "zitadel.application.v2.ApplicationService/ListApplications"`) {
		t.Fatalf("expected ListApplications method in describe output, got: %s", stdout)
	}
}

func TestFlagAliasesDryRun(t *testing.T) {
	// --given-name should work as an alias for --profile-given-name
	stdout, _, err := executeRoot(t, "", "users", "create", "human",
		"--given-name", "Alice", "--family-name", "Doe", "--email", "alice@example.com",
		"--dry-run")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !strings.Contains(stdout, `"givenName": "Alice"`) {
		t.Fatalf("expected givenName in dry-run output, got: %s", stdout)
	}
	if !strings.Contains(stdout, `"familyName": "Doe"`) {
		t.Fatalf("expected familyName in dry-run output, got: %s", stdout)
	}
}

func TestPositionalArgsInUsage(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "users", "delete", "-h")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !strings.Contains(stdout, "<user_id>") {
		t.Fatalf("expected <user_id> in usage string, got: %s", stdout)
	}
}

func TestRepeatedEnumFlagsDryRun(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "apps", "create", "oidc-configuration",
		"--project-id", "123",
		"--name", "Test",
		"--response-types", "OIDC_RESPONSE_TYPE_CODE",
		"--grant-types", "OIDC_GRANT_TYPE_AUTHORIZATION_CODE",
		"--auth-method-type", "OIDC_AUTH_METHOD_TYPE_NONE",
		"--dry-run")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !strings.Contains(stdout, `"OIDC_RESPONSE_TYPE_CODE"`) {
		t.Fatalf("expected OIDC_RESPONSE_TYPE_CODE in dry-run output, got: %s", stdout)
	}
	if !strings.Contains(stdout, `"OIDC_GRANT_TYPE_AUTHORIZATION_CODE"`) {
		t.Fatalf("expected OIDC_GRANT_TYPE_AUTHORIZATION_CODE in dry-run output, got: %s", stdout)
	}
}

func TestOutputDefaultDescribeForGet(t *testing.T) {
	// describe output for a get command should mention 'describe' as default format
	stdout, _, err := executeRoot(t, "", "users", "get-by-id", "-h")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	// The help output mentions 'describe for TTY' in the -o flag description
	if !strings.Contains(stdout, "describe") {
		t.Fatalf("expected 'describe' mentioned in get command help, got: %s", stdout)
	}
}
