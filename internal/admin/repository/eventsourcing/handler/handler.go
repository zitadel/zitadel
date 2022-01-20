package handler

import (
	"time"

	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	query2 "github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/static"
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

	es v1.Eventstore
}

func (h *handler) Eventstore() v1.Eventstore {
	return h.es
}

func Register(configs Configs, bulkLimit, errorCount uint64, view *view.View, es v1.Eventstore, defaults systemdefaults.SystemDefaults, command *command.Commands, queries *query2.Queries, static static.Storage, localDevMode bool) []query.Handler {
	handlers := []query.Handler{
		newUser(
			handler{view, bulkLimit, configs.cycleDuration("User"), errorCount, es},
			defaults, queries),
	}
	if static != nil {
		handlers = append(handlers, newStyling(
			handler{view, bulkLimit, configs.cycleDuration("Styling"), errorCount, es},
			static,
			localDevMode))
	}
	return handlers
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

func (h *handler) LockDuration() time.Duration {
	return h.cycleDuration / 3
}

func (h *handler) QueryLimit() uint64 {
	return h.bulkLimit
}
