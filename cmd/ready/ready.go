package ready

import (
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"
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
	res, err := http.Get("http://" + net.JoinHostPort("localhost", strconv.Itoa(int(config.Port))) + "/debug/ready")
	if err != nil {
		logging.WithError(err).Warn("ready check failed")
		return false
	}
	defer res.Body.Close()
	logging.WithFields("status", res.StatusCode).Warn("ready check failed")
	return res.StatusCode == 200
}
