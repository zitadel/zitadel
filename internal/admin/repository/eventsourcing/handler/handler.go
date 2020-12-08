package handler

import (
	"time"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	iam_event "github.com/caos/zitadel/internal/iam/repository/eventsourcing"

	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/query"
	org_event "github.com/caos/zitadel/internal/org/repository/eventsourcing"
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
	UserEvents *usr_event.UserEventstore
	IamEvents  *iam_event.IAMEventstore
	OrgEvents  *org_event.OrgEventstore
}

func Register(configs Configs, bulkLimit, errorCount uint64, view *view.View, eventstore eventstore.Eventstore, repos EventstoreRepos, defaults systemdefaults.SystemDefaults) []query.Handler {
	return []query.Handler{
		&Org{handler: handler{view, bulkLimit, configs.cycleDuration("Org"), errorCount}},
		&IamMember{handler: handler{view, bulkLimit, configs.cycleDuration("IamMember"), errorCount},
			userEvents: repos.UserEvents},
		&IDPConfig{handler: handler{view, bulkLimit, configs.cycleDuration("IDPConfig"), errorCount}},
		&LabelPolicy{handler: handler{view, bulkLimit, configs.cycleDuration("LabelPolicy"), errorCount}},
		&LoginPolicy{handler: handler{view, bulkLimit, configs.cycleDuration("LoginPolicy"), errorCount}},
		&IDPProvider{handler: handler{view, bulkLimit, configs.cycleDuration("LoginPolicy"), errorCount},
			systemDefaults: defaults, iamEvents: repos.IamEvents, orgEvents: repos.OrgEvents},
		&User{handler: handler{view, bulkLimit, configs.cycleDuration("User"), errorCount},
			eventstore: eventstore, orgEvents: repos.OrgEvents, iamEvents: repos.IamEvents, systemDefaults: defaults},
		&PasswordComplexityPolicy{handler: handler{view, bulkLimit, configs.cycleDuration("PasswordComplexityPolicy"), errorCount}},
		&PasswordAgePolicy{handler: handler{view, bulkLimit, configs.cycleDuration("PasswordAgePolicy"), errorCount}},
		&PasswordLockoutPolicy{handler: handler{view, bulkLimit, configs.cycleDuration("PasswordLockoutPolicy"), errorCount}},
		&OrgIAMPolicy{handler: handler{view, bulkLimit, configs.cycleDuration("OrgIAMPolicy"), errorCount}},
		&ExternalIDP{handler: handler{view, bulkLimit, configs.cycleDuration("User"), errorCount},
			orgEvents: repos.OrgEvents, iamEvents: repos.IamEvents, systemDefaults: defaults},
	}
}

func (configs Configs) cycleDuration(viewModel string) time.Duration {
	c, ok := configs[viewModel]
	if !ok {
		return 2 * time.Second
	}
	return c.MinimumCycleDuration.Duration
}

func (h *handler) MinimumCycleDuration() time.Duration {
	return h.cycleDuration
}

func (h *handler) QueryLimit() uint64 {
	return h.bulkLimit
}
