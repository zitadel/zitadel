package cmd

import (
	"github.com/spf13/viper"
	"github.com/zitadel/zitadel/internal/secrets"
)

// LoadSecretsFromFiles processes Docker secrets from _FILE environment variables
// and loads them directly into Viper configuration.
func LoadSecretsFromFiles(v *viper.Viper) {
	secrets.ProcessDockerSecretsIntoViper(v)
}
