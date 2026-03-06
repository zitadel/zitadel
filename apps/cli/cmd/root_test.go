package cmd

import (
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

	stdinDone := make(chan struct{})
	go func() {
		if stdin != "" {
			_, _ = io.WriteString(stdinW, stdin)
		}
		_ = stdinW.Close()
		close(stdinDone)
	}()

	root := NewRootCmd()
	root.SetArgs(args)
	runErr := root.Execute()

	_ = stdoutW.Close()
	_ = stderrW.Close()
	<-stdinDone

	stdoutBytes, readOutErr := io.ReadAll(stdoutR)
	if readOutErr != nil {
		t.Fatalf("read stdout: %v", readOutErr)
	}
	stderrBytes, readErrErr := io.ReadAll(stderrR)
	if readErrErr != nil {
		t.Fatalf("read stderr: %v", readErrErr)
	}

	_ = stdinR.Close()
	_ = stdoutR.Close()
	_ = stderrR.Close()

	return string(stdoutBytes), string(stderrBytes), runErr
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
	if !strings.Contains(stdout, `"name": "create human"`) {
		t.Fatalf("expected create human metadata, got: %s", stdout)
	}
}

func TestDescribeRejectsAmbiguousBaseCommand(t *testing.T) {
	_, _, err := executeRoot(t, "", "describe", "users", "create")
	if err == nil {
		t.Fatalf("expected ambiguity error")
	}
	if !strings.Contains(err.Error(), `ambiguous command "create"`) {
		t.Fatalf("expected ambiguous command error, got: %v", err)
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
	if !strings.Contains(stdout, `"group": "human"`) {
		t.Fatalf("expected flag group 'human' in describe output, got: %s", stdout)
	}
}

func TestDescribeIncludesJSONTemplate(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "describe", "users", "get-by-id")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !strings.Contains(stdout, `"json_template"`) {
		t.Fatalf("expected json_template in describe output, got: %s", stdout)
	}
	if !strings.Contains(stdout, `"userId"`) {
		t.Fatalf("expected userId field in json_template, got: %s", stdout)
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
	// JSON template should include deeply nested fields not available as CLI flags
	if !strings.Contains(stdout, `"profile"`) {
		t.Fatalf("expected nested profile in json_template, got: %s", stdout)
	}
	if !strings.Contains(stdout, `"givenName"`) {
		t.Fatalf("expected nested givenName in json_template, got: %s", stdout)
	}
	if !strings.Contains(stdout, `"password"`) {
		t.Fatalf("expected nested password in json_template, got: %s", stdout)
	}
}

func TestVariantJSONTemplateFiltersOneofBranches(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "describe", "users", "set-email", "send-code")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !strings.Contains(stdout, `"sendCode"`) {
		t.Fatalf("expected sendCode in variant json_template, got: %s", stdout)
	}
	// Should NOT contain the other oneof branches
	if strings.Contains(stdout, `"returnCode"`) {
		t.Fatalf("variant json_template should not include returnCode, got: %s", stdout)
	}
	if strings.Contains(stdout, `"isVerified"`) {
		t.Fatalf("variant json_template should not include isVerified, got: %s", stdout)
	}
}

func TestVariantJSONTemplateReturnCodeBranch(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "describe", "users", "set-email", "return-code")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !strings.Contains(stdout, `"returnCode"`) {
		t.Fatalf("expected returnCode in variant json_template, got: %s", stdout)
	}
	if strings.Contains(stdout, `"sendCode"`) {
		t.Fatalf("variant json_template should not include sendCode, got: %s", stdout)
	}
}

func TestCreateHumanHasPasswordFlags(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "describe", "users", "create", "human")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	for _, flag := range []string{`"password"`, `"password-change-required"`, `"hashed-password-hash"`, `"hashed-password-change-required"`} {
		if !strings.Contains(stdout, flag) {
			t.Fatalf("expected password flag %s in create human describe output, got: %s", flag, stdout)
		}
	}
}

func TestPasswordFlagsDryRun(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "users", "create", "human",
		"--given-name", "Alice", "--family-name", "Doe", "--email", "alice@example.com",
		"--password", "s3cret!",
		"--dry-run")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !strings.Contains(stdout, `"password"`) {
		t.Fatalf("expected password in dry-run output, got: %s", stdout)
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
	if !strings.Contains(stdout, `"user-id"`) {
		t.Fatalf("expected user-id convenience filter flag in list-keys describe, got: %s", stdout)
	}
	if !strings.Contains(stdout, `"organization-id"`) {
		t.Fatalf("expected organization-id convenience filter flag in list-keys describe, got: %s", stdout)
	}
}

func TestListKeysUserIdDryRun(t *testing.T) {
	stdout, _, err := executeRoot(t, "", "users", "list-keys", "--user-id", "abc123", "--dry-run")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !strings.Contains(stdout, `abc123`) {
		t.Fatalf("expected user-id value in dry-run output, got: %s", stdout)
	}
	if !strings.Contains(stdout, `userIdFilter`) {
		t.Fatalf("expected userIdFilter in dry-run output, got: %s", stdout)
	}
}
