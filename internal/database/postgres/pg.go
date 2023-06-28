package postgres

import (

	//sql import
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/zitadel/zitadel/internal/database/dialect"
)

func init() {
	config := &Config{}
	dialect.Register(config, config, false)
}
