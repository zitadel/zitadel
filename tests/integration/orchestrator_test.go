//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/tests/integration/infra"
)

const (
	masterKey  = "MasterkeyNeedsToHave32Characters"
	configFile = "apps/api/test-integration-api.yaml"
	logDir     = ".artifacts"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test orchestrator in short mode")
	}

	// Go tests run with cwd set to the package directory. Change to the module
	// root so that relative paths (configFile, go list) work correctly.
	require.NoError(t, chdirToModuleRoot(), "chdir to module root")

	// Create log directory and orchestrator log file first, before anything else.
	// Nx TUI buffers both stdout and stderr of child processes, so writing to
	// os.Stderr alone is insufficient — the output only appears when the task
	// finishes. A dedicated log file allows:
	//   tail -f .artifacts/integration-orchestrator.log
	// in a separate terminal for real-time progress.
	require.NoError(t, os.MkdirAll(logDir, 0o755))
	orchLogPath := filepath.Join(logDir, "integration-orchestrator.log")
	orchLogFile, err := os.Create(orchLogPath)
	require.NoError(t, err, "create orchestrator log file")
	defer orchLogFile.Close()

	// logw writes to both stderr (for Nx non-TUI / direct invocation) and the
	// orchestrator log file (for tailing under Nx TUI).
	logw := io.MultiWriter(os.Stderr, orchLogFile)
	logf := func(format string, a ...any) {
		fmt.Fprintf(logw, "[integration] "+format+"\n", a...)
	}

	absOrchLogPath, _ := filepath.Abs(orchLogPath)
	logf("Orchestrator log: %s (tail -f to follow)", absOrchLogPath)

	// Intercept SIGTERM/SIGINT so deferred cleanup (container termination,
	// ZITADEL shutdown) runs instead of the process being killed immediately
	// by Nx or Ctrl-C.
	sigCtx, sigStop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer sigStop()

	ctx, cancel := context.WithTimeout(sigCtx, 30*time.Minute)
	defer cancel()

	// Use a separate context for cleanup that is NOT tied to signals/timeout.
	// Container.Terminate and ZITADEL.Shutdown need a live context even after
	// the main ctx is cancelled by SIGTERM or timeout.
	cleanupCtx, cleanupCancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Minute)
	defer cleanupCancel()

	// Phase 1: Start infrastructure containers
	logf("Starting Postgres container...")
	pgContainer, pgCfg, err := infra.StartPostgres(ctx, logw)
	require.NoError(t, err, "start postgres")
	defer func() {
		logf("Terminating Postgres container...")
		if err := pgContainer.Terminate(cleanupCtx); err != nil {
			logf("Postgres terminate error: %v", err)
		}
	}()
	logf("Postgres ready at %s:%s", pgCfg.Host, pgCfg.Port)

	logf("Starting Redis container...")
	redisContainer, redisCfg, err := infra.StartRedis(ctx, logw)
	require.NoError(t, err, "start redis")
	defer func() {
		logf("Terminating Redis container...")
		if err := redisContainer.Terminate(cleanupCtx); err != nil {
			logf("Redis terminate error: %v", err)
		}
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
		logf("Shutting down ZITADEL...")
		if err := instance.Shutdown(cleanupCtx); err != nil {
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

	// Integration test output log (separate from orchestrator log).
	testLogPath := filepath.Join(logDir, "integration-test-output.log")
	testLogFile, err := os.Create(testLogPath)
	require.NoError(t, err, "create integration test log file")
	defer testLogFile.Close()

	absTestLogPath, _ := filepath.Abs(testLogPath)
	logf("Integration test output: %s", absTestLogPath)

	logf("Running integration tests...")
	args := []string{
		"test",
		"-count", "1",
		"-tags", "integration",
		"-timeout", "60m",
		"-p", "1", // reliability-first default; avoids cross-package contention
		"-v",
	}

	// Forward -race to the child if the orchestrator itself was built with it.
	// This avoids a redundant full recompile when -race isn't requested.
	if raceEnabled {
		args = append(args[:1], append([]string{"-race"}, args[1:]...)...)
	}
	args = append(args, packages...)

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = "." // repo root

	// Put the child in its own process group so we can kill the entire tree
	// on cleanup, preventing orphaned go-test / go-build processes.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		// Send SIGTERM to the entire process group (negative PID).
		if cmd.Process != nil {
			return syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
		}
		return nil
	}
	cmd.WaitDelay = 10 * time.Second // fallback SIGKILL after 10s

	// Tee to both the test log file and the orchestrator log writer so
	// test output is visible via both tail targets.
	out := io.MultiWriter(testLogFile, logw)
	cmd.Stdout = out
	cmd.Stderr = out
	// Inherit all env vars including our ZITADEL_* overrides
	cmd.Env = os.Environ()

	err = cmd.Run()
	require.NoError(t, err, "integration tests failed; see %s", absTestLogPath)
	logf("All integration tests passed")
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
