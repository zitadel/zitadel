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
	es                  eventstore.Eventstore
}

type EventstoreRepos struct {
	UserEvents *usr_event.UserEventstore
	IamEvents  *iam_event.IAMEventstore
	OrgEvents  *org_event.OrgEventstore
}

func Register(configs Configs, bulkLimit, errorCount uint64, view *view.View, es eventstore.Eventstore, repos EventstoreRepos, defaults systemdefaults.SystemDefaults) []query.Handler {
	return []query.Handler{
		&Org{handler: handler{view, bulkLimit, configs.cycleDuration("Org"), errorCount, es}},
		&IAMMember{handler: handler{view, bulkLimit, configs.cycleDuration("IamMember"), errorCount, es},
			userEvents: repos.UserEvents},
		&IDPConfig{handler: handler{view, bulkLimit, configs.cycleDuration("IDPConfig"), errorCount, es}},
		&LabelPolicy{handler: handler{view, bulkLimit, configs.cycleDuration("LabelPolicy"), errorCount, es}},
		&LoginPolicy{handler: handler{view, bulkLimit, configs.cycleDuration("LoginPolicy"), errorCount, es}},
		&IDPProvider{handler: handler{view, bulkLimit, configs.cycleDuration("IDPProvider"), errorCount, es},
			systemDefaults: defaults, iamEvents: repos.IamEvents, orgEvents: repos.OrgEvents},
		&User{handler: handler{view, bulkLimit, configs.cycleDuration("User"), errorCount, es},
			eventstore: es, orgEvents: repos.OrgEvents, iamEvents: repos.IamEvents, systemDefaults: defaults},
		&PasswordComplexityPolicy{handler: handler{view, bulkLimit, configs.cycleDuration("PasswordComplexityPolicy"), errorCount, es}},
		&PasswordAgePolicy{handler: handler{view, bulkLimit, configs.cycleDuration("PasswordAgePolicy"), errorCount, es}},
		&PasswordLockoutPolicy{handler: handler{view, bulkLimit, configs.cycleDuration("PasswordLockoutPolicy"), errorCount, es}},
		&OrgIAMPolicy{handler: handler{view, bulkLimit, configs.cycleDuration("OrgIAMPolicy"), errorCount, es}},
		&ExternalIDP{handler: handler{view, bulkLimit, configs.cycleDuration("ExternalIDP"), errorCount, es},
			orgEvents: repos.OrgEvents, iamEvents: repos.IamEvents, systemDefaults: defaults},
	}
}

func subscribe(es eventstore.Eventstore, handlers []query.Handler) {
	for _, handler := range handlers {
		es.Subscribe(handler.AggregateTypes()...)
	}
}

func (configs Configs) cycleDuration(viewModel string) time.Duration {
	c, ok := configs[viewModel]
	if !ok {
		return 1 * time.Second
	}
	return c.MinimumCycleDuration.Duration
}

func (h *handler) MinimumCycleDuration() time.Duration {
	return h.cycleDuration
}

func (h *handler) QueryLimit() uint64 {
	return h.bulkLimit
}
