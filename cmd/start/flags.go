package start

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/zitadel/zitadel/cmd/key"
	"github.com/zitadel/zitadel/cmd/tls"
)

var (
	tlsMode *string

	portFlag           = "port"
	externalDomainFlag = "externalDomain"
	externalPortFlag   = "externalPort"

	startFlagSet = &pflag.FlagSet{}
)

func init() {
	startFlagSet.Uint16(portFlag, viper.GetUint16(portFlag), "port to run ZITADEL on")
	startFlagSet.String(externalDomainFlag, viper.GetString(externalDomainFlag), "domain ZITADEL will be exposed on")
	startFlagSet.String(externalPortFlag, viper.GetString(externalPortFlag), "port ZITADEL will be exposed on")
}

func startFlags(cmd *cobra.Command) {
	cmd.Flags().AddFlagSet(startFlagSet)
	viper.BindPFlags(startFlagSet)

	tls.AddTLSModeFlag(cmd)
	key.AddMasterKeyFlag(cmd)
}
