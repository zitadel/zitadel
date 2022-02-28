package setup

import (
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/database"
)

type Config struct {
	AdminOrg       AdminOrg
	Database       database.Config
	SystemDefaults systemdefaults.SystemDefaults
	InternalAuthZ  authz.Config
}
