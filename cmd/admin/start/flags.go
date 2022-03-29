package start

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func startFlags(cmd *cobra.Command) {
	bindUint16Flag(cmd, "port", "port to run ZITADEL on")
	bindStringFlag(cmd, "externalDomain", "domain ZITADEL will be exposed on")
	bindStringFlag(cmd, "externalPort", "port ZITADEL will be exposed on")
	bindBoolFlag(cmd, "externalSecure", "if ZITADEL will be served on HTTPS")

	cmd.PersistentFlags().String(flagMasterKey, "", "masterkey for en/decryption keys")

}

func bindStringFlag(cmd *cobra.Command, name, description string) {
	cmd.PersistentFlags().String(name, viper.GetString(name), description)
	viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
}

func bindUint16Flag(cmd *cobra.Command, name, description string) {
	cmd.PersistentFlags().Uint16(name, uint16(viper.GetUint(name)), description)
	viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
}

func bindBoolFlag(cmd *cobra.Command, name, description string) {
	cmd.PersistentFlags().Bool(name, viper.GetBool(name), description)
	viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
}
