package setup

import (
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/database"
)

type Config struct {
	Database       database.Config
	SystemDefaults systemdefaults.SystemDefaults
	InternalAuthZ  authz.Config
}

type Steps struct {
	S1DefaultInstance *DefaultInstance
}
