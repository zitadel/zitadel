package initialise

import (
	"github.com/spf13/viper"
	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/id"
)

type Config struct {
	Database  database.Config
	AdminUser database.User
	Machine   *id.Config
	Log       *logging.Config
}

func MustNewConfig(v *viper.Viper) *Config {
	config := new(Config)
	err := v.Unmarshal(config)
	logging.OnError(err).Fatal("unable to read config")

	err = config.Log.SetLogger()
	logging.OnError(err).Fatal("unable to set logger")

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
