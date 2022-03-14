package setup

import (
	"github.com/caos/zitadel/internal/database"
)

type Config struct {
	Database database.Config
}

type Steps struct {
	S1ProjectionTable *ProjectionTable
}
