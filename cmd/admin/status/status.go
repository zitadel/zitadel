package status

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

const (
	flagKeyFile = "file"
)

type Config struct {
	Port uint16
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "check zitadel status",
	}
	cmd.AddCommand(checkHealth())
	return cmd
}

func checkHealth() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "check if zitadel is accepting requests",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := new(Config)
			if err := viper.Unmarshal(config); err != nil {
				return err
			}
			return isHealthy(config)
		},
	}
}

func isHealthy(cfg *Config) error {

	client := &http.Client{
		Timeout: 50 * time.Millisecond,
	}

	resp, err := client.Get(fmt.Sprintf("http://127.0.0.1:%d/debug/healthz", cfg.Port))
	if err != nil {
		return fmt.Errorf("zitadel is not healthy: %w", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("expected status code 200, but got %d", resp.StatusCode)
	}

	return nil
}
