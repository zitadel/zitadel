//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/tests/integration/infra"
)

const (
	masterKey  = "MasterkeyNeedsToHave32Characters"
	configFile = "apps/api/test-integration-api.yaml"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test orchestrator in short mode")
	}

	// logf writes directly to stderr so progress is visible immediately.
	// t.Log/t.Logf are buffered and only printed when the test function returns,
	// which can be 30+ minutes away.
	logf := func(format string, a ...any) {
		fmt.Fprintf(os.Stderr, "[integration] "+format+"\n", a...)
	}

	// Go tests run with cwd set to the package directory. Change to the module
	// root so that relative paths (configFile, go list) work correctly.
	require.NoError(t, chdirToModuleRoot(), "chdir to module root")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Phase 1: Start infrastructure containers
	logf("Starting Postgres container...")
	pgContainer, pgCfg, err := infra.StartPostgres(ctx)
	require.NoError(t, err, "start postgres")
	defer func() {
		require.NoError(t, pgContainer.Terminate(ctx), "terminate postgres")
	}()
	logf("Postgres ready at %s:%s", pgCfg.Host, pgCfg.Port)

	logf("Starting Redis container...")
	redisContainer, redisCfg, err := infra.StartRedis(ctx)
	require.NoError(t, err, "start redis")
	defer func() {
		require.NoError(t, redisContainer.Terminate(ctx), "terminate redis")
	}()
	logf("Redis ready at %s", redisCfg.Addr)

	// Phase 2: Configure ZITADEL to use testcontainers
	t.Setenv("ZITADEL_DATABASE_POSTGRES_HOST", pgCfg.Host)
	t.Setenv("ZITADEL_DATABASE_POSTGRES_PORT", pgCfg.Port)
	t.Setenv("ZITADEL_DATABASE_POSTGRES_ADMIN_USERNAME", pgCfg.User)
	t.Setenv("ZITADEL_DATABASE_POSTGRES_ADMIN_PASSWORD", pgCfg.Password)
	t.Setenv("ZITADEL_DATABASE_POSTGRES_USER_USERNAME", "zitadel")
	t.Setenv("ZITADEL_DATABASE_POSTGRES_USER_SSL_MODE", "disable")
	t.Setenv("ZITADEL_DATABASE_POSTGRES_ADMIN_SSL_MODE", "disable")
	t.Setenv("ZITADEL_CACHES_CONNECTORS_REDIS_ENABLED", "true")
	t.Setenv("ZITADEL_CACHES_CONNECTORS_REDIS_ADDR", redisCfg.Addr)

	// Match the port expected by internal/integration/config/client.yaml
	t.Setenv("ZITADEL_PORT", "8082")
	t.Setenv("ZITADEL_EXTERNALPORT", "8082")

	// Phase 3: Start ZITADEL in-process
	// The integration build tag causes ZITADEL to start a sink server on
	// port 8081. Fail fast with a helpful message if the port is already in
	// use (e.g. by a leftover process from a previous run).
	require.NoError(t, waitForPortFree(ctx, "127.0.0.1:8081", 10*time.Second),
		"port 8081 is in use; a previous integration test run may still be running — "+
			"kill the leftover process and retry")
	logf("Starting ZITADEL (init → setup → start)...")
	instance, err := infra.StartZITADEL(ctx, []string{configFile}, masterKey)
	require.NoError(t, err, "start ZITADEL")
	defer func() {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()
		if err := instance.Shutdown(shutdownCtx); err != nil {
			logf("ZITADEL shutdown error (may be expected): %v", err)
		}
	}()

	// Phase 4: Wait for ZITADEL to be healthy
	healthCtx, healthCancel := context.WithTimeout(ctx, 5*time.Minute)
	defer healthCancel()
	logf("Waiting for ZITADEL health check...")
	require.NoError(t, infra.WaitForHealthy(healthCtx, "http://localhost:8082"), "health check")
	logf("ZITADEL is ready")

	// Phase 5: Run integration test packages
	logf("Discovering integration test packages...")
	packages, err := discoverIntegrationPackages()
	require.NoError(t, err, "discover packages")
	logf("Found %d integration test packages", len(packages))

	// Write integration test output to a dedicated log file so it can also be
	// inspected after the run (e.g. via `cat`). Live output goes to stderr above.
	logPath := filepath.Join(".artifacts", "integration-test-output.log")
	require.NoError(t, os.MkdirAll(".artifacts", 0o755))
	logFile, err := os.Create(logPath)
	require.NoError(t, err, "create integration log file")
	defer logFile.Close()

	absLogPath, _ := filepath.Abs(logPath)
	logf("Integration test output also saved to: %s", absLogPath)

	logf("Running integration tests...")
	args := []string{
		"test",
		"-race",
		"-count", "1",
		"-tags", "integration",
		"-timeout", "60m",
		"-parallel", "1",
		"-v",
	}
	args = append(args, packages...)

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = "." // repo root
	// Tee to both the log file and stderr so output is visible in real time.
	out := io.MultiWriter(logFile, os.Stderr)
	cmd.Stdout = out
	cmd.Stderr = out
	// Inherit all env vars including our ZITADEL_* overrides
	cmd.Env = os.Environ()

	err = cmd.Run()
	require.NoError(t, err, "integration tests failed; see %s", absLogPath)
}

// discoverIntegrationPackages finds all Go packages matching the integration test pattern.
func discoverIntegrationPackages() ([]string, error) {
	cmd := exec.Command("go", "list", "-tags", "integration", "./...")
	cmd.Dir = "."
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("go list: %w", err)
	}

	var packages []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if strings.Contains(line, "integration_test") || strings.Contains(line, "events_testing") {
			packages = append(packages, line)
		}
	}

	if len(packages) == 0 {
		return nil, fmt.Errorf("no integration test packages found")
	}
	return packages, nil
}

// chdirToModuleRoot walks up from the test file's directory until it finds go.mod,
// then changes the working directory to that location.
func chdirToModuleRoot() error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("runtime.Caller failed")
	}
	dir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return os.Chdir(dir)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return fmt.Errorf("go.mod not found starting from %s", filepath.Dir(filename))
		}
		dir = parent
	}
}

// waitForPortFree polls until the given address is not in use or the deadline
// is reached. It returns nil if the port is free, or a timeout error.
func waitForPortFree(ctx context.Context, addr string, maxWait time.Duration) error {
	deadline := time.Now().Add(maxWait)
	for {
		conn, err := net.DialTimeout("tcp", addr, time.Second)
		if err != nil {
			// Port is not in use — we're free to proceed.
			return nil
		}
		conn.Close()
		if time.Now().After(deadline) {
			return fmt.Errorf("timed out waiting for %s to be free after %s", addr, maxWait)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
	}
}
