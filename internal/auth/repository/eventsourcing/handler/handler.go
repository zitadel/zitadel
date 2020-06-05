package handler

import (
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/eventstore"
	iam_events "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	org_events "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	"time"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	usr_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type Configs map[string]*Config

type Config struct {
	MinimumCycleDuration types.Duration
}

type handler struct {
	view                *view.View
	bulkLimit           uint64
	cycleDuration       time.Duration
	errorCountUntilSkip uint64
}

type EventstoreRepos struct {
	UserEvents    *usr_event.UserEventstore
	ProjectEvents *proj_event.ProjectEventstore
	OrgEvents     *org_events.OrgEventstore
	IamEvents     *iam_events.IamEventstore
}

func Register(configs Configs, bulkLimit, errorCount uint64, view *view.View, eventstore eventstore.Eventstore, repos EventstoreRepos, systemDefaults sd.SystemDefaults) []spooler.Handler {
	return []spooler.Handler{
		&User{handler: handler{view, bulkLimit, configs.cycleDuration("User"), errorCount}},
		&UserSession{handler: handler{view, bulkLimit, configs.cycleDuration("UserSession"), errorCount}, userEvents: repos.UserEvents},
		&Token{handler: handler{view, bulkLimit, configs.cycleDuration("Token"), errorCount}},
		&Key{handler: handler{view, bulkLimit, configs.cycleDuration("Key"), errorCount}},
		&Application{handler: handler{view, bulkLimit, configs.cycleDuration("Application"), errorCount}},
		&Org{handler: handler{view, bulkLimit, configs.cycleDuration("Org"), errorCount}},
		&UserGrant{
			handler:       handler{view, bulkLimit, configs.cycleDuration("UserGrant"), errorCount},
			eventstore:    eventstore,
			userEvents:    repos.UserEvents,
			orgEvents:     repos.OrgEvents,
			projectEvents: repos.ProjectEvents,
			iamEvents:     repos.IamEvents,
			iamID:         systemDefaults.IamID},
	}
}

func (configs Configs) cycleDuration(viewModel string) time.Duration {
	c, ok := configs[viewModel]
	if !ok {
		return 1 * time.Second
	}
	return c.MinimumCycleDuration.Duration
}
