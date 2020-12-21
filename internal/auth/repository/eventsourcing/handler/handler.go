package handler

import (
	"time"

	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/query"
	iam_events "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	org_events "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/types"
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

	es eventstore.Eventstore
}

func (h *handler) Eventstore() eventstore.Eventstore {
	return h.es
}

type EventstoreRepos struct {
	UserEvents    *usr_event.UserEventstore
	ProjectEvents *proj_event.ProjectEventstore
	OrgEvents     *org_events.OrgEventstore
	IamEvents     *iam_events.IAMEventstore
}

func Register(configs Configs, bulkLimit, errorCount uint64, view *view.View, es eventstore.Eventstore, repos EventstoreRepos, systemDefaults sd.SystemDefaults) []query.Handler {
	return []query.Handler{
		newUser(
			handler{view, bulkLimit, configs.cycleDuration("User"), errorCount, es},
			repos.OrgEvents,
			repos.IamEvents,
			systemDefaults.IamID),
		newUserSession(
			handler{view, bulkLimit, configs.cycleDuration("UserSession"), errorCount, es},
			repos.UserEvents),
		newUserMembership(
			handler{view, bulkLimit, configs.cycleDuration("UserMembership"), errorCount, es},
			repos.OrgEvents,
			repos.ProjectEvents),
		newToken(
			handler{view, bulkLimit, configs.cycleDuration("Token"), errorCount, es},
			repos.ProjectEvents),
		newKey(
			handler{view, bulkLimit, configs.cycleDuration("Key"), errorCount, es}),
		newApplication(handler{view, bulkLimit, configs.cycleDuration("Application"), errorCount, es},
			repos.ProjectEvents),
		newOrg(
			handler{view, bulkLimit, configs.cycleDuration("Org"), errorCount, es}),
		newUserGrant(
			handler{view, bulkLimit, configs.cycleDuration("UserGrant"), errorCount, es},
			repos.ProjectEvents,
			repos.UserEvents,
			repos.OrgEvents,
			repos.IamEvents,
			systemDefaults.IamID),
		newMachineKeys(
			handler{view, bulkLimit, configs.cycleDuration("MachineKey"), errorCount, es}),
		newLoginPolicy(
			handler{view, bulkLimit, configs.cycleDuration("LoginPolicy"), errorCount, es}),
		newIDPConfig(
			handler{view, bulkLimit, configs.cycleDuration("IDPConfig"), errorCount, es}),
		newIDPProvider(
			handler{view, bulkLimit, configs.cycleDuration("IDPProvider"), errorCount, es},
			systemDefaults,
			repos.IamEvents,
			repos.OrgEvents),
		newExternalIDP(
			handler{view, bulkLimit, configs.cycleDuration("ExternalIDP"), errorCount, es},
			systemDefaults,
			repos.IamEvents,
			repos.OrgEvents),
		newPasswordComplexityPolicy(
			handler{view, bulkLimit, configs.cycleDuration("PasswordComplexityPolicy"), errorCount, es}),
		newOrgIAMPolicy(
			handler{view, bulkLimit, configs.cycleDuration("OrgIAMPolicy"), errorCount, es}),
		newProjectRole(handler{view, bulkLimit, configs.cycleDuration("ProjectRole"), errorCount, es},
			repos.ProjectEvents),
	}
}

func (configs Configs) cycleDuration(viewModel string) time.Duration {
	c, ok := configs[viewModel]
	if !ok {
		return 3 * time.Minute
	}
	return c.MinimumCycleDuration.Duration
}

func (h *handler) MinimumCycleDuration() time.Duration {
	return h.cycleDuration
}

func (h *handler) QueryLimit() uint64 {
	return h.bulkLimit
}
