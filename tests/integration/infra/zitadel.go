//go:build integration

// Package infra provides test infrastructure helpers for integration tests.
// It uses testcontainers-go to manage Postgres and Redis containers and
// can start ZITADEL in-process for the integration test orchestrator.
package infra

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zitadel/zitadel/cmd"
	"github.com/zitadel/zitadel/cmd/initialise"
	"github.com/zitadel/zitadel/cmd/setup"
	"github.com/zitadel/zitadel/cmd/start"
)

// ZITADELInstance holds a running in-process ZITADEL server and its shutdown function.
type ZITADELInstance struct {
	Server   *start.Server
	cancel   context.CancelFunc
	shutdown chan error
}

// Shutdown gracefully stops the ZITADEL server.
func (z *ZITADELInstance) Shutdown(ctx context.Context) error {
	z.cancel()
	select {
	case err := <-z.shutdown:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// StartZITADEL performs the full init → setup → start lifecycle programmatically.
// configFiles are merged on top of cmd/defaults.yaml (same as --config flags).
// masterKey is the encryption master key (same as --masterkey flag).
// Container connection details should be set via environment variables before calling.
func StartZITADEL(ctx context.Context, configFiles []string, masterKey string) (*ZITADELInstance, error) {
	if err := setupViper(configFiles); err != nil {
		return nil, fmt.Errorf("setup viper: %w", err)
	}

	// Disable TLS for integration tests
	viper.Set("tls.enabled", false)
	viper.Set("externalSecure", false)
	viper.Set("InitProjections.Enabled", true)

	// Create a minimal cobra command to satisfy NewConfig signatures.
	// NewConfig only uses cmd.Context() and cmd.SetContext().
	// Named cobraCmd to avoid shadowing the imported "cmd" package.
	cobraCmd := &cobra.Command{}
	cobraCmd.SetContext(ctx)

	// Phase 1: Initialise (create DB schema, users, grants)
	initConfig, initShutdown, err := initialise.NewConfig(cobraCmd, viper.GetViper())
	if err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}
	defer func() {
		if err != nil {
			_ = initShutdown(ctx)
		}
	}()

	if err = initialise.InitAll(cobraCmd.Context(), initConfig); err != nil {
		return nil, fmt.Errorf("init all: %w", err)
	}
	_ = initShutdown(cobraCmd.Context())

	// Phase 2: Setup (run migrations, create first instance)
	setupConfig, setupShutdown, err := setup.NewConfig(cobraCmd, viper.GetViper())
	if err != nil {
		return nil, fmt.Errorf("setup config: %w", err)
	}
	defer func() {
		if err != nil {
			_ = setupShutdown(ctx)
		}
	}()

	setupSteps, err := setup.NewSteps(cobraCmd.Context(), viper.New())
	if err != nil {
		return nil, fmt.Errorf("setup steps: %w", err)
	}

	if err = setup.Setup(cobraCmd.Context(), setupConfig, setupSteps, masterKey); err != nil {
		return nil, fmt.Errorf("setup: %w", err)
	}
	_ = setupShutdown(cobraCmd.Context())

	// Phase 3: Start server
	startConfig, startShutdown, err := start.NewConfig(cobraCmd, viper.GetViper())
	if err != nil {
		return nil, fmt.Errorf("start config: %w", err)
	}

	serverCh := make(chan *start.Server, 1)
	serverCtx, cancel := context.WithCancel(cobraCmd.Context())

	shutdownCh := make(chan error, 1)
	go func() {
		err := start.StartZitadel(serverCtx, startConfig, masterKey, serverCh)
		_ = startShutdown(serverCtx)
		shutdownCh <- err
	}()

	// Wait for the server to be ready
	select {
	case server := <-serverCh:
		return &ZITADELInstance{
			Server:   server,
			cancel:   cancel,
			shutdown: shutdownCh,
		}, nil
	case err := <-shutdownCh:
		cancel()
		if err != nil {
			return nil, fmt.Errorf("start zitadel: %w", err)
		}
		return nil, errors.New("server exited before becoming ready")
	case <-ctx.Done():
		cancel()
		return nil, ctx.Err()
	}
}

// WaitForHealthy polls the ZITADEL health endpoint until it returns 200 or the context is cancelled.
func WaitForHealthy(ctx context.Context, baseURL string) error {
	endpoint := strings.TrimRight(baseURL, "/") + "/debug/ready"
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	client := &http.Client{Timeout: 5 * time.Second}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("health check timed out: %w", ctx.Err())
		case <-ticker.C:
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
			if err != nil {
				continue
			}
			resp, err := client.Do(req) //nolint:gosec
			if err != nil {
				continue
			}
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
	}
}

// setupViper configures the global viper instance identically to cmd/zitadel.go.
func setupViper(configFiles []string) error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("ZITADEL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("yaml")

	// Use the embedded defaults.yaml from the cmd package (avoids disk path dependency).
	if err := viper.ReadConfig(bytes.NewBuffer(cmd.DefaultConfig())); err != nil {
		return fmt.Errorf("parse defaults.yaml: %w", err)
	}

	// Merge additional config files (same as --config flags)
	for _, file := range configFiles {
		viper.SetConfigFile(file)
		if err := viper.MergeInConfig(); err != nil {
			return fmt.Errorf("merge config %s: %w", file, err)
		}
	}

	return nil
}
