package ready

import (
	"crypto/tls"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/socket"
	"net"
	"net/http"
	"os"
	"strconv"
)

func New() *cobra.Command {
	return &cobra.Command{
		Use:   "ready",
		Short: "Checks if zitadel is ready",
		Long:  "Checks if zitadel is ready",
		Run: func(cmd *cobra.Command, args []string) {
			config := MustNewConfig(viper.GetViper())
			if !ready(config) {
				os.Exit(1)
			}
		},
	}
}

func ready(config *Config) bool {
	explicitErr := tryToCheckExplicitly(config)
	if explicitErr == nil {
		logging.Info("ready check passed")
		return true
	}
	socketErr := expectTrueFromSocket(socket.ReadinessQuery)
	if socketErr == nil {
		logging.Info("ready check passed")
		return true
	}
	logging.Warnf("ready check failed: %v", explicitErr)
	logging.Warnf("ready check failed: %v", socketErr)
	return false
}

func expectTrueFromSocket(query socket.SocketRequest) error {
	resp, err := query.Request()
	if err != nil {
		return fmt.Errorf("socket request error: %w", err)
	}
	if resp != socket.True {
		return fmt.Errorf("zitadel process did not respond true to a readiness query")
	}
	return nil
}

func tryToCheckExplicitly(config *Config) error {
	scheme := "https"
	if !config.TLS.Enabled {
		scheme = "http"
	}
	// Checking the TLS cert is not in the scope of the readiness check
	httpClient := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	res, err := httpClient.Get(scheme + "://" + net.JoinHostPort("localhost", strconv.Itoa(int(config.Port))) + "/debug/ready")
	if err != nil {
		return fmt.Errorf("url error: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil
}
