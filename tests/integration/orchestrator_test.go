//go:build integration

package integration_test

import (
	"context"
	"fmt"
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
	masterKey       = "MasterkeyNeedsToHave32Characters"
	configFile      = "apps/api/test-integration-api.yaml"
	logDir          = ".artifacts/api/integration"
	coverProfileOut = logDir + "/coverage.cov"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test orchestrator in short mode")
	}

	// Go tests run with cwd set to the package directory. Change to the module
	// root so that relative paths (configFile, go list) work correctly.
	require.NoError(t, chdirToModuleRoot(), "chdir to module root")

	// Create log directory and log files first, before anything else.
	require.NoError(t, os.MkdirAll(logDir, 0o755))

	// Orchestrator log file — allows: tail -f .artifacts/api/integration/orchestrator.log
	orchLogPath := filepath.Join(logDir, "orchestrator.log")
	orchLogFile, err := os.Create(orchLogPath)
	require.NoError(t, err, "create orchestrator log file")
	defer orchLogFile.Close()

	// Server log file — ZITADEL's slog/logrus output goes here.
	serverLogPath := filepath.Join(logDir, "server.log")
	serverLogFile, err := os.Create(serverLogPath)
	require.NoError(t, err, "create server log file")
	defer serverLogFile.Close()

	// testLogWriter routes output through t.Log for real-time streaming
	// with go test -v. Direct os.Stderr writes are buffered by Go's test
	// framework until the test returns, but t.Log is streamed immediately
	// (since Go 1.14). The orchLogFile is written as a secondary persistent
	// destination for CI artifact upload.
	logw := &testLogWriter{t: t, prefix: "[integration] ", file: orchLogFile}
	logf := func(format string, a ...any) {
		msg := fmt.Sprintf("[integration] "+format, a...)
		t.Log(msg)
		fmt.Fprintln(orchLogFile, msg)
	}

	absOrchLogPath, _ := filepath.Abs(orchLogPath)
	logf("Orchestrator log: %s (tail -f to follow)", absOrchLogPath)
	logf("Server log: %s", serverLogPath)

	// Intercept SIGTERM/SIGINT so deferred cleanup (container termination,
	// ZITADEL shutdown) runs instead of the process being killed immediately
	// by Nx or Ctrl-C.
	sigCtx, sigStop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer sigStop()

	// The orchestrator timeout must exceed the child go-test timeout (60m) so
	// that the parent context never expires before the child finishes.
	// Pattern: orchestrator = child timeout + grace period.
	ctx, cancel := context.WithTimeout(sigCtx, 65*time.Minute)
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
	require.NoError(t, waitForPortFree(ctx, "127.0.0.1:8082", 10*time.Second),
		"port 8082 is in use; a previous integration test run may still be running — "+
			"kill the leftover process and retry")

	// Redirect fd 2 (stderr) at the OS level to the server log file.
	// A Go-level os.Stderr swap is insufficient because logrus (used by
	// ZITADEL internals) captures os.Stderr at init time. By redirecting
	// the actual file descriptor, ALL stderr output — including logrus,
	// slog, and any C library writes — goes to the server log file.
	origStderrFd, err := syscall.Dup(syscall.Stderr)
	require.NoError(t, err, "dup stderr")
	require.NoError(t, syscall.Dup2(int(serverLogFile.Fd()), syscall.Stderr), "dup2 stderr")
	defer func() {
		_ = syscall.Dup2(origStderrFd, syscall.Stderr)
		_ = syscall.Close(origStderrFd)
	}()

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
	// INTEGRATION_PACKAGES overrides automatic discovery with a comma-separated
	// list of package paths (e.g. for quick local testing of a single package).
	var packages []string
	var skippedPkgs []string
	if pkgOverride := os.Getenv("INTEGRATION_PACKAGES"); pkgOverride != "" {
		packages = strings.Split(pkgOverride, ",")
		logf("Using %d packages from INTEGRATION_PACKAGES override", len(packages))
	} else {
		logf("Discovering integration test packages...")
		var err error
		packages, skippedPkgs, err = discoverIntegrationPackages()
		require.NoError(t, err, "discover packages")
		logf("Found %d integration test packages (%d v2beta packages skipped)", len(packages), len(skippedPkgs))
		for _, s := range skippedPkgs {
			logf("  skipping v2beta: %s", s)
		}
	}

	logf("Running integration tests...")
	// Allow overriding parallelism via INTEGRATION_PARALLELISM env var.
	// Default: 4. Each package creates an isolated ZITADEL instance
	// (random subdomain), so cross-package state is safe. Higher values reduce wall
	// time at the cost of DB/CPU load.
	parallelism := os.Getenv("INTEGRATION_PARALLELISM")
	if parallelism == "" {
		parallelism = "4"
	}
	logf("Package parallelism: %s (override with INTEGRATION_PARALLELISM)", parallelism)
	// Set up GOCOVERDIR so the child go-test process writes binary coverage
	// data. After the run we convert it to the text profile CI uploads.
	coverDir, err := os.MkdirTemp("", "integration-covdata-*")
	require.NoError(t, err, "create coverage dir")
	defer os.RemoveAll(coverDir)
	logf("Coverage data dir: %s", coverDir)

	args := []string{
		"test",
		"-count", "1",
		"-cover",
		"-tags", "integration",
		"-timeout", "60m",
		"-p", parallelism,
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

	configureChildCommand(cmd)

	// Route all test output through the single orchestrator log writer.
	cmd.Stdout = logw
	cmd.Stderr = logw
	// Inherit all env vars including our ZITADEL_* overrides.
	// Set GOCOVERDIR so instrumented test binaries write coverage data.
	cmd.Env = append(os.Environ(), "GOCOVERDIR="+coverDir)

	err = cmd.Run()
	require.NoError(t, err, "integration tests failed; see %s", absOrchLogPath)
	logf("All integration tests passed")

	// Convert binary coverage data to a text profile for CI upload.
	logf("Converting coverage data to %s...", coverProfileOut)
	covCmd := exec.CommandContext(ctx, "go", "tool", "covdata", "textfmt", "-i="+coverDir, "-o="+coverProfileOut)
	covCmd.Dir = "."
	covCmd.Stdout = logw
	covCmd.Stderr = logw
	if covErr := covCmd.Run(); covErr != nil {
		logf("WARNING: coverage conversion failed: %v (coverage file may be missing)", covErr)
	} else {
		logf("Coverage profile written to %s", coverProfileOut)
	}
}

// discoverIntegrationPackages finds all Go packages matching the integration test pattern.
// It lists the entire module (./...) so that integration-tagged packages outside internal/
// (e.g. backend/v3/storage/database/events_testing, cmd/setup/integration_test) are included.
// It returns the packages to run and the packages that were skipped (v2beta).
func discoverIntegrationPackages() (packages []string, skipped []string, err error) {
	cmd := exec.Command("go", "list", "-tags", "integration", "./...")
	cmd.Dir = "."
	out, err := cmd.Output()
	if err != nil {
		return nil, nil, fmt.Errorf("go list: %w", err)
	}

	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if strings.Contains(line, "integration_test") || strings.Contains(line, "events_testing") {
			// Skip v2beta packages — these APIs are being removed.
			if strings.Contains(line, "v2beta") {
				skipped = append(skipped, line)
				continue
			}
			packages = append(packages, line)
		}
	}

	if len(packages) == 0 {
		return nil, nil, fmt.Errorf("no integration test packages found")
	}
	return packages, skipped, nil
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

// testLogWriter implements io.Writer by routing each line through t.Log
// for real-time streaming with go test -v. Output is also written to a
// persistent log file for CI artifact upload.
type testLogWriter struct {
	t      *testing.T
	prefix string
	file   *os.File
}

func (w *testLogWriter) Write(p []byte) (n int, err error) {
	w.file.Write(p) //nolint:errcheck
	for _, line := range strings.Split(strings.TrimRight(string(p), "\n"), "\n") {
		if line != "" {
			w.t.Log(w.prefix + line)
		}
	}
	return len(p), nil
}
