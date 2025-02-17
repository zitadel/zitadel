package start

import (
	"github.com/spf13/viper"

	"github.com/zitadel/zitadel/backend/storage/database/dialect"
)

type Config struct {
	Database dialect.Config
}

func (c Config) Hooks() []viper.DecoderConfigOption {
	return c.Database.Hooks()
}
