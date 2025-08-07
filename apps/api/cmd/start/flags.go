package start

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/key"
	"github.com/zitadel/zitadel/cmd/tls"
)

var (
	startFlagSet = &pflag.FlagSet{}
)

func init() {
	startFlagSet.Uint16("port", 0, "port to run ZITADEL on")
	startFlagSet.String("externalDomain", "", "domain ZITADEL will be exposed on")
	startFlagSet.String("externalPort", "", "port ZITADEL will be exposed on")
}

func startFlags(cmd *cobra.Command) {
	cmd.Flags().AddFlagSet(startFlagSet)
	logging.OnError(
		viper.BindPFlags(startFlagSet),
	).Fatal("start flags")

	tls.AddTLSModeFlag(cmd)
	key.AddMasterKeyFlag(cmd)
}
