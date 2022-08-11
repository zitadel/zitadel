package cockroach

import (
	"github.com/zitadel/zitadel/internal/database/dialect"
)

func init() {
	config := &Config{}
	dialect.Register(config, config, true)
}
