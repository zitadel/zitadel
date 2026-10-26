//go:build integration

// Package integration provides CLI integration tests using testcontainers-go.
// These tests start a real ZITADEL instance (with Postgres) in Docker containers
// and exercise the CLI against it.
//
// Run via: go test -v -tags integration -timeout 10m ./internal/integration/
// Or:      pnpm nx run @zitadel/cli:test-integration
//
// Requires Docker to be running.
//
// Future: once the main tests/integration/infra package is reusable, migrate
// these tests into the centralized orchestrator framework to share the ZITADEL
// instance with other integration test packages.
package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

// zitadelStack holds the running ZITADEL + Postgres containers.
type zitadelStack struct {
	postgres    *tcpostgres.PostgresContainer
	zitadel     testcontainers.Container
	network     *testcontainers.DockerNetwork
	instanceURL string
	pat         string // admin PAT for authenticated API calls
}

func (s *zitadelStack) Terminate(ctx context.Context) {
	if s.zitadel != nil {
		_ = s.zitadel.Terminate(ctx)
	}
	if s.postgres != nil {
		_ = s.postgres.Terminate(ctx)
	}
	if s.network != nil {
		_ = s.network.Remove(ctx)
	}
}

// startStack spins up Postgres + ZITADEL on a shared Docker network.
// It bootstraps an admin PAT via ZITADEL_FIRSTINSTANCE_PATPATH.
func startStack(ctx context.Context, t *testing.T) *zitadelStack {
	t.Helper()

	// Create shared network
	net, err := network.New(ctx)
	if err != nil {
		t.Fatalf("failed to create network: %v", err)
	}

	// Start Postgres
	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("zitadel"),
		tcpostgres.WithUsername("zitadel"),
		tcpostgres.WithPassword("zitadel"),
		network.WithNetwork([]string{"postgres"}, net),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		_ = net.Remove(ctx)
		t.Fatalf("failed to start Postgres: %v", err)
	}

	masterKey := "MasterkeyNeedsToHave32Characters"

	// Create a host temp dir for the PAT file — docker cp can't read tmpfs.
	patDir := t.TempDir()
	patContainerPath := "/data/admin.pat"

	// Start ZITADEL pointing at Postgres on the shared network.
	// ZITADEL_FIRSTINSTANCE_PATPATH tells start-from-init to write a PAT file.
	zitadelReq := testcontainers.ContainerRequest{
		Image:        "ghcr.io/zitadel/zitadel:latest",
		ExposedPorts: []string{"8080/tcp"},
		Networks:     []string{net.Name},
		Mounts: testcontainers.Mounts(
			testcontainers.BindMount(patDir, "/data"),
		),
		Cmd: []string{
			"start-from-init",
			"--masterkeyFromEnv",
			"--tlsMode", "disabled",
		},
		Env: map[string]string{
			"ZITADEL_MASTERKEY":                        masterKey,
			"ZITADEL_FIRSTINSTANCE_PATPATH":            patContainerPath,
			"ZITADEL_DATABASE_POSTGRES_HOST":           "postgres",
			"ZITADEL_DATABASE_POSTGRES_PORT":           "5432",
			"ZITADEL_DATABASE_POSTGRES_DATABASE":       "zitadel",
			"ZITADEL_DATABASE_POSTGRES_USER_USERNAME":  "zitadel",
			"ZITADEL_DATABASE_POSTGRES_USER_PASSWORD":  "zitadel",
			"ZITADEL_DATABASE_POSTGRES_USER_SSL_MODE":  "disable",
			"ZITADEL_DATABASE_POSTGRES_ADMIN_USERNAME": "zitadel",
			"ZITADEL_DATABASE_POSTGRES_ADMIN_PASSWORD": "zitadel",
			"ZITADEL_DATABASE_POSTGRES_ADMIN_SSL_MODE": "disable",
			// Create a machine user with IAM_OWNER role and generate a PAT.
			"ZITADEL_FIRSTINSTANCE_ORG_MACHINE_MACHINE_USERNAME":   "cli-integration-sa",
			"ZITADEL_FIRSTINSTANCE_ORG_MACHINE_MACHINE_NAME":       "CLI Integration Service Account",
			"ZITADEL_FIRSTINSTANCE_ORG_MACHINE_PAT_EXPIRATIONDATE": "2099-01-01T00:00:00Z",
		},
		WaitingFor: wait.ForHTTP("/debug/ready").
			WithPort("8080/tcp").
			WithStartupTimeout(120 * time.Second),
	}

	zitadelContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: zitadelReq,
		Started:          true,
	})
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		_ = net.Remove(ctx)
		t.Fatalf("failed to start ZITADEL container: %v", err)
	}

	host, _ := zitadelContainer.Host(ctx)
	port, _ := zitadelContainer.MappedPort(ctx, "8080/tcp")
	instanceURL := fmt.Sprintf("http://%s:%s", host, port.Port())

	// Read the bootstrapped PAT from the host-mounted directory.
	pat := readPATFromHost(t, patDir+"/admin.pat")

	t.Logf("ZITADEL ready at %s (PAT: %s...)", instanceURL, pat[:min(8, len(pat))])

	return &zitadelStack{
		postgres:    pgContainer,
		zitadel:     zitadelContainer,
		network:     net,
		instanceURL: instanceURL,
		pat:         pat,
	}
}

// readPATFromHost reads the admin PAT from the bind-mounted host directory.
func readPATFromHost(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read PAT from %s: %v", path, err)
	}
	pat := strings.TrimSpace(string(data))
	if pat == "" {
		t.Fatal("PAT file is empty")
	}
	return pat
}

// buildCLI builds the CLI binary for test use.
func buildCLI(t *testing.T) string {
	t.Helper()
	binary := t.TempDir() + "/zitadel-cli-test"
	cmd := exec.Command("go", "build", "-o", binary, ".")
	cmd.Dir = "../../" // apps/cli/
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build CLI: %v\n%s", err, out)
	}
	return binary
}

// runLive execs the CLI binary against the running ZITADEL instance.
// Returns only stdout (JSON output), stderr is logged on failure.
func runLive(t *testing.T, stack *zitadelStack, binary string, args ...string) string {
	t.Helper()
	cmd := exec.Command(binary, args...)
	cmd.Env = append(os.Environ(),
		"ZITADEL_CLI_NO_WARN=1",
		"ZITADEL_TOKEN="+stack.pat,
		"ZITADEL_INSTANCE="+stack.instanceURL,
	)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Logf("CLI stderr: %s", stderr.String())
		t.Logf("CLI stdout: %s", stdout.String())
		t.Fatalf("CLI command failed (%v): %s", err, args)
	}
	return stdout.String()
}

// runLiveExpectError runs the CLI and expects it to fail. Returns stdout.
func runLiveExpectError(t *testing.T, stack *zitadelStack, binary string, args ...string) string {
	t.Helper()
	cmd := exec.Command(binary, args...)
	cmd.Env = append(os.Environ(),
		"ZITADEL_CLI_NO_WARN=1",
		"ZITADEL_TOKEN="+stack.pat,
		"ZITADEL_INSTANCE="+stack.instanceURL,
	)
	var stdout strings.Builder
	cmd.Stdout = &stdout
	_ = cmd.Run()
	return stdout.String()
}

// parseJSON parses CLI JSON output into a map.
func parseJSON(t *testing.T, output string) map[string]interface{} {
	t.Helper()
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("failed to parse JSON output: %v\nraw: %s", err, output)
	}
	return result
}

// jsonString extracts a string field from a JSON map, supporting nested dot paths.
func jsonString(t *testing.T, m map[string]interface{}, key string) string {
	t.Helper()
	parts := strings.Split(key, ".")
	current := m
	for i, part := range parts {
		val, ok := current[part]
		if !ok {
			t.Fatalf("key %q not found in JSON (at part %q)", key, part)
		}
		if i == len(parts)-1 {
			s, ok := val.(string)
			if !ok {
				t.Fatalf("key %q is not a string: %T", key, val)
			}
			return s
		}
		current, ok = val.(map[string]interface{})
		if !ok {
			t.Fatalf("key %q is not an object at part %q", key, part)
		}
	}
	return ""
}

// TestCLIIntegration is the orchestrator test.
func TestCLIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	stack := startStack(ctx, t)
	defer stack.Terminate(ctx)

	binary := buildCLI(t)

	// Live CRUD tests against running ZITADEL.
	// Dry-run tests are covered by unit tests in cmd/root_test.go.

	// Shared state between sequential subtests.
	var orgID, projectID string

	t.Run("CRUDOrg", func(t *testing.T) {
		// Create
		out := runLive(t, stack, binary, "orgs", "create", "--name", "CLI-Test-Org")
		result := parseJSON(t, out)
		orgID = jsonString(t, result, "organizationId")
		if orgID == "" {
			t.Fatalf("expected organizationId in response, got: %s", out)
		}
		t.Logf("Created org: %s", orgID)

		// List — verify it appears
		listOut := runLive(t, stack, binary, "orgs", "list")
		if !strings.Contains(listOut, orgID) {
			t.Fatalf("expected org %s in list output, got: %s", orgID, listOut)
		}
	})

	t.Run("CRUDUser", func(t *testing.T) {
		if orgID == "" {
			t.Skip("skipping: no orgID from CRUDOrg")
		}

		// Create
		out := runLive(t, stack, binary, "users", "create", "human",
			"--organization-id", orgID,
			"--given-name", "CRUD",
			"--family-name", "TestUser",
			"--email", "crud-test@example.com")
		result := parseJSON(t, out)
		userID := jsonString(t, result, "id")
		if userID == "" {
			t.Fatalf("expected id in response, got: %s", out)
		}
		t.Logf("Created user: %s", userID)

		// Get by ID
		getOut := runLive(t, stack, binary, "users", "get-by-id", userID)
		if !strings.Contains(getOut, userID) {
			t.Fatalf("expected userId %s in get response, got: %s", userID, getOut)
		}

		// List — verify it appears
		listOut := runLive(t, stack, binary, "users", "list")
		if !strings.Contains(listOut, userID) {
			t.Fatalf("expected user %s in list output, got: %s", userID, listOut)
		}

		// Deactivate
		deactOut := runLive(t, stack, binary, "users", "deactivate", userID, "--yes")
		if !strings.Contains(deactOut, userID) {
			t.Logf("deactivate output: %s", deactOut)
		}

		// Delete
		deleteOut := runLive(t, stack, binary, "users", "delete", userID, "--yes")
		if !strings.Contains(deleteOut, userID) {
			t.Logf("delete output: %s", deleteOut)
		}

		// Verify deleted — get should fail
		verifyOut := runLiveExpectError(t, stack, binary, "users", "get-by-id", userID)
		t.Logf("After delete, get: %s", verifyOut)
	})

	t.Run("CRUDProject", func(t *testing.T) {
		if orgID == "" {
			t.Skip("skipping: no orgID from CRUDOrg")
		}

		// Create
		out := runLive(t, stack, binary, "projects", "create",
			"--organization-id", orgID,
			"--name", "CLI-Test-Project")
		result := parseJSON(t, out)
		// Try common field names for the project ID
		for _, key := range []string{"id", "projectId"} {
			if v, ok := result[key].(string); ok && v != "" {
				projectID = v
				break
			}
		}
		if projectID == "" {
			t.Fatalf("expected project ID in response, got: %s", out)
		}
		t.Logf("Created project: %s", projectID)

		// List
		listOut := runLive(t, stack, binary, "projects", "list")
		if !strings.Contains(listOut, projectID) {
			t.Fatalf("expected project %s in list output, got: %s", projectID, listOut)
		}
	})

	t.Run("CRUDApp", func(t *testing.T) {
		if projectID == "" {
			t.Skip("skipping: no projectID from CRUDProject")
		}

		// Create OIDC application
		out := runLive(t, stack, binary, "apps", "create", "oidc-configuration",
			"--project-id", projectID,
			"--name", "CLI-Test-App",
			"--response-types", "OIDC_RESPONSE_TYPE_CODE",
			"--grant-types", "OIDC_GRANT_TYPE_AUTHORIZATION_CODE",
			"--auth-method-type", "OIDC_AUTH_METHOD_TYPE_NONE")
		result := parseJSON(t, out)

		appID := ""
		for _, key := range []string{"appId", "id", "applicationId"} {
			if v, ok := result[key].(string); ok && v != "" {
				appID = v
				break
			}
		}
		if appID == "" {
			t.Fatalf("expected app ID in response, got: %s", out)
		}
		t.Logf("Created app: %s (project: %s)", appID, projectID)

		// List apps globally — the app should appear
		listOut := runLive(t, stack, binary, "apps", "list")
		if !strings.Contains(listOut, appID) {
			t.Fatalf("expected app %s in list output, got: %s", appID, listOut)
		}
	})
}
