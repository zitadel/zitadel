package initialise

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/database"
	"github.com/spf13/viper"
)

type Config struct {
	Database  database.Config
	AdminUser database.User
}

func MustNewConfig(v *viper.Viper) *Config {
	config := new(Config)
	err := viper.Unmarshal(config)
	logging.OnError(err).Fatal("unable to read config")
	return config
}

func adminConfig(config *Config) database.Config {
	adminConfig := config.Database
	adminConfig.Username = config.AdminUser.Username
	adminConfig.Password = config.AdminUser.Password
	adminConfig.SSL.Cert = config.AdminUser.SSL.Cert
	adminConfig.SSL.Key = config.AdminUser.SSL.Key
	if config.AdminUser.SSL.RootCert != "" {
		adminConfig.SSL.RootCert = config.AdminUser.SSL.RootCert
	}
	if config.AdminUser.SSL.Mode != "" {
		adminConfig.SSL.Mode = config.AdminUser.SSL.Mode
	}
	//use default database because the zitadel database might not exist
	adminConfig.Database = ""

	return adminConfig
}
