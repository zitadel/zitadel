package setup

import (
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	org_event "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	usr_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type Setup struct {
	repos       EventstoreRepos
	iamID       string
	setUpConfig types.IAMSetUp
}

type EventstoreRepos struct {
	OrgEvents     *org_event.OrgEventstore
	UserEvents    *usr_event.UserEventstore
	ProjectEvents *proj_event.ProjectEventstore
}

func StartSetup(sd systemdefaults.SystemDefaults, repos EventstoreRepos) *Setup {
	return &Setup{
		repos:       repos,
		iamID:       sd.IamId,
		setUpConfig: sd.SetUp,
	}
}

func (s *Setup) Execute() error {
	return nil
}
