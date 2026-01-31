package ready

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
)

func New() *cobra.Command {
	return &cobra.Command{
		Use:   "ready",
		Short: "Checks if zitadel is ready",
		Long:  "Checks if zitadel is ready",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Overwrite context with ready stream for logging
			cmd.SetContext(logging.NewCtx(cmd.Context(), logging.StreamReady))
			defer func() {
				logging.OnError(cmd.Context(), err).Error("zitadel ready command failed")
			}()
			config, shutdown, err := newConfig(cmd.Context(), viper.GetViper())
			if err != nil {
				return err
			}
			// Set logger again to include changes from config
			cmd.SetContext(logging.NewCtx(cmd.Context(), logging.StreamReady))
			defer shutdown(cmd.Context())
			if ready(cmd.Context(), config) {
				return nil
			}
			return errors.New("not ready")
		},
	}
}

func ready(ctx context.Context, config *Config) bool {
	scheme := "https"
	if !config.TLS.Enabled {
		scheme = "http"
	}
	// Checking the TLS cert is not in the scope of the readiness check
	httpClient := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	res, err := httpClient.Get(scheme + "://" + net.JoinHostPort("localhost", strconv.Itoa(int(config.Port))) + "/debug/ready")
	if err != nil {
		logging.WithError(ctx, err).Warn("get request failed")
		return false
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		logging.Warn(ctx, "get request failed", "status", res.StatusCode)
		return false
	}
	return true
}
